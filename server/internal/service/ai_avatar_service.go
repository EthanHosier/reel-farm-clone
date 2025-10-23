package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ethanhosier/reel-farm/db"
	"github.com/ethanhosier/reel-farm/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AIAvatarService struct {
	repo             *repository.AIAvatarRepository
	s3Client         *s3.Client
	uploader         *manager.Uploader
	bucketName       string
	tempDir          string
	cloudfrontDomain string
	cloudfrontSigner *sign.URLSigner
}

func NewAIAvatarService(repo *repository.AIAvatarRepository, bucketName string) (*AIAvatarService, error) {
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(s3Client)

	// Create temp directory
	tempDir := "/tmp/video-processing"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Get CloudFront configuration from environment
	cloudfrontDomain := os.Getenv("CLOUDFRONT_DOMAIN")
	cloudfrontKeyPairID := os.Getenv("CLOUDFRONT_KEY_PAIR_ID")
	cloudfrontPrivateKey := os.Getenv("CLOUDFRONT_PRIVATE_KEY")

	if cloudfrontDomain == "" || cloudfrontKeyPairID == "" || cloudfrontPrivateKey == "" {
		return nil, fmt.Errorf("CloudFront configuration missing: CLOUDFRONT_DOMAIN, CLOUDFRONT_KEY_PAIR_ID, and CLOUDFRONT_PRIVATE_KEY are required")
	}

	// Parse the private key
	privKey, err := sign.LoadPEMPrivKey(strings.NewReader(cloudfrontPrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse CloudFront private key: %w", err)
	}

	// Create CloudFront signer
	cloudfrontSigner := sign.NewURLSigner(cloudfrontKeyPairID, privKey)

	return &AIAvatarService{
		repo:             repo,
		s3Client:         s3Client,
		uploader:         uploader,
		bucketName:       bucketName,
		tempDir:          tempDir,
		cloudfrontDomain: cloudfrontDomain,
		cloudfrontSigner: cloudfrontSigner,
	}, nil
}

// GetAllVideos retrieves all AI avatar videos
func (s *AIAvatarService) GetAllVideos(ctx context.Context) ([]*db.AiAvatarVideo, error) {
	return s.repo.GetAllVideos(ctx)
}

// GetVideoByID retrieves a specific video by ID
func (s *AIAvatarService) GetVideoByID(ctx context.Context, id uuid.UUID) (*db.AiAvatarVideo, error) {
	return s.repo.GetVideoByID(ctx, id)
}

// CreateVideo creates a new video record
func (s *AIAvatarService) CreateVideo(ctx context.Context, params *db.CreateVideoParams) (*db.AiAvatarVideo, error) {
	return s.repo.CreateVideo(ctx, params)
}

// UpdateVideo updates an existing video
func (s *AIAvatarService) UpdateVideo(ctx context.Context, params *db.UpdateVideoParams) (*db.AiAvatarVideo, error) {
	return s.repo.UpdateVideo(ctx, params)
}

// DeleteVideo deletes a video by ID
func (s *AIAvatarService) DeleteVideo(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteVideo(ctx, id)
}

// VideoExists checks if a video exists by ID
func (s *AIAvatarService) VideoExists(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repo.VideoExists(ctx, id)
}

// GetUserGeneratedVideosByUserID retrieves all user-generated videos for a specific user
func (s *AIAvatarService) GetUserGeneratedVideosByUserID(ctx context.Context, userID uuid.UUID) ([]*db.UserGeneratedVideo, error) {
	return s.repo.GetUserGeneratedVideosByUserID(ctx, userID)
}

// ProcessVideoWithTextOverlay downloads a video, adds text overlay, and uploads the result
func (s *AIAvatarService) ProcessVideoWithTextOverlay(ctx context.Context, userID, aiAvatarVideoID uuid.UUID, videoURL, overlayText string) (*db.UserGeneratedVideo, error) {
	// Generate unique filenames
	videoID := uuid.New()
	videoFilename := fmt.Sprintf("%s.mp4", videoID.String())
	thumbnailFilename := fmt.Sprintf("%s.jpg", videoID.String())

	// Download the original video
	originalVideoPath := filepath.Join(s.tempDir, fmt.Sprintf("original_%s.mp4", videoID.String()))
	if err := s.downloadVideo(ctx, videoURL, originalVideoPath); err != nil {
		return nil, fmt.Errorf("failed to download video: %w", err)
	}
	defer os.Remove(originalVideoPath)

	// Process video with text overlay
	processedVideoPath := filepath.Join(s.tempDir, videoFilename)
	if err := s.addTextOverlay(originalVideoPath, overlayText, processedVideoPath); err != nil {
		return nil, fmt.Errorf("failed to add text overlay: %w", err)
	}
	defer os.Remove(processedVideoPath)

	// Extract thumbnail
	thumbnailPath := filepath.Join(s.tempDir, thumbnailFilename)
	if err := s.extractThumbnail(processedVideoPath, thumbnailPath); err != nil {
		return nil, fmt.Errorf("failed to extract thumbnail: %w", err)
	}
	defer os.Remove(thumbnailPath)

	// Upload processed video to S3
	videoKey := fmt.Sprintf("user-generated-videos/videos/%s", videoFilename)
	if err := s.uploadFile(processedVideoPath, videoKey); err != nil {
		return nil, fmt.Errorf("failed to upload processed video: %w", err)
	}

	// Upload thumbnail to S3
	thumbnailKey := fmt.Sprintf("user-generated-videos/thumbnails/%s", thumbnailFilename)
	if err := s.uploadFile(thumbnailPath, thumbnailKey); err != nil {
		return nil, fmt.Errorf("failed to upload thumbnail: %w", err)
	}

	// Create database record
	status := "completed"
	userGeneratedVideo, err := s.repo.CreateUserGeneratedVideo(ctx, &db.CreateUserGeneratedVideoParams{
		ID:                     videoID,
		UserID:                 pgtype.UUID{Bytes: userID, Valid: true},
		AiAvatarVideoID:        pgtype.UUID{Bytes: aiAvatarVideoID, Valid: true},
		OverlayText:            overlayText,
		GeneratedVideoFilename: videoFilename,
		ThumbnailFilename:      thumbnailFilename,
		Status:                 &status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create database record: %w", err)
	}

	return userGeneratedVideo, nil
}

// downloadVideo downloads a video from URL to local path using Go HTTP client
func (s *AIAvatarService) downloadVideo(ctx context.Context, url, outputPath string) error {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download video: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download video: HTTP %d", resp.StatusCode)
	}

	// Create output file
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Copy response body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write video file: %w", err)
	}

	return nil
}

