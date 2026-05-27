package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// SlogAdapter implements Logger interface using slog.
type SlogAdapter struct {
	logger *slog.Logger
}

// NewSlogAdapter creates a new slog-based logger.
func NewSlogAdapter(config Config) Logger {
	return newSlogAdapter(config, os.Stdout)
}

// newSlogAdapter is internal factory for testing with custom writer.
func newSlogAdapter(config Config, writer io.Writer) Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: convertLevel(config.Level),
	}

	switch config.Format {
	case FormatJSON:
		handler = slog.NewJSONHandler(writer, opts)
	case FormatText:
		handler = slog.NewTextHandler(writer, opts)
	default:
		handler = slog.NewTextHandler(writer, opts)
	}

	return SlogAdapter{
		logger: slog.New(handler),
	}
}

// Debug logs a debug message with optional key-value pairs.
func (s SlogAdapter) Debug(msg string, keysAndValues ...any) {
	s.logger.Debug(msg, keysAndValues...)
}

// Info logs an informational message with optional key-value pairs.
func (s SlogAdapter) Info(msg string, keysAndValues ...any) {
	s.logger.Info(msg, keysAndValues...)
}

// Warn logs a warning message with optional key-value pairs.
func (s SlogAdapter) Warn(msg string, keysAndValues ...any) {
	s.logger.Warn(msg, keysAndValues...)
}

func (s SlogAdapter) Error(msg string, keysAndValues ...any) {
	s.logger.Error(msg, keysAndValues...)
}

// With returns a new Logger with additional context fields.
func (s SlogAdapter) With(keysAndValues ...any) Logger {
	return SlogAdapter{
		logger: s.logger.With(keysAndValues...),
	}
}

// WithContext returns a new Logger with context.
func (s SlogAdapter) WithContext(_ context.Context) Logger {
	return SlogAdapter{
		logger: s.logger.With(), // Can extract context values here if needed
	}
}

// convertLevel converts our Level to slog.Level.
func convertLevel(level Level) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// NewDevelopment creates a logger configured for development
// Uses text format and debug level.
func NewDevelopment() Logger {
	return NewSlogAdapter(Config{
		Level:  LevelDebug,
		Format: FormatText,
	})
}

// NewProduction creates a logger configured for production
// Uses JSON format and info level.
func NewProduction() Logger {
	return NewSlogAdapter(Config{
		Level:  LevelInfo,
		Format: FormatJSON,
	})
}
