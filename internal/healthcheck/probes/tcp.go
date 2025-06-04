package probes

import (
	"net"
	"net/url"
	"time"
)

// TCPProbe checks backend health using TCP
type TCPProbe struct {
	url     *url.URL
	timeout time.Duration
}

// NewTCPProbe creates a new TCP health check probe
func NewTCPProbe(url *url.URL, timeout time.Duration) *TCPProbe {
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	return &TCPProbe{
		url:     url,
		timeout: timeout,
	}
}

// Check performs a health check
func (p *TCPProbe) Check() bool {
	// Get host and port
	host := p.url.Host

	// Connect to host
	conn, err := net.DialTimeout("tcp", host, p.timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}
