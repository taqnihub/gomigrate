package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// These flags are available on ALL commands (persistent flags).
var (
	cfgFile    string // --config path
	driver     string // --driver override
	dsn        string // --dsn override
	migDir     string // --dir override
	verbose    bool   // --verbose
	dryRun     bool   // --dry-run
	noInteract bool   // --no-interactive
)

// Version info (set at build time via ldflags).
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// rootCmd is the top-level `gomigrate` command.
// When run without a subcommand, it will launch the interactive TUI (coming in Step 9).
var rootCmd = &cobra.Command{
	Use:   "gomigrate",
	Short: "A beautiful database migration CLI for Go developers",
	Long: `GoMigrate is a easy friendly wrapper around golang-migrate that makes
database migrations simple with short commands, config files, and an
interactive TUI.

Quick start:
  gomigrate init                  # Set up your project
  gomigrate create add_users      # Create a migration
  gomigrate up                    # Apply all pending migrations
  gomigrate status                # See what's applied
  gomigrate down                  # Revert the last migration

Learn more: https://github.com/taqnihub/gomigrate`,

	// If no subcommand is given, we'll launch the TUI later.
	// For now, just show help.
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to GoMigrate!")
		fmt.Println("Run 'gomigrate --help' to see available commands.")
		fmt.Println("(Interactive TUI coming in a later step.)")
	},
}

// Execute runs the root command. Called from main().
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// init() registers persistent flags and subcommands.
// Go automatically calls init() when the package loads.
func init() {
	// Persistent flags — available on all subcommands
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to config file (default: .gomigrate.yml)")
	rootCmd.PersistentFlags().StringVarP(&driver, "driver", "d", "", "database driver (mysql or postgres)")
	rootCmd.PersistentFlags().StringVar(&dsn, "dsn", "", "full database DSN (overrides other connection flags)")
	rootCmd.PersistentFlags().StringVar(&migDir, "dir", "", "migrations directory")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show detailed output")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "show what would happen without executing")
	rootCmd.PersistentFlags().BoolVar(&noInteract, "no-interactive", false, "disable TUI, plain output (for CI/CD)")

	// Version flag
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, BuildDate)
}
