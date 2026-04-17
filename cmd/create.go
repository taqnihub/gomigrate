package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/taqnihub/gomigrate/internal/tui"
)

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new migration file",
	Long: `Create a pair of up/down SQL migration files with a timestamp version.

Examples:
  gomigrate create add_users_table
  gomigrate create "add index on email"`,

	Args: cobra.MinimumNArgs(1),
	Run:  runCreate,
}

func runCreate(cmd *cobra.Command, args []string) {
	name := strings.Join(args, " ")

	engine, _, err := newEngine()
	if err != nil {
		exitWithError(err)
	}
	defer engine.Close()

	tui.Title("✨ Creating Migration")

	upPath, downPath, err := engine.Create(name)
	if err != nil {
		exitWithError(fmt.Errorf("create failed: %w", err))
	}

	// Display paths — prefer relative to cwd, fall back to absolute
	upDisplay := displayPath(upPath)
	downDisplay := displayPath(downPath)

	tui.Success("Created migration files")
	tui.Newline()
	tui.KeyValue("UP", tui.Path(upDisplay))
	tui.KeyValue("DOWN", tui.Path(downDisplay))

	tui.Newline()
	tui.Muted("Next: edit the files with your SQL, then run %s", tui.Code("gomigrate up"))
}

// displayPath returns a user-friendly path — relative to cwd if possible,
// otherwise the absolute path. Never returns empty.
func displayPath(p string) string {
	if p == "" {
		return "(unknown)"
	}

	cwd, err := os.Getwd()
	if err != nil {
		return p // fallback to whatever we have
	}

	rel, err := filepath.Rel(cwd, p)
	if err != nil || rel == "" {
		return p // fallback to absolute
	}

	// If the relative path goes up too many levels, it's cleaner to show absolute
	if strings.HasPrefix(rel, "../..") {
		return p
	}

	return rel
}

func init() {
	rootCmd.AddCommand(createCmd)
}