// addTextOverlay adds text overlay to video using FFmpeg
func (s *AIAvatarService) addTextOverlay(inputPath, text, outputPath string) error {
	// Wrap text if it's too long (approximately 35 characters per line for 36px font)
	wrappedLines := s.wrapTextToLines(text, 35)

	// Create a temporary text file with the wrapped text
	tempTextFile := filepath.Join(s.tempDir, fmt.Sprintf("text_%d.txt", time.Now().UnixNano()))
	joinedText := strings.Join(wrappedLines, "\n")
	err := os.WriteFile(tempTextFile, []byte(joinedText), 0644)
	if err != nil {
		return fmt.Errorf("failed to create temporary text file: %w", err)
	}
	defer os.Remove(tempTextFile)

	// Check if font file exists
	fontPath := "./TikTokDisplay-Medium.ttf"
	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		return fmt.Errorf("font file not found at %s", fontPath)
	}

	// FFmpeg command to add text overlay
	videoFilter := fmt.Sprintf("drawtext=textfile=%s:fontfile=%s:fontsize=36:fontcolor=white:x=(w-text_w)/2:y=(h-text_h)/2:borderw=3:bordercolor=black:text_align=center:line_spacing=16", tempTextFile, fontPath)

	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-vf", videoFilter,
		"-c:a", "copy", // Copy audio without re-encoding
		"-y", // Overwrite output file if it exists
		outputPath,
	)

	// Stream FFmpeg output in real time
	log.Printf("🎬 Starting FFmpeg processing...")
	log.Printf("📝 Command: %s", cmd.String())

	// Create pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Stream stdout and stderr in real time
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			log.Printf("FFmpeg stdout: %s", scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			log.Printf("FFmpeg stderr: %s", scanner.Text())
		}
	}()

	// Wait for the command to complete
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	log.Printf("✅ FFmpeg processing completed successfully")
	return nil
}

// extractThumbnail extracts thumbnail from video
func (s *AIAvatarService) extractThumbnail(videoPath, thumbnailPath string) error {
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", "00:00:01", // Extract frame at 1 second
		"-vframes", "1",
		"-q:v", "2", // High quality
		"-y", // Overwrite output file if it exists
		thumbnailPath,
	)

	// Stream FFmpeg output in real time
	log.Printf("🖼️ Starting thumbnail extraction...")
	log.Printf("📝 Command: %s", cmd.String())

	// Create pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Stream stdout and stderr in real time
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			log.Printf("FFmpeg thumbnail stdout: %s", scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			log.Printf("FFmpeg thumbnail stderr: %s", scanner.Text())
		}
	}()

	// Wait for the command to complete
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("ffmpeg thumbnail extraction failed: %w", err)
	}

	log.Printf("✅ Thumbnail extraction completed successfully")
	return nil
}

// uploadFile uploads a file to S3
func (s *AIAvatarService) uploadFile(filePath, key string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = s.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

// wrapTextToLines wraps text to fit within a specified number of characters per line
func (s *AIAvatarService) wrapTextToLines(text string, maxCharsPerLine int) []string {
	words := strings.Fields(text)
	var lines []string
	var currentLine string

	for _, word := range words {
		if len(currentLine)+len(word)+1 <= maxCharsPerLine {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word // Single space between words
			}
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = word
			} else {
				// Word is longer than maxCharsPerLine, add it anyway
				lines = append(lines, word)
			}
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// GenerateSignedURL creates a signed CloudFront URL for user-generated videos
func (s *AIAvatarService) GenerateSignedURL(path string, expiresIn time.Duration) (string, error) {
	// Create the base CloudFront URL
	baseURL := fmt.Sprintf("https://%s/%s", s.cloudfrontDomain, path)

	// Sign the URL with expiration time
	expiresAt := time.Now().Add(expiresIn)
	signedURL, err := s.cloudfrontSigner.Sign(baseURL, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to sign CloudFront URL: %w", err)
	}

	return signedURL, nil
}
