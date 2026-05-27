package database

import "context"

// Logger defines the interface for structured logging required by this package.
// Consumers must supply their own implementation (dependency inversion).
type Logger interface {
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	With(keysAndValues ...any) Logger
	WithContext(ctx context.Context) Logger
}

// noopLogger is a Logger that silently discards every message.
type noopLogger struct{}

func (noopLogger) Debug(_ string, _ ...any)               {}
func (noopLogger) Info(_ string, _ ...any)                {}
func (noopLogger) Warn(_ string, _ ...any)                {}
func (noopLogger) Error(_ string, _ ...any)               {}
func (n noopLogger) With(_ ...any) Logger                 { return n }
func (n noopLogger) WithContext(_ context.Context) Logger { return n }

// NewNoopLogger returns a Logger that discards all output. Useful for tests.
func NewNoopLogger() Logger { return noopLogger{} }
