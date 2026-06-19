package codes

// Job-module wire codes returned by /api/v1/jobs/* endpoints. Mirrors the
// per-sentinel branches in modules/jobs/errors.go::handleJobError on the
// server.
const (
	JobNotFound          = "JOB_NOT_FOUND"
	JobExecutionNotFound = "JOB_EXECUTION_NOT_FOUND"
	JobOperationNotFound = "JOB_OPERATION_NOT_FOUND"

	// JobExecutionExpired is returned (410 Gone) when an execution exists but is
	// older than JobLogsRetentionHours, so it has aged out of the user-facing
	// list/get surface. Distinct from JobExecutionNotFound so clients can tell
	// "aged out" from "never existed". The row is retained in the DB for billing.
	JobExecutionExpired = "JOB_EXECUTION_EXPIRED"

	JobDeploymentInProgress = "JOB_DEPLOYMENT_IN_PROGRESS"
	JobAlreadySuspended     = "JOB_ALREADY_SUSPENDED"
	JobNotSuspended         = "JOB_NOT_SUSPENDED"
	JobQuotaExceeded        = "JOB_QUOTA_EXCEEDED"
	JobInsufficientBalance  = "JOB_INSUFFICIENT_BALANCE"

	// Validation flavours separate from the generic VALIDATION_FAILED so
	// clients (Terraform, CLI) can branch on the precise field at fault.
	JobScheduleInvalid    = "JOB_SCHEDULE_INVALID"
	JobScheduleTooFrequent = "JOB_SCHEDULE_TOO_FREQUENT"
	JobTimezoneInvalid    = "JOB_TIMEZONE_INVALID"
	JobKindInvalid        = "JOB_KIND_INVALID"

	// JobKindUnsupported is returned (400) when a create requests a job kind
	// that is recognised but currently disabled for new jobs. Distinct from
	// JobKindInvalid (an unrecognised/malformed kind) so clients can tell
	// "this kind no longer accepts new jobs" from "this kind doesn't exist".
	// app_attached is disabled for creation; existing app_attached jobs remain
	// fully manageable.
	JobKindUnsupported    = "JOB_KIND_UNSUPPORTED"

	JobAppRequired        = "JOB_APP_REQUIRED"
	JobAppNotFound        = "JOB_APP_NOT_FOUND"
	JobImageRequired      = "JOB_IMAGE_REQUIRED"
	JobConcurrencyUnsupported = "JOB_CONCURRENCY_UNSUPPORTED"

	// Standalone-job image validation (registry manifest lookup at create),
	// mirroring the apps image check. A malformed reference reuses
	// JobValidationFailed.
	JobImageNotFound            = "JOB_IMAGE_NOT_FOUND"
	JobImageUnauthorized        = "JOB_IMAGE_UNAUTHORIZED"
	JobImageRegistryUnreachable = "JOB_IMAGE_REGISTRY_UNREACHABLE"

	JobInvalidPricingSlug = "JOB_INVALID_PRICING_SLUG"

	JobValidationFailed   = "JOB_VALIDATION_FAILED"
	JobUnauthorized       = "JOB_UNAUTHORIZED"
	JobInvalidID          = "JOB_INVALID_ID"
	JobInvalidRequestBody = "JOB_INVALID_REQUEST_BODY"
	JobInternalError      = "JOB_INTERNAL_ERROR"

	// Per-execution observability (metrics/logs) internal-error fallbacks. The
	// transport/validation failures reuse the cross-cutting codes
	// INVALID_TIME_RANGE / INVALID_LOG_FILTER / METRICS_BACKEND_UNAVAILABLE /
	// LOGS_BACKEND_UNAVAILABLE; these cover the unmapped 500 default.
	JobMetricsInternalError = "JOB_METRICS_INTERNAL_ERROR"
	JobLogsInternalError    = "JOB_LOGS_INTERNAL_ERROR"

	// JobInvalidStatsWindow is returned by GET /api/v1/jobs/:id/stats when the
	// from/to window or granularity is invalid (unparseable bound, from >= to,
	// unknown granularity, or a window that would produce too many buckets).
	JobInvalidStatsWindow = "JOB_INVALID_STATS_WINDOW"

	// JobOperationFailed is the error_code persisted on a job_operations row
	// (and surfaced when polling the operation) when an async deployment
	// operation fails. It is an outcome code on the operation record rather
	// than an HTTP-response code.
	JobOperationFailed = "JOB_OPERATION_FAILED"
)
