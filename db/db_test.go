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

func TestCreateTables_Successfully(t *testing.T) {
	err := Init(":memory:")
	assert.NoError(t, err)

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
