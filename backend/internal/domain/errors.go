package domain

import "fmt"

// Custom domain errors
type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Error codes
const (
	ErrCodeInvalidInput    = "INVALID_INPUT"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeDuplicate       = "DUPLICATE"
	ErrCodeUnauthorized    = "UNAUTHORIZED"
	ErrCodeInternalError   = "INTERNAL_ERROR"
	ErrCodeInvalidGeofence = "INVALID_GEOFENCE"
	ErrCodeInvalidCoords   = "INVALID_COORDINATES"
)

// Common domain errors
var (
	ErrGeofenceNotFound  = &DomainError{Code: ErrCodeNotFound, Message: "geofence not found"}
	ErrVehicleNotFound   = &DomainError{Code: ErrCodeNotFound, Message: "vehicle not found"}
	ErrAlertNotFound     = &DomainError{Code: ErrCodeNotFound, Message: "alert rule not found"}
	ErrInvalidCoords     = &DomainError{Code: ErrCodeInvalidCoords, Message: "invalid coordinates provided"}
	ErrVehicleExists     = &DomainError{Code: ErrCodeDuplicate, Message: "vehicle already exists"}
	ErrInvalidGeometry   = &DomainError{Code: ErrCodeInvalidGeofence, Message: "invalid geofence geometry"}
)

// NewDomainError creates a new domain error
func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}
