package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/techforge-hq/go-kit/logger"
)

type failingChecker struct{}

func (failingChecker) HealthCheck(_ context.Context) error {
	return errors.New("connection refused")
}

func TestHandleHealth(t *testing.T) {
	s := &Server{
		config:   Config{Version: "0.1.0"},
		logger:   logger.NewNoop(),
		checkers: make(map[string]HealthChecker),
	}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	s.handleHealth(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response HealthResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	assert.Equal(t, healthStatusHealthy, response.Status)
}

func TestHandleHealth_Degraded(t *testing.T) {
	s := &Server{
		config: Config{Version: "0.1.0"},
		logger: logger.NewNoop(),
		checkers: map[string]HealthChecker{
			"database": failingChecker{},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	s.handleHealth(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)

	var response HealthResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	assert.Equal(t, healthStatusDegraded, response.Status)
	assert.Contains(t, response.Services["database"], healthStatusUnhealthy)
}

func TestHandlePing(t *testing.T) {
	s := &Server{
		logger: logger.NewNoop(),
	}

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	s.handlePing(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response PingResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	assert.Equal(t, "pong", response.Message)
	assert.NotEmpty(t, response.Time)
}

func TestHealthResponse_Structure(t *testing.T) {
	response := HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
		Time:    "1h30m45s",
		Services: map[string]string{
			"database": "healthy",
			"logger":   "healthy",
			"config":   "healthy",
		},
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	var unmarshaled HealthResponse
	require.NoError(t, json.Unmarshal(data, &unmarshaled))
	assert.Equal(t, response, unmarshaled)
}

func TestPingResponse_Structure(t *testing.T) {
	response := PingResponse{
		Message: "pong",
		Time:    "2023-12-14T10:30:00Z",
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	var unmarshaled PingResponse
	require.NoError(t, json.Unmarshal(data, &unmarshaled))
	assert.Equal(t, response, unmarshaled)
}

func TestHealthStatus_Constants(t *testing.T) {
	assert.Equal(t, "healthy", healthStatusHealthy)
	assert.Equal(t, "unhealthy", healthStatusUnhealthy)
	assert.Equal(t, "degraded", healthStatusDegraded)
}
