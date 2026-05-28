package di

import (
	"context"

	"github.com/samber/do/v2"
)

// Scope is an alias for do.Scope for convenient access.
type Scope = do.Scope

// Injector is an alias for do.Injector for convenient access.
type Injector = do.Injector

// Provider is an alias for do.Provider for convenient access.
type Provider[T any] = do.Provider[T]

// RootScope is an alias for do.RootScope for convenient access.
type RootScope = do.RootScope

// ShutdownReport is an alias for do.ShutdownReport for convenient access.
type ShutdownReport = do.ShutdownReport

// New creates a new root scope (DI container).
func New() *do.RootScope {
	return do.New()
}

// NewWithOpts creates a new root scope with options.
func NewWithOpts(opts *do.InjectorOpts) *do.RootScope {
	return do.NewWithOpts(opts)
}

// Provide registers a lazy-loaded service by type.
func Provide[T any](i do.Injector, provider do.Provider[T]) {
	do.Provide(i, provider)
}

// ProvideNamed registers a lazy-loaded service by name.
func ProvideNamed[T any](i do.Injector, name string, provider do.Provider[T]) {
	do.ProvideNamed(i, name, provider)
}

// ProvideValue registers an eager-loaded service value by type.
func ProvideValue[T any](i do.Injector, value T) {
	do.ProvideValue(i, value)
}

// ProvideNamedValue registers an eager-loaded service value by name.
func ProvideNamedValue[T any](i do.Injector, name string, value T) {
	do.ProvideNamedValue(i, name, value)
}

// ProvideTransient registers a transient service (new instance each invocation).
func ProvideTransient[T any](i do.Injector, provider do.Provider[T]) {
	do.ProvideTransient(i, provider)
}

// ProvideNamedTransient registers a transient service by name.
func ProvideNamedTransient[T any](i do.Injector, name string, provider do.Provider[T]) {
	do.ProvideNamedTransient(i, name, provider)
}

// Invoke retrieves a service by type.
func Invoke[T any](i do.Injector) (T, error) {
	return do.Invoke[T](i)
}

// MustInvoke retrieves a service by type, panics on error.
func MustInvoke[T any](i do.Injector) T {
	return do.MustInvoke[T](i)
}

// InvokeNamed retrieves a service by name.
func InvokeNamed[T any](i do.Injector, name string) (T, error) {
	return do.InvokeNamed[T](i, name)
}

// MustInvokeNamed retrieves a service by name, panics on error.
func MustInvokeNamed[T any](i do.Injector, name string) T {
	return do.MustInvokeNamed[T](i, name)
}

// Override replaces an existing service by type.
func Override[T any](i do.Injector, provider do.Provider[T]) {
	do.Override(i, provider)
}

// OverrideNamed replaces an existing service by name.
func OverrideNamed[T any](i do.Injector, name string, provider do.Provider[T]) {
	do.OverrideNamed(i, name, provider)
}

// OverrideValue replaces an existing service with a value.
func OverrideValue[T any](i do.Injector, value T) {
	do.OverrideValue(i, value)
}

// OverrideNamedValue replaces an existing service by name with a value.
func OverrideNamedValue[T any](i do.Injector, name string, value T) {
	do.OverrideNamedValue(i, name, value)
}

// Shutdown gracefully shuts down all services in the container.
func Shutdown(i do.Injector) error {
	report := i.Shutdown()
	if report != nil && !report.Succeed {
		return report
	}
	return nil
}

// ShutdownWithContext gracefully shuts down all services with context support.
func ShutdownWithContext(ctx context.Context, i do.Injector) error {
	report := i.ShutdownWithContext(ctx)
	if report != nil && !report.Succeed {
		return report
	}
	return nil
}

// HealthCheck runs health checks on all services.
func HealthCheck(i do.Injector) map[string]error {
	return i.HealthCheck()
}

// HealthCheckWithContext runs health checks with context support.
func HealthCheckWithContext(ctx context.Context, i do.Injector) map[string]error {
	return i.HealthCheckWithContext(ctx)
}
