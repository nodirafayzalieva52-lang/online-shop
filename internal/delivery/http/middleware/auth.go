package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"shop/internal/delivery/http/dto"
	"shop/internal/domain"
	"shop/pkg/jwt"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

type AuthMiddleware struct {
	jwtService *jwt.Service
}

func NewAuthMiddleware(jwtService *jwt.Service) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
		Error: dto.APIError{
			Code:    code,
			Message: message,
		},
	})
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "authorization header is required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid authorization header format")
			return
		}

		tokenString := parts[1]
		claims, err := m.jwtService.ParseToken(tokenString)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, RoleKey, domain.Role(claims.Role))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole enforces role-based access control middleware.
func (m *AuthMiddleware) RequireRole(roles ...domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := GetRoleFromContext(r.Context())
			if !ok {
				writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized access")
				return
			}
			allowed := false
			for _, reqRole := range roles {
				if role == reqRole {
					allowed = true
					break
				}
			}
			if !allowed {
				writeJSONError(w, http.StatusForbidden, "FORBIDDEN", "access denied: insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

func GetRoleFromContext(ctx context.Context) (domain.Role, bool) {
	role, ok := ctx.Value(RoleKey).(domain.Role)
	return role, ok
}
