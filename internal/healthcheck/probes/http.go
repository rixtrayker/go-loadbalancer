package probes

import (
	"net/http"
	"net/url"
	"time"
)

// HTTPProbe checks backend health using HTTP
type HTTPProbe struct {
	url     *url.URL
	path    string
	method  string
	timeout time.Duration
	client  *http.Client
}

// NewHTTPProbe creates a new HTTP health check probe
func NewHTTPProbe(url *url.URL, path, method string, timeout time.Duration) *HTTPProbe {
	if method == "" {
		method = http.MethodGet
	}

	if timeout == 0 {
		timeout = 5 * time.Second
	}

	return &HTTPProbe{
		url:     url,
		path:    path,
		method:  method,
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Check performs a health check
func (p *HTTPProbe) Check() bool {
	// Create request URL
	reqURL := *p.url
	reqURL.Path = p.path

	// Create request
	req, err := http.NewRequest(p.method, reqURL.String(), nil)
	if err != nil {
		return false
	}

	// Add health check header
	req.Header.Set("User-Agent", "Go-LoadBalancer-HealthCheck")

	// Send request
	resp, err := p.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Check status code
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
