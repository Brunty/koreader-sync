package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/brunty/koreader-sync-server/crypto"
	"github.com/brunty/koreader-sync-server/dao"
	"github.com/brunty/koreader-sync-server/db"
	"github.com/brunty/koreader-sync-server/types"
	"github.com/stretchr/testify/assert"
)

func TestStoreSyncProgress_SuccessfulUpdateProgress(t *testing.T) {
	db.Init(":memory:")
	db.CreateTables()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	now := time.Now()
	password, _ := crypto.HashPassword("test-password-here")
	user := types.User{
		Username:  "test-username-here",
		Password:  password,
		CreatedAt: now,
	}
	userID, err := dao.StoreUser(user)
	assert.NoError(t, err)

	progress := types.Progress{
		UserID:     *userID,
		Document:   "document-here",
		Progress:   "progress-here",
		Percentage: 0.34,
		Device:     "device-here",
		DeviceID:   "device-id-here",
		Timestamp:  now,
	}

	_, err = dao.StoreProgress(progress)

	assert.NoError(t, err)

	reqBody := &types.SyncProgressRequest{
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
	handler := http.HandlerFunc(StoreSyncProgress)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedRsp := &types.StatusResponse{Status: "sync stored"}
	actualRsp := &types.StatusResponse{}
	json.Unmarshal(rr.Body.Bytes(), &actualRsp)
	assert.Equal(t, expectedRsp, actualRsp)

	progressFromDB, err := dao.SelectProgress(*userID, "document-here")

	assert.NoError(t, err)
	assert.Equal(t, reqBody.Percentage, progressFromDB.Percentage)
	assert.Equal(t, reqBody.Progress, progressFromDB.Progress)
}
