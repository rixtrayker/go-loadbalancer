package configs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads configuration from the specified file
func LoadConfig(path string) (*Config, error) {
	// Expand path if it contains ~
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		path = filepath.Join(homeDir, path[2:])
	}

	// Read configuration file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply environment variable overrides
	applyEnvironmentOverrides(&config)

	return &config, nil
}

// applyEnvironmentOverrides applies configuration overrides from environment variables
func applyEnvironmentOverrides(config *Config) {
	// Example: Override server address
	if addr := os.Getenv("LB_SERVER_ADDRESS"); addr != "" {
		config.Server.Address = addr
	}

	// Add more environment variable overrides as needed
}
