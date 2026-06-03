package httprequest

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

// Read decodes an HTTP request into target.
//
// Datastar requests are decoded with datastar.ReadSignals. Other requests are
// decoded as JSON from the request body.
func Read(r *http.Request, target any) error {
	if IsDatastarRequest(r) {
		if err := datastar.ReadSignals(r, target); err != nil {
			return fmt.Errorf("read datastar signals: %w", err)
		}
		return nil
	}

	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return fmt.Errorf("decode json request body: %w", err)
	}
	return nil
}
