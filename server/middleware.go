package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/techforge-hq/go-kit/httpresponse"
	"github.com/techforge-hq/go-kit/logger"
)

var defaultCORSAllowedHeaders = []string{
	"Origin",
	"Content-Type",
	"Accept",
	"Authorization",
}

func resolveCORSAllowedHeaders(custom []string) string {
	if len(custom) == 0 {
		return strings.Join(defaultCORSAllowedHeaders, ", ")
	}
	return strings.Join(custom, ", ")
}

func chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func recoverMiddleware(log logger.Logger, debug bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic recovered",
						"method", r.Method,
						"uri", r.URL.RequestURI(),
						"panic", rec,
					)

					p := httpresponse.ErrInternalServerError.WithInstance(r.URL.Path)
					if debug {
						p = p.WithDetail(fmtRecover(rec))
					}
					httpresponse.Problem(w, r, p)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func fmtRecover(rec any) string {
	switch v := rec.(type) {
	case error:
		return v.Error()
	case string:
		return v
	default:
		return http.StatusText(http.StatusInternalServerError)
	}
}

func requestLogMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(rw, r)

			elapsed := time.Since(start)
			if rw.status >= http.StatusInternalServerError {
				log.Error("request failed",
					"method", r.Method,
					"uri", r.URL.RequestURI(),
					"status", rw.status,
					"latency", elapsed,
				)
				return
			}

			log.Info("request completed",
				"method", r.Method,
				"uri", r.URL.RequestURI(),
				"status", rw.status,
				"latency", elapsed,
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Flush() {
	if flusher, ok := rw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

func corsMiddleware(allowedOrigins, allowedHeaders []string) func(http.Handler) http.Handler {
	headers := resolveCORSAllowedHeaders(allowedHeaders)
	if len(allowedOrigins) == 0 {
		return permissiveCORSMiddleware(headers)
	}
	return strictCORSMiddleware(allowedOrigins, headers)
}

func permissiveCORSMiddleware(allowedHeaders string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, PUT, PATCH, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func strictCORSMiddleware(allowedOrigins []string, allowedHeaders string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowed[origin] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				if _, ok := allowed[origin]; ok {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Vary", "Origin")
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, PUT, PATCH, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
