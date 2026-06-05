package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/techforge-hq/go-kit/logger"
)

func TestResolveCORSAllowedHeaders_Default(t *testing.T) {
	got := resolveCORSAllowedHeaders(nil)
	assert.Equal(t, "Origin, Content-Type, Accept, Authorization", got)
}

func TestResolveCORSAllowedHeaders_Custom(t *testing.T) {
	got := resolveCORSAllowedHeaders([]string{"X-CSRF-Token", "X-Request-ID"})
	assert.Equal(t, "X-CSRF-Token, X-Request-ID", got)
}

func TestCORSMiddleware_Permissive_DefaultHeaders(t *testing.T) {
	handler := corsMiddleware(nil, nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/sign-out", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin, Content-Type, Accept, Authorization", rec.Header().Get("Access-Control-Allow-Headers"))
}

func TestCORSMiddleware_Permissive_CustomHeaders(t *testing.T) {
	handler := corsMiddleware(nil, []string{"X-CSRF-Token", "X-Request-ID"})(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/sign-out", nil)
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "X-CSRF-Token, X-Request-ID", rec.Header().Get("Access-Control-Allow-Headers"))
}

func TestCORSMiddleware_Strict_DefaultHeaders(t *testing.T) {
	handler := corsMiddleware([]string{"https://app.example.com"}, nil)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/sign-out", nil)
	req.Header.Set("Origin", "https://app.example.com")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "https://app.example.com", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin, Content-Type, Accept, Authorization", rec.Header().Get("Access-Control-Allow-Headers"))
}

func TestCORSMiddleware_Strict_CustomHeaders(t *testing.T) {
	handler := corsMiddleware([]string{"https://app.example.com"}, []string{"X-CSRF-Token"})(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/sign-out", nil)
	req.Header.Set("Origin", "https://app.example.com")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "X-CSRF-Token", rec.Header().Get("Access-Control-Allow-Headers"))
}

func TestNewServer_CORS_CustomHeaders(t *testing.T) {
	s := NewServer(
		Config{
			Port:                0,
			ServiceName:         "test",
			CORSAllowedHeaders:  []string{"X-CSRF-Token"},
		},
		logger.NewNoop(),
		nil,
	)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/sign-out", nil)
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	rec := httptest.NewRecorder()
	s.http.Handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "X-CSRF-Token", rec.Header().Get("Access-Control-Allow-Headers"))
}
