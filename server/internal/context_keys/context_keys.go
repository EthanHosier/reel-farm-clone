package context_keys

import "context"

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserEmailKey contextKey = "user_email"
)

// SetUserID adds the user ID to the context
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// SetUserEmail adds the user email to the context
func SetUserEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, UserEmailKey, email)
}

// GetUserID retrieves the user ID from the context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetUserEmail retrieves the user email from the context
func GetUserEmail(ctx context.Context) string {
	if email, ok := ctx.Value(UserEmailKey).(string); ok {
		return email
	}
	return ""
}
