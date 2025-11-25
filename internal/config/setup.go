package config

import (
	"alexandria/internal/logger"
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
		logger.Log.Error("failed to get config path", "error", err)
		return err
	}

	configDir := filepath.Dir(configPath)
	logger.Log.Debug("ensuring config directory exists", "path", configDir)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		logger.Log.Error("failed to create config directory", "error", err, "path", configDir)
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return nil
}

// Load reads the config file and returns the configuration
// If the file doesn't exist, it creates a default config with SQLite
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		logger.Log.Error("failed to get config path", "error", err)
		return nil, err
	}

	logger.Log.Debug("loading config", "path", configPath)

	// If config doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logger.Log.Debug("config file not found, creating default", "path", configPath)
		defaultConfig := &Config{DatabaseType: DBTypeSQLite}
		if err := Save(defaultConfig); err != nil {
			logger.Log.Error("failed to create default config", "error", err)
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		logger.Log.Info("created default config", "database_type", DBTypeSQLite)
		return defaultConfig, nil
	}

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		logger.Log.Error("failed to read config file", "error", err, "path", configPath)
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		logger.Log.Error("failed to parse config file", "error", err)
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	logger.Log.Debug("config loaded", "database_type", config.DatabaseType)
	return &config, nil
}

// Save writes the configuration to the config file
func Save(config *Config) error {
	if err := ensureConfigDir(); err != nil {
		return err
	}

	configPath, err := getConfigPath()
	if err != nil {
		logger.Log.Error("failed to get config path", "error", err)
		return err
	}

	logger.Log.Debug("saving config", "database_type", config.DatabaseType, "path", configPath)

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		logger.Log.Error("failed to marshal config", "error", err)
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		logger.Log.Error("failed to write config file", "error", err, "path", configPath)
		return fmt.Errorf("failed to write config file: %w", err)
	}

	logger.Log.Debug("config saved successfully", "path", configPath)
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
	logger.Log.Debug("switching database type", "from", "current", "to", dbType)

	if dbType != DBTypeSQLite && dbType != DBTypeTurso {
		logger.Log.Error("invalid database type requested", "type", dbType)
		return fmt.Errorf("invalid database type: %s (must be %s or %s)",
			dbType, DBTypeSQLite, DBTypeTurso)
	}

	config := &Config{DatabaseType: dbType}
	if err := Save(config); err != nil {
		logger.Log.Error("failed to save config during database switch", "error", err, "type", dbType)
		return fmt.Errorf("failed to switch database: %w", err)
	}

	logger.Log.Info("database type switched", "database_type", dbType)
	return nil
}
