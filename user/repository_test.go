package user

import (
	"testing"
	"time"

	"github.com/brunty/koreader-sync-server/db"
	"github.com/stretchr/testify/assert"
)

func setupInMemoryDb() {
	db.Init(":memory:")
	db.CreateTables()
}

func TestStoreAndSelectUser(t *testing.T) {
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	repo := NewUserRepository(db.DBCon)

	now := time.Now()

	user := User{
		Username:  "test-username-here",
		Password:  "test-password-here",
		CreatedAt: now,
	}

	_, err := repo.Store(user)

	assert.NoError(t, err)

	userFromDb, err := repo.SelectByUsername("test-username-here")

	assert.NoError(t, err)
	assert.Equal(t, user.Username, userFromDb.Username)
	assert.Equal(t, user.Password, userFromDb.Password)
}

func TestSelectUserNotFound(t *testing.T) {
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	repo := NewUserRepository(db.DBCon)

	userFromDb, err := repo.SelectByUsername("test-username-here")

	assert.NoError(t, err)
	assert.Nil(t, userFromDb)
}
