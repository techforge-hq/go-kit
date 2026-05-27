package httpresponse

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProblemDetail(t *testing.T) {
	pd := NewProblemDetail(TypeBadRequest, "Bad Request", http.StatusBadRequest)

	assert.Equal(t, TypeBadRequest, pd.Type)
	assert.Equal(t, "Bad Request", pd.Title)
	assert.Equal(t, http.StatusBadRequest, pd.Status)
	assert.Empty(t, pd.Detail)
	assert.Empty(t, pd.Instance)
	assert.Nil(t, pd.Extensions)
}

func TestProblemDetail_WithDetail(t *testing.T) {
	pd := ErrBadRequest.WithDetail("invalid input")

	assert.Equal(t, "invalid input", pd.Detail)
	assert.Empty(t, ErrBadRequest.Detail)
}

func TestProblemDetail_WithInstance(t *testing.T) {
	pd := ErrNotFound.WithInstance("/api/users/1")

	assert.Equal(t, "/api/users/1", pd.Instance)
	assert.Empty(t, ErrNotFound.Instance)
}

func TestProblemDetail_WithExtension(t *testing.T) {
	pd := ErrUnprocessableEntity.
		WithExtension("field", "email").
		WithExtension("reason", "invalid format")

	assert.Equal(t, "email", pd.Extensions["field"])
	assert.Equal(t, "invalid format", pd.Extensions["reason"])
	assert.Nil(t, ErrUnprocessableEntity.Extensions)
}

func TestProblemDetail_Error(t *testing.T) {
	t.Run("with detail", func(t *testing.T) {
		pd := ErrBadRequest.WithDetail("something specific")
		assert.Equal(t, "something specific", pd.Error())
	})

	t.Run("without detail", func(t *testing.T) {
		pd := NewProblemDetail(TypeBadRequest, "Bad Request", http.StatusBadRequest)
		assert.Equal(t, "Bad Request", pd.Error())
	})
}

func TestProblemDetail_Chaining(t *testing.T) {
	pd := ErrNotFound.
		WithDetail("user not found").
		WithInstance("/api/users/42").
		WithExtension("user_id", 42)

	assert.Equal(t, TypeNotFound, pd.Type)
	assert.Equal(t, "Not Found", pd.Title)
	assert.Equal(t, http.StatusNotFound, pd.Status)
	assert.Equal(t, "user not found", pd.Detail)
	assert.Equal(t, "/api/users/42", pd.Instance)
	assert.Equal(t, 42, pd.Extensions["user_id"])
}

func TestProblemTypeFromStatus(t *testing.T) {
	tests := []struct {
		status   int
		expected string
	}{
		{http.StatusBadRequest, TypeBadRequest},
		{http.StatusUnauthorized, TypeUnauthorized},
		{http.StatusForbidden, TypeForbidden},
		{http.StatusNotFound, TypeNotFound},
		{http.StatusMethodNotAllowed, TypeMethodNotAllowed},
		{http.StatusConflict, TypeConflict},
		{http.StatusUnprocessableEntity, TypeUnprocessableEntity},
		{http.StatusServiceUnavailable, TypeServiceUnavailable},
		{http.StatusInternalServerError, TypeInternalServerError},
		{http.StatusTeapot, TypeInternalServerError}, // unknown status defaults
	}

	for _, tt := range tests {
		t.Run(http.StatusText(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.expected, ProblemTypeFromStatus(tt.status))
		})
	}
}
