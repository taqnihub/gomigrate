package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new migration file",
	Long: `Create a pair of up/down SQL migration files with a timestamp version.

The name is sanitized (lowercased, spaces become underscores).

Examples:
  gomigrate create add_users_table
  gomigrate create "add index on email"
  gomigrate create remove_legacy_columns`,

	Args: cobra.MinimumNArgs(1),
	Run:  runCreate,
}

func runCreate(cmd *cobra.Command, args []string) {
	// Join all args as the name (allows "add users table" without quotes)
	name := strings.Join(args, " ")

	engine, _, err := newEngine()
	if err != nil {
		exitWithError(err)
	}
	defer engine.Close()

	upPath, downPath, err := engine.Create(name)
	if err != nil {
		exitWithError(fmt.Errorf("create failed: %w", err))
	}

	printSuccess("Created migration files:")
	fmt.Printf("  UP:   %s\n", upPath)
	fmt.Printf("  DOWN: %s\n", downPath)
	fmt.Println()
	printInfo("Edit the files to add your SQL, then run: gomigrate up")
}

func init() {
	rootCmd.AddCommand(createCmd)
}
