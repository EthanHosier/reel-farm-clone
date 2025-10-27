package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ethanhosier/reel-farm/db"
	"github.com/ethanhosier/reel-farm/internal/api"
	"github.com/ethanhosier/reel-farm/internal/context_keys"
	"github.com/ethanhosier/reel-farm/internal/service"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// APIServer implements the generated ServerInterface
type APIServer struct {
	userService         *service.UserService
	subscriptionService *service.SubscriptionService
	hookService         *service.HookService
	aiAvatarService     *service.AIAvatarService
}

// NewAPIServer creates a new API server handler
func NewAPIServer(userService *service.UserService, subscriptionService *service.SubscriptionService, hookService *service.HookService, aiAvatarService *service.AIAvatarService) *APIServer {
	return &APIServer{
		userService:         userService,
		subscriptionService: subscriptionService,
		hookService:         hookService,
		aiAvatarService:     aiAvatarService,
	}
}

// generateCloudFrontURL creates a CloudFront URL for a given path
func (s *APIServer) generateCloudFrontURL(path string) (string, error) {
	cloudfrontDomain := os.Getenv("CLOUDFRONT_DOMAIN")
	if cloudfrontDomain == "" {
		return "", fmt.Errorf("CLOUDFRONT_DOMAIN is not set")
	}
	return fmt.Sprintf("https://%s/%s", cloudfrontDomain, path), nil
}

// toVideoAPIResponse converts a database video to API response format
func (s *APIServer) toVideoAPIResponse(video *db.AiAvatarVideo) (*api.AIAvatarVideo, error) {
	videoPath := fmt.Sprintf("ai-avatar/videos/%s", video.Filename)
	thumbnailPath := fmt.Sprintf("ai-avatar/thumbnails/%s", video.ThumbnailFilename)

	videoURL, err := s.generateCloudFrontURL(videoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate video URL: %w", err)
	}
	thumbnailURL, err := s.generateCloudFrontURL(thumbnailPath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate thumbnail URL: %w", err)
	}

	return &api.AIAvatarVideo{
		Id:           openapi_types.UUID(video.ID),
		Title:        video.Title,
		VideoUrl:     videoURL,
		ThumbnailUrl: thumbnailURL,
		UpdatedAt:    video.UpdatedAt,
	}, nil
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

// GetAIAvatarVideos handles GET /ai-avatar/videos
func (s *APIServer) GetAIAvatarVideos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	videos, err := s.aiAvatarService.GetAllVideos(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve videos: " + err.Error(),
		})
		return
	}

	// Convert videos to API response format
	var videoResponses []api.AIAvatarVideo
	for _, video := range videos {
		videoResponse, err := s.toVideoAPIResponse(video)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to convert video to API response: " + err.Error(),
			})
			return
		}
		videoResponses = append(videoResponses, *videoResponse)
	}

	response := api.AIAvatarVideosResponse{
		Videos: videoResponses,
	}

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

