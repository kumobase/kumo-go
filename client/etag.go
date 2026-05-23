package client

import (
	"net/http"
	"strings"
)

// readETag extracts the ETag header from a response, normalising to the
// W/"…" weak form the server emits. Returns "" if absent or malformed.
//
// Callers use the returned tag to pass back via IfMatch on a subsequent
// PATCH for optimistic-concurrency protection.
func readETag(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	tag := strings.TrimSpace(resp.Header.Get(http.CanonicalHeaderKey("ETag")))
	if tag == "" {
		return ""
	}
	// Server always emits weak ETags as W/"…"; accept strong ETags too in
	// case a future server change relaxes it.
	if !(strings.HasPrefix(tag, `W/"`) || strings.HasPrefix(tag, `"`)) {
		return ""
	}
	return tag
}
