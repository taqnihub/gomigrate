package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/taqnihub/gomigrate/internal/tui"
)

var forceCmd = &cobra.Command{
	Use:   "force <version>",
	Short: "Force set the migration version (fix dirty state)",
	Long: `Force set the migration version without running any migrations.

Use this when a migration failed partway and left the database in a
"dirty" state. You'll need to manually fix the database first, then
use 'force' to tell gomigrate what version it's at.`,

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

	tui.Title("🔧 Forcing Migration Version")

	if err := engine.Force(version); err != nil {
		exitWithError(fmt.Errorf("force failed: %w", err))
	}

	tui.Success("Version forced to %d", version)
	tui.Newline()
	fmt.Println(tui.WarningBox(
		"No SQL was executed.\n" +
			"Make sure your database schema matches this version!",
	))
}

func init() {
	rootCmd.AddCommand(forceCmd)
}
