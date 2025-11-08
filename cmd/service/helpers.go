package service

import (
	"fmt"
	"time"
)

// parseDuration is a helper to parse duration strings and return an error.
func parseDuration(durationStr, fieldName string) (time.Duration, error) {
	d, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s duration: %w", fieldName, err)
	}
	return d, nil
}
