package sync_progress

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/brunty/koreader-sync-server/crypto"
	"github.com/brunty/koreader-sync-server/db"
	"github.com/brunty/koreader-sync-server/handlers"
	userpackage "github.com/brunty/koreader-sync-server/user"
	"github.com/stretchr/testify/assert"
)

func TestReadSyncProgress_Successfully(t *testing.T) {
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	syncProgressRepo := NewSyncProgressRepository(db.DBCon)
	userRepo := userpackage.NewUserRepository(db.DBCon)

	syncHandler := NewSyncProgressHandler(syncProgressRepo)

	now := time.Now()
	password, _ := crypto.BcryptHashPassword("test-password-here")
	user := userpackage.User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	userID, err := userRepo.Store(user)
	assert.NoError(t, err)

	progress := SyncProgress{
		UserID:     *userID,
		Document:   "document-here",
		Progress:   "progress-here",
		Percentage: 0.34,
		Device:     "device-here",
		DeviceID:   "device-id-here",
		Timestamp:  now,
	}

	_, err = syncProgressRepo.Store(progress)

	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/syncs/progress/document-here", nil)
	req.Header.Add("x-auth-user", "test-username-here")
	req.Header.Add("x-auth-key", "test-password-here")
	req.SetPathValue("document", "document-here")
	req = req.WithContext(context.WithValue(req.Context(), "user", *userID))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(syncHandler.ReadSyncProgress)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedRsp := &ReadSyncProgressResponse{
		Document:   "document-here",
		Progress:   "progress-here",
		Percentage: 0.34,
		Device:     "device-here",
		DeviceID:   "device-id-here",
		Timestamp:  now.Unix(),
	}
	actualRsp := &ReadSyncProgressResponse{}
	json.Unmarshal(rr.Body.Bytes(), &actualRsp)
	assert.Equal(t, expectedRsp, actualRsp)
}

func TestReadSyncProgress_SyncNotFoundInDB(t *testing.T) {
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	syncProgressRepo := NewSyncProgressRepository(db.DBCon)
	userRepo := userpackage.NewUserRepository(db.DBCon)

	syncHandler := NewSyncProgressHandler(syncProgressRepo)

	now := time.Now()
	password, _ := crypto.BcryptHashPassword("test-password-here")
	user := userpackage.User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	userID, err := userRepo.Store(user)
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/syncs/progress/document-here", nil)
	req.Header.Add("x-auth-user", "test-username-here")
	req.Header.Add("x-auth-key", "test-password-here")
	req.SetPathValue("document", "document-here")
	req = req.WithContext(context.WithValue(req.Context(), "user", *userID))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(syncHandler.ReadSyncProgress)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	expectedRsp := &ReadSyncProgressResponse{
		Document:   "",
		Progress:   "",
		Percentage: 0,
		Device:     "",
		DeviceID:   "",
		Timestamp:  0,
	}
	actualRsp := &ReadSyncProgressResponse{}
	json.Unmarshal(rr.Body.Bytes(), &actualRsp)
	assert.Equal(t, expectedRsp, actualRsp)
}

func TestGetSyncProgress_SyncNotFoundNoURLParam(t *testing.T) {
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	syncProgressRepo := NewSyncProgressRepository(db.DBCon)
	userRepo := userpackage.NewUserRepository(db.DBCon)

	syncHandler := NewSyncProgressHandler(syncProgressRepo)

	now := time.Now()
	password, _ := crypto.BcryptHashPassword("test-password-here")
	user := userpackage.User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	userID, err := userRepo.Store(user)
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/syncs/progress/document-here", nil)
	req.Header.Add("x-auth-user", "test-username-here")
	req.Header.Add("x-auth-key", "test-password-here")
	req = req.WithContext(context.WithValue(req.Context(), "user", *userID))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(syncHandler.ReadSyncProgress)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	expectedRsp := &ReadSyncProgressResponse{
		Document:   "",
		Progress:   "",
		Percentage: 0,
		Device:     "",
		DeviceID:   "",
		Timestamp:  0,
	}
	actualRsp := &ReadSyncProgressResponse{}
	json.Unmarshal(rr.Body.Bytes(), &actualRsp)
	assert.Equal(t, expectedRsp, actualRsp)
}

func TestStoreSyncProgress_SuccessfulUpdateProgress(t *testing.T) {
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	syncProgressRepo := NewSyncProgressRepository(db.DBCon)
	userRepo := userpackage.NewUserRepository(db.DBCon)

	syncHandler := NewSyncProgressHandler(syncProgressRepo)

	now := time.Now()
	password, _ := crypto.BcryptHashPassword("test-password-here")
	user := userpackage.User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	userID, err := userRepo.Store(user)
	assert.NoError(t, err)

	progress := SyncProgress{
		UserID:     *userID,
		Document:   "document-here",
		Progress:   "progress-here",
		Percentage: 0.34,
		Device:     "device-here",
		DeviceID:   "device-id-here",
		Timestamp:  now,
	}

	_, err = syncProgressRepo.Store(progress)

	assert.NoError(t, err)

	reqBody := &StoreSyncProgressRequest{
		DeviceID:   "device-here",
		Progress:   "new-progress-here",
		Document:   "document-here",
		Percentage: 0.45,
		Device:     "device-id-here",
	}
	jsonBody, _ := json.Marshal(reqBody)
	body := strings.NewReader(string(jsonBody))
	req, _ := http.NewRequest("PUT", "/syncs/progress", body)
	req.Header.Add("x-auth-user", "test-username-here")
	req.Header.Add("x-auth-key", "test-password-here")
	req = req.WithContext(context.WithValue(req.Context(), "user", *userID))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(syncHandler.StoreSyncProgress)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedRsp := &handlers.StatusResponse{Status: "sync stored"}
	actualRsp := &handlers.StatusResponse{}
	json.Unmarshal(rr.Body.Bytes(), &actualRsp)
	assert.Equal(t, expectedRsp, actualRsp)

	progressFromDB, err := syncProgressRepo.SelectByUserIDAndDocument(*userID, "document-here")

	assert.NoError(t, err)
	assert.Equal(t, reqBody.Percentage, progressFromDB.Percentage)
	assert.Equal(t, reqBody.Progress, progressFromDB.Progress)
}
