package sync_progress

import (
	"errors"
	"time"
)

type SyncProgress struct {
	ID         int64
	UserID     int64
	Document   string
	Progress   string
	Percentage float64
	Device     string
	DeviceID   string
	Timestamp  time.Time
}

type StoreSyncProgressRequest struct {
	DeviceID   string  `json:"device_id"`
	Progress   string  `json:"progress"`
	Document   string  `json:"document"`
	Percentage float64 `json:"percentage"`
	Device     string  `json:"device"`
}

func (req *StoreSyncProgressRequest) Validate() error {
	var err error
	if req.DeviceID == "" {
		err = errors.Join(err, errors.New("device_id is required"))
	}
	if req.Progress == "" {
		err = errors.Join(err, errors.New("progress is required"))
	}
	if req.Document == "" {
		err = errors.Join(err, errors.New("document is required"))
	}
	if req.Device == "" {
		err = errors.Join(err, errors.New("device is required"))
	}

	return err
}

func (req *StoreSyncProgressRequest) MarshalToSyncProgress(userID int64) SyncProgress {
	return SyncProgress{
		UserID:     userID,
		Document:   req.Document,
		Progress:   req.Progress,
		Percentage: req.Percentage,
		Device:     req.Device,
		DeviceID:   req.DeviceID,
		Timestamp:  time.Now(),
	}
}

func (progress SyncProgress) MarshalToReadResponse() ReadSyncProgressResponse {
	return ReadSyncProgressResponse{
		DeviceID:   progress.DeviceID,
		Progress:   progress.Progress,
		Document:   progress.Document,
		Percentage: progress.Percentage,
		Device:     progress.Device,
		Timestamp:  progress.Timestamp.Unix(),
	}
}

type ReadSyncProgressResponse struct {
	DeviceID   string  `json:"device_id"`
	Progress   string  `json:"progress"`
	Document   string  `json:"document"`
	Percentage float64 `json:"percentage"`
	Device     string  `json:"device"`
	Timestamp  int64   `json:"timestamp"`
}
