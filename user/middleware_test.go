package user

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brunty/koreader-sync-server/crypto"
	"github.com/brunty/koreader-sync-server/db"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_PassesToNextHandlerOnSuccessfulAuthentication(t *testing.T) {
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	userRepo := NewUserRepository(db.DBCon)

	now := time.Now()

	password, _ := crypto.HashPassword("test-password-here")
	user := User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	_, err := userRepo.Store(user)
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
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	userRepo := NewUserRepository(db.DBCon)

	now := time.Now()

	password, _ := crypto.HashPassword("test-password-here")
	user := User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	_, err := userRepo.Store(user)
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
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	userRepo := NewUserRepository(db.DBCon)

	now := time.Now()

	password, _ := crypto.HashPassword("test-password-here")
	user := User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	_, err := userRepo.Store(user)
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
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	userRepo := NewUserRepository(db.DBCon)

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
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	userRepo := NewUserRepository(db.DBCon)

	now := time.Now()

	password, _ := crypto.HashPassword("test-password-here")
	user := User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	_, err := userRepo.Store(user)
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
