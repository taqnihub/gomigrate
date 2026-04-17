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

	// Find pending migrations BEFORE running
	beforeVersion, _, _ := engine.Version()
	pending, err := pendingMigrations(engine, beforeVersion)
	if err != nil {
		exitWithError(fmt.Errorf("failed to list pending: %w", err))
	}

	if len(pending) == 0 {
		tui.Newline()
		tui.Info("Database is up to date — no pending migrations")
		return
	}

	// Determine how many to apply (affects what we display)
	toApply := len(pending)
	if len(args) > 0 {
		n, parseErr := strconv.Atoi(args[0])
		if parseErr != nil || n <= 0 {
			exitWithError(fmt.Errorf("invalid number %q: must be a positive integer", args[0]))
		}
		if n < toApply {
			toApply = n
		}
	}

	// Show what will be applied
	tui.Newline()
	tui.Muted("Will apply %d migration(s):", toApply)
	for i := 0; i < toApply; i++ {
		fmt.Printf("    %s  %d  %s\n",
			tui.Dim("→"),
			pending[i].Version,
			pending[i].Name,
		)
	}
	tui.Newline()

	// Run it
	start := time.Now()
	if len(args) == 0 {
		err = engine.Up()
	} else {
		err = engine.UpN(toApply)
	}

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			tui.Info("Database is up to date")
			return
		}
		exitWithError(fmt.Errorf("migration failed: %w", err))
	}

	// Show what actually got applied
	afterVersion, dirty, _ := engine.Version()
	elapsed := time.Since(start).Round(time.Millisecond)

	// Figure out which migrations were applied (between before and after)
	var applied []migrate.MigrationFile
	for _, m := range pending {
		if m.Version > beforeVersion && m.Version <= afterVersion {
			applied = append(applied, m)
		}
	}

	// Print each applied migration with a checkmark
	for _, m := range applied {
		icon := lipgloss.NewStyle().Foreground(tui.ColorSuccess).Bold(true).Render(tui.IconSuccess)
		fmt.Printf("    %s  %d  %s\n", icon, m.Version, m.Name)
	}

	tui.Newline()

	if dirty {
		tui.Warning("Applied but state is dirty at version %d", afterVersion)
		tui.Muted("Run %s to clean up", tui.Code("gomigrate force <version>"))
		return
	}

	tui.Success("Applied %d migration(s) in %s", len(applied), elapsed)
	tui.KeyValue("Version", fmt.Sprintf("%d → %d", beforeVersion, afterVersion))
}

func init() {
	rootCmd.AddCommand(upCmd)
}
