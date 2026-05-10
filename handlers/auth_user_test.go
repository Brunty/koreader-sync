package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthUser(t *testing.T) {
	// So, this test is a basic one because the auth user endpoint doesn't actually really do anything
	// Because it's protected by middleware.AuthMiddleware it doesn't actually need to do anything except return a
	// success message because if it's passed the middleware, it's auth'd
	req, _ := http.NewRequest("GET", "/users/auth", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AuthUser)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
