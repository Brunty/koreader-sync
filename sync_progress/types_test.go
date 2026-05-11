package sync_progress

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStoreSyncProgressRequest_ValidatesSuccessfully(t *testing.T) {
	req := &StoreSyncProgressRequest{
		DeviceID:   "device-id-here",
		Progress:   "progress-here",
		Document:   "document-here",
		Percentage: 0.34,
		Device:     "device-here",
	}

	err := req.Validate()

	assert.NoError(t, err)
}

func TestStoreSyncProgressRequest_ValidatesFieldsMissing(t *testing.T) {
	req := &StoreSyncProgressRequest{
		DeviceID: "",
		Progress: "",
		Document: "",
		Device:   "",
	}

	err := req.Validate()

	assert.Error(t, err, "device_id is required")
	assert.Error(t, err, "progress is required")
	assert.Error(t, err, "document is required")
	assert.Error(t, err, "device is required")
}

func TestStoreSyncProgressRequest_MarshalToSyncProgress(t *testing.T) {
	req := &StoreSyncProgressRequest{
		DeviceID:   "device-id-here",
		Progress:   "progress-here",
		Document:   "document-here",
		Percentage: 0.34,
		Device:     "device-here",
	}

	progress, err := req.MarshalToSyncProgress(56)

	assert.NoError(t, err)
	assert.Equal(t, req.DeviceID, progress.DeviceID)
	assert.Equal(t, req.Progress, progress.Progress)
	assert.Equal(t, req.Document, progress.Document)
	assert.Equal(t, req.Percentage, progress.Percentage)
	assert.Equal(t, req.Device, progress.Device)
}

func TestSyncProgress_MarshalToReadResponse(t *testing.T) {
	progress := &SyncProgress{
		ID:         1234, // irrelevant to output testing, but here for completeness
		UserID:     2345, // irrelevant to output testing, but here for completeness
		Document:   "document-here",
		Progress:   "progress-here",
		Percentage: 0.78,
		Device:     "device-here",
		DeviceID:   "device-id-here",
		Timestamp:  time.Date(2026, 05, 10, 20, 34, 58, 651387237, time.UTC),
	}

	rsp, err := progress.MarshalToReadResponse()

	assert.NoError(t, err)
	assert.Equal(t, progress.DeviceID, rsp.DeviceID)
	assert.Equal(t, progress.Progress, rsp.Progress)
	assert.Equal(t, progress.Document, rsp.Document)
	assert.Equal(t, progress.Percentage, rsp.Percentage)
	assert.Equal(t, progress.Device, rsp.Device)
	assert.Equal(t, int64(1778445298), rsp.Timestamp)
}

func TestSyncProgress_MarshalNilToReadResponse(t *testing.T) {
	progress := &SyncProgress{}
	progress = nil

	rsp, err := progress.MarshalToReadResponse()

	assert.NoError(t, err)
	assert.Equal(t, "", rsp.DeviceID)
	assert.Equal(t, "", rsp.Progress)
	assert.Equal(t, "", rsp.Document)
	assert.Equal(t, 0.0, rsp.Percentage)
	assert.Equal(t, "", rsp.Device)
	assert.Equal(t, int64(0), rsp.Timestamp)
}
