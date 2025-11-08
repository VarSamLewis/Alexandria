package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"alexandria/internal/database"
	"alexandria/internal/ticket"

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
		// Validate that at least one identifier is provided
		if viewID == "" && viewTitle == "" {
			return fmt.Errorf("either --id or --title must be provided")
		}

		// Validate that project is provided
		if viewProject == "" {
			return fmt.Errorf("project is required")
		}

		// Parse ID if provided
		var ticketID int64
		if viewID != "" {
			var err error
			ticketID, err = strconv.ParseInt(viewID, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid ID format: %s (must be a number)", viewID)
			}
		}

		// Get database connection
		db := database.GetDB()
		if db == nil {
			return fmt.Errorf("database not initialized")
		}

		// Create a ticket instance to populate
		t := &ticket.Ticket{}

		// Call the View method
		if err := t.View(db, viewProject, ticketID, viewTitle); err != nil {
			return fmt.Errorf("failed to view ticket: %w", err)
		}

		// Output as JSON
		jsonData, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
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
