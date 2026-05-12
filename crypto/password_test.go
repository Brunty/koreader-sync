package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHashing(t *testing.T) {
	password := "password123" // so secure amirite?!
	hashedPassword, err := BcryptHashPassword(password)

	assert.NoErrorf(t, err, "Should have no error from hashing the crypto")
	assert.True(t, BcryptCheckPasswordHash(password, hashedPassword))
}
