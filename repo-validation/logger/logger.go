package logger

import (
	"github.com/charmbracelet/log"
)

// Logger defines the interface for logging
type Logger interface {
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	Debug(msg string, keyvals ...interface{})
}

// CharmLogger implements Logger using charmbracelet/log
type CharmLogger struct {
	logger *log.Logger
}

// NewCharmLogger creates a new CharmLogger
func NewCharmLogger() *CharmLogger {
	return &CharmLogger{
		logger: log.NewWithOptions(log.Options{
			Level:           log.InfoLevel,
			ReportTimestamp: true,
		}),
	}
}

// Info logs an info message
func (l *CharmLogger) Info(msg string, keyvals ...interface{}) {
	l.logger.Info(msg, keyvals...)
}

// Warn logs a warning message
func (l *CharmLogger) Warn(msg string, keyvals ...interface{}) {
	l.logger.Warn(msg, keyvals...)
}

// Error logs an error message
func (l *CharmLogger) Error(msg string, keyvals ...interface{}) {
	l.logger.Error(msg, keyvals...)
}

// Debug logs a debug message
func (l *CharmLogger) Debug(msg string, keyvals ...interface{}) {
	l.logger.Debug(msg, keyvals...)
}

// SetLevel sets the log level
func (l *CharmLogger) SetLevel(level log.Level) {
	l.logger.SetLevel(level)
}

// WithJSON configures the logger to output JSON
func (l *CharmLogger) WithJSON() *CharmLogger {
	l.logger.SetFormatter(log.JSONFormatter)
	return l
}
