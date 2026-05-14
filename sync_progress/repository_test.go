package sync_progress

import (
	"testing"
	"time"

	"github.com/brunty/koreader-sync-server/db"
	"github.com/stretchr/testify/assert"
)

func TestStoreAndSelectProgress(t *testing.T) {
	_ = db.Init(":memory:")
	db.SetupTables()

	repo := NewSyncProgressRepository(db.DBCon)

	now := time.Now()

	progress := SyncProgress{
		UserID:     1,
		Document:   "document-here",
		Progress:   "progress-here",
		Percentage: 0.35,
		Device:     "device-here",
		DeviceID:   "device-id-here",
		Timestamp:  now,
	}

	_, err := repo.Store(t.Context(), progress)

	assert.NoError(t, err)

	progressFromDb, err := repo.SelectByUserIDAndDocument(t.Context(), 1, "document-here")

	assert.NoError(t, err)

	assert.Equal(t, progress.UserID, progressFromDb.UserID)
	assert.Equal(t, progress.Document, progressFromDb.Document)
	assert.Equal(t, progress.Progress, progressFromDb.Progress)
	assert.Equal(t, progress.Percentage, progressFromDb.Percentage)
	assert.Equal(t, progress.Device, progressFromDb.Device)
	assert.Equal(t, progress.DeviceID, progressFromDb.DeviceID)
}

func TestSelectProgressNotFound(t *testing.T) {
	_ = db.Init(":memory:")
	db.SetupTables()

	repo := NewSyncProgressRepository(db.DBCon)

	progressFromDb, err := repo.SelectByUserIDAndDocument(t.Context(), 1, "document-here")

	assert.NoError(t, err)
	assert.Nil(t, progressFromDb)
}

func TestSelectByUserIDAndDocument_DBError(t *testing.T) {
	_ = db.Init(":memory:")
	db.SetupTables()

	repo := NewSyncProgressRepository(db.DBCon)
	// Close the DB before trying to access it will cause an error which we want to test
	_ = db.DBCon.Close()

	progress, err := repo.SelectByUserIDAndDocument(t.Context(), 1, "document-here")
	assert.Error(t, err)
	assert.Nil(t, progress)
}

func TestStoreProgress_DBError(t *testing.T) {
	_ = db.Init(":memory:")
	db.SetupTables()

	repo := NewSyncProgressRepository(db.DBCon)
	// Close the DB before trying to access it will cause an error which we want to test
	_ = db.DBCon.Close()

	progress := SyncProgress{
		UserID:     1,
		Document:   "document-here",
		Progress:   "progress-here",
		Percentage: 0.35,
		Device:     "device-here",
		DeviceID:   "device-id-here",
		Timestamp:  time.Now(),
	}

	id, err := repo.Store(t.Context(), progress)
	assert.Error(t, err)
	assert.Nil(t, id)
}
