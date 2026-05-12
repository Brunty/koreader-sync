package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brunty/koreader-sync-server/crypto"
	"github.com/brunty/koreader-sync-server/db"
	"github.com/brunty/koreader-sync-server/handlers"
	"github.com/stretchr/testify/assert"
)

func TestAuthUser(t *testing.T) {
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	userRepo := NewUserRepository(db.DBCon)

	userHandler := NewUserHandler(userRepo)

	// So, this test is basic because the auth user endpoint doesn't actually really do anything
	// It's protected by middleware.AuthMiddleware it doesn't actually need to do anything except return a
	// success message because if it's passed the middleware, it is authorized
	req, _ := http.NewRequest("GET", "/users/auth", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userHandler.AuthUser)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedRsp := &handlers.StatusResponse{Status: "authorized"}
	actualRsp := &handlers.StatusResponse{}
	json.Unmarshal(rr.Body.Bytes(), &actualRsp)
	assert.Equal(t, expectedRsp, actualRsp)
}

func TestCreateUser_Successfully(t *testing.T) {
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	userRepo := NewUserRepository(db.DBCon)

	userHandler := NewUserHandler(userRepo)

	reqBody := &CreateUserRequest{
		Username: "username-here",
		Password: "password-here",
	}

	jsonBody, _ := json.Marshal(reqBody)
	body := strings.NewReader(string(jsonBody))
	req, _ := http.NewRequest("POST", "/users/create", body)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userHandler.CreateUser)

	handler.ServeHTTP(rr, req)

	user, err := userRepo.SelectByUsername(t.Context(), "username-here")

	// Check the user was created
	assert.NoError(t, err)
	assert.Equal(t, "username-here", user.Username)
	assert.True(t, crypto.BcryptCheckPasswordHash("password-here", user.Password))

	assert.Equal(t, http.StatusCreated, rr.Code)

	expectedRsp := &handlers.StatusResponse{Status: "user created"}
	actualRsp := &handlers.StatusResponse{}
	json.Unmarshal(rr.Body.Bytes(), &actualRsp)
	assert.Equal(t, expectedRsp, actualRsp)
}

func TestCreateUser_FailsBlankUserDetails(t *testing.T) {
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	userRepo := NewUserRepository(db.DBCon)

	userHandler := NewUserHandler(userRepo)

	reqBody := &CreateUserRequest{
		Username: "",
		Password: "",
	}

	jsonBody, _ := json.Marshal(reqBody)
	body := strings.NewReader(string(jsonBody))
	req, _ := http.NewRequest("POST", "/users/create", body)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userHandler.CreateUser)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedRsp := &handlers.ErrorResponse{Error: "username is required\npassword is required"}
	actualRsp := &handlers.ErrorResponse{}
	json.Unmarshal(rr.Body.Bytes(), &actualRsp)
	assert.Equal(t, expectedRsp, actualRsp)
}

func TestCreateUser_FailsDuplicateUser(t *testing.T) {
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	userRepo := NewUserRepository(db.DBCon)

	userHandler := NewUserHandler(userRepo)

	password, _ := crypto.BcryptHashPassword("original-password-here")
	user := User{
		Username: "username-here",
		Password: password,
	}

	_, err := userRepo.Store(t.Context(), user)

	assert.NoError(t, err)

	reqBody := &CreateUserRequest{
		Username: "username-here",
		Password: "new-password-here",
	}

	jsonBody, _ := json.Marshal(reqBody)
	body := strings.NewReader(string(jsonBody))
	req, _ := http.NewRequest("POST", "/users/create", body)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userHandler.CreateUser)

	handler.ServeHTTP(rr, req)

	userFromDb, err := userRepo.SelectByUsername(t.Context(), "username-here")

	// Check the original user in the DB wasn't touched
	assert.NoError(t, err)
	assert.Equal(t, "username-here", userFromDb.Username)
	assert.False(t, crypto.BcryptCheckPasswordHash("new-password-here", userFromDb.Password))

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expectedRsp := &handlers.ErrorResponse{Error: "username is already taken"}
	actualRsp := &handlers.ErrorResponse{}
	json.Unmarshal(rr.Body.Bytes(), &actualRsp)
	assert.Equal(t, expectedRsp, actualRsp)
}

func TestCreateUser_FailsMarshalingError(t *testing.T) {
	userHandler := NewUserHandler(nil)

	reqBody := &CreateUserRequest{
		Username: "username-here",
		Password: strings.Repeat("a", 73),
	}
	jsonBody, _ := json.Marshal(reqBody)
	body := strings.NewReader(string(jsonBody))
	req, _ := http.NewRequest("POST", "/users/create", body)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userHandler.CreateUser)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	expectedRsp := &handlers.ErrorResponse{Error: "something went wrong"}
	actualRsp := &handlers.ErrorResponse{}
	json.Unmarshal(rr.Body.Bytes(), &actualRsp)
	assert.Equal(t, expectedRsp, actualRsp)
}
