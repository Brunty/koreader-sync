package auth

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/brunty/koreader-sync-server/crypto"
	"github.com/brunty/koreader-sync-server/handlers"
	userpackage "github.com/brunty/koreader-sync-server/user"
)

type AuthMiddleware struct {
	repo userpackage.UserRepository
}

func NewAuthMiddleware(repo userpackage.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{repo: repo}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("auth middleware start")

		username := r.Header.Get("x-auth-user")
		password := r.Header.Get("x-auth-key")

		if username == "" || password == "" {
			slog.Info("auth middleware fail, username or password empty")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(&handlers.ErrorResponse{Error: "unauthorized"})
			// don't continue down the middleware chain as they need to be authorized
			return
		}

		user, err := m.repo.SelectByUsername(r.Context(), username)
		if user == nil || err != nil {
			slog.Info("auth middleware fail, user not found")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(&handlers.ErrorResponse{Error: "unauthorized"})
			// don't continue down the middleware chain as they need to be authorized
			return
		}

		passwordsMatch := crypto.BcryptCheckPasswordHash(password, user.Password)

		if !passwordsMatch {
			slog.Info("auth middleware fail, password incorrect")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(&handlers.ErrorResponse{Error: "unauthorized"})
			// don't continue down the middleware chain as they need to be authorized
			return
		}

		// I did consider using crypto.BcryptNeedsRehash here to update passwords if crypto.BcryptCost changes, but it
		// adds overhead into every request as the auth is checking passwords on everyone due to not using something
		// like a session token, so just handle that manually in the DB if needed

		// Attach the user ID to the context so they're usable elsewhere
		ctx := context.WithValue(r.Context(), "user", user.Id)

		slog.Debug("auth middleware success for user", slog.Int("user ID", int(user.Id)))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
