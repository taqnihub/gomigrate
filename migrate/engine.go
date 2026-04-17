package migrate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/taqnihub/gomigrate/config"
)

// Engine is the main migration runner.
// It wraps golang-migrate/migrate with a simpler API.
type Engine struct {
	config  *config.Config
	migrate *migrate.Migrate
}

// New creates a new migration engine from a config.
// It validates the config, builds the DSN, and initializes golang-migrate.
func New(cfg *config.Config) (*Engine, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Build the source URL (where migration files live)
	absPath, err := filepath.Abs(cfg.MigrationsDir)
	if err != nil {
		return nil, fmt.Errorf("invalid migrations_dir: %w", err)
	}

	// Auto-create the migrations directory if it doesn't exist.
	// This is user-friendly — no need to manually mkdir before first use.
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create migrations dir %q: %w", absPath, err)
	}

	sourceURL := fmt.Sprintf("file://%s", absPath)

	// Build the database DSN
	dsn, err := cfg.DSN()
	if err != nil {
		return nil, fmt.Errorf("failed to build DSN: %w", err)
	}

	// Initialize golang-migrate
	m, err := migrate.New(sourceURL, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrate: %w", err)
	}

	return &Engine{
		config:  cfg,
		migrate: m,
	}, nil
}

// Close releases all resources held by the engine.
// Always call this when done (use defer engine.Close()).
func (e *Engine) Close() error {
	sourceErr, dbErr := e.migrate.Close()
	if sourceErr != nil {
		return fmt.Errorf("failed to close source: %w", sourceErr)
	}
	if dbErr != nil {
		return fmt.Errorf("failed to close database: %w", dbErr)
	}
	return nil
}
