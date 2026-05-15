package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brunty/koreader-sync-server/crypto"
	"github.com/brunty/koreader-sync-server/db"
	user2 "github.com/brunty/koreader-sync-server/user"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_PassesToNextHandlerOnSuccessfulAuthentication(t *testing.T) {
	_ = db.Init(":memory:")
	db.SetupTables()

	userRepo := user2.NewUserRepository(db.DBCon)

	now := time.Now()

	password, _ := crypto.BcryptHashPassword("test-password-here")
	user := user2.User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	_, err := userRepo.Store(t.Context(), user)
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/users/auth", nil)
	req.Header.Add("x-auth-user", "test-username-here")
	req.Header.Add("x-auth-key", "test-password-here")

	rr := httptest.NewRecorder()
	handler := NewAuthMiddleware(userRepo).Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_ReturnsUnauthorizedIfUsernameIsBlank(t *testing.T) {
	_ = db.Init(":memory:")
	db.SetupTables()

	userRepo := user2.NewUserRepository(db.DBCon)

	now := time.Now()

	password, _ := crypto.BcryptHashPassword("test-password-here")
	user := user2.User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	_, err := userRepo.Store(t.Context(), user)
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/users/auth", nil)
	req.Header.Add("x-auth-key", "test-password-here")

	rr := httptest.NewRecorder()
	handler := NewAuthMiddleware(userRepo).Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_ReturnsUnauthorizedIfPasswordIsBlank(t *testing.T) {
	_ = db.Init(":memory:")
	db.SetupTables()

	userRepo := user2.NewUserRepository(db.DBCon)

	now := time.Now()

	password, _ := crypto.BcryptHashPassword("test-password-here")
	user := user2.User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	_, err := userRepo.Store(t.Context(), user)
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/users/auth", nil)
	req.Header.Add("x-auth-user", "test-username-here")

	rr := httptest.NewRecorder()
	handler := NewAuthMiddleware(userRepo).Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_ReturnsUnauthorizedIfUserIsNotFound(t *testing.T) {
	_ = db.Init(":memory:")
	db.SetupTables()

	userRepo := user2.NewUserRepository(db.DBCon)

	req, _ := http.NewRequest("GET", "/users/auth", nil)
	req.Header.Add("x-auth-user", "test-username-here")
	req.Header.Add("x-auth-key", "test-password-here")

	rr := httptest.NewRecorder()
	handler := NewAuthMiddleware(userRepo).Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_ReturnsUnauthorizedIfPasswordIsIncorrect(t *testing.T) {
	_ = db.Init(":memory:")
	db.SetupTables()

	userRepo := user2.NewUserRepository(db.DBCon)

	now := time.Now()

	password, _ := crypto.BcryptHashPassword("test-password-here")
	user := user2.User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	_, err := userRepo.Store(t.Context(), user)
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/users/auth", nil)
	req.Header.Add("x-auth-user", "test-username-here")
	req.Header.Add("x-auth-key", "incorrect-password")

	rr := httptest.NewRecorder()
	handler := NewAuthMiddleware(userRepo).Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
