package cheek

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Import SQLite driver for database/sql
	"github.com/stretchr/testify/assert"
)

// TestInitDB tests the InitDB function, including the cleanup logic.
func TestInitDB(t *testing.T) {
	// Create an in-memory SQLite database
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	// Create the log table without the UNIQUE constraint temporarily
	_, err = db.Exec(`CREATE TABLE log_temp (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        job TEXT,
        triggered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		triggered_by TEXT,
        duration INTEGER,
        status INTEGER,
        message TEXT
    )`)
	if err != nil {
		t.Fatalf("Failed to create temporary log table: %v", err)
	}

	// Insert conflicting data into the temporary log table
	_, err = db.Exec(`
		INSERT INTO log_temp (job, triggered_at, triggered_by, duration, status, message) VALUES
		('job1', '2023-10-01 10:00:00', 'user1', 120, 1, 'Success'),
		('job1', '2023-10-01 10:00:00', 'user1', 150, 1, 'Success'),  -- Duplicate
		('job1', '2023-10-01 11:00:00', 'user2', 90, 0, 'Failed')
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Create the actual log table with the UNIQUE constraint
	_, err = db.Exec(`CREATE TABLE log (
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
		t.Fatalf("Failed to create log table: %v", err)
	}

	// Move data from the temporary table to the log table
	_, err = db.Exec(`
		INSERT INTO log (job, triggered_at, triggered_by, duration, status, message)
		SELECT job, triggered_at, triggered_by, duration, status, message 
		FROM log_temp
		WHERE true
		ON CONFLICT (job, triggered_at, triggered_by) DO NOTHING
	`)

	if err != nil {
		t.Fatalf("Failed to transfer data to log table: %v", err)
	}

	//Drop the temporary table
	_, err = db.Exec("DROP TABLE log_temp;")
	if err != nil {
		t.Fatalf("Failed to drop temporary log table: %v", err)
	}

	// Call the InitDB function
	err = InitDB(db)
	assert.NoError(t, err, "InitDB should not return an error")

	// Check if cleanup worked correctly
	var cleanedCount int
	err = db.Get(&cleanedCount, "SELECT COUNT(*) FROM log")
	assert.NoError(t, err, "Querying the log table should not return an error")
	assert.Equal(t, 2, cleanedCount, "There should be 2 unique records in the log table after cleanup")

}
