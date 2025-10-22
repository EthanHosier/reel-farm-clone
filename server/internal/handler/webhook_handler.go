package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ethanhosier/reel-farm/internal/service"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
)

// WebhookHandler handles Stripe webhook events
type WebhookHandler struct {
	subscriptionService *service.SubscriptionService
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(subscriptionService *service.SubscriptionService) *WebhookHandler {
	return &WebhookHandler{
		subscriptionService: subscriptionService,
	}
}

// ServeHTTP handles webhook requests
func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Get webhook secret from environment
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if webhookSecret == "" {
		http.Error(w, "STRIPE_WEBHOOK_SECRET not configured", http.StatusInternalServerError)
		return
	}

	// Verify webhook signature
	// In webhook_handler.go, replace the webhook.ConstructEvent call with:
	event, err := webhook.ConstructEventWithOptions(body, r.Header.Get("Stripe-Signature"), webhookSecret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
	if err != nil {
		fmt.Printf("Webhook signature verification failed: %v\n", err)
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	// Handle the event
	switch event.Type {
	case "customer.subscription.created":
		err = h.handleSubscriptionCreated(event)
	case "customer.subscription.updated":
		err = h.handleSubscriptionUpdated(event)
	case "customer.subscription.deleted":
		err = h.handleSubscriptionDeleted(event)
	case "invoice.payment_succeeded":
		err = h.handlePaymentSucceeded(event)
	case "invoice.payment_failed":
		err = h.handlePaymentFailed(event)
	default:
		fmt.Printf("Unhandled event type: %s\n", event.Type)
	}

	if err != nil {
		fmt.Printf("Error handling webhook event %s: %v\n", event.Type, err)
		http.Error(w, "Failed to process webhook", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"received": true})
}

// handleSubscriptionCreated handles when a subscription is created
func (h *WebhookHandler) handleSubscriptionCreated(event stripe.Event) error {
	var subscription stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	fmt.Printf("Subscription created: %s for customer: %s\n", subscription.ID, subscription.Customer.ID)

	// Get user ID from subscription metadata
	userIDStr := ""
	if subscription.Metadata != nil {
		if val, exists := subscription.Metadata["user_id"]; exists {
			userIDStr = val
		}
	}

	if userIDStr == "" {
		return fmt.Errorf("no user ID found in subscription metadata")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse user ID: %w", err)
	}

	// Get user account
	userAccount, err := h.subscriptionService.GetUserAccount(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("failed to get user account: %w", err)
	}

	// Update user to 'pro' plan
	planStartedAt := time.Unix(subscription.Created, 0)
	planEndsAt := time.Unix(subscription.CurrentPeriodEnd, 0)

	err = h.subscriptionService.UpdateUserPlan(context.Background(), userAccount.ID, "pro", planStartedAt, &planEndsAt)
	if err != nil {
		return fmt.Errorf("failed to update user plan: %w", err)
	}

	fmt.Printf("✅ User %s upgraded to pro plan\n", userAccount.ID)
	return nil
}

// handleSubscriptionUpdated handles when a subscription is updated
func (h *WebhookHandler) handleSubscriptionUpdated(event stripe.Event) error {
	var subscription stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	fmt.Printf("Subscription updated: %s status: %s\n", subscription.ID, subscription.Status)

	// Get user ID from subscription metadata
	userIDStr := ""
	if subscription.Metadata != nil {
		if val, exists := subscription.Metadata["user_id"]; exists {
			userIDStr = val
		}
	}

	if userIDStr == "" {
		return fmt.Errorf("no user ID found in subscription metadata")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse user ID: %w", err)
	}

	// Get user account
	userAccount, err := h.subscriptionService.GetUserAccount(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("failed to get user account: %w", err)
	}

	// Handle different subscription statuses
	switch subscription.Status {
	case stripe.SubscriptionStatusActive:
		// Ensure user is on pro plan
		planStartedAt := time.Unix(subscription.Created, 0)
		planEndsAt := time.Unix(subscription.CurrentPeriodEnd, 0)

		err = h.subscriptionService.UpdateUserPlan(context.Background(), userAccount.ID, "pro", planStartedAt, &planEndsAt)
		if err != nil {
			return fmt.Errorf("failed to update user plan to pro: %w", err)
		}
		fmt.Printf("✅ User %s subscription reactivated\n", userAccount.ID)

	case stripe.SubscriptionStatusCanceled, stripe.SubscriptionStatusUnpaid, stripe.SubscriptionStatusPastDue:
		// Downgrade to free plan but keep credits
		err = h.subscriptionService.UpdateUserPlan(context.Background(), userAccount.ID, "free", time.Now(), nil)
		if err != nil {
			return fmt.Errorf("failed to downgrade user to free plan: %w", err)
		}
		fmt.Printf("⚠️ User %s subscription canceled/downgraded to free\n", userAccount.ID)
	}

	return nil
}

// handleSubscriptionDeleted handles when a subscription is deleted/canceled
func (h *WebhookHandler) handleSubscriptionDeleted(event stripe.Event) error {
	var subscription stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	fmt.Printf("Subscription deleted: %s for customer: %s\n", subscription.ID, subscription.Customer.ID)

	// TODO: Downgrade user account to 'free' plan
	// TODO: Set plan_ends_at to null
	// TODO: Keep existing credits (no expiration)

	return nil
}

// handlePaymentSucceeded handles when a payment succeeds
func (h *WebhookHandler) handlePaymentSucceeded(event stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return fmt.Errorf("failed to unmarshal invoice: %w", err)
	}

	fmt.Printf("Payment succeeded for invoice: %s, customer: %s\n", invoice.ID, invoice.Customer.ID)

	// Get user ID from client reference ID (from the subscription)
	if invoice.Subscription == nil {
		return fmt.Errorf("no subscription found in invoice")
	}

	// Get subscription to access metadata
	subscription, err := h.subscriptionService.GetSubscriptionByID(context.Background(), invoice.Subscription.ID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	userIDStr := ""
	if subscription.Metadata != nil {
		if val, exists := subscription.Metadata["user_id"]; exists {
			userIDStr = val
		}
	}

	if userIDStr == "" {
		return fmt.Errorf("no user ID found in subscription metadata")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse user ID: %w", err)
	}

	// Get user account
	userAccount, err := h.subscriptionService.GetUserAccount(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("failed to get user account: %w", err)
	}

	// Add 500 credits for monthly subscription
	err = h.subscriptionService.AddCreditsToUser(context.Background(), userAccount.ID, 500)
	if err != nil {
		return fmt.Errorf("failed to add monthly credits: %w", err)
	}

	// Update plan end date to next billing period
	planEndsAt := time.Unix(subscription.CurrentPeriodEnd, 0)
	err = h.subscriptionService.UpdateUserPlan(context.Background(), userAccount.ID, "pro", time.Now(), &planEndsAt)
	if err != nil {
		fmt.Printf("Warning: failed to update plan end date: %v\n", err)
	}

	fmt.Printf("✅ User %s received 500 monthly credits\n", userAccount.ID)
	return nil
}

// handlePaymentFailed handles when a payment fails
func (h *WebhookHandler) handlePaymentFailed(event stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return fmt.Errorf("failed to unmarshal invoice: %w", err)
	}

	fmt.Printf("Payment failed for invoice: %s, customer: %s\n", invoice.ID, invoice.Customer.ID)

	// TODO: Handle payment failure
	// TODO: Maybe send notification email
	// TODO: Don't downgrade immediately (Stripe will retry)

	return nil
}
