package repository

import (
	"context"
	"database/sql"

	"github.com/ethanhosier/reel-farm/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AIAvatarRepository struct {
	queries *db.Queries
}

func NewAIAvatarRepository(pool *pgxpool.Pool) *AIAvatarRepository {
	return &AIAvatarRepository{
		queries: db.New(pool),
	}
}

// GetAllVideos retrieves all AI avatar videos
func (r *AIAvatarRepository) GetAllVideos(ctx context.Context) ([]*db.AiAvatarVideo, error) {
	return r.queries.GetAllVideos(ctx)
}

// GetVideoByID retrieves a specific video by ID
func (r *AIAvatarRepository) GetVideoByID(ctx context.Context, id uuid.UUID) (*db.AiAvatarVideo, error) {
	return r.queries.GetVideoByID(ctx, id)
}

// CreateVideo creates a new video record
func (r *AIAvatarRepository) CreateVideo(ctx context.Context, params *db.CreateVideoParams) (*db.AiAvatarVideo, error) {
	return r.queries.CreateVideo(ctx, params)
}

// UpdateVideo updates an existing video
func (r *AIAvatarRepository) UpdateVideo(ctx context.Context, params *db.UpdateVideoParams) (*db.AiAvatarVideo, error) {
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
		// Check if it's a "no rows" error
		if err == sql.ErrNoRows || err.Error() == "no rows in result set" {
			// Video not found, it doesn't exist
			return false, nil
		}
		// Some other database error
		return false, err
	}
	return true, nil
}
