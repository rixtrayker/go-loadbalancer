package configs

import (
	"time"
)

// Config represents the application configuration
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	BackendPools []BackendPoolConfig `yaml:"backend_pools"`
	RoutingRules []RoutingRuleConfig `yaml:"routing_rules"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
}

// ServerConfig contains server-specific configuration
type ServerConfig struct {
	Address      string `yaml:"address"`
	TLSCert      string `yaml:"tls_cert"`
	TLSKey       string `yaml:"tls_key"`
	AdminEnable  bool   `yaml:"admin_enable"`
	AdminPath    string `yaml:"admin_path"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	IdleTimeout  int    `yaml:"idle_timeout"`
	CorsEnabled  bool   `yaml:"cors_enabled"`
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

// MonitoringConfig represents monitoring configuration
type MonitoringConfig struct {
	Prometheus PrometheusConfig `yaml:"prometheus"`
	Tracing    TracingConfig    `yaml:"tracing"`
	Logging    LoggingConfig    `yaml:"logging"`
	Metrics    MetricsConfig    `yaml:"metrics"`
	Alerts     AlertsConfig     `yaml:"alerts"`
}

// PrometheusConfig contains Prometheus-specific configuration
type PrometheusConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
	Port    int    `yaml:"port"`
}

// TracingConfig contains OpenTelemetry tracing configuration
type TracingConfig struct {
	Enabled        bool    `yaml:"enabled"`
	ServiceName    string  `yaml:"service_name"`
	ServiceVersion string  `yaml:"service_version"`
	Environment    string  `yaml:"environment"`
	Endpoint       string  `yaml:"endpoint"`
	SamplingRate   float64 `yaml:"sampling_rate"`
	Protocol       string  `yaml:"protocol"`
	Secure         bool    `yaml:"secure"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level         string `yaml:"level"`
	Format        string `yaml:"format"`
	Output        string `yaml:"output"`
	IncludeTraceID bool   `yaml:"include_trace_id"`
	IncludeSpanID bool   `yaml:"include_span_id"`
}

// MetricsConfig contains metrics retention and aggregation settings
type MetricsConfig struct {
	RetentionPeriod     string `yaml:"retention_period"`
	AggregationInterval string `yaml:"aggregation_interval"`
	MaxSeries          int    `yaml:"max_series"`
}

// AlertsConfig contains alerting thresholds
type AlertsConfig struct {
	ErrorRateThreshold     float64 `yaml:"error_rate_threshold"`
	LatencyThreshold       float64 `yaml:"latency_threshold"`
	HealthCheckFailures    int     `yaml:"health_check_failures"`
	ResourceUsageThreshold float64 `yaml:"resource_usage_threshold"`
}
