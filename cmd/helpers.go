package cmd

import (
	"fmt"
	"os"

	"github.com/taqnihub/gomigrate/config"
	"github.com/taqnihub/gomigrate/internal/tui"
	"github.com/taqnihub/gomigrate/migrate"
)

// loadConfig loads the config with CLI flag overrides applied.
func loadConfig() (*config.Config, error) {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return nil, err
	}

	if driver != "" {
		cfg.Driver = driver
	}
	if migDir != "" {
		cfg.MigrationsDir = migDir
	}

	return cfg, nil
}

// newEngine loads config and creates an engine.
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

// exitWithError prints an error using lipgloss styling and exits.
func exitWithError(err error) {
	tui.Error("%v", err)
	os.Exit(1)
}

// These are now just aliases to the tui package for backward compatibility.
// Commands can use either these or tui.X directly.
func printSuccess(format string, args ...interface{}) {
	tui.Success(format, args...)
}

func printInfo(format string, args ...interface{}) {
	tui.Info(format, args...)
}

func printWarn(format string, args ...interface{}) {
	tui.Warning(format, args...)
}
