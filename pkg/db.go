package cheek

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/glebarez/go-sqlite"
)

func OpenDB(dbPath string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := InitDB(db); err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}

	return db, nil
}

func InitDB(db *sqlx.DB) error {
	// Create the log table if it doesn't exist
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS log (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        job TEXT,
        triggered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		triggered_by TEXT,
        duration INTEGER,
        status INTEGER,
        message TEXT,
		UNIQUE(job, triggered_at, triggered_by)
    )`)
	if err != nil {
		return fmt.Errorf("create log table: %w", err)
	}

	// Perform cleanup to remove old, non-conforming records
	_, err = db.Exec(`
		DELETE FROM log
		WHERE id NOT IN (
			SELECT MIN(id)
			FROM log
			GROUP BY job, triggered_at, triggered_by
		);
	`)
	if err != nil {
		return fmt.Errorf("cleanup old log records: %w", err)
	}

	return nil
}
