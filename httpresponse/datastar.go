package httpresponse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/starfederation/datastar-go/datastar"
)

const (
	// HeaderDatastarRequest identifies requests issued by Datastar actions.
	HeaderDatastarRequest = "Datastar-Request"
)

// IsDatastarRequest reports whether the request was issued by a Datastar action.
func IsDatastarRequest(r *http.Request) bool {
	if r == nil {
		return false
	}

	values, ok := r.Header[http.CanonicalHeaderKey(HeaderDatastarRequest)]
	if !ok {
		return false
	}
	if len(values) == 0 {
		return true
	}

	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		return !strings.EqualFold(value, "false")
	}

	return true
}

// ReadDatastarSignals extracts Datastar signals from the request.
func ReadDatastarSignals(r *http.Request, signals any) error {
	if err := datastar.ReadSignals(r, signals); err != nil {
		return fmt.Errorf("read datastar signals: %w", err)
	}
	return nil
}

// Signals sends signals to Datastar callers and JSON to regular HTTP callers.
func Signals(w http.ResponseWriter, r *http.Request, signals any, opts ...datastar.PatchSignalsOption) error {
	if !IsDatastarRequest(r) {
		OK(w, signals)
		return nil
	}

	b, err := json.Marshal(signals)
	if err != nil {
		return fmt.Errorf("marshal datastar signals: %w", err)
	}

	sse := datastar.NewSSE(w, r)
	if err := sse.PatchSignals(b, opts...); err != nil {
		return fmt.Errorf("patch datastar signals: %w", err)
	}

	return nil
}
