package user

import (
	"testing"
	"time"

	"github.com/brunty/koreader-sync-server/db"
	"github.com/stretchr/testify/assert"
)

func setupInMemoryDb() {
	db.Init(":memory:")
	db.SetupTables()
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

	_, err := repo.Store(t.Context(), user)

	assert.NoError(t, err)

	userFromDb, err := repo.SelectByUsername(t.Context(), "test-username-here")

	assert.NoError(t, err)
	assert.Equal(t, user.Username, userFromDb.Username)
	assert.Equal(t, user.Password, userFromDb.Password)
}

func TestSelectUserNotFound(t *testing.T) {
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	repo := NewUserRepository(db.DBCon)

	userFromDb, err := repo.SelectByUsername(t.Context(), "test-username-here")

	assert.NoError(t, err)
	assert.Nil(t, userFromDb)
}

func TestStoreUser_Duplicate(t *testing.T) {
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	repo := NewUserRepository(db.DBCon)

	user := User{
		Username:  "test-user",
		Password:  "test-password",
		CreatedAt: time.Now(),
	}

	_, err := repo.Store(t.Context(), user)
	assert.NoError(t, err)

	_, err = repo.Store(t.Context(), user)
	assert.Error(t, err)
}

func TestSelectByUsername_DBError(t *testing.T) {
	setupInMemoryDb()

	repo := NewUserRepository(db.DBCon)
	db.DBCon.Close()

	user, err := repo.SelectByUsername(t.Context(), "test-user")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestStoreUser_DBError(t *testing.T) {
	setupInMemoryDb()

	repo := NewUserRepository(db.DBCon)
	db.DBCon.Close()

	user := User{
		Username:  "test-user",
		Password:  "test-password",
		CreatedAt: time.Now(),
	}

	id, err := repo.Store(t.Context(), user)
	assert.Error(t, err)
	assert.Nil(t, id)
}

func TestUpdateUser_Success(t *testing.T) {
	setupInMemoryDb()
	defer db.EmptyTables()
	defer db.DBCon.Close()

	repo := NewUserRepository(db.DBCon)

	user := User{
		Username:  "test-user",
		Password:  "original-password",
		CreatedAt: time.Now(),
	}

	_, err := repo.Store(t.Context(), user)
	assert.NoError(t, err)

	user.Password = "updated-password"
	id, err := repo.Update(t.Context(), user)
	assert.NoError(t, err)
	assert.NotNil(t, id)

	userFromDb, err := repo.SelectByUsername(t.Context(), "test-user")
	assert.NoError(t, err)
	assert.Equal(t, "updated-password", userFromDb.Password)
}

func TestUpdateUser_DBError(t *testing.T) {
	setupInMemoryDb()

	repo := NewUserRepository(db.DBCon)
	db.DBCon.Close()

	user := User{
		Username: "test-user",
		Password: "test-password",
	}

	id, err := repo.Update(t.Context(), user)
	assert.Error(t, err)
	assert.Nil(t, id)
}
