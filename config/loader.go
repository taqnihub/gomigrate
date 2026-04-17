package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Load reads configuration from (in priority order, highest first):
//  1. CLI flags (handled by cobra, not here)
//  2. Environment variables (prefix: GOMIGRATE_)
//  3. .gomigrate.yml in current directory
//  4. .gomigrate.yml in home directory
//  5. Default values
//
// Returns a validated Config ready to use.
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Start with defaults
	setDefaults(v)

	// Enable environment variable support
	// e.g., GOMIGRATE_PASSWORD=secret -> config.Password = "secret"
	v.SetEnvPrefix("GOMIGRATE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Configure where to look for config files
	if configPath != "" {
		// Explicit path provided via --config flag
		v.SetConfigFile(configPath)
	} else {
		// Default search paths
		v.SetConfigName(".gomigrate")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")     // current directory
		v.AddConfigPath("$HOME") // home directory

		// Also check for .yml extension (not just .yaml)
		if home, err := os.UserHomeDir(); err == nil {
			v.AddConfigPath(home)
		}
	}

	// Read the config file (if it exists)
	if err := v.ReadInConfig(); err != nil {
		// It's OK if the file doesn't exist — we can use env vars + defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal into our Config struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// setDefaults registers default values with Viper.
// These are used when a value isn't set in YAML, env, or flags.
func setDefaults(v *viper.Viper) {
	defaults := Default()

	v.SetDefault("driver", defaults.Driver)
	v.SetDefault("host", defaults.Host)
	v.SetDefault("port", defaults.Port)
	v.SetDefault("migrations_dir", defaults.MigrationsDir)
	v.SetDefault("ssl_mode", defaults.SSLMode)
	v.SetDefault("timezone", defaults.Timezone)
	v.SetDefault("lock_timeout", defaults.LockTimeout)
}

// ConfigFilePath returns the path to the active config file if one was found.
// Useful for debugging or displaying in error messages.
func ConfigFilePath() string {
	if path := viper.ConfigFileUsed(); path != "" {
		abs, err := filepath.Abs(path)
		if err == nil {
			return abs
		}
		return path
	}
	return ""
}
