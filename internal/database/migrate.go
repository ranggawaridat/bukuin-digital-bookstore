package database

import (
    "database/sql"
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "strings"
)

func Migrate(db *sql.DB) error {
    if err := ensureSchemaMigrationsTable(db); err != nil {
        return err
    }

    files, err := filepath.Glob("migrations/*.sql")
    if err != nil {
        return err
    }

    sort.Strings(files)

    for _, file := range files {
        migrationName := filepath.Base(file)

        applied, err := migrationAlreadyApplied(db, migrationName)
        if err != nil {
            return err
        }

        if applied {
            continue
        }

        content, err := os.ReadFile(file)
        if err != nil {
            return fmt.Errorf("failed to read migration %s: %w", file, err)
        }

        _, err = db.Exec(string(content))
        if err != nil {
            if isDuplicateColumnError(err) {
                if recordErr := markMigrationApplied(db, migrationName); recordErr != nil {
                    return fmt.Errorf("failed to record migration %s: %w", file, recordErr)
                }
                continue
            }

            return fmt.Errorf("failed to execute migration %s: %w", file, err)
        }

        if err := markMigrationApplied(db, migrationName); err != nil {
            return fmt.Errorf("failed to record migration %s: %w", file, err)
        }
    }

    return nil
}

func ensureSchemaMigrationsTable(db *sql.DB) error {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE,
            applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        )
    `)

    return err
}

func migrationAlreadyApplied(db *sql.DB, name string) (bool, error) {
    var count int

    err := db.QueryRow(
        `SELECT COUNT(*) FROM schema_migrations WHERE name = ?`,
        name,
    ).Scan(&count)

    if err != nil {
        return false, err
    }

    return count > 0, nil
}

func markMigrationApplied(db *sql.DB, name string) error {
    _, err := db.Exec(
        `INSERT INTO schema_migrations (name) VALUES (?)`,
        name,
    )

    return err
}

func isDuplicateColumnError(err error) bool {
    return strings.Contains(strings.ToLower(err.Error()), "duplicate column name")
}
