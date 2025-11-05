package database

import (
	"database/sql"
	"fmt"
)

// InitSchema creates all necessary tables and indexes if they don't exist
func InitSchema(db *sql.DB) error {
	// Create tables in order (parent tables first, then child tables)
	schemas := []string{
		createTicketsTable,
		createTicketTagsTable,
		createTicketFilesTable,
		createTicketCommentsTable,
		createTicketsIndexes,
	}

	for _, schema := range schemas {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("failed to execute schema: %w", err)
		}
	}

	return nil
}

const createTicketsTable = `
CREATE TABLE IF NOT EXISTS tickets (
    id TEXT PRIMARY KEY,
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
    ticket_id TEXT NOT NULL,
    tag TEXT NOT NULL,
    FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
    PRIMARY KEY (ticket_id, tag)
);`

const createTicketFilesTable = `
CREATE TABLE IF NOT EXISTS ticket_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ticket_id TEXT NOT NULL,
    file_path TEXT NOT NULL,
    FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE
);`

const createTicketCommentsTable = `
CREATE TABLE IF NOT EXISTS ticket_comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ticket_id TEXT NOT NULL,
    comment_text TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE
);`

const createTicketsIndexes = `
CREATE INDEX IF NOT EXISTS idx_tickets_project ON tickets(project);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_tickets_priority ON tickets(priority);
CREATE INDEX IF NOT EXISTS idx_tickets_type ON tickets(type);`
