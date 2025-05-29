package security

import (
	"net/http"
	"strings"
)

// Permission represents a specific permission
type Permission string

const (
	PermissionHealthRead     Permission = "health:read"
	PermissionHealthWrite    Permission = "health:write"
	PermissionMetricsRead    Permission = "metrics:read"
	PermissionDependencyRead Permission = "dependency:read"
	PermissionAdminAccess    Permission = "admin:access"
)

// Role represents a user role with associated permissions
type Role struct {
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}

// User represents a user with roles
type User struct {
	ID    string `json:"id"`
	Roles []Role `json:"roles"`
}

// RBACPolicy holds the role-based access control policy
type RBACPolicy struct {
	Users map[string]User `json:"users"`
	Roles map[string]Role `json:"roles"`
}

// DefaultRBACPolicy returns a default RBAC policy
func DefaultRBACPolicy() *RBACPolicy {
	return &RBACPolicy{
		Users: map[string]User{
			"admin": {
				ID: "admin",
				Roles: []Role{
					{Name: "admin", Permissions: []Permission{
						PermissionHealthRead,
						PermissionHealthWrite,
						PermissionMetricsRead,
						PermissionDependencyRead,
						PermissionAdminAccess,
					}},
				},
			},
			"service": {
				ID: "service",
				Roles: []Role{
					{Name: "service", Permissions: []Permission{
						PermissionHealthRead,
					}},
				},
			},
		},
		Roles: map[string]Role{
			"admin": {
				Name: "admin",
				Permissions: []Permission{
					PermissionHealthRead,
					PermissionHealthWrite,
					PermissionMetricsRead,
					PermissionDependencyRead,
					PermissionAdminAccess,
				},
			},
			"service": {
				Name: "service",
				Permissions: []Permission{
					PermissionHealthRead,
				},
			},
		},
	}
}

// RBACMiddleware validates user permissions for requests
func RBACMiddleware(policy *RBACPolicy) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientID := GetClientIdentity(r.Context())
			if clientID == "" {
				http.Error(w, "Client identity required", http.StatusUnauthorized)
				return
			}

			requiredPermission := getRequiredPermission(r)
			if requiredPermission == "" {
				next.ServeHTTP(w, r)
				return
			}

			if !policy.HasPermission(clientID, requiredPermission) {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// HasPermission checks if a user has a specific permission
func (p *RBACPolicy) HasPermission(userID string, permission Permission) bool {
	user, exists := p.Users[userID]
	if !exists {
		return false
	}

	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			if perm == permission {
				return true
			}
		}
	}

	return false
}

// getRequiredPermission determines the required permission based on the request
func getRequiredPermission(r *http.Request) Permission {
	path := strings.TrimPrefix(r.URL.Path, "/")
	method := r.Method

	switch {
	case strings.HasPrefix(path, "health"):
		if method == "GET" {
			return PermissionHealthRead
		}
		return PermissionHealthWrite
	case strings.HasPrefix(path, "metrics"):
		return PermissionMetricsRead
	case strings.HasPrefix(path, "dependencies"):
		return PermissionDependencyRead
	case strings.HasPrefix(path, "admin"):
		return PermissionAdminAccess
	default:
		return ""
	}
}
