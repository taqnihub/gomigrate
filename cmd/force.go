package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var forceCmd = &cobra.Command{
	Use:   "force <version>",
	Short: "Force set the migration version (fix dirty state)",
	Long: `Force set the migration version without running any migrations.

Use this when a migration failed partway and left the database in a
"dirty" state. You'll need to manually fix the database first, then
use 'force' to tell gomigrate what version it's at.

Examples:
  gomigrate force 3        # Set version to 3 (no SQL executed)
  gomigrate force 0        # Reset version to none`,

	Args: cobra.ExactArgs(1),
	Run:  runForce,
}

func runForce(cmd *cobra.Command, args []string) {
	version, err := strconv.Atoi(args[0])
	if err != nil || version < 0 {
		exitWithError(fmt.Errorf("invalid version %q: must be a non-negative integer", args[0]))
	}

	engine, _, err := newEngine()
	if err != nil {
		exitWithError(err)
	}
	defer engine.Close()

	if err := engine.Force(version); err != nil {
		exitWithError(fmt.Errorf("force failed: %w", err))
	}

	printSuccess("Migration version forced to %d", version)
	printWarn("No SQL was executed — make sure your database schema matches!")
}

func init() {
	rootCmd.AddCommand(forceCmd)
}
