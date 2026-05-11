package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brunty/koreader-sync-server/types"
	"github.com/stretchr/testify/assert"
)

func TestAuthUser(t *testing.T) {
	// So, this test is basic because the auth user endpoint doesn't actually really do anything
	// It's protected by middleware.AuthMiddleware it doesn't actually need to do anything except return a
	// success message because if it's passed the middleware, it is authorized
	req, _ := http.NewRequest("GET", "/users/auth", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AuthUser)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedRsp := &types.StatusResponse{Status: "authorized"}
	actualRsp := &types.StatusResponse{}
	json.Unmarshal(rr.Body.Bytes(), &actualRsp)
	assert.Equal(t, expectedRsp, actualRsp)
}
