package azpays

import "fmt"

// Error represents an API error returned by AzPays.
type Error struct {
	// StatusCode is the HTTP status code.
	StatusCode int `json:"status_code"`

	// Message is the error message from the API.
	Message string `json:"message"`

	// RequestID is the unique request identifier for debugging.
	RequestID string `json:"request_id,omitempty"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("azpays: %d %s (request_id: %s)", e.StatusCode, e.Message, e.RequestID)
	}
	return fmt.Sprintf("azpays: %d %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error is a 404 Not Found.
func IsNotFound(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode == 404
	}
	return false
}

// IsUnauthorized returns true if the error is a 401 Unauthorized.
func IsUnauthorized(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode == 401
	}
	return false
}

// IsForbidden returns true if the error is a 403 Forbidden.
func IsForbidden(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode == 403
	}
	return false
}

// IsRateLimited returns true if the error is a 429 Too Many Requests.
func IsRateLimited(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode == 429
	}
	return false
}

// IsBadRequest returns true if the error is a 400 Bad Request.
func IsBadRequest(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode == 400
	}
	return false
}
