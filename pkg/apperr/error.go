package apperr

import (
	"fmt"
	"runtime"
	"strings"
)

// AppError represents a structured application error with status, message, detail, and stacktrace.
type AppError struct {
	Status     Status `json:"status"`
	Message    string `json:"message"`
	Detail     string `json:"detail,omitempty"`
	Stacktrace string `json:"stacktrace,omitempty"`
	Cause      error  `json:"-"`
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Status, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Status, e.Message)
}

// Unwrap returns the underlying cause of the error.
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithDetail adds detail information to the error.
func (e *AppError) WithDetail(detail string) *AppError {
	e.Detail = detail
	return e
}

// WithCause wraps an underlying error.
func (e *AppError) WithCause(err error) *AppError {
	e.Cause = err
	if e.Detail == "" && err != nil {
		e.Detail = err.Error()
	}
	return e
}

// captureStacktrace captures the current stacktrace, skipping the specified number of frames.
func captureStacktrace(skip int) string {
	const maxDepth = 32
	var pcs [maxDepth]uintptr
	n := runtime.Callers(skip, pcs[:])
	if n == 0 {
		return ""
	}

	frames := runtime.CallersFrames(pcs[:n])
	var sb strings.Builder

	for {
		frame, more := frames.Next()
		// Skip runtime and internal frames
		if strings.Contains(frame.File, "runtime/") {
			if !more {
				break
			}
			continue
		}

		fmt.Fprintf(&sb, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)

		if !more {
			break
		}
	}

	return sb.String()
}

// newAppError creates a new AppError with the given status and message, capturing stacktrace.
func newAppError(status Status, message string) *AppError {
	return &AppError{
		Status:     status,
		Message:    message,
		Stacktrace: captureStacktrace(3), // Skip: captureStacktrace, newAppError, and the constructor
	}
}

// --- Constructor functions for each status ---

// BadGateway creates a new BadGateway error.
func BadGateway(message string) *AppError {
	return newAppError(StatusBadGateway, message)
}

// BadRequest creates a new BadRequest error.
func BadRequest(message string) *AppError {
	return newAppError(StatusBadRequest, message)
}

// Conflict creates a new Conflict error.
func Conflict(message string) *AppError {
	return newAppError(StatusConflict, message)
}

// Forbidden creates a new Forbidden error.
func Forbidden(message string) *AppError {
	return newAppError(StatusForbidden, message)
}

// InternalServer creates a new InternalServer error.
func InternalServer(message string) *AppError {
	return newAppError(StatusInternalServer, message)
}

// MethodNotAllowed creates a new MethodNotAllowed error.
func MethodNotAllowed(message string) *AppError {
	return newAppError(StatusMethodNotAllowed, message)
}

// NotFound creates a new NotFound error.
func NotFound(message string) *AppError {
	return newAppError(StatusNotFound, message)
}

// NotImplemented creates a new NotImplemented error.
func NotImplemented(message string) *AppError {
	return newAppError(StatusNotImplemented, message)
}

// ServiceUnavailable creates a new ServiceUnavailable error.
func ServiceUnavailable(message string) *AppError {
	return newAppError(StatusServiceUnavailable, message)
}

// Timeout creates a new Timeout error.
func Timeout(message string) *AppError {
	return newAppError(StatusTimeout, message)
}

// TooManyRequests creates a new TooManyRequests error.
func TooManyRequests(message string) *AppError {
	return newAppError(StatusTooManyRequests, message)
}

// Unauthorized creates a new Unauthorized error.
func Unauthorized(message string) *AppError {
	return newAppError(StatusUnauthorized, message)
}

// UnprocessableEntity creates a new UnprocessableEntity error.
func UnprocessableEntity(message string) *AppError {
	return newAppError(StatusUnprocessableEntity, message)
}

// Wrap wraps an existing error into an AppError with the given status and message.
// If the error is already an AppError, it preserves the original stacktrace.
func Wrap(err error, status Status, message string) *AppError {
	if err == nil {
		return nil
	}

	// If it's already an AppError, preserve the stacktrace but update status/message
	if appErr, ok := err.(*AppError); ok {
		return &AppError{
			Status:     status,
			Message:    message,
			Detail:     appErr.Detail,
			Stacktrace: appErr.Stacktrace,
			Cause:      appErr,
		}
	}

	appErr := newAppError(status, message)
	appErr.Cause = err
	appErr.Detail = err.Error()
	return appErr
}

// Is checks if the error is an AppError with the given status.
func Is(err error, status Status) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Status == status
	}
	return false
}

// AsAppError attempts to convert an error to an AppError.
// If the error is not an AppError, it wraps it as an InternalServer error.
func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return InternalServer("An unexpected error occurred").WithCause(err)
}
