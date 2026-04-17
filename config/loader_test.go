package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoad_FromFile verifies loading from a YAML file.
func TestLoad_FromFile(t *testing.T) {
	// Create a temporary config file
	dir := t.TempDir() // auto-cleaned up after test
	configPath := filepath.Join(dir, "test.yml")

	content := `driver: postgres
host: db.example.com
port: 5432
database: mydb
user: admin
password: s3cret
migrations_dir: /tmp/migrations
ssl_mode: require
`

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	// Load it
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify fields
	if cfg.Driver != "postgres" {
		t.Errorf("Driver: got %q, want %q", cfg.Driver, "postgres")
	}
	if cfg.Host != "db.example.com" {
		t.Errorf("Host: got %q, want %q", cfg.Host, "db.example.com")
	}
	if cfg.Port != 5432 {
		t.Errorf("Port: got %d, want %d", cfg.Port, 5432)
	}
	if cfg.Database != "mydb" {
		t.Errorf("Database: got %q, want %q", cfg.Database, "mydb")
	}
	if cfg.User != "admin" {
		t.Errorf("User: got %q, want %q", cfg.User, "admin")
	}
	if cfg.Password != "s3cret" {
		t.Errorf("Password: got %q, want %q", cfg.Password, "s3cret")
	}
	if cfg.SSLMode != "require" {
		t.Errorf("SSLMode: got %q, want %q", cfg.SSLMode, "require")
	}
}

// TestLoad_EnvOverride verifies env vars override file values.
func TestLoad_EnvOverride(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "test.yml")

	content := `driver: mysql
host: localhost
port: 3306
database: dev_db
user: root
password: default
migrations_dir: ./m
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	// Set env vars to override
	t.Setenv("GOMIGRATE_DATABASE", "production_db")
	t.Setenv("GOMIGRATE_PASSWORD", "production_secret")

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// These should come from env vars, not file
	if cfg.Database != "production_db" {
		t.Errorf("Database: env override failed. got %q, want %q", cfg.Database, "production_db")
	}
	if cfg.Password != "production_secret" {
		t.Errorf("Password: env override failed. got %q, want %q", cfg.Password, "production_secret")
	}

	// This should still come from the file (no env override)
	if cfg.User != "root" {
		t.Errorf("User: got %q, want %q", cfg.User, "root")
	}
}

// TestLoad_MissingFile returns an error when the explicit file doesn't exist.
func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/to/config.yml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

// TestLoad_InvalidConfig fails validation when required fields are missing.
func TestLoad_InvalidConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "test.yml")

	// Missing database name
	content := `driver: mysql
host: localhost
port: 3306
user: root
migrations_dir: ./m
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("expected validation error, got nil")
	}
}
