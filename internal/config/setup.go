package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	DatabaseType string `json:"database_type"` // "sqlite" or "turso"
}

// DBType constants
const (
	DBTypeSQLite = "sqlite"
	DBTypeTurso  = "turso"
)

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, "Alexandria", ".config", "config.json"), nil
}

// ensureConfigDir creates the config directory if it doesn't exist
func ensureConfigDir() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return nil
}

// Load reads the config file and returns the configuration
// If the file doesn't exist, it creates a default config with SQLite
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := &Config{DatabaseType: DBTypeSQLite}
		if err := Save(defaultConfig); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return defaultConfig, nil
	}

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Save writes the configuration to the config file
func Save(config *Config) error {
	if err := ensureConfigDir(); err != nil {
		return err
	}

	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetCurrentDBType returns the currently configured database type
func GetCurrentDBType() (string, error) {
	config, err := Load()
	if err != nil {
		return "", err
	}
	return config.DatabaseType, nil
}

// SwitchDB changes the database type to the specified type
func SwitchDB(dbType string) error {
	if dbType != DBTypeSQLite && dbType != DBTypeTurso {
		return fmt.Errorf("invalid database type: %s (must be %s or %s)",
			dbType, DBTypeSQLite, DBTypeTurso)
	}

	config := &Config{DatabaseType: dbType}
	if err := Save(config); err != nil {
		return fmt.Errorf("failed to switch database: %w", err)
	}

	return nil
}
