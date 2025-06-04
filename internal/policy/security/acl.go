package security

import (
	"net"
	"net/http"

	"github.com/amr/go-loadbalancer/internal/policy"
)

// ACL implements access control list policy
type ACL struct {
	policy.BasePolicy
	allowedIPs    []*net.IPNet
	blockedIPs    []*net.IPNet
	allowedHosts  []string
	blockedHosts  []string
}

// NewACL creates a new ACL policy
func NewACL() *ACL {
	return &ACL{
		BasePolicy: policy.BasePolicy{
			PolicyName: "acl",
		},
		allowedIPs:    make([]*net.IPNet, 0),
		blockedIPs:    make([]*net.IPNet, 0),
		allowedHosts:  make([]string, 0),
		blockedHosts:  make([]string, 0),
	}
}

// SetAllowedIPs sets the allowed IP ranges
func (a *ACL) SetAllowedIPs(ips []string) error {
	a.allowedIPs = make([]*net.IPNet, 0, len(ips))
	for _, ipStr := range ips {
		_, ipNet, err := net.ParseCIDR(ipStr)
		if err != nil {
			return err
		}
		a.allowedIPs = append(a.allowedIPs, ipNet)
	}
	return nil
}

// Apply implements the Policy interface
func (a *ACL) Apply(req *http.Request, _ *http.Response) error {
	// Check IP address
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		ip = req.RemoteAddr
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ErrInvalidIP
	}

	// Check if IP is blocked
	for _, blocked := range a.blockedIPs {
		if blocked.Contains(parsedIP) {
			return ErrAccessDenied
		}
	}

	// If allowed IPs are specified, check if IP is allowed
	if len(a.allowedIPs) > 0 {
		allowed := false
		for _, allowedIP := range a.allowedIPs {
			if allowedIP.Contains(parsedIP) {
				allowed = true
				break
			}
		}
		if !allowed {
			return ErrAccessDenied
		}
	}

	// Check host
	host := req.Host
	if len(a.blockedHosts) > 0 {
		for _, blocked := range a.blockedHosts {
			if host == blocked {
				return ErrAccessDenied
			}
		}
	}

	if len(a.allowedHosts) > 0 {
		allowed := false
		for _, allowedHost := range a.allowedHosts {
			if host == allowedHost {
				allowed = true
				break
			}
		}
		if !allowed {
			return ErrAccessDenied
		}
	}

	return nil
}

// AddAllowedIP adds an allowed IP or CIDR
func (a *ACL) AddAllowedIP(cidr string) error {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	a.allowedIPs = append(a.allowedIPs, ipnet)
	return nil
}

// AddBlockedIP adds a blocked IP or CIDR
func (a *ACL) AddBlockedIP(cidr string) error {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	a.blockedIPs = append(a.blockedIPs, ipnet)
	return nil
}

// AddAllowedHost adds an allowed host
func (a *ACL) AddAllowedHost(host string) {
	a.allowedHosts = append(a.allowedHosts, host)
}

// AddBlockedHost adds a blocked host
func (a *ACL) AddBlockedHost(host string) {
	a.blockedHosts = append(a.blockedHosts, host)
}

// Errors
var (
	ErrAccessDenied = &ACLError{msg: "access denied"}
	ErrInvalidIP    = &ACLError{msg: "invalid IP address"}
)

// ACLError represents an ACL error
type ACLError struct {
	msg string
}

func (e *ACLError) Error() string {
	return e.msg
}
