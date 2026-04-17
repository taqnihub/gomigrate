package cmd

import (
	"fmt"
	"os"

	"github.com/taqnihub/gomigrate/config"
	"github.com/taqnihub/gomigrate/migrate"
)

// loadConfig loads the config with CLI flag overrides applied.
// Commands call this instead of config.Load() directly.
func loadConfig() (*config.Config, error) {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return nil, err
	}

	// Apply CLI flag overrides (flags take highest priority)
	if driver != "" {
		cfg.Driver = driver
	}
	if migDir != "" {
		cfg.MigrationsDir = migDir
	}
	// Note: --dsn flag handling would require restructuring DSN logic;
	// we'll skip it for now (power users can still use env vars)

	return cfg, nil
}

// newEngine is a convenience function that loads config and creates an engine.
// Commands call this at the start to get a ready-to-use engine.
func newEngine() (*migrate.Engine, *config.Config, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("config error: %w", err)
	}

	engine, err := migrate.New(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("engine error: %w", err)
	}

	return engine, cfg, nil
}

// exitWithError prints an error to stderr and exits with status 1.
// Use this for fatal errors in commands.
func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "✗ %v\n", err)
	os.Exit(1)
}

// printSuccess prints a success message with a green checkmark.
func printSuccess(format string, args ...interface{}) {
	fmt.Printf("✓ "+format+"\n", args...)
}

// printInfo prints an info message.
func printInfo(format string, args ...interface{}) {
	fmt.Printf("ℹ "+format+"\n", args...)
}

// printWarn prints a warning.
func printWarn(format string, args ...interface{}) {
	fmt.Printf("⚠ "+format+"\n", args...)
}
