package configs

import (
	"time"
)

// Config represents the application configuration
type Config struct {
	Server       ServerConfig      `yaml:"server"`
	BackendPools []BackendPoolConfig `yaml:"backend_pools"`
	RoutingRules []RoutingRuleConfig `yaml:"routing_rules"`
}

// ServerConfig contains server-specific configuration
type ServerConfig struct {
	Address     string `yaml:"address"`
	TLSCert     string `yaml:"tls_cert"`
	TLSKey      string `yaml:"tls_key"`
	AdminEnable bool   `yaml:"admin_enable"`
	AdminPath   string `yaml:"admin_path"`
}

// BackendPoolConfig represents a group of backend servers
type BackendPoolConfig struct {
	Name        string           `yaml:"name"`
	Algorithm   string           `yaml:"algorithm"`
	Backends    []BackendConfig  `yaml:"backends"`
	HealthCheck HealthCheckConfig `yaml:"health_check"`
}

// BackendConfig represents a single backend server
type BackendConfig struct {
	URL    string `yaml:"url"`
	Weight int    `yaml:"weight"`
}

// HealthCheckConfig defines health check parameters
type HealthCheckConfig struct {
	Path     string        `yaml:"path"`
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
	Method   string        `yaml:"method"`
}

// RoutingRuleConfig defines how requests are routed
type RoutingRuleConfig struct {
	Match      MatchConfig      `yaml:"match"`
	TargetPool string           `yaml:"target_pool"`
	Policies   []PolicyConfig   `yaml:"policies"`
}

// MatchConfig defines criteria for matching requests
type MatchConfig struct {
	Host    string            `yaml:"host"`
	Path    string            `yaml:"path"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
}

// PolicyConfig defines policies to apply to matched requests
type PolicyConfig struct {
	RateLimit string `yaml:"rate_limit"`
	Transform string `yaml:"transform"`
	ACL       string `yaml:"acl"`
}
