package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSyncProgressRequest_ValidatesFieldsMissing(t *testing.T) {
	req := &SyncProgressRequest{
		DeviceId:   "device-id-here",
		Progress:   "progress-here",
		Document:   "document-here",
		Percentage: 0.34,
		Device:     "device-here",
	}

	err := req.Validate()

	assert.NoError(t, err)
}

func TestSyncProgressRequest_ValidatesSuccessfully(t *testing.T) {
	req := &SyncProgressRequest{
		DeviceId: "",
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

func TestSyncProgressRequest_MarshalToProgress(t *testing.T) {
	req := &SyncProgressRequest{
		DeviceId:   "device-id-here",
		Progress:   "progress-here",
		Document:   "document-here",
		Percentage: 0.34,
		Device:     "device-here",
	}

	progress, err := req.MarshalToProgress(56)

	assert.NoError(t, err)
	assert.Equal(t, req.DeviceId, progress.DeviceId)
	assert.Equal(t, req.Progress, progress.Progress)
	assert.Equal(t, req.Document, progress.Document)
	assert.Equal(t, req.Percentage, progress.Percentage)
	assert.Equal(t, req.Device, progress.Device)
}

func TestProgress_MarshalToResponse(t *testing.T) {
	progress := &Progress{
		Id:         1234, // irrelevant to output testing, but here for completeness
		UserId:     2345, // irrelevant to output testing, but here for completeness
		Document:   "document-here",
		Progress:   "progress-here",
		Percentage: 0.78,
		Device:     "device-here",
		DeviceId:   "device-id-here",
		Timestamp:  time.Date(2026, 05, 10, 20, 34, 58, 651387237, time.UTC),
	}

	rsp, err := progress.MarshalToResponse()

	assert.NoError(t, err)
	assert.Equal(t, progress.DeviceId, rsp.DeviceId)
	assert.Equal(t, progress.Progress, rsp.Progress)
	assert.Equal(t, progress.Document, rsp.Document)
	assert.Equal(t, progress.Percentage, rsp.Percentage)
	assert.Equal(t, progress.Device, rsp.Device)
	assert.Equal(t, int64(1778445298), rsp.Timestamp)

}
