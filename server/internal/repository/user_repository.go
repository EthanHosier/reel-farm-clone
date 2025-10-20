package repository

import (
	"context"
	"fmt"

	"github.com/ethanhosier/reel-farm/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
