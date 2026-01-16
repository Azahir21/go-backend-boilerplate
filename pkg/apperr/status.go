package apperr

import "net/http"

// Status represents a standardized error status code.
type Status string

// Error status constants
const (
	StatusBadGateway          Status = "bad_gateway"
	StatusBadRequest          Status = "bad_request"
	StatusConflict            Status = "conflict"
	StatusForbidden           Status = "forbidden"
	StatusInternalServer      Status = "internal_server"
	StatusMethodNotAllowed    Status = "method_not_allowed"
	StatusNotFound            Status = "not_found"
	StatusNotImplemented      Status = "not_implemented"
	StatusServiceUnavailable  Status = "service_unavailable"
	StatusTimeout             Status = "timeout"
	StatusTooManyRequests     Status = "too_many_requests"
	StatusUnauthorized        Status = "unauthorized"
	StatusUnprocessableEntity Status = "unprocessable_entity"
)

// statusHTTPMap maps Status to HTTP status codes.
var statusHTTPMap = map[Status]int{
	StatusBadGateway:          http.StatusBadGateway,
	StatusBadRequest:          http.StatusBadRequest,
	StatusConflict:            http.StatusConflict,
	StatusForbidden:           http.StatusForbidden,
	StatusInternalServer:      http.StatusInternalServerError,
	StatusMethodNotAllowed:    http.StatusMethodNotAllowed,
	StatusNotFound:            http.StatusNotFound,
	StatusNotImplemented:      http.StatusNotImplemented,
	StatusServiceUnavailable:  http.StatusServiceUnavailable,
	StatusTimeout:             http.StatusGatewayTimeout,
	StatusTooManyRequests:     http.StatusTooManyRequests,
	StatusUnauthorized:        http.StatusUnauthorized,
	StatusUnprocessableEntity: http.StatusUnprocessableEntity,
}

// HTTPCode returns the HTTP status code for the given status.
func (s Status) HTTPCode() int {
	if code, ok := statusHTTPMap[s]; ok {
		return code
	}
	return http.StatusInternalServerError
}

// String returns the string representation of the status.
func (s Status) String() string {
	return string(s)
}

// IsValid checks if the status is a valid defined status.
func (s Status) IsValid() bool {
	_, ok := statusHTTPMap[s]
	return ok
}
