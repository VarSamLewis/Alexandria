package database

import (
	"alexandria/internal/logger"
	"database/sql"
	"fmt"
)

// InitSchema creates all necessary tables and indexes if they don't exist
func InitSchema(db *sql.DB) error {
	logger.Log.Debug("initializing database schema")

	// Create tables in order (parent tables first, then child tables)
	schemas := []struct {
		name   string
		script string
	}{
		{"tickets table", createTicketsTable},
		{"ticket_tags table", createTicketTagsTable},
		{"ticket_files table", createTicketFilesTable},
		{"ticket_comments table", createTicketCommentsTable},
		{"users table", createUsersTable}
		{"indexes", createTicketsIndexes},
	}

	for _, schema := range schemas {
		logger.Log.Debug("creating schema", "name", schema.name)
		if _, err := db.Exec(schema.script); err != nil {
			logger.Log.Error("failed to execute schema", "error", err, "name", schema.name)
			return fmt.Errorf("failed to execute schema: %w", err)
		}
	}

	logger.Log.Debug("database schema initialized successfully")
	return nil
}

const createTicketsTable = `
CREATE TABLE IF NOT EXISTS tickets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project TEXT NOT NULL,
    type TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    critical_path BOOLEAN DEFAULT 0,
    status TEXT NOT NULL,
    priority TEXT NOT NULL,
    created_by TEXT,
    assigned_to TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);`

const createTicketTagsTable = `
CREATE TABLE IF NOT EXISTS ticket_tags (
    ticket_id INTEGER NOT NULL,
    tag TEXT NOT NULL,
		FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
    PRIMARY KEY (ticket_id, tag)
);`

const createTicketFilesTable = `
CREATE TABLE IF NOT EXISTS ticket_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ticket_id INTEGER NOT NULL,
    file_path TEXT NOT NULL,
    FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE
);`

const createTicketCommentsTable = `
CREATE TABLE IF NOT EXISTS ticket_comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ticket_id INTEGER NOT NULL,
    comment_text TEXT NOT NULL,
    created_at DATETIME NOT NULL,
		FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE:w http.ResponseWriter, r *http.Request
);`

const createUsersTable = `
  CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      username TEXT NOT NULL UNIQUE,
      email TEXT NOT NULL UNIQUE COLLATE NOCASE,
      hashed_password TEXT NOT NULL,
      fullname TEXT NOT NULL,
      role TEXT NOT NULL CHECK(role IN ('admin', 'user', 'viewer')) DEFAULT 'user',
      created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`


const createTicketsIndexes = `
CREATE INDEX IF NOT EXISTS idx_tickets_project ON tickets(project);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_tickets_priority ON tickets(priority);
CREATE INDEX IF NOT EXISTS idx_tickets_type ON tickets(type);
CREATE INDEX IF NOT EXISTS idx_tickets_type ON tickets(type);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(LOWER(email));`
