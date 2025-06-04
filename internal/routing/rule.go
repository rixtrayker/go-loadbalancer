package routing

import (
	"net/http"
	"regexp"

	"github.com/rixtrayker/go-loadbalancer/configs"
)

// Rule represents a routing rule
type Rule struct {
	HostPattern *regexp.Regexp
	PathPattern *regexp.Regexp
	Method      string
	HeaderRules map[string]*regexp.Regexp
	TargetPool  string
	Policies    []configs.PolicyConfig
}

// NewRule creates a new routing rule from config
func NewRule(config configs.RoutingRuleConfig) *Rule {
	rule := &Rule{
		Method:     config.Match.Method,
		TargetPool: config.TargetPool,
		Policies:   config.Policies,
	}

	// Compile host pattern
	if config.Match.Host != "" {
		pattern := wildcardToRegexp(config.Match.Host)
		rule.HostPattern = regexp.MustCompile(pattern)
	}

	// Compile path pattern
	if config.Match.Path != "" {
		pattern := wildcardToRegexp(config.Match.Path)
		rule.PathPattern = regexp.MustCompile(pattern)
	}

	// Compile header patterns
	if len(config.Match.Headers) > 0 {
		rule.HeaderRules = make(map[string]*regexp.Regexp)
		for k, v := range config.Match.Headers {
			rule.HeaderRules[k] = regexp.MustCompile(v)
		}
	}

	return rule
}

// Matches checks if a request matches this rule
func (r *Rule) Matches(req *http.Request) bool {
	// Check host
	if r.HostPattern != nil && !r.HostPattern.MatchString(req.Host) {
		return false
	}

	// Check path
	if r.PathPattern != nil && !r.PathPattern.MatchString(req.URL.Path) {
		return false
	}

	// Check method
	if r.Method != "" && r.Method != req.Method {
		return false
	}

	// Check headers
	if r.HeaderRules != nil {
		for k, v := range r.HeaderRules {
			headerVal := req.Header.Get(k)
			if headerVal == "" || !v.MatchString(headerVal) {
				return false
			}
		}
	}

	return true
}

// wildcardToRegexp converts a wildcard pattern to a regexp pattern
func wildcardToRegexp(pattern string) string {
	// Replace * with .*
	return "^" + regexp.QuoteMeta(pattern) + "$"
}
