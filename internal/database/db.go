package database

import (
	"alexandria/internal/config"
	"alexandria/internal/logger"
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
		logger.Log.Error("failed to load database config", "error", err)
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger.Log.Debug("initializing database", "type", cfg.DatabaseType, "path", dbPath)

	// Get the appropriate connection factory
	factory, exists := connectionFactories[cfg.DatabaseType]
	if !exists {
		logger.Log.Error("unsupported database type", "type", cfg.DatabaseType)
		return fmt.Errorf("unsupported database type: %s", cfg.DatabaseType)
	}

	// Create the connection using the factory
	logger.Log.Debug("creating database connection", "type", cfg.DatabaseType)
	db, err = factory(dbPath)
	if err != nil {
		logger.Log.Error("failed to connect to database", "error", err, "type", cfg.DatabaseType)
		return fmt.Errorf("failed to connect to %s database: %w", cfg.DatabaseType, err)
	}

	// Initialize schema
	logger.Log.Debug("initializing database schema")
	if err := InitSchema(db); err != nil {
		logger.Log.Error("failed to initialize schema", "error", err)
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	logger.Log.Info("database connection established", "type", cfg.DatabaseType)
	return nil
}

// newSQLiteConnection creates a new SQLite database connection
func newSQLiteConnection(dbPath string) (*sql.DB, error) {
	if dbPath == "" {
		var err error
		dbPath, err = getDefaultDBPath()
		if err != nil {
			logger.Log.Error("failed to get default database path", "error", err)
			return nil, fmt.Errorf("failed to get default database path: %w", err)
		}
	}

	logger.Log.Debug("connecting to SQLite", "path", dbPath)

	// Ensure the directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		logger.Log.Error("failed to create database directory", "error", err, "path", dbDir)
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Log.Error("failed to open SQLite database", "error", err, "path", dbPath)
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		conn.Close()
		logger.Log.Error("failed to ping SQLite database", "error", err)
		return nil, fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	// Enable foreign key constraints (disabled by default in SQLite)
	if _, err := conn.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		conn.Close()
		logger.Log.Error("failed to enable foreign keys", "error", err)
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	logger.Log.Debug("SQLite connection established", "path", dbPath)
	return conn, nil
}

// newTursoConnection creates a new Turso database connection
func newTursoConnection(dbPath string) (*sql.DB, error) {
	tursoURL := os.Getenv("TURSO_URL")
	tursoToken := os.Getenv("TURSO_AUTH_TOKEN")

	if tursoURL == "" {
		logger.Log.Error("Turso URL not set")
		return nil, fmt.Errorf("TURSO_URL environment variable is not set")
	}
	if tursoToken == "" {
		logger.Log.Error("Turso auth token not set")
		return nil, fmt.Errorf("TURSO_AUTH_TOKEN environment variable is not set")
	}

	logger.Log.Debug("connecting to Turso", "url", tursoURL)

	// Construct the connection string for libsql
	// Format: libsql://host?authToken=token
	connStr := fmt.Sprintf("%s?authToken=%s", tursoURL, tursoToken)

	// Open database connection
	conn, err := sql.Open("libsql", connStr)
	if err != nil {
		logger.Log.Error("failed to open Turso database", "error", err)
		return nil, fmt.Errorf("failed to open Turso database: %w", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		conn.Close()
		logger.Log.Error("failed to ping Turso database", "error", err)
		return nil, fmt.Errorf("failed to ping Turso database: %w", err)
	}

	logger.Log.Debug("Turso connection established", "url", tursoURL)
	return conn, nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// Close closes the database connection
func Close() error {
	if db != nil {
		logger.Log.Debug("closing database connection")
		if err := db.Close(); err != nil {
			logger.Log.Error("failed to close database", "error", err)
			return err
		}
		logger.Log.Debug("database connection closed")
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
