package httpresponse

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOK(t *testing.T) {
	w := httptest.NewRecorder()
	OK(w, map[string]string{"name": "test"})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp Response
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	data, ok := resp.Data.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "test", data["name"])
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	Created(w, map[string]string{"id": "123"})

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp Response
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	data, ok := resp.Data.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "123", data["id"])
}

func TestNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	NoContent(w)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestWithMeta(t *testing.T) {
	w := httptest.NewRecorder()
	data := []string{"a", "b"}
	meta := map[string]int{"total": 2}
	WithMeta(w, http.StatusOK, data, meta)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotNil(t, resp.Data)
	assert.NotNil(t, resp.Meta)
}

func TestProblemResponse(t *testing.T) {
	w := httptest.NewRecorder()
	problem := ErrBadRequest.WithDetail("invalid input").WithInstance("/api/users")
	Problem(w, problem)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, ContentTypeProblemJSON, w.Header().Get("Content-Type"))

	var pd ProblemDetail
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &pd))
	assert.Equal(t, TypeBadRequest, pd.Type)
	assert.Equal(t, "Bad Request", pd.Title)
	assert.Equal(t, http.StatusBadRequest, pd.Status)
	assert.Equal(t, "invalid input", pd.Detail)
	assert.Equal(t, "/api/users", pd.Instance)
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	BadRequest(w, r, "missing field")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var pd ProblemDetail
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &pd))
	assert.Equal(t, "missing field", pd.Detail)
	assert.Equal(t, "/api/users", pd.Instance)
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/users/1", nil)
	NotFound(w, r, "user not found")

	assert.Equal(t, http.StatusNotFound, w.Code)

	var pd ProblemDetail
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &pd))
	assert.Equal(t, "user not found", pd.Detail)
}

func TestUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/secret", nil)
	Unauthorized(w, r, "invalid token")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/admin", nil)
	Forbidden(w, r, "insufficient permissions")

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestConflict(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/users", nil)
	Conflict(w, r, "email already exists")

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestUnprocessableEntity(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/users", nil)
	UnprocessableEntity(w, r, "invalid email format")

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestInternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	InternalServerError(w, r, "database connection failed")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
