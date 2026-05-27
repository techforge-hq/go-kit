package httpresponse

import (
	"encoding/json"
	"net/http"
)

// Response represents a successful API response.
type Response struct {
	Data any `json:"data,omitempty"`
	Meta any `json:"meta,omitempty"`
}

// JSON writes a JSON response with the given status code and payload.
func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// OK sends a 200 OK response with the given data.
func OK(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, Response{Data: data})
}

// Created sends a 201 Created response with the given data.
func Created(w http.ResponseWriter, data any) {
	JSON(w, http.StatusCreated, Response{Data: data})
}

// NoContent sends a 204 No Content response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// WithMeta sends a JSON response with data and metadata.
func WithMeta(w http.ResponseWriter, status int, data, meta any) {
	JSON(w, status, Response{Data: data, Meta: meta})
}

// Problem sends an RFC 9457 problem detail response.
func Problem(w http.ResponseWriter, problem ProblemDetail) {
	w.Header().Set("Content-Type", ContentTypeProblemJSON)
	w.WriteHeader(problem.Status)
	json.NewEncoder(w).Encode(problem)
}

// BadRequest sends a 400 Bad Request problem response.
func BadRequest(w http.ResponseWriter, r *http.Request, detail string) {
	Problem(w, ErrBadRequest.WithDetail(detail).WithInstance(r.URL.Path))
}

// NotFound sends a 404 Not Found problem response.
func NotFound(w http.ResponseWriter, r *http.Request, detail string) {
	Problem(w, ErrNotFound.WithDetail(detail).WithInstance(r.URL.Path))
}

// Unauthorized sends a 401 Unauthorized problem response.
func Unauthorized(w http.ResponseWriter, r *http.Request, detail string) {
	Problem(w, ErrUnauthorized.WithDetail(detail).WithInstance(r.URL.Path))
}

// Forbidden sends a 403 Forbidden problem response.
func Forbidden(w http.ResponseWriter, r *http.Request, detail string) {
	Problem(w, ErrForbidden.WithDetail(detail).WithInstance(r.URL.Path))
}

// Conflict sends a 409 Conflict problem response.
func Conflict(w http.ResponseWriter, r *http.Request, detail string) {
	Problem(w, ErrConflict.WithDetail(detail).WithInstance(r.URL.Path))
}

// UnprocessableEntity sends a 422 Unprocessable Entity problem response.
func UnprocessableEntity(w http.ResponseWriter, r *http.Request, detail string) {
	Problem(w, ErrUnprocessableEntity.WithDetail(detail).WithInstance(r.URL.Path))
}

// InternalServerError sends a 500 Internal Server Error problem response.
func InternalServerError(w http.ResponseWriter, r *http.Request, detail string) {
	Problem(w, ErrInternalServerError.WithDetail(detail).WithInstance(r.URL.Path))
}
