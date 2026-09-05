package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	version := "001_init"

	var appliedVersion string

	err = db.QueryRow(
		`
		SELECT version
		FROM schema_migrations
		WHERE version = ?
		`,
		version,
	).Scan(&appliedVersion)

	if err == nil {
		return nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	sqlFile, err := os.ReadFile(
		"./migrations/001_init.sql",
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(sqlFile))
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`
		INSERT INTO schema_migrations (version)
		VALUES (?)
		`,
		version,
	)
	if err != nil {
		return err
	}

	fmt.Println(
		"migration applied:",
		version,
	)

	return nil
}
