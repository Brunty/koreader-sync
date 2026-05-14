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
	"time"

	"github.com/brunty/koreader-sync-server/crypto"
	"github.com/brunty/koreader-sync-server/db"
	userpackage "github.com/brunty/koreader-sync-server/user"
	"github.com/stretchr/testify/assert"
)

const TestDBName = "/tmp/koreader-sync/data.test.db.sqlite3"

// These tests are more like acceptance tests as I want to call via HTTP so that it's not just invoking handlers, that
// way we can check things like auth middleware works properly on appropriate routes and the whole flows work
// These aren't testing full logic of the handlers (the handler tests are better for that) these are to ensure that
// things work within the HTTP context and setting the user on the request context is handled correctly

func setupAcceptanceTests(t *testing.T) *httptest.Server {
	t.Cleanup(func() {
		_ = os.Mkdir("/tmp/koreader-sync", 0755)
		_ = os.Remove(TestDBName)
	})

	err := db.Init(TestDBName)
	assert.NoError(t, err)

	db.SetupTables()

	mux := &ServeMux{http.NewServeMux()}
	mux.RegisterRoutes()

	return httptest.NewServer(mux)
}

func marshalResponseBodyToMap(t *testing.T, rsp *http.Response) map[string]any {
	body, err := io.ReadAll(rsp.Body)
	slog.Debug(string(body))
	assert.NoError(t, err)

	var result map[string]any
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

	for desc, tt := range endpoints {
		t.Run(desc, func(t *testing.T) {
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

	createBody := `{"username":"test-user","password":"test-pass"}`
	createReq, _ := http.NewRequest("POST", server.URL+"/users/create", strings.NewReader(createBody))
	createRsp, err := client.Do(createReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, createRsp.StatusCode)

	authReq, _ := http.NewRequest("GET", server.URL+"/users/auth", nil)
	authReq.Header.Set("x-auth-user", "test-user")
	authReq.Header.Set("x-auth-key", "test-pass")
	authRsp, err := client.Do(authReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, authRsp.StatusCode)

	body := marshalResponseBodyToMap(t, authRsp)
	assert.Equal(t, "authorized", body["status"])
}

func TestUserSyncProgressFlow(t *testing.T) {
	server := setupAcceptanceTests(t)
	defer server.Close()

	repo := userpackage.NewUserRepository(db.DBCon)

	password, _ := crypto.BcryptHashPassword("test-pass")

	user := userpackage.User{
		Username:  "test-user",
		Password:  password,
		CreatedAt: time.Now(),
	}

	_, err := repo.Store(t.Context(), user)

	client := &http.Client{}

	syncProgressBody := `{"device_id":"my-device","progress":"/some/progress[2]/here","document":"super-book","percentage":0.123,"device":"my-reader"}`
	createReq, _ := http.NewRequest("PUT", server.URL+"/syncs/progress", strings.NewReader(syncProgressBody))
	createReq.Header.Set("x-auth-user", "test-user")
	createReq.Header.Set("x-auth-key", "test-pass")
	createRsp, err := client.Do(createReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, createRsp.StatusCode)

	getSyncProgressReq, _ := http.NewRequest("GET", server.URL+"/syncs/progress/super-book", nil)
	getSyncProgressReq.Header.Set("x-auth-user", "test-user")
	getSyncProgressReq.Header.Set("x-auth-key", "test-pass")
	authRsp, err := client.Do(getSyncProgressReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, authRsp.StatusCode)

	body := marshalResponseBodyToMap(t, authRsp)
	assert.Equal(t, "my-device", body["device_id"])
	assert.Equal(t, "/some/progress[2]/here", body["progress"])
	assert.Equal(t, "super-book", body["document"])
	assert.Equal(t, 0.123, body["percentage"])
	assert.Equal(t, "my-reader", body["device"])
	assert.NotEqual(t, 0, body["timestamp"])
}
