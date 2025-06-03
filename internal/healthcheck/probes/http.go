package probes

import (
	"context"
	"net/http"
	"time"
)

// HTTPProbe represents an HTTP health check probe
type HTTPProbe struct {
	client      *http.Client
	path        string
	expectedStatus int
}

// NewHTTPProbe creates a new HTTP probe
func NewHTTPProbe(timeout time.Duration, path string, expectedStatus int) *HTTPProbe {
	return &HTTPProbe{
		client: &http.Client{
			Timeout: timeout,
		},
		path: path,
		expectedStatus: expectedStatus,
	}
}

// Check performs an HTTP health check
func (p *HTTPProbe) Check(ctx context.Context, url string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url+p.path, nil)
	if err != nil {
		return false, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == p.expectedStatus, nil
}

// SetTimeout updates the probe's timeout
func (p *HTTPProbe) SetTimeout(timeout time.Duration) {
	p.client.Timeout = timeout
}

// SetPath updates the health check path
func (p *HTTPProbe) SetPath(path string) {
	p.path = path
}

// SetExpectedStatus updates the expected HTTP status code
func (p *HTTPProbe) SetExpectedStatus(status int) {
	p.expectedStatus = status
} 