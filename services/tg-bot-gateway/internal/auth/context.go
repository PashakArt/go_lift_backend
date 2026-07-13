package auth

import "context"

type contextKey string

const (
	userIDKey   contextKey = "userID"
	tenantIDKey contextKey = "tenantID"
)

func ContextWithAuth(ctx context.Context, userID, tenantID string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

func UserIDFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(userIDKey).(string); ok {
		return val
	}
	return ""
}

func TenantIDFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(tenantIDKey).(string); ok {
		return val
	}
	return ""
}
