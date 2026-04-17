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

var downAll bool

var downCmd = &cobra.Command{
	Use:   "down [n]",
	Short: "Revert applied migrations",
	Long: `Revert the last N applied migrations (default: 1).

Examples:
  gomigrate down           # Revert the last migration
  gomigrate down 3         # Revert the last 3 migrations
  gomigrate down --all     # Revert ALL migrations (destructive!)`,

	Args: cobra.MaximumNArgs(1),
	Run:  runDown,
}

func runDown(cmd *cobra.Command, args []string) {
	engine, cfg, err := newEngine()
	if err != nil {
		exitWithError(err)
	}
	defer engine.Close()

	// Title changes based on --all
	if downAll {
		tui.Title("⚠️  Reverting ALL Migrations")
		fmt.Println(tui.WarningBox(
			"This will drop all tables managed by gomigrate.\n" +
				"All data in those tables will be lost.",
		))
	} else {
		tui.Title("⬇  Reverting Migrations")
	}

	tui.KeyValue("Database", fmt.Sprintf("%s on %s:%d/%s",
		cfg.Driver, cfg.Host, cfg.Port, cfg.Database))

	beforeVersion, _, _ := engine.Version()
	start := time.Now()

	if downAll {
		err = engine.DownAll()
	} else {
		n := 1
		if len(args) == 1 {
			parsed, parseErr := strconv.Atoi(args[0])
			if parseErr != nil || parsed <= 0 {
				exitWithError(fmt.Errorf("invalid number %q: must be a positive integer", args[0]))
			}
			n = parsed
		}
		err = engine.Down(n)
	}

	tui.Newline()

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			tui.Info("No applied migrations to revert")
			return
		}
		exitWithError(fmt.Errorf("revert failed: %w", err))
	}

	afterVersion, dirty, _ := engine.Version()
	elapsed := time.Since(start).Round(time.Millisecond)

	if dirty {
		tui.Warning("Reverted but state is dirty at version %d", afterVersion)
		return
	}

	reverted := beforeVersion - afterVersion
	if afterVersion == 0 {
		tui.Success("Reverted %d migration(s) — database is now empty (%s)", reverted, elapsed)
	} else {
		tui.Success("Reverted %d migration(s) in %s", reverted, elapsed)
		tui.KeyValue("Version", fmt.Sprintf("%d → %d", beforeVersion, afterVersion))
	}
}

func init() {
	downCmd.Flags().BoolVar(&downAll, "all", false, "revert ALL migrations (destructive)")
	rootCmd.AddCommand(downCmd)
}
