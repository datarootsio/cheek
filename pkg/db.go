package cheek

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := InitDB(db); err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}

	return db, nil
}

// assert that log table exists
func InitDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS log (
		id INTEGER PRIMARY KEY,
		job TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		message TEXT
	)`)
	if err != nil {
		return fmt.Errorf("create log table: %w", err)
	}

	return nil
}
