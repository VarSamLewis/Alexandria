package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var name string
var enthusiastic bool

var greetCmd = &cobra.Command{
	Use:   "greet",
	Short: "Greet someone",
	Long:  `Greet someone with a friendly message. You can customize the greeting with flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		greeting := fmt.Sprintf("Hello, %s!", name)
		if enthusiastic {
			greeting += " How wonderful to meet you!"
		}
		fmt.Println(greeting)
	},
}

func init() {
	rootCmd.AddCommand(greetCmd)

	// Add flags for the greet command
	greetCmd.Flags().StringVarP(&name, "name", "n", "World", "Name of the person to greet")
	greetCmd.Flags().BoolVarP(&enthusiastic, "enthusiastic", "e", false, "Make the greeting more enthusiastic")
}
