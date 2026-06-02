package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/techforge-hq/go-kit/httpresponse"
	"github.com/techforge-hq/go-kit/logger"
)

func TestNewServer_BaseRoutes(t *testing.T) {
	s := NewServer(Config{Port: 0, ServiceName: "test"}, logger.NewNoop(), nil)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	s.http.Handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestNewServer_CustomRouteWithErrorMiddleware(t *testing.T) {
	s := NewServer(Config{Port: 0, ServiceName: "test"}, logger.NewNoop(), func(mux *http.ServeMux) {
		mux.Handle("GET /api/items", httpresponse.ErrorMiddleware(httpresponse.ErrorMiddlewareConfig{}, func(w http.ResponseWriter, r *http.Request) error {
			return httpresponse.ErrNotFound.WithDetail("missing")
		}))
	})

	req := httptest.NewRequest(http.MethodGet, "/api/items", nil)
	rec := httptest.NewRecorder()
	s.http.Handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	require.Equal(t, httpresponse.ContentTypeProblemJSON, rec.Header().Get("Content-Type"))
}

func TestNewServer_DatastarProblemResponse(t *testing.T) {
	s := NewServer(Config{Port: 0, ServiceName: "test"}, logger.NewNoop(), func(mux *http.ServeMux) {
		mux.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
			httpresponse.BadRequest(w, r, "invalid password")
		})
	})

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.Header.Set(httpresponse.HeaderDatastarRequest, "true")
	rec := httptest.NewRecorder()
	s.http.Handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), "event: datastar-patch-signals")
	assert.Contains(t, rec.Body.String(), `"detail":"invalid password"`)
}

func TestRecoverMiddleware(t *testing.T) {
	log := logger.NewNoop()
	handler := recoverMiddleware(log, true)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	}))

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
