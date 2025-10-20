package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ethanhosier/reel-farm/internal/service"

	"github.com/google/uuid"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUserAccount handles GET /users/{id}
func (h *UserHandler) GetUserAccount(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path
	idStr := r.URL.Path[len("/users/"):]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get user account
	userAccount, err := h.userService.GetUserAccount(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get user account", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userAccount)
}
