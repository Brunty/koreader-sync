package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/brunty/koreader-sync-server/db"
	"github.com/stretchr/testify/assert"
)

// These tests are more like acceptance tests as I want to call via HTTP so that it's not just invoking handlers, that
// way we can check things like auth middleware works properly on appropriate routes and the whole flows work

func setupAcceptanceTests(t *testing.T) *httptest.Server {
	t.Cleanup(func() {
		_ = os.Remove("./data/data.test.db.sqlite3")
	})

	err := db.Init("./data/data.test.db.sqlite3")
	assert.NoError(t, err)

	db.SetupTables()

	mux := &ServeMux{http.NewServeMux()}
	mux.RegisterRoutes()

	return httptest.NewServer(mux)
}

func marshalResponseBodyToMap(t *testing.T, rsp *http.Response) map[string]string {
	body, err := io.ReadAll(rsp.Body)
	slog.Debug(string(body))
	assert.NoError(t, err)

	var result map[string]string
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	return result
}

func TestHome(t *testing.T) {
	server := setupAcceptanceTests(t)
	defer server.Close()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/", server.URL), nil)

	client := &http.Client{}
	rsp, err := client.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rsp.StatusCode)

	body := marshalResponseBodyToMap(t, rsp)
	assert.Equal(t, "This is a server compatible with the KOReader Sync Protocol", body["status"])
}

func TestNotFound(t *testing.T) {
	server := setupAcceptanceTests(t)
	defer server.Close()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/bad-url-here", server.URL), nil)

	client := &http.Client{}
	rsp, err := client.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, rsp.StatusCode)

	body := marshalResponseBodyToMap(t, rsp)
	assert.Equal(t, "not found", body["error"])
}

func TestUnauthorizedEndpoints(t *testing.T) {
	server := setupAcceptanceTests(t)
	defer server.Close()

	client := &http.Client{}

	endpoints := map[string]struct {
		method string
		url    string
	}{
		"user auth": {
			"GET",
			"/users/auth",
		},
		"sync document progress": {
			"PUT",
			"/syncs/progress",
		},
		"get document progress": {
			"GET",
			"/syncs/progress/document-here",
		},
	}

	for testName, tt := range endpoints {
		t.Run(testName, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, fmt.Sprintf("%s%s", server.URL, tt.url), nil)
			authRsp, err := client.Do(req)
			assert.NoError(t, err)

			assert.Equal(t, http.StatusUnauthorized, authRsp.StatusCode)

			body := marshalResponseBodyToMap(t, authRsp)
			assert.Equal(t, "unauthorized", body["error"])
		})
	}
}

func TestUserCreateAndAuthFlow(t *testing.T) {
	server := setupAcceptanceTests(t)
	defer server.Close()

	client := &http.Client{}

	createBody := `{"username":"test-user","password":"test- pass"}`
	createReq, _ := http.NewRequest("POST", server.URL+"/users/create", strings.NewReader(createBody))
	createRsp, err := client.Do(createReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, createRsp.StatusCode)

	authReq, _ := http.NewRequest("GET", server.URL+"/users/auth", nil)
	authReq.Header.Set("x-auth-user", "test-user")
	authReq.Header.Set("x-auth-key", "test- pass")
	authRsp, err := client.Do(authReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, authRsp.StatusCode)

	body := marshalResponseBodyToMap(t, authRsp)
	assert.Equal(t, "authorized", body["status"])
}
