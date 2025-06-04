package transform

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Apply applies transformation rules to a request
func Apply(transformStr string, r *http.Request) error {
	// Parse transform string (e.g., "add-header:X-Forwarded-Host:example.com,remove-header:Referer")
	rules := strings.Split(transformStr, ",")

	for _, rule := range rules {
		parts := strings.SplitN(rule, ":", 3)
		if len(parts) < 2 {
			continue
		}

		action := strings.ToLower(parts[0])
		
		switch action {
		case "add-header":
			if len(parts) != 3 {
				return errors.New("invalid add-header format")
			}
			r.Header.Add(parts[1], parts[2])
			
		case "set-header":
			if len(parts) != 3 {
				return errors.New("invalid set-header format")
			}
			r.Header.Set(parts[1], parts[2])
			
		case "remove-header":
			r.Header.Del(parts[1])
			
		case "rewrite-path":
			if len(parts) != 3 {
				return errors.New("invalid rewrite-path format")
			}
			r.URL.Path = strings.Replace(r.URL.Path, parts[1], parts[2], 1)
			
		case "add-query":
			if len(parts) != 3 {
				return errors.New("invalid add-query format")
			}
			q := r.URL.Query()
			q.Add(parts[1], parts[2])
			r.URL.RawQuery = q.Encode()
			
		default:
			return errors.New("unknown transform action: " + action)
		}
	}

	return nil
}

// HeaderTransformer transforms HTTP headers
type HeaderTransformer struct {
	AddHeaders    map[string]string `json:"add_headers"`
	RemoveHeaders []string          `json:"remove_headers"`
}

// NewHeaderTransformer creates a new header transformer from JSON
func NewHeaderTransformer(jsonStr string) (*HeaderTransformer, error) {
	var ht HeaderTransformer
	if err := json.Unmarshal([]byte(jsonStr), &ht); err != nil {
		return nil, err
	}
	return &ht, nil
}

// Apply applies header transformations to a request
func (ht *HeaderTransformer) Apply(r *http.Request) {
	// Add headers
	for name, value := range ht.AddHeaders {
		r.Header.Add(name, value)
	}

	// Remove headers
	for _, name := range ht.RemoveHeaders {
		r.Header.Del(name)
	}
}
