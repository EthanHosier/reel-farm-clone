package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethanhosier/reel-farm/db"
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
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

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

// GetUserByBillingCustomerID retrieves a user by Stripe customer ID
func (s *SubscriptionService) GetUserByBillingCustomerID(ctx context.Context, customerID string) (*db.UserAccount, error) {
	return s.userRepo.GetUserByBillingCustomerID(ctx, customerID)
}

// UpdateUserPlan updates a user's subscription plan
func (s *SubscriptionService) UpdateUserPlan(ctx context.Context, userID uuid.UUID, plan string, planStartedAt time.Time, planEndsAt *time.Time) error {
	return s.userRepo.UpdateUserPlan(ctx, userID, plan, planStartedAt, planEndsAt)
}

// AddCreditsToUser adds credits to a user's account
func (s *SubscriptionService) AddCreditsToUser(ctx context.Context, userID uuid.UUID, credits int32) error {
	return s.userRepo.AddCreditsToUser(ctx, userID, credits)
}

// GetSubscriptionByID retrieves a Stripe subscription by ID
func (s *SubscriptionService) GetSubscriptionByID(ctx context.Context, subscriptionID string) (*stripe.Subscription, error) {
	subscription, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription %s: %w", subscriptionID, err)
	}
	return subscription, nil
}

// GetUserAccount retrieves a user account by ID
func (s *SubscriptionService) GetUserAccount(ctx context.Context, userID uuid.UUID) (*db.UserAccount, error) {
	return s.userRepo.GetUserAccount(ctx, userID)
}
