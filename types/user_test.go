package types

import (
	"testing"

	"github.com/brunty/koreader-sync-server/crypto"
	"github.com/stretchr/testify/assert"
)

func TestUserType_ValidatesBlankUsername(t *testing.T) {
	req := &CreateUserRequest{
		Username: "",
		Password: "Password",
	}

	err := req.Validate()

	assert.Error(t, err, "username is required")
}

func TestUserType_ValidatesBlankPassword(t *testing.T) {
	req := &CreateUserRequest{
		Username: "username",
		Password: "",
	}

	err := req.Validate()

	assert.Error(t, err, "password is required")
}

func TestUserType_ValidatesBothBlankFields(t *testing.T) {
	req := &CreateUserRequest{
		Username: "",
		Password: "",
	}

	err := req.Validate()

	assert.Error(t, err, "username is required")
	assert.Error(t, err, "password is required")
}

func TestUserType_ValidatesSuccessfully(t *testing.T) {
	req := &CreateUserRequest{
		Username: "username",
		Password: "password",
	}

	err := req.Validate()

	assert.NoError(t, err)
}

func TestCreateUserRequest_MarshalToUser(t *testing.T) {
	req := &CreateUserRequest{
		Username: "username",
		Password: "password",
	}

	user, err := req.MarshalToUser()

	assert.NoError(t, err)

	assert.Equal(t, req.Username, user.Username)
	assert.True(t, crypto.CheckPasswordHash(req.Password, user.Password))
}
