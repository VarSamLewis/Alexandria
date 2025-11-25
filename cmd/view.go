package cmd

import (
	"alexandria/internal/database"
	"alexandria/internal/logger"
	"alexandria/internal/ticket"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	viewID      string
	viewTitle   string
	viewProject string
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a single ticket's details",
	Long:  `View the full details of a ticket by ID or title.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Log.Debug("viewing ticket", "id", viewID, "title", viewTitle, "project", viewProject)

		// Validate that at least one identifier is provided
		if viewID == "" && viewTitle == "" {
			logger.Log.Error("validation failed", "error", "no identifier provided")
			return fmt.Errorf("either --id or --title must be provided")
		}

		// Validate that project is provided
		if viewProject == "" {
			logger.Log.Error("validation failed", "error", "project is required")
			return fmt.Errorf("project is required")
		}

		// Parse ID if provided
		var ticketID int64
		if viewID != "" {
			var err error
			ticketID, err = strconv.ParseInt(viewID, 10, 64)
			if err != nil {
				logger.Log.Error("failed to parse ticket ID", "error", err, "id", viewID)
				return fmt.Errorf("invalid ID format: %s (must be a number)", viewID)
			}
			logger.Log.Debug("parsed ticket ID", "id", ticketID)
		}

		// Get database connection
		db := database.GetDB()
		if db == nil {
			logger.Log.Error("database not initialized")
			return fmt.Errorf("database not initialized")
		}

		// Create a ticket instance to populate
		t := &ticket.Ticket{}

		// Call the View method
		logger.Log.Debug("calling view method", "project", viewProject, "id", ticketID, "title", viewTitle)
		if err := t.View(db, viewProject, ticketID, viewTitle); err != nil {
			logger.Log.Error("failed to view ticket", "error", err, "project", viewProject)
			return fmt.Errorf("failed to view ticket: %w", err)
		}

		logger.Log.Info("ticket retrieved successfully", "id", t.ID, "title", t.Title)

		// Output as JSON
		jsonData, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			logger.Log.Error("failed to marshal ticket", "error", err)
			return fmt.Errorf("failed to marshal ticket: %w", err)
		}

		fmt.Println("Ticket details:")
		fmt.Println(string(jsonData))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)

	viewCmd.Flags().StringVarP(&viewID, "id", "i", "", "Ticket ID to view")
	viewCmd.Flags().StringVarP(&viewTitle, "title", "t", "", "Ticket title to view")
	viewCmd.Flags().StringVarP(&viewProject, "project", "p", "", "Project name (required)")
	viewCmd.MarkFlagRequired("project")
}
