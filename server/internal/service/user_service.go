package service

import (
	"context"

	"github.com/ethanhosier/reel-farm/db"
	"github.com/ethanhosier/reel-farm/internal/repository"

	"github.com/google/uuid"
)

// UserService handles user business logic
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserAccount retrieves a user account by ID
func (s *UserService) GetUserAccount(ctx context.Context, id uuid.UUID) (*db.UserAccount, error) {
	return s.userRepo.GetUserAccount(ctx, id)
}
