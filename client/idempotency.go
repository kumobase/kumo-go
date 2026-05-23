package client

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

// writeOpts is the resolved configuration for a single write call
// (POST/PATCH/DELETE). Populated by applying WriteOption funcs in order.
type writeOpts struct {
	idempotencyKey string // explicit override; empty means "auto-generate"
	ifMatch        string // ETag value to send as If-Match; empty means "no header"
}

// WriteOption configures a single write call. Pass to the per-resource
// service method:
//
//	app, err := client.Apps().Create(ctx, &req,
//	    client.WithIdempotencyKey("tf-resource-abc"),
//	)
type WriteOption func(*writeOpts)

// WithIdempotencyKey overrides the SDK's auto-generated Idempotency-Key for
// this call. Use it when your own pipeline owns idempotency identity —
// e.g. Terraform binding the key to a resource's stable id, or a cron job
// that wants to dedupe on a window timestamp.
//
// IMPORTANT: when you supply your own key you also own retry semantics for
// that call. Calling Create twice with the same key + same body replays
// the cached response. Same key + DIFFERENT body returns 409
// IDEMPOTENCY_KEY_CONFLICT (do not retry — it's a caller bug).
func WithIdempotencyKey(key string) WriteOption {
	return func(o *writeOpts) { o.idempotencyKey = key }
}

// IfMatch sends If-Match: <etag> with the request. The server returns 412
// ETAG_MISMATCH if the resource has been modified since the caller
// captured the tag — surface that as ErrETagMismatch via errors.Is.
//
// Empty etag is treated as "no header" (omit it from the request rather
// than sending If-Match: "").
func IfMatch(etag string) WriteOption {
	return func(o *writeOpts) { o.ifMatch = etag }
}

// resolveWriteOpts applies opts and fills in defaults (auto-generated
// Idempotency-Key when no override was supplied).
func resolveWriteOpts(opts []WriteOption) (writeOpts, error) {
	out := writeOpts{}
	for _, opt := range opts {
		opt(&out)
	}
	if out.idempotencyKey == "" {
		key, err := newIdempotencyKey()
		if err != nil {
			return writeOpts{}, err
		}
		out.idempotencyKey = key
	}
	return out, nil
}

// newIdempotencyKey returns a fresh random 128-bit hex string. Equivalent
// in entropy to UUIDv4 but avoids the github.com/google/uuid dep. Hex form
// is friendly to log greps and matches what curl users typically generate.
func newIdempotencyKey() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", errors.Join(errors.New("kumo: failed to generate idempotency key"), err)
	}
	return hex.EncodeToString(b[:]), nil
}
