package cmd

import (
	"alexandria/internal/config"
	"alexandria/internal/database"
	"alexandria/internal/logger"
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
			logger.Log.Debug("showing database status")
			return showDatabaseStatus()
		}

		dbType := args[0]
		logger.Log.Debug("switching database", "target", dbType)

		// Validate database type
		if dbType != config.DBTypeSQLite && dbType != config.DBTypeTurso {
			logger.Log.Error("invalid database type", "type", dbType)
			return fmt.Errorf("invalid database type: %s (must be sqlite or turso)", dbType)
		}

		// If switching to Turso, validate environment variables
		if dbType == config.DBTypeTurso {
			logger.Log.Debug("validating Turso environment variables")
			if err := validateTursoEnv(); err != nil {
				logger.Log.Error("Turso environment validation failed", "error", err)
				return err
			}
		}

		// Switch the database
		logger.Log.Debug("updating config", "database_type", dbType)
		if err := config.SwitchDB(dbType); err != nil {
			logger.Log.Error("failed to update config", "error", err, "database_type", dbType)
			return fmt.Errorf("failed to switch database: %w", err)
		}

		fmt.Printf("Switching to %s database...\n", dbType)

		// Close current database connection
		logger.Log.Debug("closing current database connection")
		if err := database.Close(); err != nil {
			logger.Log.Error("failed to close database", "error", err)
			return fmt.Errorf("failed to close current database: %w", err)
		}

		// Reinitialize database with new configuration
		logger.Log.Debug("initializing new database connection", "database_type", dbType)
		if err := database.Init(""); err != nil {
			logger.Log.Error("failed to initialize database", "error", err, "database_type", dbType)
			return fmt.Errorf("failed to connect to %s database: %w", dbType, err)
		}

		logger.Log.Info("successfully switched database", "database_type", dbType)
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
		logger.Log.Error("Turso environment variable missing", "variable", "TURSO_URL")
		return fmt.Errorf("TURSO_URL environment variable is not set")
	}
	if tursoToken == "" {
		logger.Log.Error("Turso environment variable missing", "variable", "TURSO_AUTH_TOKEN")
		return fmt.Errorf("TURSO_AUTH_TOKEN environment variable is not set")
	}

	logger.Log.Debug("Turso environment variables validated", "url", tursoURL)
	return nil
}

// showDatabaseStatus displays the current database configuration
func showDatabaseStatus() error {
	cfg, err := config.Load()
	if err != nil {
		logger.Log.Error("failed to load config", "error", err)
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger.Log.Debug("loaded database config", "type", cfg.DatabaseType)

	fmt.Println("Current Database Configuration:")
	fmt.Println("================================")
	fmt.Printf("Database Type: %s\n", cfg.DatabaseType)

	if cfg.DatabaseType == config.DBTypeTurso {
		tursoURL := os.Getenv("TURSO_URL")
		if tursoURL != "" {
			fmt.Printf("Turso URL: %s\n", tursoURL)
		} else {
			fmt.Println("Warning: TURSO_URL environment variable is not set")
			logger.Log.Warn("Turso URL not set")
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
			logger.Log.Warn("Turso auth token not set")
		}
	}

	return nil
}
