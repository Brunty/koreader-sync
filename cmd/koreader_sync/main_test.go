package main

import (
	"context"
	"errors"
	"testing"

	"github.com/brunty/koreader-sync-server/crypto"
	userpackage "github.com/brunty/koreader-sync-server/user"
	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	selectByUsernameFn func(ctx context.Context, username string) (*userpackage.User, error)
	updateFn           func(ctx context.Context, user userpackage.User) (*int64, error)
}

func (m *mockUserRepo) SelectByUsername(ctx context.Context, username string) (*userpackage.User, error) {
	return m.selectByUsernameFn(ctx, username)
}

func (m *mockUserRepo) Store(ctx context.Context, user userpackage.User) (*int64, error) {
	return nil, nil
}

func (m *mockUserRepo) Update(ctx context.Context, user userpackage.User) (*int64, error) {
	if m.updateFn == nil {
		return nil, nil
	}
	return m.updateFn(ctx, user)
}

func TestMd5(t *testing.T) {
	assert.Equal(t, "5d41402abc4b2a76b9719d911017c592", md5("hello"))
	assert.Equal(t, "e99a18c428cb38d5f260853678922e03", md5("abc123"))
	assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", md5(""))
}

func TestChangePassword_Success(t *testing.T) {
	repo := &mockUserRepo{
		selectByUsernameFn: func(ctx context.Context, username string) (*userpackage.User, error) {
			return &userpackage.User{
				Id:       1,
				Username: "test-user",
				Password: "old-hash",
			}, nil
		},
		updateFn: func(ctx context.Context, user userpackage.User) (*int64, error) {
			assert.True(t, crypto.BcryptCheckPasswordHash(md5("new-password"), user.Password))
			return int64Ptr(1), nil
		},
	}

	err := changePassword(context.Background(), repo, "test-user", "new-password")
	assert.NoError(t, err)
}

func TestChangePassword_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{
		selectByUsernameFn: func(ctx context.Context, username string) (*userpackage.User, error) {
			return nil, nil
		},
	}

	err := changePassword(context.Background(), repo, "nonexistent", "password")
	assert.EqualError(t, err, "user not found")
}

func TestChangePassword_UpdateError(t *testing.T) {
	repo := &mockUserRepo{
		selectByUsernameFn: func(ctx context.Context, username string) (*userpackage.User, error) {
			return &userpackage.User{
				Id:       1,
				Username: "test-user",
				Password: "old-hash",
			}, nil
		},
		updateFn: func(ctx context.Context, user userpackage.User) (*int64, error) {
			return nil, errors.New("db error")
		},
	}

	err := changePassword(context.Background(), repo, "test-user", "new-password")
	assert.ErrorContains(t, err, "error storing user")
}

func int64Ptr(n int64) *int64 {
	return &n
}
