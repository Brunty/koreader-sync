package user

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type UserRepository interface {
	SelectByUsername(ctx context.Context, username string) (*User, error)
	Store(ctx context.Context, user User) (*int64, error)
	Update(ctx context.Context, user User) (*int64, error)
}

type UserRepositorySQLite struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositorySQLite{db: db}
}

func (r *UserRepositorySQLite) SelectByUsername(ctx context.Context, username string) (*User, error) {
	var user = User{}
	user.Username = username

	query := "SELECT id, password, created_at FROM users WHERE username = $1"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.Id, &user.Password, &user.CreatedAt)

	if err != nil {
		// if there are no rows, that's fine, just return nil, nil
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositorySQLite) Store(ctx context.Context, user User) (*int64, error) {
	res, err := r.db.ExecContext(ctx, "INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)", user.Username, user.Password, time.Now())
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()

	return &id, nil
}

func (r *UserRepositorySQLite) Update(ctx context.Context, user User) (*int64, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE users SET password = ? WHERE username = ?", user.Password, user.Username)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()

	return &id, nil
}
