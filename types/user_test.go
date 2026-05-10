package types

import (
	"testing"

	"github.com/brunty/koreader-sync-server/crypto"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserRequest_ValidatesFieldsMissing(t *testing.T) {
	req := &CreateUserRequest{
		Username: "",
		Password: "",
	}

	err := req.Validate()

	assert.Error(t, err, "username is required")
	assert.Error(t, err, "password is required")
}

func TestCreateUserRequest_ValidatesSuccessfully(t *testing.T) {
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
