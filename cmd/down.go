package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/taqnihub/gomigrate/migrate"
)

var downAll bool // --all flag

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

	printInfo("Connected to %s://%s:%d/%s", cfg.Driver, cfg.Host, cfg.Port, cfg.Database)

	start := time.Now()

	// Handle --all flag
	if downAll {
		printWarn("Reverting ALL migrations — this will drop all tables!")
		err = engine.DownAll()
	} else {
		// Default: revert 1, or N if specified
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

	// Handle result
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			printInfo("No applied migrations to revert")
			return
		}
		exitWithError(fmt.Errorf("revert failed: %w", err))
	}

	version, dirty, _ := engine.Version()
	elapsed := time.Since(start).Round(time.Millisecond)

	if dirty {
		printWarn("Revert succeeded but state is dirty at version %d", version)
	} else if version == 0 {
		printSuccess("All migrations reverted (%s)", elapsed)
	} else {
		printSuccess("Migrations reverted — now at version %d (%s)", version, elapsed)
	}
}

func init() {
	downCmd.Flags().BoolVar(&downAll, "all", false, "revert ALL migrations (destructive)")
	rootCmd.AddCommand(downCmd)
}
