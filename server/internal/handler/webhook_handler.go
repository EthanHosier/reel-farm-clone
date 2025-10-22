package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ethanhosier/reel-farm/internal/service"
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

	// Process the subscription creation business logic
	err := h.subscriptionService.ProcessSubscriptionCreated(context.Background(), &subscription)
	if err != nil {
		return fmt.Errorf("failed to process subscription creation: %w", err)
	}

	fmt.Printf("✅ Subscription %s processed successfully\n", subscription.ID)
	return nil
}

// handleSubscriptionUpdated handles when a subscription is updated
func (h *WebhookHandler) handleSubscriptionUpdated(event stripe.Event) error {
	var subscription stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	fmt.Printf("Subscription updated: %s status: %s\n", subscription.ID, subscription.Status)

	// Process the subscription update business logic
	err := h.subscriptionService.ProcessSubscriptionUpdated(context.Background(), &subscription)
	if err != nil {
		return fmt.Errorf("failed to process subscription update: %w", err)
	}

	fmt.Printf("✅ Subscription %s updated successfully\n", subscription.ID)
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

	// Process the payment success business logic
	err := h.subscriptionService.ProcessPaymentSucceeded(context.Background(), &invoice)
	if err != nil {
		return fmt.Errorf("failed to process payment success: %w", err)
	}

	fmt.Printf("✅ Payment for invoice %s processed successfully\n", invoice.ID)
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
