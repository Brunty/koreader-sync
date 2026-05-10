package types

import (
	"errors"
	"time"

	"github.com/brunty/koreader-sync-server/crypto"
)

type User struct {
	Id        int64
	Username  string
	Password  string
	CreatedAt time.Time
}

type Progress struct {
	Id         int64
	UserId     int64
	Document   string
	Progress   string
	Percentage float64
	Device     string
	DeviceId   string
	Timestamp  time.Time
}

type StatusResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req *CreateUserRequest) Validate() error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (req *CreateUserRequest) MarshalToUser() (User, error) {
	hashedPw, err := crypto.HashPassword(req.Password)
	if err != nil {
		return User{}, err
	}

	return User{
		Username:  req.Username,
		Password:  hashedPw,
		CreatedAt: time.Now(),
	}, nil
}

type SyncProgressRequest struct {
	DeviceId   string  `json:"device_id"`
	Progress   string  `json:"progress"`
	Document   string  `json:"document"`
	Percentage float64 `json:"percentage"`
	Device     string  `json:"device"`
}

func (req *SyncProgressRequest) Validate() error {
	if req.DeviceId == "" {
		return errors.New("device_id is required")
	}
	if req.Progress == "" {
		return errors.New("progress is required")
	}
	if req.Document == "" {
		return errors.New("document is required")
	}
	if req.Device == "" {
		return errors.New("device is required")
	}

	return nil
}

func (req *SyncProgressRequest) MarshalToProgress(userId int64) (Progress, error) {
	return Progress{
		UserId:     userId,
		Document:   req.Document,
		Progress:   req.Progress,
		Percentage: req.Percentage,
		Device:     req.Device,
		DeviceId:   req.DeviceId,
		Timestamp:  time.Now(),
	}, nil
}

func (progress *Progress) MarshalToResponse() (ProgressResponse, error) {
	return ProgressResponse{
		DeviceId:   progress.DeviceId,
		Progress:   progress.Progress,
		Document:   progress.Document,
		Percentage: progress.Percentage,
		Device:     progress.Device,
		Timestamp:  progress.Timestamp.Unix(),
	}, nil
}

type ProgressResponse struct {
	DeviceId   string  `json:"device_id"`
	Progress   string  `json:"progress"`
	Document   string  `json:"document"`
	Percentage float64 `json:"percentage"`
	Device     string  `json:"device"`
	Timestamp  int64   `json:"timestamp"`
}
