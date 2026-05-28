package di

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type greeter struct {
	msg string
}

func TestProvideValueAndInvoke(t *testing.T) {
	scope := New()
	ProvideValue(scope, greeter{msg: "hello"})

	got, err := Invoke[greeter](scope)
	require.NoError(t, err)
	assert.Equal(t, "hello", got.msg)
}

func TestProvideAndMustInvoke(t *testing.T) {
	scope := New()
	Provide(scope, func(Injector) (string, error) {
		return "world", nil
	})

	assert.Equal(t, "world", MustInvoke[string](scope))
}

func TestShutdown(t *testing.T) {
	scope := New()
	ProvideValue(scope, 42)

	err := Shutdown(scope)
	assert.NoError(t, err)
}

func TestShutdownWithContext(t *testing.T) {
	scope := New()
	ProvideValue(scope, true)

	err := ShutdownWithContext(context.Background(), scope)
	assert.NoError(t, err)
}

func TestHealthCheck(t *testing.T) {
	scope := New()
	ProvideValue(scope, struct{}{})

	errs := HealthCheck(scope)
	require.Len(t, errs, 1)
	for name, err := range errs {
		assert.NoError(t, err, "service %s", name)
	}
}
