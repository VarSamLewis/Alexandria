package cmd

import (
	"alexandria/internal/database"
	"alexandria/internal/logger"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/common-nighthawk/go-figure"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "alexandria",
	Short: "A simple ticket management CLI",
	Long:  `A command-line tool for managing tickets and tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show banner only when running root command with no arguments
		myBanner := figure.NewFigure("Alexandria", "", true)
		myBanner.Print()
		fmt.Println() // blank line for spacing
		cmd.Help()
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize logger based on verbose flag
		logger.Init(verbose)
		logger.Log.Debug("starting alexandria", "verbose", verbose)

		// Load .env file from project root (ignore error if file doesn't exist)
		_ = godotenv.Load()
		logger.Log.Debug("loaded environment variables")

		// Initialize database with default path
		logger.Log.Debug("initializing database connection")
		if err := database.Init(""); err != nil {
			logger.Log.Error("failed to initialize database", "error", err)
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		logger.Log.Info("database initialized successfully")
		return nil
	},
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}

func Execute() {
	// Ensure database is closed on exit
	defer func() {
		if err := database.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close database: %v\n", err)
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
