package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
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

	Args: cobra.MaximumNArgs(1), // accept 0 or 1 argument
	Run:  runUp,
}

func runUp(cmd *cobra.Command, args []string) {
	engine, cfg, err := newEngine()
	if err != nil {
		exitWithError(err)
	}
	defer engine.Close()

	printInfo("Connected to %s://%s:%d/%s", cfg.Driver, cfg.Host, cfg.Port, cfg.Database)

	start := time.Now()

	// Determine how many migrations to apply
	if len(args) == 0 {
		// No argument: apply ALL pending
		err = engine.Up()
	} else {
		// Argument given: apply N
		n, parseErr := strconv.Atoi(args[0])
		if parseErr != nil || n <= 0 {
			exitWithError(fmt.Errorf("invalid number %q: must be a positive integer", args[0]))
		}
		err = engine.UpN(n)
	}

	// Handle the result
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			printInfo("No pending migrations — database is up to date")
			return
		}
		exitWithError(fmt.Errorf("migration failed: %w", err))
	}

	// Show new version
	version, dirty, _ := engine.Version()
	elapsed := time.Since(start).Round(time.Millisecond)

	if dirty {
		printWarn("Migration succeeded but state is dirty at version %d", version)
	} else {
		printSuccess("Migrations applied — now at version %d (%s)", version, elapsed)
	}
}

func init() {
	rootCmd.AddCommand(upCmd)
}
