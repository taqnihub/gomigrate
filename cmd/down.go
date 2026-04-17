package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
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

	// Find applied migrations BEFORE reverting
	beforeVersion, _, _ := engine.Version()
	applied, err := appliedMigrations(engine, beforeVersion)
	if err != nil {
		exitWithError(fmt.Errorf("failed to list applied: %w", err))
	}

	if len(applied) == 0 {
		tui.Newline()
		tui.Info("No applied migrations to revert")
		return
	}

	// Determine how many to revert
	toRevert := 1
	if downAll {
		toRevert = len(applied)
	} else if len(args) == 1 {
		parsed, parseErr := strconv.Atoi(args[0])
		if parseErr != nil || parsed <= 0 {
			exitWithError(fmt.Errorf("invalid number %q: must be a positive integer", args[0]))
		}
		toRevert = parsed
		if toRevert > len(applied) {
			toRevert = len(applied)
		}
	}

	// Pick the last N applied migrations (these are what we'll revert)
	// applied is sorted ascending, so the last ones are the most recent
	toRevertList := applied[len(applied)-toRevert:]

	// Show what will be reverted (in reverse order since we unwind newest first)
	tui.Newline()
	tui.Muted("Will revert %d migration(s):", toRevert)
	for i := len(toRevertList) - 1; i >= 0; i-- {
		fmt.Printf("    %s  %d  %s\n",
			tui.Dim("←"),
			toRevertList[i].Version,
			toRevertList[i].Name,
		)
	}
	tui.Newline()

	// Run it
	start := time.Now()
	if downAll {
		err = engine.DownAll()
	} else {
		err = engine.Down(toRevert)
	}

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			tui.Info("No applied migrations to revert")
			return
		}
		exitWithError(fmt.Errorf("revert failed: %w", err))
	}

	afterVersion, dirty, _ := engine.Version()
	elapsed := time.Since(start).Round(time.Millisecond)

	// Print each reverted migration with a checkmark (newest first)
	for i := len(toRevertList) - 1; i >= 0; i-- {
		icon := lipgloss.NewStyle().Foreground(tui.ColorSuccess).Bold(true).Render(tui.IconSuccess)
		fmt.Printf("    %s  %d  %s\n", icon, toRevertList[i].Version, toRevertList[i].Name)
	}

	tui.Newline()

	if dirty {
		tui.Warning("Reverted but state is dirty at version %d", afterVersion)
		return
	}

	if afterVersion == 0 {
		tui.Success("Reverted %d migration(s) — database is now empty (%s)", len(toRevertList), elapsed)
	} else {
		tui.Success("Reverted %d migration(s) in %s", len(toRevertList), elapsed)
		tui.KeyValue("Version", fmt.Sprintf("%d → %d", beforeVersion, afterVersion))
	}
}

// countAppliedMigrations returns how many migrations are currently applied.
// This is correct regardless of whether versions are sequential or timestamp-based.
func countAppliedMigrations(engine *migrate.Engine) (int, error) {
	currentVersion, _, err := engine.Version()
	if err != nil {
		return 0, err
	}

	// No migrations applied
	if currentVersion == 0 {
		return 0, nil
	}

	files, err := engine.List()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, f := range files {
		if f.Version <= currentVersion {
			count++
		}
	}
	return count, nil
}

func init() {
	downCmd.Flags().BoolVar(&downAll, "all", false, "revert ALL migrations (destructive)")
	rootCmd.AddCommand(downCmd)
}
