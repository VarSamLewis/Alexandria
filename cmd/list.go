package cmd

import (
	"alexandria/internal/database"
	"alexandria/internal/logger"
	"alexandria/internal/ticket"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	filterStatus     string
	filterType       string
	filterPriority   string
	filterAssignedTo string
	filterTags       string
	filterProject    string
	outputFormat     string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tickets from the database",
	Long:  `List all tickets from the database with optional filtering by status, type, priority, assigned user, or tags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Log.Debug("listing tickets", "project", filterProject, "status", filterStatus, "type", filterType, "output", outputFormat)

		// Get database connection
		db := database.GetDB()
		if db == nil {
			logger.Log.Error("database not initialized")
			return fmt.Errorf("database not initialized")
		}

		// Build filters
		filters := ticket.Filters{}

		if filterStatus != "" {
			status := ticket.Status(filterStatus)
			if !status.Valid() {
				logger.Log.Error("validation failed", "error", "invalid status", "status", filterStatus)
				return fmt.Errorf("invalid status: %s (must be: open, in-progress, or closed)", filterStatus)
			}
			filters.Status = &status
			logger.Log.Debug("applying status filter", "status", status)
		}

		if filterType != "" {
			tType := ticket.Type(filterType)
			if !tType.Valid() {
				logger.Log.Error("validation failed", "error", "invalid type", "type", filterType)
				return fmt.Errorf("invalid type: %s (must be: bug, feature, or task)", filterType)
			}
			filters.Type = &tType
			logger.Log.Debug("applying type filter", "type", tType)
		}

		if filterPriority != "" {
			priority := ticket.Priority(filterPriority)
			if !priority.Valid() {
				logger.Log.Error("validation failed", "error", "invalid priority", "priority", filterPriority)
				return fmt.Errorf("invalid priority: %s (must be: undefined, low, medium, or high)", filterPriority)
			}
			filters.Priority = &priority
			logger.Log.Debug("applying priority filter", "priority", priority)
		}

		if filterAssignedTo != "" {
			filters.AssignedTo = &filterAssignedTo
			logger.Log.Debug("applying assigned_to filter", "assigned_to", filterAssignedTo)
		}

		if filterProject != "" {
			filters.Project = &filterProject
			logger.Log.Debug("applying project filter", "project", filterProject)
		}

		if filterTags != "" {
			tagList := strings.Split(filterTags, ",")
			for i, tag := range tagList {
				tagList[i] = strings.TrimSpace(tag)
			}
			filters.Tags = tagList
			logger.Log.Debug("applying tags filter", "count", len(tagList))
		}

		// Query tickets
		logger.Log.Debug("querying tickets with filters")
		tickets, err := ticket.List(db, filters)
		if err != nil {
			logger.Log.Error("failed to list tickets", "error", err)
			return fmt.Errorf("failed to list tickets: %w", err)
		}

		logger.Log.Info("tickets retrieved", "count", len(tickets))

		if len(tickets) == 0 {
			logger.Log.Debug("no tickets found")
			fmt.Println("No tickets found.")
			return nil
		}

		// Output based on format
		logger.Log.Debug("formatting output", "format", outputFormat, "count", len(tickets))
		switch outputFormat {
		case "json":
			jsonData, err := json.MarshalIndent(tickets, "", "  ")
			if err != nil {
				logger.Log.Error("failed to marshal tickets", "error", err)
				return fmt.Errorf("failed to marshal tickets: %w", err)
			}
			fmt.Println(string(jsonData))

		case "table":
			printTicketsTable(tickets)

		case "summary":
			printTicketsSummary(tickets)

		default:
			logger.Log.Error("invalid output format", "format", outputFormat)
			return fmt.Errorf("invalid output format: %s (must be: json, table, or summary)", outputFormat)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&filterProject, "project", "", "Filter tickets by project")
	listCmd.Flags().StringVar(&filterStatus, "status", "", "Filter by status (open, in-progress, closed)")
	listCmd.Flags().StringVar(&filterType, "type", "", "Filter by type (bug, feature, task)")
	listCmd.Flags().StringVar(&filterPriority, "priority", "", "Filter by priority (undefined, low, medium, high)")
	listCmd.Flags().StringVar(&filterAssignedTo, "assigned-to", "", "Filter by assigned user")
	listCmd.Flags().StringVar(&filterTags, "tags", "", "Filter by tags (comma-separated)")
	listCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (json, table, summary)")
}

// printTicketsTable prints tickets in a table format
func printTicketsTable(tickets []ticket.Ticket) {
	// Print header
	fmt.Printf("%-6s %-18s %-10s %-10s %-10s %-35s %-13s %-12s\n",
		"ID", "PROJECT", "TYPE", "PRIORITY", "CRITICAL", "TITLE", "STATUS", "ASSIGNED TO")
	fmt.Println(strings.Repeat("-", 114))

	// Print rows
	for _, t := range tickets {
		assignedTo := "unassigned"
		if t.AssignedTo != nil {
			assignedTo = *t.AssignedTo
		}

		// Truncate title if too long
		title := t.Title
		if len(title) > 35 {
			title = title[:32] + "..."
		}

		fmt.Printf("%-6d %-18s %-10s %-10s %-10t %-35s %-13s %-12s\n",
			t.ID,
			t.Project,
			t.Type,
			t.Priority,
			t.CriticalPath,
			title,
			t.Status,
			assignedTo)
	}

	fmt.Printf("\nTotal: %d ticket(s)\n", len(tickets))
}

// printTicketsSummary prints a summary of each ticket
func printTicketsSummary(tickets []ticket.Ticket) {
	for i, t := range tickets {
		if i > 0 {
			fmt.Println()
		}

		fmt.Printf("ID: %d\n", t.ID)
		fmt.Printf("Type: %s | Priority: %s | Status: %s\n", t.Type, t.Priority, t.Status)
		fmt.Printf("Title: %s\n", t.Title)

		if t.Description != "" {
			fmt.Printf("Description: %s\n", t.Description)
		}

		if t.CriticalPath {
			fmt.Println("Critical Path: YES")
		}

		if t.AssignedTo != nil {
			fmt.Printf("Assigned To: %s\n", *t.AssignedTo)
		}

		if t.CreatedBy != nil {
			fmt.Printf("Created By: %s\n", *t.CreatedBy)
		}

		if len(t.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(t.Tags, ", "))
		}

		if len(t.Files) > 0 {
			fmt.Printf("Files: %d\n", len(t.Files))
		}

		if len(t.Comments) > 0 {
			fmt.Printf("Comments: %d\n", len(t.Comments))
		}

		fmt.Printf("Created: %s | Updated: %s\n", t.CreatedAt.Format("2006-01-02 15:04"), t.UpdatedAt.Format("2006-01-02 15:04"))
		fmt.Println(strings.Repeat("-", 80))
	}

	fmt.Printf("\nTotal: %d ticket(s)\n", len(tickets))
}


