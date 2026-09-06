package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func Migrate(db *sql.DB) error {
	files, err := filepath.Glob(
		"migrations/*.sql",
	)
	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf(
				"failed to read migration %s: %w",
				file,
				err,
			)
		}

		_, err = db.Exec(
			string(content),
		)
		if err != nil {
			return fmt.Errorf(
				"failed to execute migration %s: %w",
				file,
				err,
			)
		}
	}

	return nil
}