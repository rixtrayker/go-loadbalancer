package configs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
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

	// Load default configuration
	config := DefaultConfig()

	// Read config file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply environment variables
	applyEnvironmentVariables(config)

	// Validate configuration
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Address:      ":8080",
			ReadTimeout:  30,
			WriteTimeout: 30,
			IdleTimeout:  60,
			CorsEnabled:  false,
		},
		Monitoring: MonitoringConfig{
			Prometheus: PrometheusConfig{
				Enabled: false,
				Path:    "/metrics",
				Port:    9090,
			},
			Tracing: TracingConfig{
				Enabled:        false,
				ServiceName:    "go-loadbalancer",
				ServiceVersion: "1.0.0",
				Environment:    "development",
				Endpoint:       "localhost:4317",
				SamplingRate:   1.0,
				Protocol:       "grpc",
				Secure:         false,
			},
			Logging: LoggingConfig{
				Level:          "info",
				Format:         "json",
				Output:         "stdout",
				IncludeTraceID: true,
				IncludeSpanID:  true,
			},
		},
	}
}

// applyEnvironmentVariables applies environment variables to the configuration
func applyEnvironmentVariables(config *Config) {
	// Server configuration
	if addr := os.Getenv("LB_SERVER_ADDRESS"); addr != "" {
		config.Server.Address = addr
	}

	// Logging configuration
	if level := os.Getenv("LB_LOG_LEVEL"); level != "" {
		config.Monitoring.Logging.Level = level
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
