package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

	// Get current DB version
	currentVersion, dirty, err := engine.Version()
	if err != nil {
		exitWithError(fmt.Errorf("failed to get version: %w", err))
	}

	// Get all migration files from disk
	files, err := engine.List()
	if err != nil {
		exitWithError(fmt.Errorf("failed to list migrations: %w", err))
	}

	// Print header
	fmt.Printf("Migration Status — %s://%s:%d/%s\n\n", cfg.Driver, cfg.Host, cfg.Port, cfg.Database)

	if len(files) == 0 {
		printInfo("No migration files found in %s", cfg.MigrationsDir)
		printInfo("Run 'gomigrate create <name>' to create your first migration")
		return
	}

	// Print table header
	fmt.Printf("%-16s  %-40s  %-10s\n", "VERSION", "NAME", "STATUS")
	fmt.Printf("%-16s  %-40s  %-10s\n", "-------", "----", "------")

	applied := 0
	pending := 0

	for _, f := range files {
		status := "pending"
		if f.Version <= currentVersion {
			status = "applied"
			applied++
		} else {
			pending++
		}

		fmt.Printf("%-16d  %-40s  %-10s\n", f.Version, truncate(f.Name, 40), status)
	}

	fmt.Println()

	if dirty {
		printWarn("Database is in a DIRTY state at version %d", currentVersion)
		printWarn("This means a migration failed partway. Use 'gomigrate force <version>' to reset.")
	}

	printInfo("Total: %d | Applied: %d | Pending: %d", len(files), applied, pending)
}

// truncate shortens a string to max length, adding "..." if cut.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