// GenerateHooks handles POST /hooks/generate
func (s *APIServer) GenerateHooks(w http.ResponseWriter, r *http.Request) {
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
	var req api.GenerateHooksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
		return
	}

	// Validate request
	if req.Prompt == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "missing_prompt",
			Message: "prompt is required",
		})
		return
	}

	if req.NumHooks < 1 || req.NumHooks > 10 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "invalid_num_hooks",
			Message: "num_hooks must be between 1 and 10",
		})
		return
	}

	// Generate hooks
	hooks, err := s.hookService.GenerateHooks(r.Context(), userID, req.Prompt, int(req.NumHooks))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "hook_generation_failed",
			Message: err.Error(),
		})
		return
	}

	// Return hooks
	response := api.GenerateHooksResponse{
		Hooks: hooks,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetHooks handles GET /hooks
func (s *APIServer) GetHooks(w http.ResponseWriter, r *http.Request, params api.GetHooksParams) {
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

	// Set default pagination values
	limit := int32(20)
	offset := int32(0)

	if params.Limit != nil {
		limit = int32(*params.Limit)
	}
	if params.Offset != nil {
		offset = int32(*params.Offset)
	}

	// Get hooks from service
	hooks, totalCount, err := s.hookService.GetHooks(r.Context(), userID, limit, offset)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "failed_to_get_hooks",
			Message: "Failed to retrieve hooks",
		})
		return
	}

	// Ensure hooks is never nil - initialize as empty slice if nil
	response := api.GetHooksResponse{
		Hooks:      hooks,
		TotalCount: int(totalCount),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteHook handles DELETE /hooks/{hookId}
func (s *APIServer) DeleteHook(w http.ResponseWriter, r *http.Request, hookId openapi_types.UUID) {
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

	// Parse hook ID
	hookUUID, err := uuid.Parse(hookId.String())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "invalid_hook_id",
			Message: "Invalid hook ID format",
		})
		return
	}

	// Delete the hook
	err = s.hookService.DeleteHook(r.Context(), hookUUID, userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "hook_not_found",
			Message: "Hook not found or doesn't belong to user",
		})
		return
	}

	// Return success response
	response := map[string]string{
		"message": "Hook deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteHooksBulk handles DELETE /hooks/bulk
func (s *APIServer) DeleteHooksBulk(w http.ResponseWriter, r *http.Request) {
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
	var requestBody struct {
		HookIds []string `json:"hook_ids"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
		return
	}

	// Validate that hook_ids is not empty
	if len(requestBody.HookIds) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "invalid_request",
			Message: "hook_ids array cannot be empty",
		})
		return
	}

	// Parse hook IDs
	hookIDs := make([]uuid.UUID, len(requestBody.HookIds))
	for i, idStr := range requestBody.HookIds {
		hookID, err := uuid.Parse(idStr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ErrorResponse{
				Error:   "invalid_hook_id",
				Message: fmt.Sprintf("Invalid hook ID format: %s", idStr),
			})
			return
		}
		hookIDs[i] = hookID
	}

	// Delete the hooks
	deletedHooks, err := s.hookService.DeleteHooks(r.Context(), hookIDs, userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "failed_to_delete_hooks",
			Message: "Failed to delete hooks",
		})
		return
	}

	// Extract deleted hook IDs
	deletedIDs := make([]string, len(deletedHooks))
	for i, hook := range deletedHooks {
		deletedIDs[i] = hook.Id.String()
	}

	// Return success response
	response := map[string]interface{}{
		"message":       "Successfully deleted hooks",
		"deleted_count": len(deletedHooks),
		"deleted_ids":   deletedIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateUserGeneratedVideo handles POST /user-generated-videos
func (s *APIServer) CreateUserGeneratedVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract user ID from context
	userIDStr := context_keys.GetUserID(r.Context())
	if userIDStr == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "unauthorized",
			Message: "User ID not found in context",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Invalid user ID format",
		})
		return
	}

	// Parse request body
	var req api.CreateUserGeneratedVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
		})
		return
	}

	// Parse AI avatar video ID
	aiAvatarVideoID := uuid.UUID(req.AiAvatarVideoId)

	// Get the AI avatar video to get its URL
	aiAvatarVideo, err := s.aiAvatarService.GetVideoByID(r.Context(), aiAvatarVideoID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "video_not_found",
			Message: "AI avatar video not found",
		})
		return
	}

	// Generate CloudFront URL for the video
	videoPath := fmt.Sprintf("ai-avatar/videos/%s", aiAvatarVideo.Filename)
	videoURL, err := s.generateCloudFrontURL(videoPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to generate video URL",
		})
		return
	}

	// Process video with text overlay
	userGeneratedVideo, err := s.aiAvatarService.ProcessVideoWithTextOverlay(r.Context(), userID, aiAvatarVideoID, videoURL, req.OverlayText)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "processing_error",
			Message: "Failed to process video with text overlay" + err.Error(),
		})
		return
	}

	// Generate signed CloudFront URLs for the generated video (24 hour expiration)
	generatedVideoPath := fmt.Sprintf("user-generated-videos/videos/%s", userGeneratedVideo.GeneratedVideoFilename)
	generatedThumbnailPath := fmt.Sprintf("user-generated-videos/thumbnails/%s", userGeneratedVideo.ThumbnailFilename)

	generatedVideoURL, err := s.aiAvatarService.GenerateSignedURL(generatedVideoPath, 24*time.Hour)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to generate signed video URL",
		})
		return
	}

	generatedThumbnailURL, err := s.aiAvatarService.GenerateSignedURL(generatedThumbnailPath, 24*time.Hour)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to generate signed thumbnail URL",
		})
		return
	}

	// Create response
	response := api.UserGeneratedVideoResponse{
		Video: api.UserGeneratedVideo{
			Id:              openapi_types.UUID(userGeneratedVideo.ID),
			UserId:          openapi_types.UUID(userGeneratedVideo.UserID.Bytes),
			AiAvatarVideoId: openapi_types.UUID(userGeneratedVideo.AiAvatarVideoID.Bytes),
			OverlayText:     userGeneratedVideo.OverlayText,
			VideoUrl:        generatedVideoURL,
			ThumbnailUrl:    generatedThumbnailURL,
			Status:          api.UserGeneratedVideoStatus(*userGeneratedVideo.Status),
			CreatedAt:       userGeneratedVideo.CreatedAt,
		},
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetUserGeneratedVideos handles GET /user-generated-videos
func (s *APIServer) GetUserGeneratedVideos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract user ID from context
	userIDStr := context_keys.GetUserID(r.Context())
	if userIDStr == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "unauthorized",
			Message: "User ID not found in context",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Invalid user ID format",
		})
		return
	}

	// Get user-generated videos from service
	userGeneratedVideos, err := s.aiAvatarService.GetUserGeneratedVideosByUserID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve user-generated videos",
		})
		return
	}

	// Convert to API response format
	var videoResponses []api.UserGeneratedVideo
	for _, video := range userGeneratedVideos {
		// Generate signed CloudFront URLs for each video (24 hour expiration)
		videoPath := fmt.Sprintf("user-generated-videos/videos/%s", video.GeneratedVideoFilename)
		thumbnailPath := fmt.Sprintf("user-generated-videos/thumbnails/%s", video.ThumbnailFilename)

		videoURL, err := s.aiAvatarService.GenerateSignedURL(videoPath, 24*time.Hour)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to generate signed video URL",
			})
			return
		}

		thumbnailURL, err := s.aiAvatarService.GenerateSignedURL(thumbnailPath, 24*time.Hour)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to generate signed thumbnail URL",
			})
			return
		}

		videoResponses = append(videoResponses, api.UserGeneratedVideo{
			Id:              openapi_types.UUID(video.ID),
			UserId:          openapi_types.UUID(video.UserID.Bytes),
			AiAvatarVideoId: openapi_types.UUID(video.AiAvatarVideoID.Bytes),
			OverlayText:     video.OverlayText,
			VideoUrl:        videoURL,
			ThumbnailUrl:    thumbnailURL,
			Status:          api.UserGeneratedVideoStatus(*video.Status),
			CreatedAt:       video.CreatedAt,
		})
	}

	response := api.UserGeneratedVideosResponse{
		Videos: videoResponses,
	}

	json.NewEncoder(w).Encode(response)
}
