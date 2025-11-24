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

// ConnectionFactory is a function that creates a database connection
type ConnectionFactory func(dbPath string) (*sql.DB, error)

// Registry of database connection factories
var connectionFactories = map[string]ConnectionFactory{
	config.DBTypeSQLite: newSQLiteConnection,
	config.DBTypeTurso:  newTursoConnection,
}

// RegisterConnectionFactory allows registering new database types
// This makes it easy to add support for PostgreSQL, MySQL, etc. later
func RegisterConnectionFactory(dbType string, factory ConnectionFactory) {
	connectionFactories[dbType] = factory
}

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

	// Get the appropriate connection factory
	factory, exists := connectionFactories[cfg.DatabaseType]
	if !exists {
		return fmt.Errorf("unsupported database type: %s", cfg.DatabaseType)
	}

	// Create the connection using the factory
	db, err = factory(dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to %s database: %w", cfg.DatabaseType, err)
	}

	// Initialize schema
	if err := InitSchema(db); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	return nil
}

// newSQLiteConnection creates a new SQLite database connection
func newSQLiteConnection(dbPath string) (*sql.DB, error) {
	if dbPath == "" {
		var err error
		dbPath, err = getDefaultDBPath()
		if err != nil {
			return nil, fmt.Errorf("failed to get default database path: %w", err)
		}
	}

	// Ensure the directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	// Enable foreign key constraints (disabled by default in SQLite)
	if _, err := conn.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return conn, nil
}

// newTursoConnection creates a new Turso database connection
func newTursoConnection(dbPath string) (*sql.DB, error) {
	tursoURL := os.Getenv("TURSO_URL")
	tursoToken := os.Getenv("TURSO_AUTH_TOKEN")

	if tursoURL == "" {
		return nil, fmt.Errorf("TURSO_URL environment variable is not set")
	}
	if tursoToken == "" {
		return nil, fmt.Errorf("TURSO_AUTH_TOKEN environment variable is not set")
	}

	// Construct the connection string for libsql
	// Format: libsql://host?authToken=token
	connStr := fmt.Sprintf("%s?authToken=%s", tursoURL, tursoToken)

	// Open database connection
	conn, err := sql.Open("libsql", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open Turso database: %w", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping Turso database: %w", err)
	}

	return conn, nil
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
