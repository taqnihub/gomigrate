package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initForce bool // --force flag to overwrite existing config

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gomigrate in the current directory",
	Long: `Create a .gomigrate.yml config file and migrations directory.

This sets up everything needed to start using gomigrate in a new project.`,
	Args: cobra.NoArgs,
	Run:  runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	configPath := ".gomigrate.yml"

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil && !initForce {
		exitWithError(fmt.Errorf("%s already exists (use --force to overwrite)", configPath))
	}

	// Write a default config
	defaultConfig := `# GoMigrate configuration
# Docs: https://github.com/taqnihub/gomigrate

# Database driver: mysql or postgres
driver: mysql

# Connection details
host: localhost
port: 3306
database: myapp
user: root
password: password

# Where migration files are stored
migrations_dir: ./db/migrations

# Optional: SSL mode (disable, require, verify-full)
ssl_mode: disable

# Optional: timezone
timezone: UTC

# Optional: how long to wait for migration lock
lock_timeout: 15s
`

	if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		exitWithError(fmt.Errorf("failed to write config: %w", err))
	}

	// Create migrations directory
	migrationsDir := "./db/migrations"
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		exitWithError(fmt.Errorf("failed to create migrations dir: %w", err))
	}

	absConfig, _ := filepath.Abs(configPath)
	absMigrations, _ := filepath.Abs(migrationsDir)

	printSuccess("GoMigrate initialized!")
	fmt.Printf("  Config:         %s\n", absConfig)
	fmt.Printf("  Migrations dir: %s\n\n", absMigrations)
	printInfo("Next steps:")
	fmt.Println("  1. Edit .gomigrate.yml with your database credentials")
	fmt.Println("  2. Run 'gomigrate create <name>' to create your first migration")
	fmt.Println("  3. Run 'gomigrate up' to apply it")
}

func init() {
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "overwrite existing config")
	rootCmd.AddCommand(initCmd)
}
