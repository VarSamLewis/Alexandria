package database

import (
	"alexandria/internal/config"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var db *sql.DB

// Init initializes the database connection and creates the schema
// It loads .env file, checks config for database type, and connects accordingly
func Init(dbPath string) error {
	// Load .env file from project root (ignore error if file doesn't exist)
	_ = godotenv.Load()

	// Get database type from config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Connect to appropriate database based on config
	switch cfg.DatabaseType {
	case config.DBTypeSQLite:
		if err := initSQLite(dbPath); err != nil {
			return err
		}
	case config.DBTypeTurso:
		if err := initTurso(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown database type: %s", cfg.DatabaseType)
	}

	// Initialize schema
	if err := InitSchema(db); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	return nil
}

// initSQLite initializes a local SQLite database connection
func initSQLite(dbPath string) error {
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
		return fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	// Enable foreign key constraints (disabled by default in SQLite)
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return nil
}

// initTurso initializes a Turso database connection
func initTurso() error {
	tursoURL := os.Getenv("TURSO_URL")
	tursoToken := os.Getenv("TURSO_AUTH_TOKEN")

	if tursoURL == "" {
		return fmt.Errorf("TURSO_URL environment variable is not set")
	}
	if tursoToken == "" {
		return fmt.Errorf("TURSO_AUTH_TOKEN environment variable is not set")
	}

	// Construct the connection string for libsql
	// Format: libsql://host?authToken=token
	connStr := fmt.Sprintf("%s?authToken=%s", tursoURL, tursoToken)

	// Open database connection
	var err error
	db, err = sql.Open("libsql", connStr)
	if err != nil {
		return fmt.Errorf("failed to open Turso database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping Turso database: %w", err)
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
