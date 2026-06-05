package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/techforge-hq/go-kit/logger"
)

// Config holds server configuration.
type Config struct {
	Port        int
	ServiceName string
	Version     string
	Debug       bool
	// CORSAllowedOrigins lists browser origins allowed for credentialed cross-origin API calls.
	// Empty enables permissive CORS (fine for local dev).
	CORSAllowedOrigins []string
	// CORSAllowedHeaders lists headers permitted in Access-Control-Allow-Headers for
	// preflight responses. Empty uses defaultCORSAllowedHeaders. Add custom auth/CSRF
	// headers here (e.g. "X-CSRF-Token", "X-Request-ID") when consumers send them.
	CORSAllowedHeaders []string
}

// HealthChecker reports dependency health for /health.
type HealthChecker interface {
	HealthCheck(ctx context.Context) error
}

// RouteFunc registers application routes on the server mux.
type RouteFunc func(mux *http.ServeMux)

// Server is an HTTP server with health endpoints and shared middleware.
type Server struct {
	config   Config
	logger   logger.Logger
	checkers map[string]HealthChecker
	Mux      *http.ServeMux
	http     *http.Server
}

// NewServer creates a configured HTTP server. Pass nil for registerRoutes if none are needed.
func NewServer(config Config, log logger.Logger, registerRoutes RouteFunc) *Server {
	if config.Version == "" {
		config.Version = "0.1.0"
	}

	mux := http.NewServeMux()
	s := &Server{
		config:   config,
		logger:   log.With("component", "server"),
		checkers: make(map[string]HealthChecker),
		Mux:      mux,
	}

	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /ping", s.handlePing)

	if registerRoutes != nil {
		registerRoutes(mux)
	}

	handler := chain(
		mux,
		corsMiddleware(config.CORSAllowedOrigins, config.CORSAllowedHeaders),
		requestLogMiddleware(s.logger),
		recoverMiddleware(s.logger, config.Debug),
	)

	s.http = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: handler,
	}

	s.logger.Info("HTTP server initialized",
		"port", config.Port,
		"service", config.ServiceName,
	)

	return s
}

// RegisterHealthChecker adds a named dependency health check.
func (s *Server) RegisterHealthChecker(name string, checker HealthChecker) {
	s.checkers[name] = checker
}

// Start listens until ctx is cancelled, then shuts down gracefully.
func (s *Server) Start(ctx context.Context) error {
	s.logger.WithContext(ctx).Info("starting HTTP server",
		"address", s.http.Addr,
		"service", s.config.ServiceName,
	)

	errCh := make(chan error, 1)
	go func() {
		if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			s.logger.WithContext(ctx).Error("server error", "error", err)
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("server listen failed: %w", err)
		}
	}

	s.logger.WithContext(ctx).Info("shutting down HTTP server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.http.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.WithContext(ctx).Info("shutting down HTTP server")

	if err := s.http.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}
