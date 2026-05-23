// Package codes contains the stable wire-protocol error codes returned in
// the Code field of every Kumo API error response.
//
// Consumers should branch on these constants rather than parsing the
// human-readable Message field (which may evolve between releases). New
// codes are added under minor SDK bumps; existing codes are never renamed
// or removed (would be a wire-breaking change).
package codes

// Cross-cutting codes emitted by shared pkg/common helpers in the server.
// These can appear on any endpoint, so the module-specific code packages
// (codes/apps, codes/secrets, …) intentionally do not redeclare them.
const (
	// IdempotencyKeyConflict — the supplied Idempotency-Key matched an
	// existing record whose request body differs from this attempt. The
	// client either reused the key by mistake or changed the body between
	// retries. Choose a new key.
	IdempotencyKeyConflict = "IDEMPOTENCY_KEY_CONFLICT"

	// IdempotencyInProgress — the supplied Idempotency-Key was registered by
	// a previous attempt that is still running on the server. Retry after a
	// short delay (200..2000ms).
	IdempotencyInProgress = "IDEMPOTENCY_IN_PROGRESS"

	// ETagMismatch — the If-Match header on a PATCH did not match the
	// resource's current ETag. Re-fetch the resource and retry.
	ETagMismatch = "ETAG_MISMATCH"

	// ValidationFailed — generic request-body validation rejection. The
	// Data field carries a ValidationErrorsResponse listing field-level
	// failures.
	ValidationFailed = "VALIDATION_FAILED"

	// InvalidFilterCombination — mutually exclusive list filters were
	// supplied together (e.g. both app_id and attached on /volumes).
	InvalidFilterCombination = "INVALID_FILTER_COMBINATION"
)
