package errors

import (
	"fmt"
	"net/http"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"
	// ErrorTypeBackend represents backend-related errors
	ErrorTypeBackend ErrorType = "BACKEND_ERROR"
	// ErrorTypePolicy represents policy-related errors
	ErrorTypePolicy ErrorType = "POLICY_ERROR"
	// ErrorTypeRouting represents routing-related errors
	ErrorTypeRouting ErrorType = "ROUTING_ERROR"
	// ErrorTypeHealthCheck represents health check related errors
	ErrorTypeHealthCheck ErrorType = "HEALTH_CHECK_ERROR"
)

// LoadBalancerError represents a custom error type for the load balancer
type LoadBalancerError struct {
	Type    ErrorType
	Message string
	Err     error
	Code    int
}

func (e *LoadBalancerError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(message string, err error) *LoadBalancerError {
	return &LoadBalancerError{
		Type:    ErrorTypeValidation,
		Message: message,
		Err:     err,
		Code:    http.StatusBadRequest,
	}
}

// NewBackendError creates a new backend error
func NewBackendError(message string, err error) *LoadBalancerError {
	return &LoadBalancerError{
		Type:    ErrorTypeBackend,
		Message: message,
		Err:     err,
		Code:    http.StatusBadGateway,
	}
}

// NewPolicyError creates a new policy error
func NewPolicyError(message string, err error) *LoadBalancerError {
	return &LoadBalancerError{
		Type:    ErrorTypePolicy,
		Message: message,
		Err:     err,
		Code:    http.StatusForbidden,
	}
}

// NewRoutingError creates a new routing error
func NewRoutingError(message string, err error) *LoadBalancerError {
	return &LoadBalancerError{
		Type:    ErrorTypeRouting,
		Message: message,
		Err:     err,
		Code:    http.StatusNotFound,
	}
}

// NewHealthCheckError creates a new health check error
func NewHealthCheckError(message string, err error) *LoadBalancerError {
	return &LoadBalancerError{
		Type:    ErrorTypeHealthCheck,
		Message: message,
		Err:     err,
		Code:    http.StatusServiceUnavailable,
	}
} 