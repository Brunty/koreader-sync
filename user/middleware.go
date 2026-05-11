package user

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/brunty/koreader-sync-server/crypto"
	"github.com/brunty/koreader-sync-server/handlers"
)

type AuthMiddleware struct {
	repo UserRepository
}

func NewAuthMiddleware(repo UserRepository) *AuthMiddleware {
	return &AuthMiddleware{repo: repo}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Auth middleware invoked")

		username := r.Header.Get("x-auth-user")
		password := r.Header.Get("x-auth-key")

		if username == "" || password == "" {
			slog.Debug("Auth middleware fail, username or password empty")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(&handlers.StatusResponse{Status: "unauthorized"})
			// don't continue down the middleware chain as they need to be authorized
			return
		}

		user, err := m.repo.SelectByUsername(username)
		if user == nil || err != nil {
			slog.Debug("Auth middleware fail, user not found")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(&handlers.StatusResponse{Status: "unauthorized"})
			// don't continue down the middleware chain as they need to be authorized
			return
		}

		passwordsMatch := crypto.CheckPasswordHash(password, user.Password)

		if !passwordsMatch {
			slog.Debug("Auth middleware fail, password incorrect")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(&handlers.StatusResponse{Status: "unauthorized"})
			// don't continue down the middleware chain as they need to be authorized
			return
		}

		// Attach the user ID to the context so they're usable elsewhere
		ctx := context.WithValue(r.Context(), "user", user.Id)

		slog.Debug("Auth middleware success for user", slog.Int("user ID", int(user.Id)))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
