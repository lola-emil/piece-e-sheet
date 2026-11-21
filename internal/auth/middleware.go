package auth

import (
	"api/internal/respond"
	"context"
	"net/http"
	"strings"
)

type ctxKey struct{}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ctxKey{}, userID)
}

func UserID(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(ctxKey{}).(string)
	return uid, ok
}

func UserIDFromRequest(r *http.Request) (string, bool) {
	return UserID(r.Context())
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				tokenString = parts[1]
			} else {
				respond.Error(w, http.StatusUnauthorized, "Invalid authorization format")
				return
			}
		}

		if tokenString == "" {
			cookie, err := r.Cookie("access_token")
			if err == nil && cookie.Value != "" {
				tokenString = cookie.Value
			}
		}

		if tokenString == "" {
			respond.Error(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		claims, err := ValidateTokenClaims(tokenString)
		if err != nil {
			respond.Error(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		if claims.Type == TokenTypeRefresh {
			respond.Error(w, http.StatusUnauthorized, "Access token required")
			return
		}

		next.ServeHTTP(w, r.WithContext(WithUserID(r.Context(), claims.UserID)))
	})
}
