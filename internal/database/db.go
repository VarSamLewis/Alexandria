package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Init initializes the database connection and creates the schema
// If dbPath is empty, uses the default location: ~/work/alexandria/tickets.db
func Init(dbPath string) error {
	if dbPath == "" {
		var err error
		dbPath, err = getDefaultDBPath()
		if err != nil {
			return fmt.Errorf("failed to get default database path: %w", err)
		}
	}

	// Ensure the directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable foreign key constraints (disabled by default in SQLite)
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Initialize schema
	if err := InitSchema(db); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// Close closes the database connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// getDefaultDBPath returns the default database path: ~/work/DB/Alexandria/tickets.db
func getDefaultDBPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, "work", "DB", "Alexandria", "tickets.db"), nil
}

// GetDBPath returns the current database path (for informational purposes)
/*func GetDBPath() (string, error) {
	return getDefaultDBPath()
}
*/
