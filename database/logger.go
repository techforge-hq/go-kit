package database

import "context"

// Logger is a structured logger whose With and WithContext return the same type T.
// This lets callers pass their own logger implementation without an adapter, as long
// as chaining methods return the concrete logger type (F-bounded constraint).
type Logger[T any] interface {
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	With(keysAndValues ...any) T
	WithContext(ctx context.Context) T
}

// noopLogger is a Logger that silently discards every message.
type noopLogger struct{}

func (noopLogger) Debug(_ string, _ ...any)               {}
func (noopLogger) Info(_ string, _ ...any)                {}
func (noopLogger) Warn(_ string, _ ...any)                {}
func (noopLogger) Error(_ string, _ ...any)               {}
func (n noopLogger) With(_ ...any) noopLogger             { return n }
func (n noopLogger) WithContext(_ context.Context) noopLogger { return n }

// NewNoopLogger returns a Logger that discards all output. Useful for tests.
func NewNoopLogger() noopLogger { return noopLogger{} }
