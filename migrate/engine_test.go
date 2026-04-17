package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/taqnihub/gomigrate/config"
)

// TestNew_InvalidConfig verifies that invalid configs are rejected.
func TestNew_InvalidConfig(t *testing.T) {
	cfg := &config.Config{
		// Missing required fields
		Driver: "mysql",
	}

	_, err := New(cfg)
	if err == nil {
		t.Error("expected error for invalid config, got nil")
	}
	if !strings.Contains(err.Error(), "invalid config") {
		t.Errorf("error should mention 'invalid config': %v", err)
	}
}

// TestNew_UnsupportedDriver verifies we reject unknown drivers.
func TestNew_UnsupportedDriver(t *testing.T) {
	cfg := &config.Config{
		Driver:        "mongodb",
		Host:          "localhost",
		Port:          27017,
		Database:      "test",
		User:          "root",
		MigrationsDir: t.TempDir(),
	}

	_, err := New(cfg)
	if err == nil {
		t.Error("expected error for unsupported driver, got nil")
	}
}

// TestNew_CreatesMigrationsDir verifies the directory is auto-created.
func TestNew_CreatesMigrationsDir(t *testing.T) {
	// Use a path that doesn't exist yet
	dir := filepath.Join(t.TempDir(), "nested", "path", "migrations")

	cfg := &config.Config{
		Driver:        "mysql",
		Host:          "localhost",
		Port:          3306,
		Database:      "test",
		User:          "root",
		Password:      "wrong", // intentionally wrong
		MigrationsDir: dir,
	}

	// This will fail to connect (no real DB), but it should create the dir first
	_, _ = New(cfg)

	// Verify the directory was created
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("migrations dir was not auto-created: %s", dir)
	}
}
