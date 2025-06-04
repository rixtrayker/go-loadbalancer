package policy

import (
	"net/http"
)

// Policy defines the interface for request/response policies
type Policy interface {
	// Apply applies the policy to the request/response
	Apply(req *http.Request, resp *http.Response) error
	// Name returns the name of the policy
	Name() string
}

// BasePolicy provides common functionality for policies
type BasePolicy struct {
	PolicyName string
}

// Name returns the name of the policy
func (b *BasePolicy) Name() string {
	return b.PolicyName
}

// PolicyChain represents a chain of policies to be applied
type PolicyChain struct {
	policies []Policy
}

// NewPolicyChain creates a new policy chain
func NewPolicyChain() *PolicyChain {
	return &PolicyChain{
		policies: make([]Policy, 0),
	}
}

// AddPolicy adds a policy to the chain
func (pc *PolicyChain) AddPolicy(policy Policy) {
	pc.policies = append(pc.policies, policy)
}

// RemovePolicy removes a policy from the chain
func (pc *PolicyChain) RemovePolicy(name string) {
	for i, p := range pc.policies {
		if p.Name() == name {
			pc.policies = append(pc.policies[:i], pc.policies[i+1:]...)
			return
		}
	}
}

// Apply applies all policies in the chain
func (pc *PolicyChain) Apply(req *http.Request, resp *http.Response) error {
	for _, policy := range pc.policies {
		if err := policy.Apply(req, resp); err != nil {
			return err
		}
	}
	return nil
}
