package policy

import (
	"fmt"
	"net/http"

	"github.com/rixtrayker/go-loadbalancer/configs"
	"github.com/rixtrayker/go-loadbalancer/internal/policy/ratelimit"
	"github.com/rixtrayker/go-loadbalancer/internal/policy/security"
	"github.com/rixtrayker/go-loadbalancer/internal/policy/transform"
)

// Apply applies a policy to a request
func Apply(policy configs.PolicyConfig, r *http.Request) error {
	// Apply rate limiting
	if policy.RateLimit != "" {
		if err := ratelimit.Apply(policy.RateLimit, r); err != nil {
			return fmt.Errorf("rate limit policy failed: %w", err)
		}
	}

	// Apply ACL
	if policy.ACL != "" {
		if err := security.Apply(policy.ACL, r); err != nil {
			return fmt.Errorf("ACL policy failed: %w", err)
		}
	}

	// Apply transformations
	if policy.Transform != "" {
		if err := transform.Apply(policy.Transform, r); err != nil {
			return fmt.Errorf("transform policy failed: %w", err)
		}
	}

	return nil
}
