package service

import (
	"context"

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
