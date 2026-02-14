package middleware

import (
	"context"
	"net/http"

	"auth-proxy/internal/token"
)

type contextKey string

const UserKey contextKey = "user"

func Auth(tokenService *token.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr, err := tokenService.GetTokenFromRequest(r)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			claims, err := tokenService.ValidateToken(tokenStr)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
