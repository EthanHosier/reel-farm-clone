package repository

import (
	"context"

	"github.com/ethanhosier/reel-farm/db"
	"github.com/google/uuid"
)

type AIAvatarRepository struct {
	queries *db.Queries
}

func NewAIAvatarRepository(queries *db.Queries) *AIAvatarRepository {
	return &AIAvatarRepository{
		queries: queries,
	}
}

// GetAllVideos retrieves all AI avatar videos
func (r *AIAvatarRepository) GetAllVideos(ctx context.Context) ([]*db.Video, error) {
	return r.queries.GetAllVideos(ctx)
}

// GetVideoByID retrieves a specific video by ID
func (r *AIAvatarRepository) GetVideoByID(ctx context.Context, id uuid.UUID) (*db.Video, error) {
	return r.queries.GetVideoByID(ctx, id)
}

// CreateVideo creates a new video record
func (r *AIAvatarRepository) CreateVideo(ctx context.Context, params *db.CreateVideoParams) (*db.Video, error) {
	return r.queries.CreateVideo(ctx, params)
}

// UpdateVideo updates an existing video
func (r *AIAvatarRepository) UpdateVideo(ctx context.Context, params *db.UpdateVideoParams) (*db.Video, error) {
	return r.queries.UpdateVideo(ctx, params)
}

// DeleteVideo deletes a video by ID
func (r *AIAvatarRepository) DeleteVideo(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteVideo(ctx, id)
}

// VideoExists checks if a video exists by ID
func (r *AIAvatarRepository) VideoExists(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := r.queries.GetVideoByID(ctx, id)
	if err != nil {
		// If video not found, it doesn't exist
		return false, nil
	}
	return true, nil
}
