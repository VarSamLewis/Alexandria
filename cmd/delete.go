package cmd

import (
	"alexandria/internal/database"
	"alexandria/internal/logger"
	"alexandria/internal/ticket"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	deleteID      string
	deleteTitle   string
	deleteProject string
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a ticket from the database",
	Long:  `Delete a ticket and all its related data (tags, files, comments) by ID or title.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Log.Debug("deleting ticket", "id", deleteID, "title", deleteTitle, "project", deleteProject)

		// Validate that at least one identifier is provided
		if deleteID == "" && deleteTitle == "" {
			logger.Log.Error("validation failed", "error", "no identifier provided")
			return fmt.Errorf("either --id or --title must be provided")
		}

		// Validate that project is provided
		if deleteProject == "" {
			logger.Log.Error("validation failed", "error", "project is required")
			return fmt.Errorf("project is required")
		}

		// Parse ID if provided
		var ticketID int64
		if deleteID != "" {
			var err error
			ticketID, err = strconv.ParseInt(deleteID, 10, 64)
			if err != nil {
				logger.Log.Error("failed to parse ticket ID", "error", err, "id", deleteID)
				return fmt.Errorf("invalid ID format: %s (must be a number)", deleteID)
			}
			logger.Log.Debug("parsed ticket ID", "id", ticketID)
		}

		// Get database connection
		db := database.GetDB()
		if db == nil {
			logger.Log.Error("database not initialized")
			return fmt.Errorf("database not initialized")
		}

		// Create a ticket instance for deletion
		t := &ticket.Ticket{}

		// Call the Delete method
		logger.Log.Debug("calling delete method", "project", deleteProject, "id", ticketID, "title", deleteTitle)
		if err := t.Delete(db, deleteProject, ticketID, deleteTitle); err != nil {
			logger.Log.Error("failed to delete ticket", "error", err, "project", deleteProject)
			return fmt.Errorf("failed to delete ticket: %w", err)
		}

		// Success message
		if ticketID != 0 {
			logger.Log.Info("ticket deleted successfully", "id", ticketID, "project", deleteProject)
			fmt.Printf("Successfully deleted ticket with ID: %d from project: %s\n", ticketID, deleteProject)
		} else {
			logger.Log.Info("ticket deleted successfully", "title", deleteTitle, "project", deleteProject)
			fmt.Printf("Successfully deleted ticket with title: %s from project: %s\n", deleteTitle, deleteProject)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&deleteID, "id", "i", "", "Ticket ID to delete")
	deleteCmd.Flags().StringVarP(&deleteTitle, "title", "t", "", "Ticket title to delete")
	deleteCmd.Flags().StringVarP(&deleteProject, "project", "p", "", "Project name (required)")
	if err := createCmd.MarkFlagRequired("project"); err != nil {
		panic(err)
	}

}
