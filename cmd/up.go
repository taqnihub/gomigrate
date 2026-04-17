package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/taqnihub/gomigrate/internal/tui"
	"github.com/taqnihub/gomigrate/migrate"
)

var upCmd = &cobra.Command{
	Use:   "up [n]",
	Short: "Apply pending migrations",
	Long: `Apply pending migrations to the database.

Examples:
  gomigrate up          # Apply all pending migrations
  gomigrate up 1        # Apply only the next pending migration
  gomigrate up 3        # Apply the next 3 pending migrations`,

	Args: cobra.MaximumNArgs(1),
	Run:  runUp,
}

func runUp(cmd *cobra.Command, args []string) {
	engine, cfg, err := newEngine()
	if err != nil {
		exitWithError(err)
	}
	defer engine.Close()

	// Header
	tui.Title("⚡ Applying Migrations")
	tui.KeyValue("Database", fmt.Sprintf("%s on %s:%d/%s",
		cfg.Driver, cfg.Host, cfg.Port, cfg.Database))

	// Get initial version for reporting
	beforeVersion, _, _ := engine.Version()

	start := time.Now()

	// Determine how many to apply
	if len(args) == 0 {
		err = engine.Up()
	} else {
		n, parseErr := strconv.Atoi(args[0])
		if parseErr != nil || n <= 0 {
			exitWithError(fmt.Errorf("invalid number %q: must be a positive integer", args[0]))
		}
		err = engine.UpN(n)
	}

	// Handle result
	tui.Newline()

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			tui.Info("Database is up to date — no pending migrations")
			return
		}
		exitWithError(fmt.Errorf("migration failed: %w", err))
	}

	afterVersion, dirty, _ := engine.Version()
	elapsed := time.Since(start).Round(time.Millisecond)

	if dirty {
		tui.Warning("Applied but state is dirty at version %d", afterVersion)
		tui.Muted("Run '%s' to clean up", tui.Code("gomigrate force <version>"))
		return
	}

	applied := afterVersion - beforeVersion
	tui.Success("Applied %d migration(s) in %s", applied, elapsed)
	tui.KeyValue("Version", fmt.Sprintf("%d → %d", beforeVersion, afterVersion))
}

func init() {
	rootCmd.AddCommand(upCmd)
}
