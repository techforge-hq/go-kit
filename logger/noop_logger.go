package logger

import "context"

// NoopLogger is a logger that does nothing.
// Useful for testing when you don't want log output.
type NoopLogger struct{}

// NewNoop creates a new no-op logger.
func NewNoop() Logger {
	return NoopLogger{}
}

// Debug logs a debug message (no-op implementation).
func (n NoopLogger) Debug(_ string, _ ...interface{}) {}

// Info logs an informational message (no-op implementation).
func (n NoopLogger) Info(_ string, _ ...interface{}) {}

// Warn logs a warning message (no-op implementation).
func (n NoopLogger) Warn(_ string, _ ...interface{}) {}

// Error logs an error message (no-op implementation).
func (n NoopLogger) Error(_ string, _ ...interface{}) {}

// With returns a new Logger with additional context fields (no-op implementation).
func (n NoopLogger) With(_ ...interface{}) Logger {
	return n
}

// WithContext returns a new Logger with context (no-op implementation).
func (n NoopLogger) WithContext(_ context.Context) Logger {
	return n
}
