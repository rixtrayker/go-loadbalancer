package routing

import (
	"net/http"
	"strings"
)

// Rule represents a routing rule
type Rule struct {
	Host    string
	Path    string
	Methods []string
	Headers map[string]string
}

// NewRule creates a new routing rule
func NewRule(host, path string, methods []string, headers map[string]string) *Rule {
	return &Rule{
		Host:    host,
		Path:    path,
		Methods: methods,
		Headers: headers,
	}
}

// Matches checks if a request matches this rule
func (r *Rule) Matches(req *http.Request) bool {
	// Check host
	if r.Host != "" && !strings.EqualFold(req.Host, r.Host) {
		return false
	}

	// Check path
	if r.Path != "" && !strings.HasPrefix(req.URL.Path, r.Path) {
		return false
	}

	// Check methods
	if len(r.Methods) > 0 {
		methodMatch := false
		for _, method := range r.Methods {
			if req.Method == method {
				methodMatch = true
				break
			}
		}
		if !methodMatch {
			return false
		}
	}

	// Check headers
	for key, value := range r.Headers {
		if req.Header.Get(key) != value {
			return false
		}
	}

	return true
}

// WithHost sets the host for the rule
func (r *Rule) WithHost(host string) *Rule {
	r.Host = host
	return r
}

// WithPath sets the path for the rule
func (r *Rule) WithPath(path string) *Rule {
	r.Path = path
	return r
}

// WithMethods sets the allowed methods for the rule
func (r *Rule) WithMethods(methods []string) *Rule {
	r.Methods = methods
	return r
}

// WithHeader adds a header requirement to the rule
func (r *Rule) WithHeader(key, value string) *Rule {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[key] = value
	return r
} 