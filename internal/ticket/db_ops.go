package ticket

import (
	"database/sql"
	"fmt"
	"time"
)
// Create inserts a new ticket into the database
func (t *Ticket) Create(db *sql.DB, project string) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert the main ticket record
	insertTicketQuery := `
		INSERT INTO tickets (
			id, project, type, title, description, critical_path,
			status, priority, created_by, assigned_to, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.Exec(
		insertTicketQuery,
		t.ID,
		project,
		t.Type,
		t.Title,
		t.Description,
		t.CriticalPath,
		t.Status,
		t.Priority,
		t.CreatedBy,
		t.AssignedTo,
		t.CreatedAt,
		t.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert ticket: %w", err)
	}

	// Insert tags
	if len(t.Tags) > 0 {
		insertTagQuery := `INSERT INTO ticket_tags (ticket_id, tag) VALUES (?, ?)`
		for _, tag := range t.Tags {
			if _, err := tx.Exec(insertTagQuery, t.ID, tag); err != nil {
				return fmt.Errorf("failed to insert tag: %w", err)
			}
		}
	}

	// Insert files
	if len(t.Files) > 0 {
		insertFileQuery := `INSERT INTO ticket_files (ticket_id, file_path) VALUES (?, ?)`
		for _, file := range t.Files {
			if _, err := tx.Exec(insertFileQuery, t.ID, file); err != nil {
				return fmt.Errorf("failed to insert file: %w", err)
			}
		}
	}

	// Insert comments
	if len(t.Comments) > 0 {
		insertCommentQuery := `INSERT INTO ticket_comments (ticket_id, comment_text, created_at) VALUES (?, ?, ?)`
		for _, comment := range t.Comments {
			if _, err := tx.Exec(insertCommentQuery, t.ID, comment, time.Now()); err != nil {
				return fmt.Errorf("failed to insert comment: %w", err)
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// List retrieves tickets from the database based on the provided filters
func List(db *sql.DB, filters Filters) ([]Ticket, error) {
	query := `
		SELECT DISTINCT t.id, t.project, t.type, t.title, t.description,
		       t.critical_path, t.status, t.priority, t.created_by,
		       t.assigned_to, t.created_at, t.updated_at
		FROM tickets t
		LEFT JOIN ticket_tags tt ON t.id = tt.ticket_id
		WHERE 1=1`

	var args []interface{}

	// Build dynamic query based on filters
	if filters.Status != nil {
		query += " AND t.status = ?"
		args = append(args, *filters.Status)
	}

	if filters.Type != nil {
		query += " AND t.type = ?"
		args = append(args, *filters.Type)
	}

	if filters.Priority != nil {
		query += " AND t.priority = ?"
		args = append(args, *filters.Priority)
	}

	if filters.AssignedTo != nil {
		query += " AND t.assigned_to = ?"
		args = append(args, *filters.AssignedTo)
	}

	// Filter by tags if provided
	if len(filters.Tags) > 0 {
		query += " AND tt.tag IN ("
		for i := range filters.Tags {
			if i > 0 {
				query += ","
			}
			query += "?"
			args = append(args, filters.Tags[i])
		}
		query += ")"
	}

	query += " ORDER BY t.created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tickets: %w", err)
	}
	defer rows.Close()

	var tickets []Ticket
	ticketMap := make(map[string]*Ticket)

	for rows.Next() {
		var t Ticket
		var project string

		err := rows.Scan(
			&t.ID,
			&project,
			&t.Type,
			&t.Title,
			&t.Description,
			&t.CriticalPath,
			&t.Status,
			&t.Priority,
			&t.CreatedBy,
			&t.AssignedTo,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ticket: %w", err)
		}

		// Avoid duplicates from JOIN
		if _, exists := ticketMap[t.ID]; !exists {
			ticketMap[t.ID] = &t
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Convert map to slice and load related data
	for _, ticket := range ticketMap {
		// Load tags
		tags, err := loadTags(db, ticket.ID)
		if err != nil {
			return nil, err
		}
		ticket.Tags = tags

		// Load files
		files, err := loadFiles(db, ticket.ID)
		if err != nil {
			return nil, err
		}
		ticket.Files = files

		// Load comments
		comments, err := loadComments(db, ticket.ID)
		if err != nil {
			return nil, err
		}
		ticket.Comments = comments

		tickets = append(tickets, *ticket)
	}

	return tickets, nil
}

// loadTags loads tags for a specific ticket
func loadTags(db *sql.DB, ticketID string) ([]string, error) {
	rows, err := db.Query("SELECT tag FROM ticket_tags WHERE ticket_id = ?", ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to load tags: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

// loadFiles loads files for a specific ticket
func loadFiles(db *sql.DB, ticketID string) ([]string, error) {
	rows, err := db.Query("SELECT file_path FROM ticket_files WHERE ticket_id = ?", ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to load files: %w", err)
	}
	defer rows.Close()

	var files []string
	for rows.Next() {
		var file string
		if err := rows.Scan(&file); err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, rows.Err()
}

// loadComments loads comments for a specific ticket
func loadComments(db *sql.DB, ticketID string) ([]string, error) {
	rows, err := db.Query("SELECT comment_text FROM ticket_comments WHERE ticket_id = ? ORDER BY created_at", ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to load comments: %w", err)
	}
	defer rows.Close()

	var comments []string
	for rows.Next() {
		var comment string
		if err := rows.Scan(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}
  // Delete removes a ticket from the database
func (t *Ticket) Delete(db *sql.DB, project string, id string, title string) error {
  // Start a transaction
      tx, err := db.Begin()
      if err != nil {
          return fmt.Errorf("failed to begin transaction: %w", err)
      }
      defer tx.Rollback()

      var ticketID string

      // Determine which identifier to use
      if id != "" {
          ticketID = id
      } else if title != "" {
          err = tx.QueryRow("SELECT id FROM tickets WHERE title = ? AND project = ?", title, project).Scan(&ticketID)
          if err == sql.ErrNoRows {
              return fmt.Errorf("no ticket found with title '%s'", title)
          }
          if err != nil {
              return fmt.Errorf("failed to find ticket: %w", err)
          }
      } else {
          return fmt.Errorf("either id or title must be provided")
      }

      // Delete from all tables using ticketID
      if _, err := tx.Exec("DELETE FROM ticket_tags WHERE ticket_id = ?", ticketID); err != nil {
          return fmt.Errorf("failed to delete tags: %w", err)
      }

      if _, err := tx.Exec("DELETE FROM ticket_files WHERE ticket_id = ?", ticketID); err != nil {
          return fmt.Errorf("failed to delete files: %w", err)
      }

      if _, err := tx.Exec("DELETE FROM ticket_comments WHERE ticket_id = ?", ticketID); err != nil {
          return fmt.Errorf("failed to delete comments: %w", err)
      }

      if _, err := tx.Exec("DELETE FROM tickets WHERE id = ? AND project = ?", ticketID, project); err != nil {
          return fmt.Errorf("failed to delete ticket: %w", err)
      }

      return tx.Commit()
  }
