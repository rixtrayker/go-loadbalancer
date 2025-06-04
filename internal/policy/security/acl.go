package security

import (
	"errors"
	"net"
	"net/http"
	"strings"
	"sync"
)

// ACL implements access control lists
type ACL struct {
	allowList []string
	denyList  []string
	mutex     sync.RWMutex
}

var (
	// Global ACL instance
	globalACL = NewACL()
)

// NewACL creates a new ACL
func NewACL() *ACL {
	return &ACL{
		allowList: make([]string, 0),
		denyList:  make([]string, 0),
	}
}

// Apply applies ACL rules to a request
func Apply(aclStr string, r *http.Request) error {
	// Parse ACL string (e.g., "allow:192.168.1.0/24,deny:10.0.0.1")
	rules := strings.Split(aclStr, ",")
	
	// Get client IP
	clientIP := getClientIP(r)
	ip := net.ParseIP(clientIP)
	if ip == nil {
		return errors.New("invalid client IP")
	}

	// Check each rule
	for _, rule := range rules {
		parts := strings.SplitN(rule, ":", 2)
		if len(parts) != 2 {
			continue
		}

		action := strings.ToLower(parts[0])
		cidr := parts[1]

		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			// Try as a single IP
			singleIP := net.ParseIP(cidr)
			if singleIP == nil {
				continue
			}
			
			// Check exact IP match
			if ip.Equal(singleIP) {
				if action == "deny" {
					return errors.New("access denied by ACL")
				}
				return nil // Explicitly allowed
			}
		} else {
			// Check CIDR match
			if ipNet.Contains(ip) {
				if action == "deny" {
					return errors.New("access denied by ACL")
				}
				return nil // Explicitly allowed
			}
		}
	}

	// Default to allow if no matching rules
	return nil
}

// AddAllowRule adds an allow rule to the ACL
func (a *ACL) AddAllowRule(cidr string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.allowList = append(a.allowList, cidr)
}

// AddDenyRule adds a deny rule to the ACL
func (a *ACL) AddDenyRule(cidr string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.denyList = append(a.denyList, cidr)
}

// Check checks if an IP is allowed by the ACL
func (a *ACL) Check(ipStr string) bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// Check deny list first
	for _, cidr := range a.denyList {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if ipNet.Contains(ip) {
			return false
		}
	}

	// If allow list is empty, allow all
	if len(a.allowList) == 0 {
		return true
	}

	// Check allow list
	for _, cidr := range a.allowList {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if ipNet.Contains(ip) {
			return true
		}
	}

	return false
}

// getClientIP extracts the client IP from a request
func getClientIP(r *http.Request) string {
	// Try X-Forwarded-For header first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Try X-Real-IP header
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	// Fall back to RemoteAddr
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
