package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

const (
	healthStatusHealthy   = "healthy"
	healthStatusUnhealthy = "unhealthy"
	healthStatusDegraded  = "degraded"
)

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status   string            `json:"status"`
	Version  string            `json:"version"`
	Time     string            `json:"time"`
	Services map[string]string `json:"services"`
}

// PingResponse represents the ping response.
type PingResponse struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	services := make(map[string]string)

	for name, checker := range s.checkers {
		if err := checker.HealthCheck(ctx); err != nil {
			services[name] = healthStatusUnhealthy + ": " + err.Error()
		} else {
			services[name] = healthStatusHealthy
		}
	}

	status := healthStatusHealthy
	for _, svcStatus := range services {
		if svcStatus != healthStatusHealthy {
			status = healthStatusDegraded
			break
		}
	}

	response := HealthResponse{
		Status:   status,
		Version:  s.config.Version,
		Time:     time.Now().UTC().Format(time.RFC3339),
		Services: services,
	}

	httpStatus := http.StatusOK
	if status != healthStatusHealthy {
		httpStatus = http.StatusServiceUnavailable
	}

	writeJSON(w, httpStatus, response)
}

func (s *Server) handlePing(w http.ResponseWriter, _ *http.Request) {
	response := PingResponse{
		Message: "pong",
		Time:    time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
