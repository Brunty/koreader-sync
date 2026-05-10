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

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req *CreateUserRequest) Validate() error {
	var err error
	if req.Username == "" {
		err = errors.Join(err, errors.New("username is required"))
	}
	if req.Password == "" {
		err = errors.Join(err, errors.New("password is required"))
	}

	return err
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
