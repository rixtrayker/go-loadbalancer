package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Loader handles loading configuration from various sources
type Loader struct {
	configPath string
}

// NewLoader creates a new configuration loader
func NewLoader(configPath string) *Loader {
	return &Loader{
		configPath: configPath,
	}
}

// Load loads the configuration from the specified source
func (l *Loader) Load() (*Config, error) {
	// First try to load from file
	config, err := l.loadFromFile()
	if err != nil {
		return nil, fmt.Errorf("failed to load config from file: %w", err)
	}

	// Then override with environment variables if present
	l.overrideFromEnv(config)

	return config, nil
}

// loadFromFile loads configuration from a JSON file
func (l *Loader) loadFromFile() (*Config, error) {
	file, err := os.Open(l.configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// overrideFromEnv overrides configuration values with environment variables
func (l *Loader) overrideFromEnv(config *Config) {
	// TODO: Implement environment variable overrides
	// This would allow configuration to be overridden by environment variables
	// following a specific naming convention (e.g., LB_SERVER_PORT, LB_BACKEND_1_URL)
} 