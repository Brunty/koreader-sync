package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brunty/koreader-sync-server/crypto"
	"github.com/brunty/koreader-sync-server/dao"
	"github.com/brunty/koreader-sync-server/db"
	"github.com/brunty/koreader-sync-server/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Successfully(t *testing.T) {
	db.Init(":memory:")
	db.CreateTables()
	defer db.DBCon.Close()

	reqBody := &types.CreateUserRequest{
		Username: "username-here",
		Password: "password-here",
	}

	jsonBody, _ := json.Marshal(reqBody)
	body := strings.NewReader(string(jsonBody))
	req, _ := http.NewRequest("GET", "/users/create", body)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateUser)

	handler.ServeHTTP(rr, req)

	user, err := dao.SelectUserByUsername("username-here")

	// Check the user was created
	assert.NoError(t, err)
	assert.Equal(t, "username-here", user.Username)
	assert.True(t, crypto.CheckPasswordHash("password-here", user.Password))

	assert.Equal(t, http.StatusCreated, rr.Code)

	expectedRsp := &types.StatusResponse{Status: "user created"}
	actualRsp := &types.StatusResponse{}
	json.Unmarshal(rr.Body.Bytes(), &actualRsp)
	assert.Equal(t, expectedRsp, actualRsp)
}

func TestCreateUser_FailsBlankUserDetails(t *testing.T) {
	db.Init(":memory:")
	db.CreateTables()
	defer db.DBCon.Close()

	reqBody := &types.CreateUserRequest{
		Username: "",
		Password: "",
	}

	jsonBody, _ := json.Marshal(reqBody)
	body := strings.NewReader(string(jsonBody))
	req, _ := http.NewRequest("GET", "/users/create", body)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateUser)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedRsp := &types.ErrorResponse{Error: "username is required\npassword is required"}
	actualRsp := &types.ErrorResponse{}
	json.Unmarshal(rr.Body.Bytes(), &actualRsp)
	assert.Equal(t, expectedRsp, actualRsp)
}

func TestCreateUser_FailsDuplicateUser(t *testing.T) {
	db.Init(":memory:")
	db.CreateTables()
	defer db.DBCon.Close()

	password, _ := crypto.HashPassword("original-password-here")
	user := types.User{
		Username: "username-here",
		Password: password,
	}

	_, err := dao.StoreUser(user)

	assert.NoError(t, err)

	reqBody := &types.CreateUserRequest{
		Username: "username-here",
		Password: "new-password-here",
	}

	jsonBody, _ := json.Marshal(reqBody)
	body := strings.NewReader(string(jsonBody))
	req, _ := http.NewRequest("GET", "/users/create", body)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateUser)

	handler.ServeHTTP(rr, req)

	userFromDb, err := dao.SelectUserByUsername("username-here")

	// Check the original user in the DB wasn't touched
	assert.NoError(t, err)
	assert.Equal(t, "username-here", userFromDb.Username)
	assert.False(t, crypto.CheckPasswordHash("new-password-here", userFromDb.Password))

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedRsp := &types.ErrorResponse{Error: "username is already taken"}
	actualRsp := &types.ErrorResponse{}
	json.Unmarshal(rr.Body.Bytes(), &actualRsp)
	assert.Equal(t, expectedRsp, actualRsp)
}
