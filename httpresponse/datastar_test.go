package httpresponse

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testSignals struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
}

func TestIsDatastarRequest(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   bool
	}{
		{name: "missing header", want: false},
		{name: "true header", header: "true", want: true},
		{name: "mixed case true header", header: "TRUE", want: true},
		{name: "false header", header: "false", want: false},
		{name: "non-empty header", header: "1", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				r.Header.Set(HeaderDatastarRequest, tt.header)
			}

			assert.Equal(t, tt.want, IsDatastarRequest(r))
		})
	}
}

func TestReadDatastarSignals_Get(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/?datastar=%7B%22message%22%3A%22hello%22%2C%22count%22%3A2%7D", nil)

	var signals testSignals
	require.NoError(t, ReadDatastarSignals(r, &signals))

	assert.Equal(t, "hello", signals.Message)
	assert.Equal(t, 2, signals.Count)
}

func TestReadDatastarSignals_Post(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"message":"hello","count":2}`))

	var signals testSignals
	require.NoError(t, ReadDatastarSignals(r, &signals))

	assert.Equal(t, "hello", signals.Message)
	assert.Equal(t, 2, signals.Count)
}

func TestSignals_JSONFallback(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	require.NoError(t, Signals(w, r, testSignals{Message: "hello", Count: 2}))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp Response
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	data, ok := resp.Data.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "hello", data["message"])
	assert.Equal(t, float64(2), data["count"])
}

func TestSignals_DatastarRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(HeaderDatastarRequest, "true")

	require.NoError(t, Signals(w, r, testSignals{Message: "hello", Count: 2}))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/event-stream", w.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
	assert.Contains(t, w.Body.String(), "event: datastar-patch-signals")
	assert.Contains(t, w.Body.String(), `data: signals {"message":"hello","count":2}`)
}
