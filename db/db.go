package db

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

var DBCon *sql.DB

func Init(dbFile string) error {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	DBCon = db

	return nil
}

func CreateTables() {
	DBCon.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	  	)
	`)

	DBCon.Exec(`
		CREATE TABLE IF NOT EXISTS progress (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			document TEXT NOT NULL,
			progress TEXT NOT NULL,
			percentage REAL NOT NULL,
			device TEXT NOT NULL,
			device_id TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id),
			UNIQUE(user_id, document)
		)
	`)
}
