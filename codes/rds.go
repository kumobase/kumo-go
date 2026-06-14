package codes

// RDS-module wire codes returned by /api/v1/rds/* endpoints.
//
// RDS (Relational Database Service) is Kumo's managed-database umbrella; the
// first engine is PostgreSQL. Codes here are engine-agnostic so MySQL and
// other engines can reuse them.
const (
	// RDSInstanceNotFound — no database instance with the given id/name exists
	// in the caller's scope (also covers cross-tenant access attempts).
	RDSInstanceNotFound = "RDS_INSTANCE_NOT_FOUND"

	// RDSFlavorNotFound — the requested flavor slug does not exist or is not
	// active in the catalogue.
	RDSFlavorNotFound = "RDS_FLAVOR_NOT_FOUND"

	// RDSFlavorDisabled — the flavor exists but is no longer offered for new
	// instances (running instances keep their pinned flavor version).
	RDSFlavorDisabled = "RDS_FLAVOR_DISABLED"

	// RDSEngineNotSupported — the requested engine/version is not offered.
	RDSEngineNotSupported = "RDS_ENGINE_NOT_SUPPORTED"

	// RDSActionInProgress — a lifecycle action (provision/scale/resize/delete)
	// is already running on this instance; the new request was rejected. Poll
	// the instance until status leaves its transient state, then retry.
	RDSActionInProgress = "RDS_ACTION_IN_PROGRESS"

	// RDSInstanceNotReady — the instance exists but is not yet in a state that
	// supports the requested operation (e.g. connection info requested before
	// the credentials secret/endpoint is published).
	RDSInstanceNotReady = "RDS_INSTANCE_NOT_READY"

	// RDSInstanceNotSuspended — a start (resume) was requested on an instance
	// that is not in the suspended state. Only a suspended database can be
	// started.
	RDSInstanceNotSuspended = "RDS_INSTANCE_NOT_SUSPENDED"

	// RDSInvalidStorageSize — the requested storage size is outside the
	// flavor/tier bounds, or a resize attempted to shrink (not allowed).
	RDSInvalidStorageSize = "RDS_INVALID_STORAGE_SIZE"

	// RDSOperationNotFound — the operation_id supplied to the polling endpoint
	// does not exist in the caller's scope.
	RDSOperationNotFound = "RDS_OPERATION_NOT_FOUND"

	// RDSInsufficientBalance — the caller's wallet cannot cover the minimum
	// up-front cost (≈1h compute + storage) of the requested instance.
	RDSInsufficientBalance = "RDS_INSUFFICIENT_BALANCE"

	// RDSUnauthorized — the request lacked a valid session/credential.
	RDSUnauthorized = "RDS_UNAUTHORIZED"

	// RDSInvalidRequestBody — the JSON body could not be parsed.
	RDSInvalidRequestBody = "RDS_INVALID_REQUEST_BODY"

	// RDSValidationError — request-body validation failed; Data carries a
	// ValidationErrorsResponse.
	RDSValidationError = "RDS_VALIDATION_ERROR"

	// RDSInvalidInstanceID — the path id-or-name segment was malformed.
	RDSInvalidInstanceID = "RDS_INVALID_INSTANCE_ID"

	// RDSInvalidPagination — page / page_size were out of range.
	RDSInvalidPagination = "RDS_INVALID_PAGINATION"

	// RDSInvalidStatusFilter — the status list filter held an unknown value.
	RDSInvalidStatusFilter = "RDS_INVALID_STATUS_FILTER"

	// RDSInternalError — unexpected server-side failure. Safe to retry.
	RDSInternalError = "RDS_INTERNAL_ERROR"
)
