package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ethanhosier/reel-farm/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// UserRepository handles user account operations
type UserRepository struct {
	queries *db.Queries
}

// NewUserRepository creates a new user repository
func NewUserRepository(conn *pgx.Conn) *UserRepository {
	return &UserRepository{
		queries: db.New(conn),
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
