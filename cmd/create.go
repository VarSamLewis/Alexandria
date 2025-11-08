package cmd

import (
	"encoding/json"
	"fmt"
	"alexandria/internal/database"
	"alexandria/internal/ticket"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	title       string
	description string
	ticketType  string
	criticalpath bool
	priority    string
	assignedTo  string
	tags        string
	createdBy   string
	project     string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new ticket",
	Long:  `Create a new ticket with the specified title, description, type, priority, and other attributes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate required fields
		if title == "" {
			return fmt.Errorf("title is required")
		}

		// Parse type
		tType := ticket.Type(ticketType)
		if !tType.Valid() {
			return fmt.Errorf("invalid type: %s (must be: bug, feature, or task)", ticketType)
		}

		// Parse priority
		tPriority := ticket.Priority(priority)
		if !tPriority.Valid() {
			return fmt.Errorf("invalid priority: %s (must be: low, medium, or high)", priority)
		} 

		// Parse tags
		var tagList []string
		if tags != "" {
			tagList = strings.Split(tags, ",")
			for i, tag := range tagList {
				tagList[i] = strings.TrimSpace(tag)
			}
		}

		// Generate simple ID (timestamp-based)
		id := strconv.FormatInt(time.Now().UnixNano(), 10)

		// Create ticket
		newTicket := ticket.Ticket{
			ID:          id,
			Type:        tType,
			Title:       title,
			Description: description,
      CriticalPath: criticalpath,  
			Status:      ticket.StatusOpen, // Default to open
			Priority:    tPriority,
			Tags:        tagList,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Handle optional fields
		if assignedTo != "" {
			newTicket.AssignedTo = &assignedTo
		}
		if createdBy != "" {
			newTicket.CreatedBy = &createdBy
		}

		// Validate project is provided
		if project == "" {
			return fmt.Errorf("project is required")
		}

		// Save ticket to database
		db := database.GetDB()
		if db == nil {
			return fmt.Errorf("database not initialized")
		}

		if err := newTicket.Create(db, project); err != nil {
			return fmt.Errorf("failed to save ticket: %w", err)
		}

		// Output as JSON
		jsonData, err := json.MarshalIndent(newTicket, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal ticket: %w", err)
		}

		fmt.Println("Ticket created and saved successfully:")
		fmt.Println(string(jsonData))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&title, "title", "t", "", "Ticket title (required)")
	createCmd.Flags().StringVarP(&description, "description", "d", "", "Ticket description")
	createCmd.Flags().StringVar(&ticketType, "type", "task", "Ticket type (bug, feature, task)")
	createCmd.Flags().StringVarP(&priority, "priority", "p", "undefined", "Ticket priority (low, medium, high)")
	createCmd.Flags().BoolVarP(&criticalpath, "criticalpath", "c", false, "Mark ticket as critical path")
	createCmd.Flags().StringVarP(&assignedTo, "assigned-to", "a", "", "Assign ticket to user")
	createCmd.Flags().StringVar(&createdBy, "created-by", "", "Ticket creator")
	createCmd.Flags().StringVar(&tags, "tags", "", "Comma-separated list of tags")
	createCmd.Flags().StringVar(&project, "project", "", "Project name (required)")
	createCmd.MarkFlagRequired("title")
	createCmd.MarkFlagRequired("project")
}
