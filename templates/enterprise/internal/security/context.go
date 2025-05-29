package security

import "context"

type contextKey string

const (
	clientIdentityKey contextKey = "client_identity"
	auditContextKey   contextKey = "audit_context"
)

// WithClientIdentity adds client identity to context
func WithClientIdentity(ctx context.Context, clientID string) context.Context {
	return context.WithValue(ctx, clientIdentityKey, clientID)
}

// GetClientIdentity retrieves client identity from context
func GetClientIdentity(ctx context.Context) string {
	if clientID, ok := ctx.Value(clientIdentityKey).(string); ok {
		return clientID
	}
	return ""
}

// AuditContext holds audit information
type AuditContext struct {
	UserID    string
	Action    string
	Resource  string
	Timestamp int64
	RequestID string
}

// WithAuditContext adds audit context
func WithAuditContext(ctx context.Context, auditCtx *AuditContext) context.Context {
	return context.WithValue(ctx, auditContextKey, auditCtx)
}

// GetAuditContext retrieves audit context
func GetAuditContext(ctx context.Context) *AuditContext {
	if auditCtx, ok := ctx.Value(auditContextKey).(*AuditContext); ok {
		return auditCtx
	}
	return nil
}
