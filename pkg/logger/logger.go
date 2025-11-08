package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogger creates a new logger instance with JSON formatting
func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
	return log
}

// NewLoggerWithLevel creates a new logger with a specific log level
func NewLoggerWithLevel(level logrus.Level) *logrus.Logger {
	log := NewLogger()
	log.SetLevel(level)
	return log
}
