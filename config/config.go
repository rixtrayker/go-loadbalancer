package config

// Config represents the main configuration structure for the load balancer
type Config struct {
	Server   ServerConfig   `json:"server"`
	Backends []BackendConfig `json:"backends"`
	Policies PolicyConfig   `json:"policies"`
}

// ServerConfig holds the main server configuration
type ServerConfig struct {
	Port            int    `json:"port"`
	Host            string `json:"host"`
	Algorithm       string `json:"algorithm"`
	ReadTimeout     int    `json:"readTimeout"`
	WriteTimeout    int    `json:"writeTimeout"`
	IdleTimeout     int    `json:"idleTimeout"`
	ShutdownTimeout int    `json:"shutdownTimeout"`
}

// BackendConfig represents a backend server configuration
type BackendConfig struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Weight   int    `json:"weight"`
	MaxConns int    `json:"maxConns"`
}

// PolicyConfig holds various policy configurations
type PolicyConfig struct {
	RateLimit  RateLimitConfig  `json:"rateLimit"`
	Security   SecurityConfig   `json:"security"`
	Transform  TransformConfig  `json:"transform"`
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled     bool `json:"enabled"`
	RequestsPer int  `json:"requestsPer"`
	Period      int  `json:"period"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	AllowedIPs []string `json:"allowedIPs"`
}

// TransformConfig holds request/response transformation configuration
type TransformConfig struct {
	AddHeaders    map[string]string `json:"addHeaders"`
	RemoveHeaders []string          `json:"removeHeaders"`
} 