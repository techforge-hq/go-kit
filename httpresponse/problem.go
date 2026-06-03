// Package httpresponse implements RFC 9457 Problem Details for HTTP APIs.
package httpresponse

import "net/http"

const (
	// ContentTypeProblemJSON is the media type for RFC 9457 problem details.
	ContentTypeProblemJSON = "application/problem+json"
)

// ProblemDetail represents an RFC 9457 Problem Details object.
type ProblemDetail struct {
	// Type is a URI reference that identifies the problem type.
	Type string `json:"type"`

	// Title is a short, human-readable summary of the problem type.
	Title string `json:"title"`

	// Status is the HTTP status code.
	Status int `json:"status"`

	// Detail is a human-readable explanation specific to this occurrence.
	Detail string `json:"detail,omitempty"`

	// Instance is a URI reference that identifies the specific occurrence.
	Instance string `json:"instance,omitempty"`

	// Extensions allows additional members to be included.
	Extensions map[string]any `json:"extensions,omitempty"`
}

// ProblemDetailSignals is the global signals payload emitted for problem responses.
type ProblemDetailSignals struct {
	// Type is a URI reference that identifies the problem type.
	Type string `json:"type"`

	// Title is a short, human-readable summary of the problem type.
	Title string `json:"title"`

	// Status is the HTTP status code.
	Status int `json:"status"`

	// Detail is a human-readable explanation specific to this occurrence.
	Detail string `json:"detail,omitempty"`

	// Instance is a URI reference that identifies the specific occurrence.
	Instance string `json:"instance,omitempty"`
}

// NewProblemDetail creates a new ProblemDetail with the given type, title, and status.
func NewProblemDetail(problemType, title string, status int) ProblemDetail {
	return ProblemDetail{
		Type:   problemType,
		Title:  title,
		Status: status,
	}
}

// WithDetail adds a detail message to the problem.
func (p ProblemDetail) WithDetail(detail string) ProblemDetail {
	p.Detail = detail
	return p
}

// WithInstance adds an instance URI to the problem.
func (p ProblemDetail) WithInstance(instance string) ProblemDetail {
	p.Instance = instance
	return p
}

// WithExtension adds a custom extension field to the problem.
func (p ProblemDetail) WithExtension(key string, value any) ProblemDetail {
	if p.Extensions == nil {
		p.Extensions = make(map[string]any)
	}
	p.Extensions[key] = value
	return p
}

// Error implements the error interface.
func (p ProblemDetail) Error() string {
	if p.Detail != "" {
		return p.Detail
	}
	return p.Title
}

// Common problem types as URI references.
const (
	TypeBadRequest          = "https://httpstatuses.com/400"
	TypeUnauthorized        = "https://httpstatuses.com/401"
	TypeForbidden           = "https://httpstatuses.com/403"
	TypeNotFound            = "https://httpstatuses.com/404"
	TypeMethodNotAllowed    = "https://httpstatuses.com/405"
	TypeUnprocessableEntity = "https://httpstatuses.com/422"
	TypeConflict            = "https://httpstatuses.com/409"
	TypeInternalServerError = "https://httpstatuses.com/500"
	TypeServiceUnavailable  = "https://httpstatuses.com/503"
)

// Pre-defined problem details for common HTTP errors.
var (
	ErrBadRequest = NewProblemDetail(
		TypeBadRequest,
		"Bad Request",
		http.StatusBadRequest,
	)

	ErrUnauthorized = NewProblemDetail(
		TypeUnauthorized,
		"Unauthorized",
		http.StatusUnauthorized,
	)

	ErrForbidden = NewProblemDetail(
		TypeForbidden,
		"Forbidden",
		http.StatusForbidden,
	)

	ErrNotFound = NewProblemDetail(
		TypeNotFound,
		"Not Found",
		http.StatusNotFound,
	)

	ErrMethodNotAllowed = NewProblemDetail(
		TypeMethodNotAllowed,
		"Method Not Allowed",
		http.StatusMethodNotAllowed,
	)

	ErrConflict = NewProblemDetail(
		TypeConflict,
		"Conflict",
		http.StatusConflict,
	)

	ErrUnprocessableEntity = NewProblemDetail(
		TypeUnprocessableEntity,
		"Unprocessable Entity",
		http.StatusUnprocessableEntity,
	)

	ErrInternalServerError = NewProblemDetail(
		TypeInternalServerError,
		"Internal Server Error",
		http.StatusInternalServerError,
	)

	ErrServiceUnavailable = NewProblemDetail(
		TypeServiceUnavailable,
		"Service Unavailable",
		http.StatusServiceUnavailable,
	)
)

// ProblemTypeFromStatus returns the problem type URI for a given HTTP status code.
func ProblemTypeFromStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return TypeBadRequest
	case http.StatusUnauthorized:
		return TypeUnauthorized
	case http.StatusForbidden:
		return TypeForbidden
	case http.StatusNotFound:
		return TypeNotFound
	case http.StatusMethodNotAllowed:
		return TypeMethodNotAllowed
	case http.StatusConflict:
		return TypeConflict
	case http.StatusUnprocessableEntity:
		return TypeUnprocessableEntity
	case http.StatusServiceUnavailable:
		return TypeServiceUnavailable
	default:
		return TypeInternalServerError
	}
}
