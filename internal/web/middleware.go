package handler

import (
	"context"
	"net/http"
	"strings"

	"tic-tac-toe/internal/domain"
)

func AuthMiddleware(
	jwtService *domain.JwtService,
	userService *domain.UserService,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "authorization required", http.StatusUnauthorized)
				return
			}

			const prefix = "Bearer "

			if !strings.HasPrefix(authHeader, prefix) {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, prefix)

			if err := jwtService.ValidateAccessToken(tokenString); err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			userUUID, err := jwtService.GetUUID(tokenString)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			user, err := userService.GetUserByUUID(userUUID)
			if err != nil {
				http.Error(w, "user not found", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				userContextKey,
				user,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
