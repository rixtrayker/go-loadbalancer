package config

import (
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Config represents the main configuration structure for the load balancer
type Config struct {
	Server   ServerConfig   `json:"server" koanf:"server"`
	Backends []BackendConfig `json:"backends" koanf:"backends"`
	Policies PolicyConfig   `json:"policies" koanf:"policies"`
}

// ServerConfig holds the main server configuration
type ServerConfig struct {
	Port            int    `json:"port" koanf:"port"`
	Host            string `json:"host" koanf:"host"`
	Algorithm       string `json:"algorithm" koanf:"algorithm"`
	ReadTimeout     int    `json:"readTimeout" koanf:"readTimeout"`
	WriteTimeout    int    `json:"writeTimeout" koanf:"writeTimeout"`
	IdleTimeout     int    `json:"idleTimeout" koanf:"idleTimeout"`
	ShutdownTimeout int    `json:"shutdownTimeout" koanf:"shutdownTimeout"`
}

// BackendConfig represents a backend server configuration
type BackendConfig struct {
	Name     string `json:"name" koanf:"name"`
	URL      string `json:"url" koanf:"url"`
	Weight   int    `json:"weight" koanf:"weight"`
	MaxConns int    `json:"maxConns" koanf:"maxConns"`
}

// PolicyConfig holds various policy configurations
type PolicyConfig struct {
	RateLimit  RateLimitConfig  `json:"rateLimit" koanf:"rateLimit"`
	Security   SecurityConfig   `json:"security" koanf:"security"`
	Transform  TransformConfig  `json:"transform" koanf:"transform"`
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled     bool `json:"enabled" koanf:"enabled"`
	RequestsPer int  `json:"requestsPer" koanf:"requestsPer"`
	Period      int  `json:"period" koanf:"period"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	AllowedIPs []string `json:"allowedIPs" koanf:"allowedIPs"`
}

// TransformConfig holds request/response transformation configuration
type TransformConfig struct {
	AddHeaders    map[string]string `json:"addHeaders" koanf:"addHeaders"`
	RemoveHeaders []string          `json:"removeHeaders" koanf:"removeHeaders"`
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	k := koanf.New(".")

	// Load from file if provided
	if configPath != "" {
		if err := k.Load(file.Provider(configPath), json.Parser()); err != nil {
			return nil, err
		}
	}

	// Load from environment variables
	// Prefix all env vars with LB_
	if err := k.Load(env.Provider("LB_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "LB_")), "_", ".", -1)
	}), nil); err != nil {
		return nil, err
	}

	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		return nil, err
	}

	return &config, nil
} 