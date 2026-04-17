package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/taqnihub/gomigrate/internal/tui"
)

var (
	initForce      bool
	initDriver     string
	initHost       string
	initPort       string
	initDatabase   string
	initUser       string
	initPassword   string
	initMigrations string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gomigrate in the current directory",
	Long: `Interactively create a .gomigrate.yml config file and migrations directory.

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

	tui.Title("🚀 Welcome to GoMigrate")
	tui.Muted("Let's set up your project. Press Enter to accept defaults in [brackets].")
	tui.Newline()

	// Set defaults that will change based on driver selection
	initDriver = "mysql"
	initPort = "3306"
	initUser = "root"
	initMigrations = "./db/migrations"

	// Step 1: Choose driver
	driverForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Database driver").
				Options(
					huh.NewOption("MySQL", "mysql"),
					huh.NewOption("PostgreSQL", "postgres"),
				).
				Value(&initDriver),
		),
	)

	if err := driverForm.Run(); err != nil {
		exitWithError(fmt.Errorf("init cancelled: %w", err))
	}

	// Adjust default port based on driver
	if initDriver == "postgres" {
		initPort = "5432"
		initUser = "postgres"
	}

	// Step 2: Collect connection details
	connectionForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Host").
				Placeholder("localhost").
				Value(&initHost),

			huh.NewInput().
				Title("Port").
				Value(&initPort).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("port is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("Database name").
				Placeholder("myapp").
				Value(&initDatabase).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("database name is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("Username").
				Value(&initUser).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("username is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Value(&initPassword),

			huh.NewInput().
				Title("Migrations directory").
				Value(&initMigrations),
		),
	)

	if err := connectionForm.Run(); err != nil {
		exitWithError(fmt.Errorf("init cancelled: %w", err))
	}

	// Apply default for host if empty
	if initHost == "" {
		initHost = "localhost"
	}

	// Build the YAML config
	configContent := fmt.Sprintf(`# GoMigrate configuration
# Docs: https://github.com/taqnihub/gomigrate

driver: %s
host: %s
port: %s
database: %s
user: %s
password: %s
migrations_dir: %s

# Optional: SSL mode (disable, require, verify-full)
ssl_mode: disable

# Optional: timezone
timezone: UTC

# Optional: how long to wait for migration lock
lock_timeout: 15s
`,
		initDriver, initHost, initPort, initDatabase,
		initUser, initPassword, initMigrations,
	)

	// Write the config
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		exitWithError(fmt.Errorf("failed to write config: %w", err))
	}

	// Create migrations directory
	if err := os.MkdirAll(initMigrations, 0755); err != nil {
		exitWithError(fmt.Errorf("failed to create migrations dir: %w", err))
	}

	// Success message
	absConfig, _ := filepath.Abs(configPath)
	absMigrations, _ := filepath.Abs(initMigrations)

	tui.Newline()
	tui.Success("GoMigrate initialized!")
	tui.Newline()
	tui.KeyValue("Config", tui.Path(absConfig))
	tui.KeyValue("Migrations", tui.Path(absMigrations))
	tui.KeyValue("Database", fmt.Sprintf("%s on %s:%s/%s",
		initDriver, initHost, initPort, initDatabase))

	tui.Section("Next Steps")
	fmt.Printf("  %s Create your first migration:\n", tui.Dim("1."))
	fmt.Printf("     %s\n", tui.Code("gomigrate create add_users_table"))
	fmt.Printf("  %s Edit the generated SQL files in %s\n",
		tui.Dim("2."), tui.Path(initMigrations))
	fmt.Printf("  %s Apply it:\n", tui.Dim("3."))
	fmt.Printf("     %s\n", tui.Code("gomigrate up"))
	tui.Newline()
}

func init() {
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "overwrite existing config")
	rootCmd.AddCommand(initCmd)
}
