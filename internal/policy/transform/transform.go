package transform

import (
	"net/http"

	"github.com/amr/go-loadbalancer/internal/policy"
)

// Transformer implements request/response transformation policy
type Transformer struct {
	policy.BasePolicy
	headerTransformations map[string]string
	queryTransformations  map[string]string
}

// NewTransformer creates a new transformer
func NewTransformer() *Transformer {
	return &Transformer{
		BasePolicy: policy.BasePolicy{
			name: "transformer",
		},
		headerTransformations: make(map[string]string),
		queryTransformations:  make(map[string]string),
	}
}

// Apply implements the Policy interface
func (t *Transformer) Apply(req *http.Request, resp *http.Response) error {
	// Transform headers
	for key, value := range t.headerTransformations {
		req.Header.Set(key, value)
	}

	// Transform query parameters
	q := req.URL.Query()
	for key, value := range t.queryTransformations {
		q.Set(key, value)
	}
	req.URL.RawQuery = q.Encode()

	return nil
}

// AddHeaderTransformation adds a header transformation
func (t *Transformer) AddHeaderTransformation(key, value string) {
	t.headerTransformations[key] = value
}

// RemoveHeaderTransformation removes a header transformation
func (t *Transformer) RemoveHeaderTransformation(key string) {
	delete(t.headerTransformations, key)
}

// AddQueryTransformation adds a query parameter transformation
func (t *Transformer) AddQueryTransformation(key, value string) {
	t.queryTransformations[key] = value
}

// RemoveQueryTransformation removes a query parameter transformation
func (t *Transformer) RemoveQueryTransformation(key string) {
	delete(t.queryTransformations, key)
} 