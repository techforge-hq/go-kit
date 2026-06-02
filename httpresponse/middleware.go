package httpresponse

import (
	"errors"
	"net/http"
)

// HandlerFunc is an HTTP handler that can return an error.
// Use with ErrorMiddleware to automatically convert returned errors
// into RFC 9457 problem detail responses.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ErrorMiddlewareConfig holds configuration for the error middleware.
type ErrorMiddlewareConfig struct {
	Debug bool
}

// ErrorMiddleware returns an http.Handler that catches errors returned by
// HandlerFunc and converts them into RFC 9457 problem detail responses.
func ErrorMiddleware(config ErrorMiddlewareConfig, handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err == nil {
			return
		}

		handleError(w, r, err, config)
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error, config ErrorMiddlewareConfig) {
	var problem ProblemDetail
	if errors.As(err, &problem) {
		Problem(w, r, problem.WithInstance(r.URL.Path))
		return
	}

	p := ErrInternalServerError.WithInstance(r.URL.Path)
	if config.Debug {
		p = p.WithDetail(err.Error())
	}

	Problem(w, r, p)
}
