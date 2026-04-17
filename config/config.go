package config

import (
	"fmt"
	"strings"
	"time"
)

// Config holds all configuration for gomigrate.
// It can be loaded from a YAML file, environment variables, or CLI flags.
type Config struct {
	// Database connection
	Driver   string `mapstructure:"driver" yaml:"driver"`
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Database string `mapstructure:"database" yaml:"database"`
	User     string `mapstructure:"user" yaml:"user"`
	Password string `mapstructure:"password" yaml:"password"`

	// Migration files
	MigrationsDir string `mapstructure:"migrations_dir" yaml:"migrations_dir"`

	// Optional settings
	SSLMode     string        `mapstructure:"ssl_mode" yaml:"ssl_mode"`
	Timezone    string        `mapstructure:"timezone" yaml:"timezone"`
	LockTimeout time.Duration `mapstructure:"lock_timeout" yaml:"lock_timeout"`
}

// DSN builds a database connection string from the config fields.
// Returns a driver-specific connection URL that golang-migrate understands.
func (c *Config) DSN() (string, error) {
	switch strings.ToLower(c.Driver) {
	case "mysql":
		// Format: mysql://user:password@tcp(host:port)/dbname
		return fmt.Sprintf(
			"mysql://%s:%s@tcp(%s:%d)/%s",
			c.User, c.Password, c.Host, c.Port, c.Database,
		), nil

	case "postgres", "postgresql":
		// Format: postgres://user:password@host:port/dbname?sslmode=disable
		sslMode := c.SSLMode
		if sslMode == "" {
			sslMode = "disable"
		}
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			c.User, c.Password, c.Host, c.Port, c.Database, sslMode,
		), nil

	default:
		return "", fmt.Errorf("unsupported driver: %q (supported: mysql, postgres)", c.Driver)
	}
}

// Validate checks if the config has all required fields and valid values.
// Returns an error with a helpful message if something is wrong.
func (c *Config) Validate() error {
	if c.Driver == "" {
		return fmt.Errorf("driver is required (mysql or postgres)")
	}

	driver := strings.ToLower(c.Driver)
	if driver != "mysql" && driver != "postgres" && driver != "postgresql" {
		return fmt.Errorf("unsupported driver %q: must be mysql or postgres", c.Driver)
	}

	if c.Host == "" {
		return fmt.Errorf("host is required")
	}

	if c.Port == 0 {
		return fmt.Errorf("port is required")
	}

	if c.Database == "" {
		return fmt.Errorf("database name is required")
	}

	if c.User == "" {
		return fmt.Errorf("user is required")
	}

	if c.MigrationsDir == "" {
		return fmt.Errorf("migrations_dir is required")
	}

	return nil
}

// Default returns a Config with sensible defaults.
// These are used as fallbacks when values aren't provided.
func Default() *Config {
	return &Config{
		Driver:        "mysql",
		Host:          "localhost",
		Port:          3306,
		MigrationsDir: "./db/migrations",
		SSLMode:       "disable",
		Timezone:      "UTC",
		LockTimeout:   15 * time.Second,
	}
}
