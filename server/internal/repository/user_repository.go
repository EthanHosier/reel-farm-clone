package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ethanhosier/reel-farm/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository handles user account operations
type UserRepository struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

// NewUserRepository creates a new user repository
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		queries: db.New(pool),
		pool:    pool,
	}
}

// GetUserAccount retrieves a user account by ID
func (r *UserRepository) GetUserAccount(ctx context.Context, id uuid.UUID) (*db.UserAccount, error) {
	userAccount, err := r.queries.GetUserAccount(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user account: %w", err)
	}
	return userAccount, nil
}

// GetUserByBillingCustomerID retrieves a user account by Stripe customer ID
func (r *UserRepository) GetUserByBillingCustomerID(ctx context.Context, customerID string) (*db.UserAccount, error) {
	userAccount, err := r.queries.GetUserByBillingCustomerID(ctx, &customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by billing customer ID: %w", err)
	}
	return userAccount, nil
}

// UpdateUserPlan updates a user's subscription plan
func (r *UserRepository) UpdateUserPlan(ctx context.Context, id uuid.UUID, plan string, planStartedAt time.Time, planEndsAt *time.Time) error {
	var pgPlanEndsAt pgtype.Timestamptz
	if planEndsAt != nil {
		pgPlanEndsAt = pgtype.Timestamptz{Time: *planEndsAt, Valid: true}
	}

	params := &db.UpdateUserPlanParams{
		ID:            id,
		Plan:          plan,
		PlanStartedAt: planStartedAt,
		PlanEndsAt:    pgPlanEndsAt,
	}

	err := r.queries.UpdateUserPlan(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update user plan: %w", err)
	}
	return nil
}

// UpdateUserBillingCustomerID updates a user's Stripe customer ID
func (r *UserRepository) UpdateUserBillingCustomerID(ctx context.Context, id uuid.UUID, customerID string) error {
	params := &db.UpdateUserBillingCustomerIDParams{
		ID:                id,
		BillingCustomerID: &customerID,
	}

	err := r.queries.UpdateUserBillingCustomerID(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update billing customer ID: %w", err)
	}
	return nil
}

// AddCreditsToUser adds credits to a user's account
func (r *UserRepository) AddCreditsToUser(ctx context.Context, id uuid.UUID, credits int32) error {
	params := &db.AddCreditsToUserParams{
		ID:      id,
		Credits: credits,
	}

	err := r.queries.AddCreditsToUser(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to add credits to user: %w", err)
	}
	return nil
}

// WithTransaction executes a function within a database transaction
func (r *UserRepository) WithTransaction(ctx context.Context, fn func(*UserRepository) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Always rollback unless committed

	// Create a new repository instance with the transaction
	txRepo := &UserRepository{
		queries: db.New(tx), // SQLC works with transactions
		pool:    r.pool,     // Keep reference to pool for potential nested transactions
	}

	if err := fn(txRepo); err != nil {
		return err // Transaction will be rolled back via defer
	}

	return tx.Commit(ctx)
}
