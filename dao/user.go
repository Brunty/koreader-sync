package dao

import (
	"database/sql"
	"errors"
	"time"

	"github.com/brunty/koreader-sync-server/db"
	"github.com/brunty/koreader-sync-server/types"
)

func SelectUserByUsername(username string) (*types.User, error) {
	var user = types.User{}
	user.Username = username

	query := "SELECT id, password, created_at FROM users WHERE username = $1"
	err := db.DBCon.QueryRow(query, username).Scan(&user.Id, &user.Password, &user.CreatedAt)

	if err != nil {
		// if there's no rows, that's fine, just return nil nil
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func StoreUser(user types.User) error {
	_, err := db.DBCon.Exec("INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)", user.Username, user.Password, time.Now())
	if err != nil {
		return err
	}

	return nil
}
