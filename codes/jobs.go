package codes

// Job-module wire codes returned by /api/v1/jobs/* endpoints. Mirrors the
// per-sentinel branches in modules/jobs/errors.go::handleJobError on the
// server.
const (
	JobNotFound          = "JOB_NOT_FOUND"
	JobExecutionNotFound = "JOB_EXECUTION_NOT_FOUND"
	JobOperationNotFound = "JOB_OPERATION_NOT_FOUND"

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
	JobAppRequired        = "JOB_APP_REQUIRED"
	JobAppNotFound        = "JOB_APP_NOT_FOUND"
	JobImageRequired      = "JOB_IMAGE_REQUIRED"
	JobConcurrencyUnsupported = "JOB_CONCURRENCY_UNSUPPORTED"

	JobInvalidPricingSlug = "JOB_INVALID_PRICING_SLUG"

	JobValidationFailed   = "JOB_VALIDATION_FAILED"
	JobUnauthorized       = "JOB_UNAUTHORIZED"
	JobInvalidID          = "JOB_INVALID_ID"
	JobInvalidRequestBody = "JOB_INVALID_REQUEST_BODY"
	JobInternalError      = "JOB_INTERNAL_ERROR"
)
