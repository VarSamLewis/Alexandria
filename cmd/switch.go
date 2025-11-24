package cmd

import (
	"alexandria/internal/config"
	"alexandria/internal/database"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var showStatus bool

var switchCmd = &cobra.Command{
	Use:   "source [sqlite|turso]",
	Short: "Switch between SQLite and Turso databases",
	Long: `Switch the active database between local SQLite and Turso cloud database.

Examples:
  alexandria source sqlite   # Switch to local SQLite database
  alexandria source turso    # Switch to Turso cloud database
  alexandria source --status # Show current database configuration`,
	Args: func(cmd *cobra.Command, args []string) error {
		if showStatus {
			return nil
		}
		if len(args) != 1 {
			return fmt.Errorf("requires exactly one argument: sqlite or turso")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// If --status flag is set, show current database info
		if showStatus {
			return showDatabaseStatus()
		}

		dbType := args[0]

		// Validate database type
		if dbType != config.DBTypeSQLite && dbType != config.DBTypeTurso {
			return fmt.Errorf("invalid database type: %s (must be sqlite or turso)", dbType)
		}

		// If switching to Turso, validate environment variables
		if dbType == config.DBTypeTurso {
			if err := validateTursoEnv(); err != nil {
				return err
			}
		}

		// Switch the database
		if err := config.SwitchDB(dbType); err != nil {
			return fmt.Errorf("failed to switch database: %w", err)
		}

		fmt.Printf("Switching to %s database...\n", dbType)

		// Close current database connection
		if err := database.Close(); err != nil {
			return fmt.Errorf("failed to close current database: %w", err)
		}

		// Reinitialize database with new configuration
		if err := database.Init(""); err != nil {
			return fmt.Errorf("failed to connect to %s database: %w", dbType, err)
		}

		fmt.Printf("Successfully switched to %s database.\n", dbType)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
	switchCmd.Flags().BoolVar(&showStatus, "status", false, "Show current database configuration")
}

// validateTursoEnv checks if required Turso environment variables are set
func validateTursoEnv() error {
	tursoURL := os.Getenv("TURSO_URL")
	tursoToken := os.Getenv("TURSO_AUTH_TOKEN")

	if tursoURL == "" {
		return fmt.Errorf("TURSO_URL environment variable is not set")
	}
	if tursoToken == "" {
		return fmt.Errorf("TURSO_AUTH_TOKEN environment variable is not set")
	}

	return nil
}

// showDatabaseStatus displays the current database configuration
func showDatabaseStatus() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("Current Database Configuration:")
	fmt.Println("================================")
	fmt.Printf("Database Type: %s\n", cfg.DatabaseType)

	if cfg.DatabaseType == config.DBTypeTurso {
		tursoURL := os.Getenv("TURSO_URL")
		if tursoURL != "" {
			fmt.Printf("Turso URL: %s\n", tursoURL)
		} else {
			fmt.Println("Warning: TURSO_URL environment variable is not set")
		}

		tursoToken := os.Getenv("TURSO_AUTH_TOKEN")
		if tursoToken != "" {
			// Show only first few characters of the token for security
			tokenPreview := tursoToken
			if len(tursoToken) > 8 {
				tokenPreview = tursoToken[:8] + "..."
			}
			fmt.Printf("Turso Token: %s\n", tokenPreview)
		} else {
			fmt.Println("Warning: TURSO_AUTH_TOKEN environment variable is not set")
		}
	}

	return nil
}
