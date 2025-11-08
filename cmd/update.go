package cmd

import (
	"fmt"
	"mycli/internal/database"
	"mycli/internal/ticket"
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
		// Validate that at least one identifier is provided
		if updateID == "" && updateFindTitle == "" {
			return fmt.Errorf("either --id or --title must be provided to identify the ticket")
		}

		// Validate that project is provided
		if updateProject == "" {
			return fmt.Errorf("project is required")
		}

		// Get database connection
		db := database.GetDB()
		if db == nil {
			return fmt.Errorf("database not initialized")
		}

		// First, fetch the existing ticket to preserve current values
		filters := ticket.Filters{}
		tickets, err := ticket.List(db, filters)
		if err != nil {
			return fmt.Errorf("failed to fetch tickets: %w", err)
		}

		// Find the ticket to update
		var existingTicket *ticket.Ticket
		for i := range tickets {
			if updateID != "" && tickets[i].ID == updateID {
				existingTicket = &tickets[i]
				break
			} else if updateFindTitle != "" && tickets[i].Title == updateFindTitle {
				existingTicket = &tickets[i]
				break
			}
		}

		if existingTicket == nil {
			if updateID != "" {
				return fmt.Errorf("ticket with ID '%s' not found", updateID)
			}
			return fmt.Errorf("ticket with title '%s' not found", updateFindTitle)
		}

		// Track if at least one field is being updated
		hasUpdates := false

		// Update only the fields that were specified
		if updateType != "" {
			tType := ticket.Type(updateType)
			if !tType.Valid() {
				return fmt.Errorf("invalid type: %s (must be: bug, feature, or task)", updateType)
			}
			existingTicket.Type = tType
			hasUpdates = true
		}

		if updateStatus != "" {
			tStatus := ticket.Status(updateStatus)
			if !tStatus.Valid() {
				return fmt.Errorf("invalid status: %s (must be: open, in-progress, or closed)", updateStatus)
			}
			existingTicket.Status = tStatus
			hasUpdates = true
		}

		if updatePriority != "" {
			tPriority := ticket.Priority(updatePriority)
			if !tPriority.Valid() {
				return fmt.Errorf("invalid priority: %s (must be: low, medium, high, or undefined)", updatePriority)
			}
			existingTicket.Priority = tPriority
			hasUpdates = true
		}

		if updateTitle != "" {
			existingTicket.Title = updateTitle
			hasUpdates = true
		}

		if updateDesc != "" {
			existingTicket.Description = updateDesc
			hasUpdates = true
		}

		if cmd.Flags().Changed("critical") {
			existingTicket.CriticalPath = *updateCritical
			hasUpdates = true
		}

		if updateAssignedTo != "" {
			existingTicket.AssignedTo = &updateAssignedTo
			hasUpdates = true
		}

		if updateCreatedBy != "" {
			existingTicket.CreatedBy = &updateCreatedBy
			hasUpdates = true
		}

		if updateTags != "" {
			tagList := strings.Split(updateTags, ",")
			for i, tag := range tagList {
				tagList[i] = strings.TrimSpace(tag)
			}
			existingTicket.Tags = tagList
			hasUpdates = true
		}

		if updateFiles != "" {
			fileList := strings.Split(updateFiles, ",")
			for i, file := range fileList {
				fileList[i] = strings.TrimSpace(file)
			}
			existingTicket.Files = fileList
			hasUpdates = true
		}

		if updateComments != "" {
			commentList := strings.Split(updateComments, ",")
			for i, comment := range commentList {
				commentList[i] = strings.TrimSpace(comment)
			}
			existingTicket.Comments = commentList
			hasUpdates = true
		}

		// Check if at least one field is being updated
		if !hasUpdates {
			return fmt.Errorf("no fields specified to update")
		}

		// Call the Update method
		if err := existingTicket.Update(db, updateProject, updateID, updateFindTitle); err != nil {
			return fmt.Errorf("failed to update ticket: %w", err)
		}

		// Success message
		if updateID != "" {
			fmt.Printf("Successfully updated ticket with ID: %s in project: %s\n", updateID, updateProject)
		} else {
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
