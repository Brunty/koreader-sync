package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brunty/koreader-sync-server/crypto"
	"github.com/brunty/koreader-sync-server/dao"
	"github.com/brunty/koreader-sync-server/db"
	"github.com/brunty/koreader-sync-server/types"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_PassesToNextHandlerOnSuccessfulAuthentication(t *testing.T) {
	db.Init(":memory:")
	db.CreateTables()
	defer db.DBCon.Close()

	now := time.Now()

	password, _ := crypto.HashPassword("test-password-here")
	user := types.User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	err := dao.StoreUser(user)
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/users/auth", nil)
	req.Header.Add("x-auth-user", "test-username-here")
	req.Header.Add("x-auth-key", "test-password-here")

	rr := httptest.NewRecorder()
	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_ReturnsUnauthorizedIfUsernameIsBlank(t *testing.T) {
	db.Init(":memory:")
	db.CreateTables()
	defer db.DBCon.Close()

	now := time.Now()

	password, _ := crypto.HashPassword("test-password-here")
	user := types.User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	err := dao.StoreUser(user)
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/users/auth", nil)
	req.Header.Add("x-auth-key", "test-password-here")

	rr := httptest.NewRecorder()
	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_ReturnsUnauthorizedIfPasswordIsBlank(t *testing.T) {
	db.Init(":memory:")
	db.CreateTables()
	defer db.DBCon.Close()

	now := time.Now()

	password, _ := crypto.HashPassword("test-password-here")
	user := types.User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	err := dao.StoreUser(user)
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/users/auth", nil)
	req.Header.Add("x-auth-user", "test-username-here")

	rr := httptest.NewRecorder()
	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_ReturnsUnauthorizedIfUserIsNotFound(t *testing.T) {
	db.Init(":memory:")
	db.CreateTables()
	defer db.DBCon.Close()

	req, _ := http.NewRequest("GET", "/users/auth", nil)
	req.Header.Add("x-auth-user", "test-username-here")
	req.Header.Add("x-auth-key", "test-password-here")

	rr := httptest.NewRecorder()
	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_ReturnsUnauthorizedIfPasswordIsIncorrect(t *testing.T) {
	db.Init(":memory:")
	db.CreateTables()
	defer db.DBCon.Close()

	now := time.Now()

	password, _ := crypto.HashPassword("test-password-here")
	user := types.User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	err := dao.StoreUser(user)
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/users/auth", nil)
	req.Header.Add("x-auth-user", "test-username-here")
	req.Header.Add("x-auth-key", "incorrect-password")

	rr := httptest.NewRecorder()
	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
