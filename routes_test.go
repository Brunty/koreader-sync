package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/brunty/koreader-sync-server/db"
	"github.com/stretchr/testify/assert"
)

// These tests are more like API tests as I want to call via HTTP so that it's not just invoking handlers, that way
// we can check things like auth middleware works properly on appropriate routes

func TestRootEndpoint(t *testing.T) {
	mux := &ServeMux{http.NewServeMux()}
	mux.RegisterRoutes()

	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var result map[string]string
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	assert.Equal(t, "This is a server compatible with the KOReader Sync Protocol", result["status"])
}

func TestAuthEndpoint_Unauthorized(t *testing.T) {
	t.Cleanup(func() {
		_ = os.Remove("./data/data.test.db.sqlite3")
	})

	err := db.Init("./data/data.test.db.sqlite3")
	assert.NoError(t, err)

	db.SetupTables()

	mux := &ServeMux{http.NewServeMux()}
	mux.RegisterRoutes()

	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/users/auth")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	slog.Debug(string(body))
	assert.NoError(t, err)

	var result map[string]string
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	assert.Equal(t, "unauthorized", result["error"])
}
