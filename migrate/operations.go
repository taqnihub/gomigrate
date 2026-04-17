package migrate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
)

// ErrNoChange is returned when there's nothing to migrate.
var ErrNoChange = migrate.ErrNoChange

// Up applies all pending migrations.
// Returns ErrNoChange if there's nothing to apply.
func (e *Engine) Up() error {
	if err := e.migrate.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return ErrNoChange
		}
		return fmt.Errorf("migration up failed: %w", err)
	}
	return nil
}

// UpN applies the next N pending migrations.
// Use UpN(1) to apply only the next one.
func (e *Engine) UpN(n int) error {
	if n <= 0 {
		return fmt.Errorf("n must be positive, got %d", n)
	}
	if err := e.migrate.Steps(n); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return ErrNoChange
		}
		return fmt.Errorf("migration up %d failed: %w", n, err)
	}
	return nil
}

// Down reverts the last N applied migrations.
// Use Down(1) to revert only the most recent one (safe default).
func (e *Engine) Down(n int) error {
	if n <= 0 {
		return fmt.Errorf("n must be positive, got %d", n)
	}
	// Steps with negative number means "go down N steps"
	if err := e.migrate.Steps(-n); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return ErrNoChange
		}
		return fmt.Errorf("migration down %d failed: %w", n, err)
	}
	return nil
}

// DownAll reverts ALL applied migrations.
// This is destructive — use with caution!
func (e *Engine) DownAll() error {
	if err := e.migrate.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return ErrNoChange
		}
		return fmt.Errorf("migration down all failed: %w", err)
	}
	return nil
}

// Force sets the migration version without running any migrations.
// Useful for fixing a "dirty" state when a migration failed mid-way.
// Passing -1 resets the version to nothing.
func (e *Engine) Force(version int) error {
	if err := e.migrate.Force(version); err != nil {
		return fmt.Errorf("force version %d failed: %w", version, err)
	}
	return nil
}

// Version returns the current migration version and whether the state is dirty.
// A dirty state means a migration failed partway and needs manual intervention.
// Returns (0, false, nil) if no migrations have been applied yet.
func (e *Engine) Version() (version uint, dirty bool, err error) {
	v, d, err := e.migrate.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			// No migrations applied yet — not an error
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get version: %w", err)
	}
	return v, d, nil
}

// MigrationFile represents a single migration on disk.
type MigrationFile struct {
	Version  uint
	Name     string
	UpPath   string
	DownPath string
}

// List returns all migration files found in the migrations directory.
// Does NOT query the database — just reads files from disk.
func (e *Engine) List() ([]MigrationFile, error) {
	dir := e.config.MigrationsDir

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations dir: %w", err)
	}

	// Migrations come in pairs: NNNNNN_name.up.sql and NNNNNN_name.down.sql
	// We group them by version.
	files := make(map[uint]*MigrationFile)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		version, migName, direction, ok := parseFilename(name)
		if !ok {
			continue // skip unrecognized files
		}

		if _, exists := files[version]; !exists {
			files[version] = &MigrationFile{
				Version: version,
				Name:    migName,
			}
		}

		fullPath := filepath.Join(dir, name)
		if direction == "up" {
			files[version].UpPath = fullPath
		} else if direction == "down" {
			files[version].DownPath = fullPath
		}
	}

	// Sort by version
	result := make([]MigrationFile, 0, len(files))
	for _, f := range files {
		result = append(result, *f)
	}

	// Simple sort by version (ascending)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Version > result[j].Version {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result, nil
}

// Create creates a new migration file pair (up and down) with the given name.
// Returns an error if a migration with the same name already exists.
// Returns the paths of the created files.
func (e *Engine) Create(name string) (upPath, downPath string, err error) {
	if name == "" {
		return "", "", fmt.Errorf("migration name is required")
	}

	// Ensure migrations dir exists
	if err := os.MkdirAll(e.config.MigrationsDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create migrations dir: %w", err)
	}

	// Sanitize name
	sanitized := sanitizeName(name)
	if sanitized == "" {
		return "", "", fmt.Errorf("migration name %q produced empty result after sanitization", name)
	}

	// Check for duplicate names
	existing, err := e.findMigrationsByName(sanitized)
	if err != nil {
		return "", "", fmt.Errorf("failed to check existing migrations: %w", err)
	}
	if len(existing) > 0 {
		return "", "", fmt.Errorf(
			"a migration named %q already exists (version %d)\n\nChoose a different name or delete the existing files:\n  - %s\n  - %s",
			sanitized,
			existing[0].Version,
			existing[0].UpPath,
			existing[0].DownPath,
		)
	}

	// Generate timestamp-based version
	version := time.Now().UTC().Format("20060102150405")

	upName := fmt.Sprintf("%s_%s.up.sql", version, sanitized)
	downName := fmt.Sprintf("%s_%s.down.sql", version, sanitized)

	upPath = filepath.Join(e.config.MigrationsDir, upName)
	downPath = filepath.Join(e.config.MigrationsDir, downName)

	// Create files with helpful header comments
	upContent := fmt.Sprintf("-- Migration: %s\n-- Created: %s\n\n-- Write your UP migration SQL here\n\n",
		sanitized, time.Now().Format(time.RFC3339))
	downContent := fmt.Sprintf("-- Migration: %s (rollback)\n-- Created: %s\n\n-- Write your DOWN migration SQL here\n\n",
		sanitized, time.Now().Format(time.RFC3339))

	if err := os.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		return "", "", fmt.Errorf("failed to write up migration: %w", err)
	}
	if err := os.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		os.Remove(upPath) // cleanup
		return "", "", fmt.Errorf("failed to write down migration: %w", err)
	}

	return upPath, downPath, nil
}

// findMigrationsByName returns existing migrations that match the given name.
// This is case-insensitive matching against the sanitized name.
func (e *Engine) findMigrationsByName(name string) ([]MigrationFile, error) {
	files, err := e.List()
	if err != nil {
		return nil, err
	}

	var matches []MigrationFile
	for _, f := range files {
		if strings.EqualFold(f.Name, name) {
			matches = append(matches, f)
		}
	}
	return matches, nil
}
