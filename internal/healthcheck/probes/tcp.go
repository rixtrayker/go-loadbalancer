package probes

import (
	"context"
	"net"
	"time"
)

// TCPProbe represents a TCP health check probe
type TCPProbe struct {
	timeout time.Duration
}

// NewTCPProbe creates a new TCP probe
func NewTCPProbe(timeout time.Duration) *TCPProbe {
	return &TCPProbe{
		timeout: timeout,
	}
}

// Check performs a TCP health check
func (p *TCPProbe) Check(ctx context.Context, url string) (bool, error) {
	dialer := net.Dialer{
		Timeout: p.timeout,
	}

	conn, err := dialer.DialContext(ctx, "tcp", url)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	return true, nil
}

// SetTimeout updates the probe's timeout
func (p *TCPProbe) SetTimeout(timeout time.Duration) {
	p.timeout = timeout
} 