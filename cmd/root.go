package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "A simple CLI application built with Cobra",
	Long:  `This is a simple CLI application demonstrating the use of Cobra framework in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to MyCLI! Use --help to see available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Add global flags here if needed
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mycli.yaml)")
}
