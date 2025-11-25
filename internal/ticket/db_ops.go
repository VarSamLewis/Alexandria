package ticket

import (
	"alexandria/internal/logger"
	"database/sql"
	"fmt"
	"time"
)
// Create inserts a new ticket into the database
func (t *Ticket) Create(db *sql.DB, project string) error {
	logger.Log.Debug("creating ticket in database", "project", project, "title", t.Title)

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		logger.Log.Error("failed to begin transaction", "error", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert the main ticket record (ID is auto-generated)
	insertTicketQuery := `
		INSERT INTO tickets (
			project, type, title, description, critical_path,
			status, priority, created_by, assigned_to, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	logger.Log.Debug("inserting ticket record")
	result, err := tx.Exec(
		insertTicketQuery,
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
		logger.Log.Error("failed to insert ticket", "error", err, "title", t.Title)
		return fmt.Errorf("failed to insert ticket: %w", err)
	}

	// Get the auto-generated ID
	id, err := result.LastInsertId()
	if err != nil {
		logger.Log.Error("failed to get inserted ID", "error", err)
		return fmt.Errorf("failed to get inserted ID: %w", err)
	}
	t.ID = id
	logger.Log.Debug("ticket record inserted", "id", t.ID)

	// Insert tags
	if len(t.Tags) > 0 {
		logger.Log.Debug("inserting tags", "count", len(t.Tags))
		insertTagQuery := `INSERT INTO ticket_tags (ticket_id, tag) VALUES (?, ?)`
		for _, tag := range t.Tags {
			if _, err := tx.Exec(insertTagQuery, t.ID, tag); err != nil {
				logger.Log.Error("failed to insert tag", "error", err, "tag", tag)
				return fmt.Errorf("failed to insert tag: %w", err)
			}
		}
	}

	// Insert files
	if len(t.Files) > 0 {
		logger.Log.Debug("inserting files", "count", len(t.Files))
		insertFileQuery := `INSERT INTO ticket_files (ticket_id, file_path) VALUES (?, ?)`
		for _, file := range t.Files {
			if _, err := tx.Exec(insertFileQuery, t.ID, file); err != nil {
				logger.Log.Error("failed to insert file", "error", err, "file", file)
				return fmt.Errorf("failed to insert file: %w", err)
			}
		}
	}

	// Insert comments
	if len(t.Comments) > 0 {
		logger.Log.Debug("inserting comments", "count", len(t.Comments))
		insertCommentQuery := `INSERT INTO ticket_comments (ticket_id, comment_text, created_at) VALUES (?, ?, ?)`
		for _, comment := range t.Comments {
			if _, err := tx.Exec(insertCommentQuery, t.ID, comment, time.Now()); err != nil {
				logger.Log.Error("failed to insert comment", "error", err)
				return fmt.Errorf("failed to insert comment: %w", err)
			}
		}
	}

	// Commit the transaction
	logger.Log.Debug("committing transaction", "ticket_id", t.ID)
	if err := tx.Commit(); err != nil {
		logger.Log.Error("failed to commit transaction", "error", err, "ticket_id", t.ID)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Log.Info("ticket created in database", "id", t.ID, "project", project, "title", t.Title)
	return nil
}

// Update modifies an existing ticket in the database
func (t *Ticket) Update(db *sql.DB, project string, id int64, title string) error {
	logger.Log.Debug("updating ticket", "project", project, "id", id, "title", title)

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		logger.Log.Error("failed to begin transaction", "error", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var ticketID int64

	// Determine which identifier to use
	if id != 0 {
		ticketID = id
		logger.Log.Debug("using ID to find ticket", "id", ticketID)
	} else if title != "" {
		logger.Log.Debug("looking up ticket by title", "title", title)
		err = tx.QueryRow("SELECT id FROM tickets WHERE title = ? AND project = ?", title, project).Scan(&ticketID)
		if err == sql.ErrNoRows {
			logger.Log.Error("ticket not found", "title", title, "project", project)
			return fmt.Errorf("no ticket found with title '%s'", title)
		}
		if err != nil {
			logger.Log.Error("failed to find ticket", "error", err)
			return fmt.Errorf("failed to find ticket: %w", err)
		}
		logger.Log.Debug("found ticket", "ticket_id", ticketID)
	} else {
		logger.Log.Error("no identifier provided for update")
		return fmt.Errorf("either id or title must be provided")
	}

	// Update the main ticket record
	updateTicketQuery := `
		UPDATE tickets SET
			type = ?, title = ?, description = ?, critical_path = ?,
			status = ?, priority = ?, assigned_to = ?, updated_at = ?
		WHERE id = ? AND project = ?`

	logger.Log.Debug("executing update query", "ticket_id", ticketID)
	result, err := tx.Exec(
		updateTicketQuery,
		t.Type,
		t.Title,
		t.Description,
		t.CriticalPath,
		t.Status,
		t.Priority,
		t.AssignedTo,
		time.Now(),
		ticketID,
		project,
	)
	if err != nil {
		logger.Log.Error("failed to update ticket", "error", err)
		return fmt.Errorf("failed to update ticket: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Log.Error("failed to check rows affected", "error", err)
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		logger.Log.Error("no rows affected by update", "ticket_id", ticketID)
		return fmt.Errorf("no ticket found with the provided identifier")
	}

	logger.Log.Debug("ticket record updated", "rows_affected", rowsAffected)

	// Update tags - delete existing and insert new ones
	logger.Log.Debug("updating tags", "ticket_id", ticketID)
	if _, err := tx.Exec("DELETE FROM ticket_tags WHERE ticket_id = ?", ticketID); err != nil {
		logger.Log.Error("failed to delete existing tags", "error", err)
		return fmt.Errorf("failed to delete existing tags: %w", err)
	}

	if len(t.Tags) > 0 {
		logger.Log.Debug("inserting new tags", "count", len(t.Tags))
		insertTagQuery := `INSERT INTO ticket_tags (ticket_id, tag) VALUES (?, ?)`
		for _, tag := range t.Tags {
			if _, err := tx.Exec(insertTagQuery, ticketID, tag); err != nil {
				logger.Log.Error("failed to insert tag", "error", err, "tag", tag)
				return fmt.Errorf("failed to insert tag: %w", err)
			}
		}
	}

	// Update files - delete existing and insert new ones
	logger.Log.Debug("updating files", "ticket_id", ticketID)
	if _, err := tx.Exec("DELETE FROM ticket_files WHERE ticket_id = ?", ticketID); err != nil {
		logger.Log.Error("failed to delete existing files", "error", err)
		return fmt.Errorf("failed to delete existing files: %w", err)
	}

	if len(t.Files) > 0 {
		logger.Log.Debug("inserting new files", "count", len(t.Files))
		insertFileQuery := `INSERT INTO ticket_files (ticket_id, file_path) VALUES (?, ?)`
		for _, file := range t.Files {
			if _, err := tx.Exec(insertFileQuery, ticketID, file); err != nil {
				logger.Log.Error("failed to insert file", "error", err, "file", file)
				return fmt.Errorf("failed to insert file: %w", err)
			}
		}
	}

	// Add new comments (don't delete existing ones)
	if len(t.Comments) > 0 {
		logger.Log.Debug("adding new comments", "count", len(t.Comments))
		insertCommentQuery := `INSERT INTO ticket_comments (ticket_id, comment_text, created_at) VALUES (?, ?, ?)`
		for _, comment := range t.Comments {
			if _, err := tx.Exec(insertCommentQuery, ticketID, comment, time.Now()); err != nil {
				logger.Log.Error("failed to insert comment", "error", err)
				return fmt.Errorf("failed to insert comment: %w", err)
			}
		}
	}

	// Commit the transaction
	logger.Log.Debug("committing update transaction", "ticket_id", ticketID)
	if err := tx.Commit(); err != nil {
		logger.Log.Error("failed to commit transaction", "error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Log.Info("ticket updated", "ticket_id", ticketID, "project", project)
	return nil
}

// List retrieves tickets from the database based on the provided filters
func List(db *sql.DB, filters Filters) ([]Ticket, error) {
	logger.Log.Debug("listing tickets", "filters", fmt.Sprintf("%+v", filters))

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

	if filters.Project != nil {
		query += " AND t.project = ?"
		args = append(args, *filters.Project)
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

	logger.Log.Debug("executing list query")
	rows, err := db.Query(query, args...)
	if err != nil {
		logger.Log.Error("failed to query tickets", "error", err)
		return nil, fmt.Errorf("failed to query tickets: %w", err)
	}
	defer rows.Close()

	var tickets []Ticket
	ticketMap := make(map[int64]*Ticket)

	for rows.Next() {
		var t Ticket

		err := rows.Scan(
			&t.ID,
			&t.Project,
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
		logger.Log.Error("error iterating ticket rows", "error", err)
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	logger.Log.Debug("found tickets", "count", len(ticketMap))

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

	logger.Log.Info("tickets listed", "count", len(tickets))
	return tickets, nil
}

// loadTags loads tags for a specific ticket
func loadTags(db *sql.DB, ticketID int64) ([]string, error) {
	logger.Log.Debug("loading tags", "ticket_id", ticketID)
	rows, err := db.Query("SELECT tag FROM ticket_tags WHERE ticket_id = ?", ticketID)
	if err != nil {
		logger.Log.Error("failed to query tags", "error", err, "ticket_id", ticketID)
		return nil, fmt.Errorf("failed to load tags: %w", err)
	}
	defer rows.Close()

	tags := []string{}
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			logger.Log.Error("failed to scan tag", "error", err)
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("error iterating tags", "error", err)
		return nil, fmt.Errorf("error iterating tags: %w", err)
	}

	logger.Log.Debug("tags loaded", "count", len(tags))
	return tags, nil
}

// loadFiles loads files for a specific ticket
func loadFiles(db *sql.DB, ticketID int64) ([]string, error) {
	logger.Log.Debug("loading files", "ticket_id", ticketID)
	rows, err := db.Query("SELECT file_path FROM ticket_files WHERE ticket_id = ?", ticketID)
	if err != nil {
		logger.Log.Error("failed to query files", "error", err, "ticket_id", ticketID)
		return nil, fmt.Errorf("failed to load files: %w", err)
	}
	defer rows.Close()

	files := []string{}
	for rows.Next() {
		var file string
		if err := rows.Scan(&file); err != nil {
			logger.Log.Error("failed to scan file", "error", err)
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("error iterating files", "error", err)
		return nil, fmt.Errorf("error iterating files: %w", err)
	}

	logger.Log.Debug("files loaded", "count", len(files))
	return files, nil
}

// loadComments loads comments for a specific ticket
func loadComments(db *sql.DB, ticketID int64) ([]string, error) {
	logger.Log.Debug("loading comments", "ticket_id", ticketID)
	rows, err := db.Query("SELECT comment_text FROM ticket_comments WHERE ticket_id = ? ORDER BY created_at", ticketID)
	if err != nil {
		logger.Log.Error("failed to query comments", "error", err, "ticket_id", ticketID)
		return nil, fmt.Errorf("failed to load comments: %w", err)
	}
	defer rows.Close()

	comments := []string{}
	for rows.Next() {
		var comment string
		if err := rows.Scan(&comment); err != nil {
			logger.Log.Error("failed to scan comment", "error", err)
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("error iterating comments", "error", err)
		return nil, fmt.Errorf("error iterating comments: %w", err)
	}

	logger.Log.Debug("comments loaded", "count", len(comments))
	return comments, nil
}
  // Delete removes a ticket from the database
func (t *Ticket) Delete(db *sql.DB, project string, id int64, title string) error {
	logger.Log.Debug("deleting ticket", "project", project, "id", id, "title", title)

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		logger.Log.Error("failed to begin transaction", "error", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var ticketID int64

	// Determine which identifier to use
	if id != 0 {
		ticketID = id
	} else if title != "" {
		err = tx.QueryRow("SELECT id FROM tickets WHERE title = ? AND project = ?", title, project).Scan(&ticketID)
		if err == sql.ErrNoRows {
			logger.Log.Error("ticket not found", "title", title, "project", project)
			return fmt.Errorf("no ticket found with title '%s'", title)
		}
		if err != nil {
			logger.Log.Error("failed to find ticket", "error", err)
			return fmt.Errorf("failed to find ticket: %w", err)
		}
	} else {
		logger.Log.Error("no identifier provided for delete")
		return fmt.Errorf("either id or title must be provided")
	}

	logger.Log.Debug("deleting ticket data", "ticket_id", ticketID)

	// Delete from all tables using ticketID
	if _, err := tx.Exec("DELETE FROM ticket_tags WHERE ticket_id = ?", ticketID); err != nil {
		logger.Log.Error("failed to delete tags", "error", err)
		return fmt.Errorf("failed to delete tags: %w", err)
	}

	if _, err := tx.Exec("DELETE FROM ticket_files WHERE ticket_id = ?", ticketID); err != nil {
		logger.Log.Error("failed to delete files", "error", err)
		return fmt.Errorf("failed to delete files: %w", err)
	}

	if _, err := tx.Exec("DELETE FROM ticket_comments WHERE ticket_id = ?", ticketID); err != nil {
		logger.Log.Error("failed to delete comments", "error", err)
		return fmt.Errorf("failed to delete comments: %w", err)
	}

	if _, err := tx.Exec("DELETE FROM tickets WHERE id = ? AND project = ?", ticketID, project); err != nil {
		logger.Log.Error("failed to delete ticket record", "error", err)
		return fmt.Errorf("failed to delete ticket: %w", err)
	}

	if err := tx.Commit(); err != nil {
		logger.Log.Error("failed to commit delete transaction", "error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Log.Info("ticket deleted", "ticket_id", ticketID, "project", project)
	return nil
}

func (t *Ticket) View(db *sql.DB, project string, id int64, title string) error {
	logger.Log.Debug("viewing ticket", "project", project, "id", id, "title", title)

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		logger.Log.Error("failed to begin transaction", "error", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var ticketID int64

	if id != 0 {
		ticketID = id
		logger.Log.Debug("using provided ID", "id", ticketID)
	} else if title != "" {
		logger.Log.Debug("resolving ticket by title", "title", title, "project", project)
		err = tx.QueryRow("SELECT id FROM tickets WHERE title = ? AND project = ?", title, project).Scan(&ticketID)
		if err == sql.ErrNoRows {
			logger.Log.Error("ticket not found by title", "title", title, "project", project)
			return fmt.Errorf("no ticket found with title '%s'", title)
		}
		if err != nil {
			logger.Log.Error("failed to find ticket by title", "error", err, "title", title)
			return fmt.Errorf("failed to find ticket: %w", err)
		}
		logger.Log.Debug("resolved ticket ID", "id", ticketID, "title", title)
	} else {
		logger.Log.Error("validation failed", "error", "neither id nor title provided")
		return fmt.Errorf("either id or title must be provided")
	}

	logger.Log.Debug("fetching ticket from database", "id", ticketID, "project", project)
	query := `SELECT id, project, type, title, description, critical_path,
                status, priority, created_by, assigned_to, created_at, updated_at
                FROM tickets WHERE id = ? AND project = ?`

	err = tx.QueryRow(query, ticketID, project).Scan(
		&t.ID,
		&t.Project,
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
	if err == sql.ErrNoRows {
		logger.Log.Error("ticket not found", "id", ticketID, "project", project)
		return fmt.Errorf("ticket not found")
	}
	if err != nil {
		logger.Log.Error("failed to fetch ticket", "error", err, "id", ticketID)
		return fmt.Errorf("failed to fetch ticket: %w", err)
	}

	logger.Log.Debug("ticket fetched successfully", "id", t.ID, "title", t.Title)

	// Load related data
	logger.Log.Debug("loading related data", "ticket_id", ticketID)
	tags, err := loadTags(db, ticketID)
	if err != nil {
		return err
	}
	t.Tags = tags
	logger.Log.Debug("loaded tags", "count", len(tags))

	files, err := loadFiles(db, ticketID)
	if err != nil {
		return err
	}
	t.Files = files
	logger.Log.Debug("loaded files", "count", len(files))

	comments, err := loadComments(db, ticketID)
	if err != nil {
		return err
	}
	t.Comments = comments
	logger.Log.Debug("loaded comments", "count", len(comments))

	logger.Log.Debug("committing view transaction")
	if err := tx.Commit(); err != nil {
		logger.Log.Error("failed to commit view transaction", "error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Log.Info("ticket viewed successfully", "id", t.ID, "title", t.Title, "project", project)
	return nil
}
