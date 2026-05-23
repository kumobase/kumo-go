package codes

// VPS-module wire codes returned by /api/v1/vps/* endpoints.
const (
	InstanceNotFound          = "INSTANCE_NOT_FOUND"
	InstanceExpired           = "INSTANCE_EXPIRED"
	InstanceNotRunning        = "INSTANCE_NOT_RUNNING"
	SSHNotReady               = "SSH_NOT_READY"
	AutoRenewAlreadyCancelled = "AUTO_RENEW_ALREADY_CANCELLED"
	ActionInProgress          = "ACTION_IN_PROGRESS"
	ActionQueued              = "ACTION_QUEUED"

	ProviderNotFound  = "PROVIDER_NOT_FOUND"
	PlanNotFound      = "PLAN_NOT_FOUND"
	ProviderDisabled  = "PROVIDER_DISABLED"
	PlanDisabled      = "PLAN_DISABLED"
	InvalidRegion     = "INVALID_REGION"
	MissingRegion     = "MISSING_REGION"

	InsufficientBalance      = "INSUFFICIENT_BALANCE"
	ProviderBalanceIssue     = "PROVIDER_BALANCE_ISSUE"
	PlatformCapacityExceeded = "PLATFORM_CAPACITY_EXCEEDED"
	ServiceUnavailable       = "SERVICE_UNAVAILABLE"
	QuotaExceeded            = "QUOTA_EXCEEDED"

	VPSUnauthorized       = "UNAUTHORIZED"
	VPSInvalidRequestBody = "INVALID_REQUEST_BODY"
	VPSValidationError    = "VALIDATION_ERROR"
	VPSInvalidServerID    = "INVALID_SERVER_ID"
	VPSInvalidPagination  = "INVALID_PAGINATION"
	VPSInvalidStatusFilter = "INVALID_STATUS_FILTER"
	VPSInvalidTimeFilter  = "INVALID_TIME_FILTER"
	VPSInternalError      = "INTERNAL_ERROR"
)
