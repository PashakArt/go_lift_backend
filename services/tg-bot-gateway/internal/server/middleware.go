package server

import (
	"net/http"
	"strings"

	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/auth"
	"github.com/PashakArt/go_lift_backend/services/tg-bot-gateway/internal/server/handlers"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Если это preflight-запрос, сразу отвечаем 200
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(jwtManager *auth.JwtManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handlers.RespondWithError(w, http.StatusUnauthorized, "Authorization header missing")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				handlers.RespondWithError(w, http.StatusUnauthorized, "Invalid authorization format")
				return
			}

			claims, err := jwtManager.Verify(parts[1])
			if err != nil {
				handlers.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := auth.ContextWithAuth(r.Context(), claims.UserId, claims.TenantId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
