package dao

import (
	"testing"
	"time"

	"github.com/brunty/koreader-sync-server/db"
	"github.com/brunty/koreader-sync-server/types"
	"github.com/stretchr/testify/assert"
)

func TestStoreAndSelectUser(t *testing.T) {
	db.Init(":memory:")
	db.CreateTables()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	now := time.Now()

	user := types.User{
		Username:  "test-username-here",
		Password:  "test-password-here",
		CreatedAt: now,
	}

	_, err := StoreUser(user)

	assert.NoError(t, err)

	userFromDb, err := SelectUserByUsername("test-username-here")

	assert.NoError(t, err)
	assert.Equal(t, user.Username, userFromDb.Username)
	assert.Equal(t, user.Password, userFromDb.Password)
}

func TestSelectUserNotFound(t *testing.T) {
	db.Init(":memory:")
	db.CreateTables()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	userFromDb, err := SelectUserByUsername("test-username-here")

	assert.NoError(t, err)
	assert.Nil(t, userFromDb)
}
