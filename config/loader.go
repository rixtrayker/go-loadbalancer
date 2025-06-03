package config

import (
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
	// First try to get config path from environment
	configPath := os.Getenv("LB_CONFIG_PATH")
	if configPath == "" {
		configPath = l.configPath
	}

	// Load configuration using koanf
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return config, nil
} 