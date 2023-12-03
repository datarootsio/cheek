package cheek

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(dbPath string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := InitDB(db); err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}

	return db, nil
}

func InitDB(db *sqlx.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS log (
        id INTEGER PRIMARY KEY,
        job TEXT,
        triggered_at DATETIME,
		triggered_by TEXT,
        duration INTEGER,
        status INTEGER,
        message TEXT
    )`)
	if err != nil {
		return fmt.Errorf("create log table: %w", err)
	}

	return nil
}
