// Package logger provides structured logging functionality with multiple levels and output formats.
package logger

import "context"

// Level represents the logging level.
type Level int

const (
	// LevelDebug is the lowest logging level.
	LevelDebug Level = iota
	// LevelInfo is the default logging level.
	LevelInfo
	// LevelWarn is for warning messages.
	LevelWarn
	// LevelError is for error messages.
	LevelError
)

// Format represents the output format.
type Format string

const (
	// FormatJSON represents JSON output format.
	FormatJSON Format = "json"
	// FormatText represents plain text output format.
	FormatText Format = "text"
)

// Config holds logger configuration.
type Config struct {
	Level  Level  // Minimum level to log
	Format Format // Output format (json or text)
}

// Logger defines the interface for structured logging.
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	With(keysAndValues ...interface{}) Logger
	WithContext(ctx context.Context) Logger
}
