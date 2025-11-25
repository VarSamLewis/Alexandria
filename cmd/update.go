package cmd

import (
	"alexandria/internal/database"
	"alexandria/internal/logger"
	"alexandria/internal/ticket"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	updateID         string
	updateFindTitle  string
	updateProject    string
	updateTitle      string
	updateDesc       string
	updateType       string
	updateStatus     string
	updatePriority   string
	updateCritical   *bool
	updateAssignedTo string
	updateCreatedBy  string
	updateTags       string
	updateFiles      string
	updateComments   string
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing ticket",
	Long:  `Update an existing ticket's fields. Specify the ticket by ID or title, then provide the fields to update.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Log.Debug("updating ticket", "id", updateID, "title", updateFindTitle, "project", updateProject)

		// Validate that at least one identifier is provided
		if updateID == "" && updateFindTitle == "" {
			logger.Log.Error("validation failed", "error", "no identifier provided")
			return fmt.Errorf("either --id or --title must be provided to identify the ticket")
		}

		// Validate that project is provided
		if updateProject == "" {
			logger.Log.Error("validation failed", "error", "project is required")
			return fmt.Errorf("project is required")
		}

		// Parse ID if provided
		var ticketID int64
		if updateID != "" {
			var err error
			ticketID, err = strconv.ParseInt(updateID, 10, 64)
			if err != nil {
				logger.Log.Error("failed to parse ticket ID", "error", err, "id", updateID)
				return fmt.Errorf("invalid ID format: %s (must be a number)", updateID)
			}
			logger.Log.Debug("parsed ticket ID", "id", ticketID)
		}

		// Get database connection
		db := database.GetDB()
		if db == nil {
			logger.Log.Error("database not initialized")
			return fmt.Errorf("database not initialized")
		}

		// First, fetch the existing ticket to preserve current values
		logger.Log.Debug("fetching existing ticket")
		filters := ticket.Filters{}
		tickets, err := ticket.List(db, filters)
		if err != nil {
			logger.Log.Error("failed to fetch tickets", "error", err)
			return fmt.Errorf("failed to fetch tickets: %w", err)
		}

		// Find the ticket to update
		var existingTicket *ticket.Ticket
		for i := range tickets {
			if ticketID != 0 && tickets[i].ID == ticketID {
				existingTicket = &tickets[i]
				break
			} else if updateFindTitle != "" && tickets[i].Title == updateFindTitle {
				existingTicket = &tickets[i]
				break
			}
		}

		if existingTicket == nil {
			if ticketID != 0 {
				logger.Log.Error("ticket not found", "id", ticketID)
				return fmt.Errorf("ticket with ID '%d' not found", ticketID)
			}
			logger.Log.Error("ticket not found", "title", updateFindTitle)
			return fmt.Errorf("ticket with title '%s' not found", updateFindTitle)
		}

		logger.Log.Debug("found existing ticket", "id", existingTicket.ID, "title", existingTicket.Title)

		// Track if at least one field is being updated
		hasUpdates := false

		// Update only the fields that were specified
		if updateType != "" {
			tType := ticket.Type(updateType)
			if !tType.Valid() {
				logger.Log.Error("validation failed", "error", "invalid type", "type", updateType)
				return fmt.Errorf("invalid type: %s (must be: bug, feature, or task)", updateType)
			}
			existingTicket.Type = tType
			hasUpdates = true
			logger.Log.Debug("updating type", "new_type", tType)
		}

		if updateStatus != "" {
			tStatus := ticket.Status(updateStatus)
			if !tStatus.Valid() {
				logger.Log.Error("validation failed", "error", "invalid status", "status", updateStatus)
				return fmt.Errorf("invalid status: %s (must be: open, in-progress, or closed)", updateStatus)
			}
			existingTicket.Status = tStatus
			hasUpdates = true
			logger.Log.Debug("updating status", "new_status", tStatus)
		}

		if updatePriority != "" {
			tPriority := ticket.Priority(updatePriority)
			if !tPriority.Valid() {
				logger.Log.Error("validation failed", "error", "invalid priority", "priority", updatePriority)
				return fmt.Errorf("invalid priority: %s (must be: low, medium, high, or undefined)", updatePriority)
			}
			existingTicket.Priority = tPriority
			hasUpdates = true
			logger.Log.Debug("updating priority", "new_priority", tPriority)
		}

		if updateTitle != "" {
			existingTicket.Title = updateTitle
			hasUpdates = true
			logger.Log.Debug("updating title", "new_title", updateTitle)
		}

		if updateDesc != "" {
			existingTicket.Description = updateDesc
			hasUpdates = true
			logger.Log.Debug("updating description")
		}

		if cmd.Flags().Changed("criticalpath") {
			existingTicket.CriticalPath = *updateCritical
			hasUpdates = true
			logger.Log.Debug("updating critical path", "critical", *updateCritical)
		}

		if updateAssignedTo != "" {
			existingTicket.AssignedTo = &updateAssignedTo
			hasUpdates = true
			logger.Log.Debug("updating assigned to", "assigned_to", updateAssignedTo)
		}

		if updateCreatedBy != "" {
			existingTicket.CreatedBy = &updateCreatedBy
			hasUpdates = true
			logger.Log.Debug("updating created by", "created_by", updateCreatedBy)
		}

		if updateTags != "" {
			tagList := strings.Split(updateTags, ",")
			for i, tag := range tagList {
				tagList[i] = strings.TrimSpace(tag)
			}
			existingTicket.Tags = tagList
			hasUpdates = true
			logger.Log.Debug("updating tags", "count", len(tagList))
		}

		if updateFiles != "" {
			fileList := strings.Split(updateFiles, ",")
			for i, file := range fileList {
				fileList[i] = strings.TrimSpace(file)
			}
			existingTicket.Files = fileList
			hasUpdates = true
			logger.Log.Debug("updating files", "count", len(fileList))
		}

		if updateComments != "" {
			commentList := strings.Split(updateComments, ",")
			for i, comment := range commentList {
				commentList[i] = strings.TrimSpace(comment)
			}
			existingTicket.Comments = commentList
			hasUpdates = true
			logger.Log.Debug("updating comments", "count", len(commentList))
		}

		// Check if at least one field is being updated
		if !hasUpdates {
			logger.Log.Error("validation failed", "error", "no fields to update")
			return fmt.Errorf("no fields specified to update")
		}

		// Call the Update method
		logger.Log.Debug("calling update method", "project", updateProject, "id", ticketID, "title", updateFindTitle)
		if err := existingTicket.Update(db, updateProject, ticketID, updateFindTitle); err != nil {
			logger.Log.Error("failed to update ticket", "error", err, "project", updateProject)
			return fmt.Errorf("failed to update ticket: %w", err)
		}

		// Success message
		if ticketID != 0 {
			logger.Log.Info("ticket updated successfully", "id", ticketID, "project", updateProject)
			fmt.Printf("Successfully updated ticket with ID: %d in project: %s\n", ticketID, updateProject)
		} else {
			logger.Log.Info("ticket updated successfully", "title", updateFindTitle, "project", updateProject)
			fmt.Printf("Successfully updated ticket with title: %s in project: %s\n", updateFindTitle, updateProject)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Flags to identify the ticket
	updateCmd.Flags().StringVarP(&updateID, "id", "i", "", "Ticket ID to update")
	updateCmd.Flags().StringVarP(&updateFindTitle, "title", "t", "", "Find ticket by title to update")
	updateCmd.Flags().StringVar(&updateProject, "project", "", "Project name (required)")

	// Flags for fields to update
	updateCmd.Flags().StringVar(&updateTitle, "new-title", "", "New title for the ticket")
	updateCmd.Flags().StringVarP(&updateDesc, "description", "d", "", "New description for the ticket")
	updateCmd.Flags().StringVar(&updateType, "type", "", "New type (bug, feature, task)")
	updateCmd.Flags().StringVar(&updateStatus, "status", "", "New status (open, in-progress, closed)")
	updateCmd.Flags().StringVarP(&updatePriority, "priority", "p", "", "New priority (low, medium, high, undefined)")
	updateCritical = updateCmd.Flags().BoolP("criticalpath", "c", false, "Mark ticket as critical path")
	updateCmd.Flags().StringVarP(&updateAssignedTo, "assigned-to", "a", "", "Assign ticket to user")
	updateCmd.Flags().StringVar(&updateCreatedBy, "created-by", "", "Update ticket creator")
	updateCmd.Flags().StringVar(&updateTags, "tags", "", "Comma-separated list of tags (replaces existing)")
	updateCmd.Flags().StringVar(&updateFiles, "files", "", "Comma-separated list of file paths (replaces existing)")
	updateCmd.Flags().StringVar(&updateComments, "comments", "", "Comma-separated list of comments to add")

	updateCmd.MarkFlagRequired("project")
}
