package dao

import (
	"testing"
	"time"

	"github.com/brunty/koreader-sync-server/db"
	"github.com/brunty/koreader-sync-server/types"
	"github.com/stretchr/testify/assert"
)

func TestStoreAndSelectProgress(t *testing.T) {
	db.Init(":memory:")
	db.CreateTables()
	defer db.DBCon.Close()

	now := time.Now()

	progress := types.Progress{
		UserID:     1,
		Document:   "document-here",
		Progress:   "progress-here",
		Percentage: 0.35,
		Device:     "device-here",
		DeviceID:   "device-id-here",
		Timestamp:  now,
	}

	err := StoreProgress(progress)

	assert.NoError(t, err)

	progressFromDb, err := SelectProgress(1, "document-here")

	assert.NoError(t, err)

	assert.Equal(t, progress.UserID, progressFromDb.UserID)
	assert.Equal(t, progress.Document, progressFromDb.Document)
	assert.Equal(t, progress.Progress, progressFromDb.Progress)
	assert.Equal(t, progress.Percentage, progressFromDb.Percentage)
	assert.Equal(t, progress.Device, progressFromDb.Device)
	assert.Equal(t, progress.DeviceID, progressFromDb.DeviceID)
}

func TestSelectProgressNotFound(t *testing.T) {
	db.Init(":memory:")
	db.CreateTables()
	defer db.DBCon.Close()

	progressFromDb, err := SelectProgress(1, "document-here")

	assert.NoError(t, err)
	assert.Nil(t, progressFromDb)
}
