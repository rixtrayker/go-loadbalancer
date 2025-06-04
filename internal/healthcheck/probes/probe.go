package probes

// Probe defines the interface for health check probes
type Probe interface {
	// Check performs a health check and returns true if healthy
	Check() bool
}
