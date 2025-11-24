package cmd

import (
	"fmt"
	"alexandria/internal/database"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "alexandria",
	Short: "A simple ticket management CLI",
	Long:  `A command-line tool for managing tickets and tasks.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load .env file from project root (ignore error if file doesn't exist)
		_ = godotenv.Load()

		// Initialize database with default path
		if err := database.Init(""); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		return nil
	},
}

func Execute() {
	// Ensure database is closed on exit
	defer database.Close()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
