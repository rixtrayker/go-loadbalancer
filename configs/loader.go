package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
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

	// Initialize koanf
	k := koanf.New(".")

	// Load default configuration
	defaultConfig := DefaultConfig()
	if err := k.Load(koanf.Provider(defaultConfig, ".", koanf.UnmarshalConf{}), nil); err != nil {
		return nil, fmt.Errorf("failed to load default config: %w", err)
	}

	// Load from YAML file
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	// Load from environment variables
	if err := k.Load(env.Provider("LB_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "LB_")), "_", ".", -1)
	}), nil); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	// Unmarshal into Config struct
	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// DefaultConfig returns the default configuration
func DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"server": map[string]interface{}{
			"address":       ":8080",
			"read_timeout":  30,
			"write_timeout": 30,
			"idle_timeout":  60,
			"cors_enabled":  false,
		},
		"monitoring": map[string]interface{}{
			"prometheus": map[string]interface{}{
				"enabled": false,
				"path":    "/metrics",
				"port":    9090,
			},
			"tracing": map[string]interface{}{
				"enabled":         false,
				"service_name":    "go-loadbalancer",
				"service_version": "1.0.0",
				"environment":     "development",
				"endpoint":        "localhost:4317",
				"sampling_rate":   1.0,
				"protocol":        "grpc",
				"secure":          false,
			},
			"logging": map[string]interface{}{
				"level":           "info",
				"format":          "json",
				"output":          "stdout",
				"include_trace_id": true,
				"include_span_id":  true,
			},
		},
	}
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// Validate server configuration
	if config.Server.Address == "" {
		return fmt.Errorf("server address is required")
	}

	// Validate backend pools
	if len(config.BackendPools) == 0 {
		return fmt.Errorf("at least one backend pool is required")
	}

	// Validate each backend pool
	poolNames := make(map[string]bool)
	for _, pool := range config.BackendPools {
		if pool.Name == "" {
			return fmt.Errorf("backend pool name is required")
		}

		if poolNames[pool.Name] {
			return fmt.Errorf("duplicate backend pool name: %s", pool.Name)
		}
		poolNames[pool.Name] = true

		if len(pool.Backends) == 0 {
			return fmt.Errorf("at least one backend is required in pool: %s", pool.Name)
		}

		for _, backend := range pool.Backends {
			if backend.URL == "" {
				return fmt.Errorf("backend URL is required in pool: %s", pool.Name)
			}
		}
	}

	// Validate routing rules
	if len(config.RoutingRules) == 0 {
		return fmt.Errorf("at least one routing rule is required")
	}

	// Validate each routing rule
	for _, rule := range config.RoutingRules {
		if rule.TargetPool == "" {
			return fmt.Errorf("target pool is required in routing rule")
		}

		if !poolNames[rule.TargetPool] {
			return fmt.Errorf("target pool does not exist: %s", rule.TargetPool)
		}
	}

	return nil
}
