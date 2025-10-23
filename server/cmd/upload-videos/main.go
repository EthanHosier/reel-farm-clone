package main

import (
	"context"
	"crypto/md5"

	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ethanhosier/reel-farm/db"
	"github.com/ethanhosier/reel-farm/internal/repository"
	"github.com/ethanhosier/reel-farm/internal/service"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type VideoUploader struct {
	s3Client   *s3.Client
	uploader   *manager.Uploader
	service    *service.AIAvatarService
	bucketName string
	tempDir    string
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Hardcoded bucket name from Terraform output
	bucketName := "reel-farm-bucket-an04kbe3"

	// Initialize uploader
	uploader, err := NewVideoUploader(bucketName)
	if err != nil {
		log.Fatalf("Failed to initialize uploader: %v", err)
	}
	defer uploader.Cleanup()

	// Find video files in videos/ directory
	videoDir := "videos"
	videoFiles, err := findVideoFiles(videoDir)
	if err != nil {
		log.Fatalf("Failed to find video files: %v", err)
	}

	if len(videoFiles) == 0 {
		fmt.Printf("❌ No video files found in directory: %s\n", videoDir)
		fmt.Println("Supported formats: .mp4, .mov, .avi, .mkv, .webm")
		fmt.Println("Create a 'videos' folder and add your video files there.")
		os.Exit(1)
	}

	fmt.Printf("🎬 Found %d video file(s) in %s/\n", len(videoFiles), videoDir)
	fmt.Printf("🪣 Bucket: %s\n", bucketName)
	fmt.Println()

	// Process each video file
	var successCount, errorCount int
	for i, videoFile := range videoFiles {
		filename := filepath.Base(videoFile)
		fmt.Printf("📹 Processing video %d/%d: %s\n", i+1, len(videoFiles), filename)

		err := uploader.UploadVideo(videoFile)
		if err != nil {
			log.Printf("❌ Failed to upload %s: %v", filename, err)
			errorCount++
			continue
		}

		fmt.Printf("✅ Successfully uploaded: %s\n", filename)
		successCount++
		fmt.Println()
	}

	// Summary
	fmt.Printf("📊 Upload Summary:\n")
	fmt.Printf("   ✅ Successful: %d\n", successCount)
	fmt.Printf("   ❌ Failed: %d\n", errorCount)
	fmt.Printf("   📁 Total: %d\n", len(videoFiles))

	if errorCount > 0 {
		fmt.Printf("\n⚠️  Some videos failed to upload. Check the errors above.\n")
		os.Exit(1)
	} else {
		fmt.Println("\n🎉 All videos processed successfully!")
	}
}

func NewVideoUploader(bucketName string) (*VideoUploader, error) {
	// Load AWS config with us-west-2 region
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client and uploader
	s3Client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(s3Client)

	// Connect to database
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create temp directory for thumbnails
	tempDir, err := os.MkdirTemp("", "video-upload-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	s, err := service.NewAIAvatarService(repository.NewAIAvatarRepository(pool), bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to create AIAvatarService: %w", err)
	}

	return &VideoUploader{
		s3Client:   s3Client,
		uploader:   uploader,
		service:    s,
		bucketName: bucketName,
		tempDir:    tempDir,
	}, nil
}

func (u *VideoUploader) Cleanup() {
	if u.tempDir != "" {
		os.RemoveAll(u.tempDir)
	}
}

func findVideoFiles(dir string) ([]string, error) {
	var videoFiles []string

	extensions := []string{".mp4", ".mov", ".avi", ".mkv", ".webm", ".m4v"}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		for _, videoExt := range extensions {
			if ext == videoExt {
				videoFiles = append(videoFiles, path)
				break
			}
		}

		return nil
	})

	return videoFiles, err
}

func (u *VideoUploader) UploadVideo(videoPath string) error {
	filename := filepath.Base(videoPath)
	title := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Generate deterministic UUID from filename hash
	hash := md5.Sum([]byte(filename))
	// Convert hash to UUID format
	id := uuid.Must(uuid.FromBytes(hash[:]))

	// Check if video already exists in database
	exists, err := u.service.VideoExists(context.Background(), id)
	if err != nil {
		return fmt.Errorf("failed to check if video exists: %w", err)
	}
	if exists {
		fmt.Printf("   ⏭️  Video already exists in database: %s\n", title)
		return nil
	}

	// Get file info
	fileInfo, err := os.Stat(videoPath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Extract thumbnail using ffmpeg
	thumbnailPath, duration, err := u.extractThumbnail(videoPath)
	if err != nil {
		return fmt.Errorf("failed to extract thumbnail: %w", err)
	}

	// Generate S3 filenames using the UUID
	videoFilename := fmt.Sprintf("%s.mp4", id.String())
	thumbnailFilename := fmt.Sprintf("%s.jpg", id.String())

	// Upload video to S3
	videoKey := fmt.Sprintf("ai-avatar/videos/%s", videoFilename)
	err = u.uploadFile(videoPath, videoKey)
	if err != nil {
		return fmt.Errorf("failed to upload video: %w", err)
	}

	// Upload thumbnail to S3
	thumbnailKey := fmt.Sprintf("ai-avatar/thumbnails/%s", thumbnailFilename)
	err = u.uploadFile(thumbnailPath, thumbnailKey)
	if err != nil {
		return fmt.Errorf("failed to upload thumbnail: %w", err)
	}

	// Save to database
	fileSize := fileInfo.Size()
	_, err = u.service.CreateVideo(context.Background(), &db.CreateVideoParams{
		ID:                id,
		Title:             title,
		Description:       nil, // No description for now
		Filename:          videoFilename,
		ThumbnailFilename: thumbnailFilename,
		Duration:          &duration,
		FileSize:          &fileSize,
	})
	if err != nil {
		return fmt.Errorf("failed to save video to database: %w", err)
	}

	return nil
}

func (u *VideoUploader) extractThumbnail(videoPath string) (string, int32, error) {
	// Generate thumbnail filename
	thumbnailPath := filepath.Join(u.tempDir, fmt.Sprintf("thumb_%s.jpg", filepath.Base(videoPath)))

	// Extract thumbnail using ffmpeg
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", "00:00:01", // Extract frame at 1 second
		"-vframes", "1",
		"-q:v", "2", // High quality
		"-y", // Overwrite output file
		thumbnailPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", 0, fmt.Errorf("ffmpeg failed: %s, error: %w", string(output), err)
	}

	// Get video duration using ffprobe
	durationCmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-show_entries", "format=duration",
		"-of", "csv=p=0",
		videoPath,
	)

	durationOutput, err := durationCmd.Output()
	if err != nil {
		return "", 0, fmt.Errorf("ffprobe failed: %w", err)
	}

	durationStr := strings.TrimSpace(string(durationOutput))
	durationFloat, err := strconv.ParseFloat(durationStr, 32)
	if err != nil {
		return "", 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	duration := int32(durationFloat)

	return thumbnailPath, duration, nil
}

func (u *VideoUploader) uploadFile(filePath, key string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = u.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(u.bucketName),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}
