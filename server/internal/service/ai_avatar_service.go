package service

import (
	"context"
	"fmt"

	"github.com/ethanhosier/reel-farm/db"
	"github.com/ethanhosier/reel-farm/internal/repository"
	"github.com/google/uuid"
)

type AIAvatarService struct {
	repo *repository.AIAvatarRepository
}

func NewAIAvatarService(repo *repository.AIAvatarRepository) *AIAvatarService {
	return &AIAvatarService{
		repo: repo,
	}
}

// GetAllVideos retrieves all AI avatar videos
func (s *AIAvatarService) GetAllVideos(ctx context.Context) ([]*db.Video, error) {
	return s.repo.GetAllVideos(ctx)
}

// GetVideoByID retrieves a specific video by ID
func (s *AIAvatarService) GetVideoByID(ctx context.Context, id uuid.UUID) (*db.Video, error) {
	return s.repo.GetVideoByID(ctx, id)
}

// CreateVideo creates a new video record
func (s *AIAvatarService) CreateVideo(ctx context.Context, params *db.CreateVideoParams) (*db.Video, error) {
	return s.repo.CreateVideo(ctx, params)
}

// UpdateVideo updates an existing video
func (s *AIAvatarService) UpdateVideo(ctx context.Context, params *db.UpdateVideoParams) (*db.Video, error) {
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

// GetVideoWithURLs retrieves a video with CloudFront URLs
func (s *AIAvatarService) GetVideoWithURLs(ctx context.Context, id uuid.UUID, cloudfrontDomain string) (*VideoWithURLs, error) {
	video, err := s.repo.GetVideoByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get video: %w", err)
	}

	return &VideoWithURLs{
		Video:        video,
		VideoURL:     fmt.Sprintf("https://%s/ai-avatar/videos/%s", cloudfrontDomain, video.Filename),
		ThumbnailURL: fmt.Sprintf("https://%s/ai-avatar/thumbnails/%s", cloudfrontDomain, video.ThumbnailFilename),
	}, nil
}

// GetAllVideosWithURLs retrieves all videos with CloudFront URLs
func (s *AIAvatarService) GetAllVideosWithURLs(ctx context.Context, cloudfrontDomain string) ([]*VideoWithURLs, error) {
	videos, err := s.repo.GetAllVideos(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get videos: %w", err)
	}

	var videosWithURLs []*VideoWithURLs
	for _, video := range videos {
		videosWithURLs = append(videosWithURLs, &VideoWithURLs{
			Video:        video,
			VideoURL:     fmt.Sprintf("https://%s/ai-avatar/videos/%s", cloudfrontDomain, video.Filename),
			ThumbnailURL: fmt.Sprintf("https://%s/ai-avatar/thumbnails/%s", cloudfrontDomain, video.ThumbnailFilename),
		})
	}

	return videosWithURLs, nil
}

// VideoWithURLs represents a video with CloudFront URLs
type VideoWithURLs struct {
	Video        *db.Video `json:"video"`
	VideoURL     string    `json:"video_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
}
