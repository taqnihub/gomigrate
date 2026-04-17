package config

import (
	"strings"
	"testing"
	"time"
)

// TestDefault verifies default values are sensible.
func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Driver != "mysql" {
		t.Errorf("Driver: got %q, want %q", cfg.Driver, "mysql")
	}
	if cfg.Host != "localhost" {
		t.Errorf("Host: got %q, want %q", cfg.Host, "localhost")
	}
	if cfg.Port != 3306 {
		t.Errorf("Port: got %d, want %d", cfg.Port, 3306)
	}
	if cfg.MigrationsDir != "./db/migrations" {
		t.Errorf("MigrationsDir: got %q, want %q", cfg.MigrationsDir, "./db/migrations")
	}
	if cfg.LockTimeout != 15*time.Second {
		t.Errorf("LockTimeout: got %v, want %v", cfg.LockTimeout, 15*time.Second)
	}
}

// TestValidate exercises the Validate() method with many cases.
// This is a "table-driven test" — idiomatic Go style for testing many inputs.
func TestValidate(t *testing.T) {
	tests := []struct {
		name    string // test case name (shown in output)
		config  Config // the config to validate
		wantErr bool   // should Validate() return an error?
		errMsg  string // if wantErr=true, the error should contain this string
	}{
		{
			name: "valid mysql config",
			config: Config{
				Driver:        "mysql",
				Host:          "localhost",
				Port:          3306,
				Database:      "test",
				User:          "root",
				MigrationsDir: "./migrations",
			},
			wantErr: false,
		},
		{
			name: "valid postgres config",
			config: Config{
				Driver:        "postgres",
				Host:          "localhost",
				Port:          5432,
				Database:      "test",
				User:          "postgres",
				MigrationsDir: "./migrations",
			},
			wantErr: false,
		},
		{
			name:    "empty driver",
			config:  Config{Host: "localhost", Port: 3306, Database: "test", User: "root", MigrationsDir: "./m"},
			wantErr: true,
			errMsg:  "driver is required",
		},
		{
			name: "unsupported driver",
			config: Config{
				Driver:        "mongodb",
				Host:          "localhost",
				Port:          27017,
				Database:      "test",
				User:          "root",
				MigrationsDir: "./m",
			},
			wantErr: true,
			errMsg:  "unsupported driver",
		},
		{
			name:    "missing host",
			config:  Config{Driver: "mysql", Port: 3306, Database: "test", User: "root", MigrationsDir: "./m"},
			wantErr: true,
			errMsg:  "host is required",
		},
		{
			name:    "missing port",
			config:  Config{Driver: "mysql", Host: "localhost", Database: "test", User: "root", MigrationsDir: "./m"},
			wantErr: true,
			errMsg:  "port is required",
		},
		{
			name:    "missing database",
			config:  Config{Driver: "mysql", Host: "localhost", Port: 3306, User: "root", MigrationsDir: "./m"},
			wantErr: true,
			errMsg:  "database name is required",
		},
		{
			name:    "missing user",
			config:  Config{Driver: "mysql", Host: "localhost", Port: 3306, Database: "test", MigrationsDir: "./m"},
			wantErr: true,
			errMsg:  "user is required",
		},
		{
			name:    "missing migrations_dir",
			config:  Config{Driver: "mysql", Host: "localhost", Port: 3306, Database: "test", User: "root"},
			wantErr: true,
			errMsg:  "migrations_dir is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()

			if tc.wantErr && err == nil {
				t.Errorf("expected error but got nil")
				return
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if tc.wantErr && !strings.Contains(err.Error(), tc.errMsg) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.errMsg)
			}
		})
	}
}

// TestDSN verifies DSN generation for each driver.
func TestDSN(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		want    string
		wantErr bool
	}{
		{
			name: "mysql dsn",
			config: Config{
				Driver:   "mysql",
				Host:     "localhost",
				Port:     3306,
				Database: "testdb",
				User:     "root",
				Password: "secret",
			},
			want: "mysql://root:secret@tcp(localhost:3306)/testdb",
		},
		{
			name: "mysql case insensitive",
			config: Config{
				Driver:   "MySQL",
				Host:     "localhost",
				Port:     3306,
				Database: "testdb",
				User:     "root",
				Password: "secret",
			},
			want: "mysql://root:secret@tcp(localhost:3306)/testdb",
		},
		{
			name: "postgres dsn with default ssl",
			config: Config{
				Driver:   "postgres",
				Host:     "localhost",
				Port:     5432,
				Database: "testdb",
				User:     "postgres",
				Password: "secret",
			},
			want: "postgres://postgres:secret@localhost:5432/testdb?sslmode=disable",
		},
		{
			name: "postgres dsn with require ssl",
			config: Config{
				Driver:   "postgres",
				Host:     "prod.example.com",
				Port:     5432,
				Database: "testdb",
				User:     "postgres",
				Password: "secret",
				SSLMode:  "require",
			},
			want: "postgres://postgres:secret@prod.example.com:5432/testdb?sslmode=require",
		},
		{
			name:    "unsupported driver",
			config:  Config{Driver: "mongodb"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.config.DSN()

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if got != tc.want {
				t.Errorf("DSN:\n  got:  %q\n  want: %q", got, tc.want)
			}
		})
	}
}
