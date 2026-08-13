package crypto

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHashing(t *testing.T) {
	password := "password123" // so secure amirite?!
	hashedPassword, err := BcryptHashPassword(password)

	assert.NoErrorf(t, err, "Should have no error from hashing the crypto")
	assert.True(t, BcryptCheckPasswordHash(password, hashedPassword))
}

func TestPasswordHashingErrorsOnStringTooLong(t *testing.T) {
	password := strings.Repeat("a", 100)
	_, err := BcryptHashPassword(password)

	assert.Error(t, err, "Should have errored with a string too long")
}

func TestBcryptNeedsRehash(t *testing.T) {
	// Has a cost of 4
	oldHashedPw := "$2y$04$S4sNf81050H0lCz8v3Hfh.tnKtTJGoZ63Z4qUCL//f284CS5i7d8q"

	// Checks against BcryptCost which wouldn't be set to 4
	assert.True(t, BcryptNeedsRehash(oldHashedPw))
}

func TestBcryptDoesNotNeedRehash(t *testing.T) {
	newHashedPassword, _ := BcryptHashPassword("password")

	assert.False(t, BcryptNeedsRehash(newHashedPassword))
}

func TestBcryptDoesNotNeedRehashIfHashError(t *testing.T) {
	assert.False(t, BcryptNeedsRehash("something bad here"))
}
