package codes

// Source-connection wire codes returned by /api/v1/source-connections/*
// endpoints. Mirrors the per-sentinel branches in
// modules/sourceconnect/errors.go::handleError on the server.
//
// Codes for the interactive install / OAuth-callback / inbound-webhook
// surfaces are intentionally NOT declared here: those endpoints are consumed
// by browsers and the git provider, never by an SDK client, so their error
// shapes are server-internal only.
const (
	// SourceConnectionNotFound — no connection with the given id exists, or
	// it isn't owned by the authenticated user.
	SourceConnectionNotFound = "SOURCE_CONNECTION_NOT_FOUND"

	// SourceConnectionForbidden — the connection exists but the caller may
	// not act on it.
	SourceConnectionForbidden = "SOURCE_CONNECTION_FORBIDDEN"

	// SourceConnectionSuspended — the provider-side installation is
	// suspended; reactivate it on the provider before retrying.
	SourceConnectionSuspended = "SOURCE_CONNECTION_SUSPENDED"

	// SourceProviderError — a call to the upstream git provider (GitHub)
	// failed. Usually transient; safe to retry.
	SourceProviderError = "SOURCE_PROVIDER_ERROR"

	// SourceConnectionInternalError — unexpected server-side failure.
	SourceConnectionInternalError = "SOURCE_CONNECTION_INTERNAL_ERROR"
)
