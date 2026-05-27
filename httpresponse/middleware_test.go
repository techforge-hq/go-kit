package httpresponse

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorMiddleware_NoError(t *testing.T) {
	handler := ErrorMiddleware(ErrorMiddlewareConfig{}, func(w http.ResponseWriter, r *http.Request) error {
		OK(w, "success")
		return nil
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestErrorMiddleware_ProblemDetailError(t *testing.T) {
	handler := ErrorMiddleware(ErrorMiddlewareConfig{}, func(w http.ResponseWriter, r *http.Request) error {
		return ErrNotFound.WithDetail("resource not found")
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/users/1", nil)
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, ContentTypeProblemJSON, w.Header().Get("Content-Type"))

	var pd ProblemDetail
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &pd))
	assert.Equal(t, "resource not found", pd.Detail)
	assert.Equal(t, "/api/users/1", pd.Instance)
}

func TestErrorMiddleware_GenericError_DebugOff(t *testing.T) {
	handler := ErrorMiddleware(ErrorMiddlewareConfig{Debug: false}, func(w http.ResponseWriter, r *http.Request) error {
		return errors.New("something went wrong")
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var pd ProblemDetail
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &pd))
	assert.Empty(t, pd.Detail)
	assert.Equal(t, "Internal Server Error", pd.Title)
}

func TestErrorMiddleware_GenericError_DebugOn(t *testing.T) {
	handler := ErrorMiddleware(ErrorMiddlewareConfig{Debug: true}, func(w http.ResponseWriter, r *http.Request) error {
		return errors.New("database connection timeout")
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var pd ProblemDetail
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &pd))
	assert.Equal(t, "database connection timeout", pd.Detail)
}
