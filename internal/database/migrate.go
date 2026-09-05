package database

import (
	"database/sql"
	"os"
)

func Migrate(db *sql.DB) error {
	sqlFile, err := os.ReadFile("./migrations/001_init.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(sqlFile))
	if err != nil {
		return err
	}

	return nil
}
