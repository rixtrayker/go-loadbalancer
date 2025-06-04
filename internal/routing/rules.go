package routing

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/rixtrayker/go-loadbalancer/configs"
)

var (
	// ErrNoMatchingRule is returned when no routing rule matches a request
	ErrNoMatchingRule = errors.New("no matching routing rule")
)

// Rule represents a routing rule
type Rule struct {
	HostPattern  *regexp.Regexp
	PathPattern  *regexp.Regexp
	Method       string
	HeaderRules  map[string]*regexp.Regexp
	TargetPool   string
	Policies     []configs.PolicyConfig
}

// NewRule creates a new routing rule from config
func NewRule(config configs.RoutingRuleConfig) *Rule {
	rule := &Rule{
		Method:      config.Match.Method,
		TargetPool:  config.TargetPool,
		Policies:    config.Policies,
		HeaderRules: make(map[string]*regexp.Regexp),
	}

	// Compile host pattern
	if config.Match.Host != "" {
		pattern := strings.Replace(config.Match.Host, ".", "\\.", -1)
		pattern = strings.Replace(pattern, "*", ".*", -1)
		rule.HostPattern = regexp.MustCompile("^" + pattern + "$")
	}

	// Compile path pattern
	if config.Match.Path != "" {
		pattern := strings.Replace(config.Match.Path, "*", ".*", -1)
		rule.PathPattern = regexp.MustCompile("^" + pattern + "$")
	}

	// Compile header patterns
	for name, pattern := range config.Match.Headers {
		rule.HeaderRules[name] = regexp.MustCompile(pattern)
	}

	return rule
}

// Matches checks if a request matches this rule
func (r *Rule) Matches(req *http.Request) bool {
	// Check host
	if r.HostPattern != nil {
		if !r.HostPattern.MatchString(req.Host) {
			return false
		}
	}

	// Check path
	if r.PathPattern != nil {
		if !r.PathPattern.MatchString(req.URL.Path) {
			return false
		}
	}

	// Check method
	if r.Method != "" && r.Method != req.Method {
		return false
	}

	// Check headers
	for name, pattern := range r.HeaderRules {
		value := req.Header.Get(name)
		if value == "" || !pattern.MatchString(value) {
			return false
		}
	}

	return true
}
