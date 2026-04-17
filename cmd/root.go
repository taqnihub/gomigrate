package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"github.com/taqnihub/gomigrate/internal/tui"
)

// Persistent flags
var (
	cfgFile    string
	driver     string
	dsn        string
	migDir     string
	verbose    bool
	dryRun     bool
	noInteract bool
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "gomigrate",
	Short: "A beautiful database migration CLI for Go developers",
	Long: `GoMigrate is a friendly wrapper around golang-migrate that makes
database migrations simple with short commands, config files, and an
interactive TUI.

Quick start:
  gomigrate init                  # Set up your project
  gomigrate create add_users      # Create a migration
  gomigrate up                    # Apply all pending migrations
  gomigrate status                # See what's applied
  gomigrate down                  # Revert the last migration

Learn more: https://github.com/taqnihub/gomigrate`,

	Run: runRoot,
}

// runRoot launches the interactive menu when no subcommand is given.
func runRoot(cmd *cobra.Command, args []string) {
	// Try to load config — if we can't, show a helpful message
	cfg, err := loadConfig()
	if err != nil {
		tui.Warning("No configuration found")
		tui.Newline()
		tui.Muted("Run %s to set up gomigrate in this directory", tui.Code("gomigrate init"))
		os.Exit(0)
	}

	// Launch the interactive menu
	menu := tui.NewMenuModel(cfg)
	p := tea.NewProgram(menu)

	finalModel, err := p.Run()
	if err != nil {
		exitWithError(fmt.Errorf("menu error: %w", err))
	}

	// Act on the user's choice
	selected := finalModel.(tui.MenuModel).Selected()
	if selected == "" {
		return // user quit
	}

	// Dispatch to the appropriate action
	dispatchMenuAction(selected)
}

// dispatchMenuAction runs the appropriate command based on menu selection.
func dispatchMenuAction(action string) {
	switch action {
	case "up":
		runUp(nil, nil)
	case "down":
		runDown(nil, nil)
	case "create":
		runCreateInteractive()
	case "status":
		runStatus(nil, nil)
	case "force":
		runForceInteractive()
	case "quit":
		return
	}
}

// runCreateInteractive prompts for a migration name, then creates it.
func runCreateInteractive() {
	input := tui.NewInputModel("Migration name", "e.g., add_users_table")
	p := tea.NewProgram(input)

	finalModel, err := p.Run()
	if err != nil {
		exitWithError(fmt.Errorf("input error: %w", err))
	}

	result := finalModel.(tui.InputModel)
	if result.Cancelled() || !result.Submitted() {
		tui.Muted("Cancelled")
		return
	}

	name := result.Value()
	runCreate(nil, []string{name})
}

// runForceInteractive prompts for a version number, then forces it.
func runForceInteractive() {
	input := tui.NewInputModel("Version number to force", "e.g., 20260417090000")
	p := tea.NewProgram(input)

	finalModel, err := p.Run()
	if err != nil {
		exitWithError(fmt.Errorf("input error: %w", err))
	}

	result := finalModel.(tui.InputModel)
	if result.Cancelled() || !result.Submitted() {
		tui.Muted("Cancelled")
		return
	}

	runForce(nil, []string{result.Value()})
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to config file (default: .gomigrate.yml)")
	rootCmd.PersistentFlags().StringVarP(&driver, "driver", "d", "", "database driver (mysql or postgres)")
	rootCmd.PersistentFlags().StringVar(&dsn, "dsn", "", "full database DSN (overrides other connection flags)")
	rootCmd.PersistentFlags().StringVar(&migDir, "dir", "", "migrations directory")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show detailed output")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "show what would happen without executing")
	rootCmd.PersistentFlags().BoolVar(&noInteract, "no-interactive", false, "disable TUI, plain output (for CI/CD)")

	cobra.OnInitialize(func() {
		if noInteract || os.Getenv("NO_COLOR") != "" {
			lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(os.Stdout, termenv.WithProfile(termenv.Ascii)))
		}
	})

	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, BuildDate)
}
