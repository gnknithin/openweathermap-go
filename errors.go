package openweathermap

import (
	"fmt"
)

// APIError represents a structured error returned by the OpenWeatherMap API endpoints.
// It implements the standard Go error interface.
type APIError struct {
	Code    int    `json:"cod"`
	Message string `json:"message"`
}

// Error formats the APIError into a human-readable string to fulfill the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("openweathermap api error: status code %d: %s", e.Code, e.Message)
}

// IsAPIError checks if a generic error is an OpenWeatherMap APIError of a specific HTTP status code.
func IsAPIError(err error, statusCode int) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Code == statusCode
	}
	return false
}
