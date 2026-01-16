package apperr

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestStatus_HTTPCode(t *testing.T) {
	tests := []struct {
		status   Status
		expected int
	}{
		{StatusBadGateway, http.StatusBadGateway},
		{StatusBadRequest, http.StatusBadRequest},
		{StatusConflict, http.StatusConflict},
		{StatusForbidden, http.StatusForbidden},
		{StatusInternalServer, http.StatusInternalServerError},
		{StatusMethodNotAllowed, http.StatusMethodNotAllowed},
		{StatusNotFound, http.StatusNotFound},
		{StatusNotImplemented, http.StatusNotImplemented},
		{StatusServiceUnavailable, http.StatusServiceUnavailable},
		{StatusTimeout, http.StatusGatewayTimeout},
		{StatusTooManyRequests, http.StatusTooManyRequests},
		{StatusUnauthorized, http.StatusUnauthorized},
		{StatusUnprocessableEntity, http.StatusUnprocessableEntity},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.HTTPCode(); got != tt.expected {
				t.Errorf("HTTPCode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStatus_IsValid(t *testing.T) {
	tests := []struct {
		status   Status
		expected bool
	}{
		{StatusBadRequest, true},
		{StatusNotFound, true},
		{Status("invalid_status"), false},
		{Status(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.expected {
				t.Errorf("IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAppError_Error(t *testing.T) {
	err := BadRequest("invalid input")
	expected := "bad_request: invalid input"
	if got := err.Error(); got != expected {
		t.Errorf("Error() = %v, want %v", got, expected)
	}
}

func TestAppError_WithCause(t *testing.T) {
	cause := errors.New("database connection failed")
	err := InternalServer("failed to process request").WithCause(cause)

	if err.Cause != cause {
		t.Errorf("Cause = %v, want %v", err.Cause, cause)
	}

	if err.Detail != cause.Error() {
		t.Errorf("Detail = %v, want %v", err.Detail, cause.Error())
	}
}

func TestAppError_WithDetail(t *testing.T) {
	err := BadRequest("validation failed").WithDetail("field 'email' is required")

	if err.Detail != "field 'email' is required" {
		t.Errorf("Detail = %v, want %v", err.Detail, "field 'email' is required")
	}
}

func TestAppError_Stacktrace(t *testing.T) {
	err := NotFound("user not found")

	if err.Stacktrace == "" {
		t.Error("Stacktrace should not be empty")
	}
}

func TestWrap(t *testing.T) {
	cause := errors.New("original error")
	wrapped := Wrap(cause, StatusBadRequest, "request failed")

	if wrapped.Status != StatusBadRequest {
		t.Errorf("Status = %v, want %v", wrapped.Status, StatusBadRequest)
	}

	if wrapped.Cause != cause {
		t.Errorf("Cause = %v, want %v", wrapped.Cause, cause)
	}

	if wrapped.Detail != cause.Error() {
		t.Errorf("Detail = %v, want %v", wrapped.Detail, cause.Error())
	}
}

func TestWrap_NilError(t *testing.T) {
	wrapped := Wrap(nil, StatusBadRequest, "request failed")
	if wrapped != nil {
		t.Errorf("Wrap(nil) should return nil, got %v", wrapped)
	}
}

func TestIs(t *testing.T) {
	err := NotFound("resource not found")

	if !Is(err, StatusNotFound) {
		t.Error("Is() should return true for matching status")
	}

	if Is(err, StatusBadRequest) {
		t.Error("Is() should return false for non-matching status")
	}
}

func TestAsAppError(t *testing.T) {
	t.Run("AppError passthrough", func(t *testing.T) {
		original := BadRequest("bad input")
		result := AsAppError(original)

		if result != original {
			t.Error("AsAppError should return the same AppError")
		}
	})

	t.Run("Regular error conversion", func(t *testing.T) {
		regular := errors.New("some error")
		result := AsAppError(regular)

		if result.Status != StatusInternalServer {
			t.Errorf("Status = %v, want %v", result.Status, StatusInternalServer)
		}

		if result.Cause != regular {
			t.Errorf("Cause = %v, want %v", result.Cause, regular)
		}
	})

	t.Run("Nil error", func(t *testing.T) {
		result := AsAppError(nil)
		if result != nil {
			t.Errorf("AsAppError(nil) should return nil, got %v", result)
		}
	})
}

func TestResponder_ProductionConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := DefaultConfig() // Production-like config
	responder := NewResponder(config)

	err := BadRequest("validation failed").WithDetail("email is required")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	responder.Respond(c, err)

	// Check that detail is not in response
	body := w.Body.String()
	if containsString(body, "email is required") {
		t.Error("Production response should not contain detail")
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestResponder_DevelopmentConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := DevelopmentConfig()
	responder := NewResponder(config)

	err := BadRequest("validation failed").WithDetail("email is required")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	responder.Respond(c, err)

	body := w.Body.String()
	if !containsString(body, "email is required") {
		t.Error("Development response should contain detail")
	}

	if !containsString(body, "stacktrace") {
		t.Error("Development response should contain stacktrace")
	}
}

func TestConfigFromEnv(t *testing.T) {
	tests := []struct {
		env            string
		expectDetail   bool
		expectStack    bool
	}{
		{"production", false, false},
		{"staging", false, false},
		{"development", true, true},
		{"local", true, true},
		{"", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			config := ConfigFromEnv(tt.env)

			if config.ShowDetail != tt.expectDetail {
				t.Errorf("ShowDetail = %v, want %v", config.ShowDetail, tt.expectDetail)
			}

			if config.ShowStacktrace != tt.expectStack {
				t.Errorf("ShowStacktrace = %v, want %v", config.ShowStacktrace, tt.expectStack)
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStringHelper(s, substr))
}

func containsStringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
