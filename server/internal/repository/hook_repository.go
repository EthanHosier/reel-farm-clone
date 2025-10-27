package repository

import (
	"context"
	"fmt"

	"github.com/ethanhosier/reel-farm/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// HookRepository handles hook operations
type HookRepository struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

// NewHookRepository creates a new hook repository
func NewHookRepository(pool *pgxpool.Pool) *HookRepository {
	return &HookRepository{
		queries: db.New(pool),
		pool:    pool,
	}
}

// CreateHooksBatch creates multiple hooks in a single database call
func (r *HookRepository) CreateHooksBatch(ctx context.Context, userID uuid.UUID, generationID uuid.UUID, prompt string, hookTexts []string, creditsUsed int32) ([]*db.Hook, error) {
	// Create hook indices array
	hookIndices := make([]int32, len(hookTexts))
	for i := range hookTexts {
		hookIndices[i] = int32(i)
	}

	params := &db.CreateHooksBatchParams{
		UserID:       pgtype.UUID{Bytes: userID, Valid: true},
		GenerationID: pgtype.UUID{Bytes: generationID, Valid: true},
		Prompt:       prompt,
		Column4:      hookTexts,
		Column5:      hookIndices,
		CreditsUsed:  creditsUsed,
	}

	hooks, err := r.queries.CreateHooksBatch(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create hooks batch: %w", err)
	}
	return hooks, nil
}

// GetHooksByUser gets hooks for a user with pagination
func (r *HookRepository) GetHooksByUser(ctx context.Context, userID uuid.UUID, limit int32, offset int32) ([]*db.Hook, error) {
	params := &db.GetHooksByUserParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		Limit:  limit,
		Offset: offset,
	}

	hooks, err := r.queries.GetHooksByUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get hooks by user: %w", err)
	}
	return hooks, nil
}

// GetHooksByGeneration gets all hooks from a specific generation
func (r *HookRepository) GetHooksByGeneration(ctx context.Context, generationID uuid.UUID) ([]*db.Hook, error) {
	hooks, err := r.queries.GetHooksByGeneration(ctx, pgtype.UUID{Bytes: generationID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get hooks by generation: %w", err)
	}
	return hooks, nil
}

// GetHookByID gets a specific hook by ID
func (r *HookRepository) GetHookByID(ctx context.Context, hookID uuid.UUID) (*db.Hook, error) {
	hook, err := r.queries.GetHookByID(ctx, hookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hook by ID: %w", err)
	}
	return hook, nil
}

// DeleteHook deletes a hook (only if it belongs to the user)
func (r *HookRepository) DeleteHook(ctx context.Context, hookID uuid.UUID, userID uuid.UUID) error {
	params := &db.DeleteHookParams{
		ID:     hookID,
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	}

	err := r.queries.DeleteHook(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to delete hook: %w", err)
	}
	return nil
}

// DeleteHooks deletes multiple hooks (only if they belong to the user)
func (r *HookRepository) DeleteHooks(ctx context.Context, hookIDs []uuid.UUID, userID uuid.UUID) ([]*db.Hook, error) {
	// Convert []uuid.UUID to []pgtype.UUID
	pgtypes := make([]pgtype.UUID, len(hookIDs))
	for i, id := range hookIDs {
		pgtypes[i] = pgtype.UUID{Bytes: id, Valid: true}
	}

	params := &db.DeleteHooksParams{
		HookIds: pgtypes,
		UserID:  pgtype.UUID{Bytes: userID, Valid: true},
	}

	hooks, err := r.queries.DeleteHooks(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to delete hooks: %w", err)
	}
	return hooks, nil
}

// GetUserHookCount gets the total number of hooks for a user
func (r *HookRepository) GetUserHookCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	count, err := r.queries.GetUserHookCount(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return 0, fmt.Errorf("failed to get user hook count: %w", err)
	}
	return count, nil
}
