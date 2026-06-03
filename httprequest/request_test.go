package httprequest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testPayload struct {
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
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			if tt.header != "" {
				r.Header.Set(HeaderDatastarRequest, tt.header)
			}

			assert.Equal(t, tt.want, IsDatastarRequest(r))
		})
	}
}

func TestRead_JSONBody(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"message":"hello","count":2}`))

	var payload testPayload
	require.NoError(t, Read(r, &payload))

	assert.Equal(t, "hello", payload.Message)
	assert.Equal(t, 2, payload.Count)
}

func TestRead_DatastarBody(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"message":"hello","count":2}`))
	r.Header.Set(HeaderDatastarRequest, "true")

	var payload testPayload
	require.NoError(t, Read(r, &payload))

	assert.Equal(t, "hello", payload.Message)
	assert.Equal(t, 2, payload.Count)
}

func TestRead_DatastarQuery(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/?datastar=%7B%22message%22%3A%22hello%22%2C%22count%22%3A2%7D", nil)
	r.Header.Set(HeaderDatastarRequest, "true")

	var payload testPayload
	require.NoError(t, Read(r, &payload))

	assert.Equal(t, "hello", payload.Message)
	assert.Equal(t, 2, payload.Count)
}

func TestRead_InvalidJSON(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{`))

	var payload testPayload
	require.Error(t, Read(r, &payload))
}
