package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ethanhosier/reel-farm/internal/api"
	"github.com/ethanhosier/reel-farm/internal/context_keys"
	"github.com/ethanhosier/reel-farm/internal/service"
	"github.com/google/uuid"
)

// APIServer implements the generated ServerInterface
type APIServer struct {
	userService *service.UserService
}

// NewAPIServer creates a new API server handler
func NewAPIServer(userService *service.UserService) *APIServer {
	return &APIServer{
		userService: userService,
	}
}

// GetHealth handles GET /health
func (s *APIServer) GetHealth(w http.ResponseWriter, r *http.Request) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Get port from environment or use default
	port := "3000"
	if r.Header.Get("X-Port") != "" {
		port = r.Header.Get("X-Port")
	}

	// Create response using generated type
	response := api.HealthResponse{
		Message: "Reel Farm API is healthy!!!",
		Status:  "ok",
		Port:    port,
	}

	// Convert to JSON and send
	json.NewEncoder(w).Encode(response)
}

// GetUserAccount handles GET /user
func (s *APIServer) GetUserAccount(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userIDStr := context_keys.GetUserID(r.Context())
	if userIDStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "unauthorized",
			Message: "User ID not found in context",
		})
		return
	}

	// Convert string to UUID
	id, err := uuid.Parse(userIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Invalid user ID format",
		})
		return
	}

	// Get user account
	userAccount, err := s.userService.GetUserAccount(r.Context(), id)
	if err != nil {
		// Return error response using generated type
		errorResponse := api.ErrorResponse{
			Error:   "user_not_found",
			Message: "Failed to get user account",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Convert database model to API model
	var planEndsAt *time.Time
	if userAccount.PlanEndsAt.Valid {
		planEndsAt = &userAccount.PlanEndsAt.Time
	}

	apiUserAccount := api.UserAccount{
		Id:                userAccount.ID,
		Plan:              userAccount.Plan,
		PlanStartedAt:     userAccount.PlanStartedAt,
		PlanEndsAt:        planEndsAt,
		BillingCustomerId: userAccount.BillingCustomerID,
		CreatedAt:         userAccount.CreatedAt,
		UpdatedAt:         userAccount.UpdatedAt,
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiUserAccount)
}
