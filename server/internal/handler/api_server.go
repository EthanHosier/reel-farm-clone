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
	userService         *service.UserService
	subscriptionService *service.SubscriptionService
}

// NewAPIServer creates a new API server handler
func NewAPIServer(userService *service.UserService, subscriptionService *service.SubscriptionService) *APIServer {
	return &APIServer{
		userService:         userService,
		subscriptionService: subscriptionService,
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
		Credits:           int(userAccount.Credits),
		BillingCustomerId: userAccount.BillingCustomerID,
		CreatedAt:         userAccount.CreatedAt,
		UpdatedAt:         userAccount.UpdatedAt,
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiUserAccount)
}

// CreateCheckoutSession handles POST /subscription/create-checkout-session
func (s *APIServer) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
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
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Invalid user ID format",
		})
		return
	}

	// Get email from context
	email := context_keys.GetUserEmail(r.Context())
	if email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "unauthorized",
			Message: "Email not found in context",
		})
		return
	}

	// Parse request body
	var req api.CreateCheckoutSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
		return
	}

	// Validate required URLs
	if req.SuccessUrl == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "missing_success_url",
			Message: "success_url is required",
		})
		return
	}

	if req.CancelUrl == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "missing_cancel_url",
			Message: "cancel_url is required",
		})
		return
	}

	// Create checkout session
	checkoutURL, err := s.subscriptionService.CreateCheckoutSession(
		r.Context(),
		userID,
		email,
		req.PriceId,
		req.SuccessUrl,
		req.CancelUrl,
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "checkout_session_failed",
			Message: "Failed to create checkout session",
		})
		return
	}

	// Return checkout URL
	response := api.CheckoutSessionResponse{
		CheckoutUrl: checkoutURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateCustomerPortalSession handles POST /subscription/customer-portal
func (s *APIServer) CreateCustomerPortalSession(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
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
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Invalid user ID format",
		})
		return
	}

	// Parse request body
	var req api.CreateCustomerPortalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
		return
	}

	// Validate required return URL
	if req.ReturnUrl == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "missing_return_url",
			Message: "return_url is required",
		})
		return
	}

	// Create customer portal session
	portalURL, err := s.subscriptionService.CreateCustomerPortalSession(
		r.Context(),
		userID,
		req.ReturnUrl,
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "portal_session_failed",
			Message: "Failed to create customer portal session",
		})
		return
	}

	// Return portal URL
	response := api.CustomerPortalResponse{
		PortalUrl: portalURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
