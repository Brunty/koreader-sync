package db

import (
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit_Successfully(t *testing.T) {
	err := Init(":memory:")

	assert.NoError(t, err)
}

func TestInit_InvalidDriver(t *testing.T) {
	err := initWithDriver("nonexistent", ":memory:")

	assert.Error(t, err)
}

// Is this necessary as table creation is implicitely tested in other tests because y'know, WE CAN PUT STUFF IN THE
// TABLES... no. It's not. But it was a fun test to write to poke around with the sqlite_schema stuff, so there's that
func TestCreateTables_Successfully(t *testing.T) {
	err := Init(":memory:")
	assert.NoError(t, err)
	defer EmptyTables()

	CreateTables()

	rows, err := DBCon.Query(`
		SELECT COUNT(name)
		  FROM sqlite_schema
		 WHERE name = @table_name
		   AND tbl_name = @table_name
	`, sql.Named("table_name", "users"))
	assert.NoError(t, err)
	defer rows.Close()
	var userTableCount int

	for rows.Next() {
		if err := rows.Scan(&userTableCount); err != nil {
			log.Fatal(err)
		}
	}

	assert.Equal(t, 1, userTableCount)

	rows, err = DBCon.Query(`
		SELECT COUNT(name)
		  FROM sqlite_schema
		 WHERE name = @table_name
		   AND tbl_name = @table_name
	`, sql.Named("table_name", "progress"))
	assert.NoError(t, err)
	defer rows.Close()
	var progressTableCount int

	for rows.Next() {
		if err := rows.Scan(&progressTableCount); err != nil {
			log.Fatal(err)
		}
	}

	assert.Equal(t, 1, progressTableCount)
}
