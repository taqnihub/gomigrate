package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/taqnihub/gomigrate/internal/tui"
	"github.com/taqnihub/gomigrate/migrate"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Long:  `Show all migrations and their status (applied or pending).`,
	Args:  cobra.NoArgs,
	Run:   runStatus,
}

func runStatus(cmd *cobra.Command, args []string) {
	engine, cfg, err := newEngine()
	if err != nil {
		exitWithError(err)
	}
	defer engine.Close()

	currentVersion, dirty, err := engine.Version()
	if err != nil {
		exitWithError(fmt.Errorf("failed to get version: %w", err))
	}

	files, err := engine.List()
	if err != nil {
		exitWithError(fmt.Errorf("failed to list migrations: %w", err))
	}

	// Header
	tui.Title("📋 Migration Status")
	tui.KeyValue("Database", fmt.Sprintf("%s on %s:%d/%s",
		cfg.Driver, cfg.Host, cfg.Port, cfg.Database))
	tui.KeyValue("Directory", tui.Path(cfg.MigrationsDir))
	tui.Newline()

	if len(files) == 0 {
		fmt.Println(tui.WarningBox(
			"No migration files found.\n\n" +
				"Run '" + tui.Code("gomigrate create <name>") + "' to create your first migration.",
		))
		return
	}

	// Render the table
	renderStatusTable(files, currentVersion)

	// Summary
	applied := 0
	pending := 0
	for _, f := range files {
		if f.Version <= currentVersion {
			applied++
		} else {
			pending++
		}
	}

	tui.Newline()
	if dirty {
		fmt.Println(tui.ErrorBox(
			fmt.Sprintf("Database is DIRTY at version %d\n", currentVersion) +
				"A migration failed partway. Fix the database manually, then run:\n" +
				tui.Code("gomigrate force <version>"),
		))
	}

	tui.Muted("%d total · %d applied · %d pending",
		len(files), applied, pending)
}

// renderStatusTable renders a pretty table of migrations using lipgloss.
func renderStatusTable(files []migrate.MigrationFile, currentVersion uint) {
	// Column widths
	const (
		versionW = 18
		nameW    = 40
		statusW  = 14
	)

	headerStyle := lipgloss.NewStyle().
		Foreground(tui.ColorMuted).
		Bold(true).
		Padding(0, 1)

	rowStyle := lipgloss.NewStyle().Padding(0, 1)

	// Divider line
	divider := tui.Dim(strings.Repeat("─", versionW+nameW+statusW+6))

	// Header row
	fmt.Println("  " +
		headerStyle.Width(versionW).Render("VERSION") +
		headerStyle.Width(nameW).Render("NAME") +
		headerStyle.Width(statusW).Render("STATUS"))

	fmt.Println("  " + divider)

	// Data rows
	for _, f := range files {
		var statusIcon, statusText string
		var statusColor lipgloss.AdaptiveColor

		if f.Version <= currentVersion {
			statusIcon = tui.IconSuccess
			statusText = "applied"
			statusColor = tui.ColorSuccess
		} else {
			statusIcon = "○"
			statusText = "pending"
			statusColor = tui.ColorWarning
		}

		status := lipgloss.NewStyle().
			Foreground(statusColor).
			Render(fmt.Sprintf("%s %s", statusIcon, statusText))

		version := fmt.Sprintf("%d", f.Version)
		name := truncate(f.Name, nameW-2)

		fmt.Println("  " +
			rowStyle.Width(versionW).Render(version) +
			rowStyle.Width(nameW).Render(name) +
			rowStyle.Width(statusW).Render(status))
	}
}

// truncate shortens a string to max length.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
