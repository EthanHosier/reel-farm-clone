package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethanhosier/reel-farm/internal/repository"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v78"
	billingportalsession "github.com/stripe/stripe-go/v78/billingportal/session"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/subscription"
)

type SubscriptionService struct {
	userRepo *repository.UserRepository
}

func NewSubscriptionService(userRepo *repository.UserRepository) *SubscriptionService {
	// Set Stripe API key
	stripeSecretKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeSecretKey == "" {
		log.Fatal("STRIPE_SECRET_KEY is not set")
	}

	stripe.Key = stripeSecretKey

	return &SubscriptionService{
		userRepo: userRepo,
	}
}

// CreateCheckoutSession creates a Stripe checkout session for subscription
func (s *SubscriptionService) CreateCheckoutSession(ctx context.Context, userID uuid.UUID, email, priceID, successURL, cancelURL string) (string, error) {
	// Get user account to check if they already have a Stripe customer ID
	userAccount, err := s.userRepo.GetUserAccount(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user account: %w", err)
	}

	var customerID string

	// If user doesn't have a Stripe customer ID, create one
	if userAccount.BillingCustomerID == nil || *userAccount.BillingCustomerID == "" {
		// Create Stripe customer
		customerParams := &stripe.CustomerParams{
			Email: stripe.String(email),
		}

		customer, err := customer.New(customerParams)
		if err != nil {
			return "", fmt.Errorf("failed to create Stripe customer: %w", err)
		}

		customerID = customer.ID

		// Update user account with Stripe customer ID
		// TODO: Implement UpdateUserAccount method in repository
		// s.userRepo.UpdateUserAccount(ctx, userID, map[string]interface{}{
		// 	"billing_customer_id": customerID,
		// })
	} else {
		customerID = *userAccount.BillingCustomerID
	}

	// Create checkout session
	params := &stripe.CheckoutSessionParams{
		Customer:          stripe.String(customerID),
		ClientReferenceID: stripe.String(userID.String()), // Use user ID as reference
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"user_id": userID.String(), // Store user ID in subscription metadata
			},
		},
	}

	session, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	return session.URL, nil
}

// CreateCustomerPortalSession creates a Stripe customer portal session
func (s *SubscriptionService) CreateCustomerPortalSession(ctx context.Context, userID uuid.UUID, returnURL string) (string, error) {
	// Get user account
	userAccount, err := s.userRepo.GetUserAccount(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user account: %w", err)
	}

	if userAccount.BillingCustomerID == nil || *userAccount.BillingCustomerID == "" {
		return "", fmt.Errorf("user does not have a Stripe customer ID")
	}

	// Create customer portal session
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(*userAccount.BillingCustomerID),
		ReturnURL: stripe.String(returnURL),
	}

	session, err := billingportalsession.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create customer portal session: %w", err)
	}

	return session.URL, nil
}

// GetSubscriptionByID retrieves a Stripe subscription by ID
func (s *SubscriptionService) GetSubscriptionByID(ctx context.Context, subscriptionID string) (*stripe.Subscription, error) {
	subscription, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription %s: %w", subscriptionID, err)
	}
	return subscription, nil
}

// ProcessSubscriptionCreated handles subscription creation business logic
func (s *SubscriptionService) ProcessSubscriptionCreated(ctx context.Context, subscription *stripe.Subscription) error {
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

	// Update user to 'pro' plan
	planStartedAt := time.Unix(subscription.Created, 0)
	planEndsAt := time.Unix(subscription.CurrentPeriodEnd, 0)

	err = s.userRepo.UpdateUserPlan(ctx, userID, "pro", planStartedAt, &planEndsAt)
	if err != nil {
		return fmt.Errorf("failed to update user plan: %w", err)
	}

	return nil
}

// ProcessSubscriptionUpdated handles subscription update business logic
func (s *SubscriptionService) ProcessSubscriptionUpdated(ctx context.Context, subscription *stripe.Subscription) error {
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

	// Handle different subscription statuses
	switch subscription.Status {
	case stripe.SubscriptionStatusActive:
		// Ensure user is on pro plan
		planStartedAt := time.Unix(subscription.Created, 0)
		planEndsAt := time.Unix(subscription.CurrentPeriodEnd, 0)

		err = s.userRepo.UpdateUserPlan(ctx, userID, "pro", planStartedAt, &planEndsAt)
		if err != nil {
			return fmt.Errorf("failed to update user plan to pro: %w", err)
		}

	case stripe.SubscriptionStatusCanceled, stripe.SubscriptionStatusUnpaid, stripe.SubscriptionStatusPastDue:
		// Downgrade to free plan but keep credits
		err = s.userRepo.UpdateUserPlan(ctx, userID, "free", time.Now(), nil)
		if err != nil {
			return fmt.Errorf("failed to downgrade user to free plan: %w", err)
		}
	}

	return nil
}

// ProcessPaymentSucceeded handles payment success business logic
func (s *SubscriptionService) ProcessPaymentSucceeded(ctx context.Context, invoice *stripe.Invoice) error {
	// Get subscription to access metadata
	if invoice.Subscription == nil {
		return fmt.Errorf("no subscription found in invoice")
	}

	subscription, err := s.GetSubscriptionByID(ctx, invoice.Subscription.ID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

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

	// Execute credit addition and plan update in a transaction
	return s.userRepo.WithTransaction(ctx, func(txRepo *repository.UserRepository) error {
		// Add 500 credits for monthly subscription
		err := txRepo.AddCreditsToUser(ctx, userID, 500)
		if err != nil {
			return fmt.Errorf("failed to add monthly credits: %w", err)
		}

		// Update plan end date to next billing period
		planEndsAt := time.Unix(subscription.CurrentPeriodEnd, 0)
		err = txRepo.UpdateUserPlan(ctx, userID, "pro", time.Now(), &planEndsAt)
		if err != nil {
			return fmt.Errorf("failed to update plan end date: %w", err)
		}

		return nil
	})
}
