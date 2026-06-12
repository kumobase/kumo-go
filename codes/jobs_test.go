package codes

import "testing"

// Wire codes are a public contract. These assert the exact string values
// so an accidental rename is caught here before release.
func TestJobsCodeValues(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"JobNotFound", JobNotFound, "JOB_NOT_FOUND"},
		{"JobExecutionNotFound", JobExecutionNotFound, "JOB_EXECUTION_NOT_FOUND"},
		{"JobOperationNotFound", JobOperationNotFound, "JOB_OPERATION_NOT_FOUND"},
		{"JobDeploymentInProgress", JobDeploymentInProgress, "JOB_DEPLOYMENT_IN_PROGRESS"},
		{"JobAlreadySuspended", JobAlreadySuspended, "JOB_ALREADY_SUSPENDED"},
		{"JobNotSuspended", JobNotSuspended, "JOB_NOT_SUSPENDED"},
		{"JobQuotaExceeded", JobQuotaExceeded, "JOB_QUOTA_EXCEEDED"},
		{"JobInsufficientBalance", JobInsufficientBalance, "JOB_INSUFFICIENT_BALANCE"},
		{"JobScheduleInvalid", JobScheduleInvalid, "JOB_SCHEDULE_INVALID"},
		{"JobScheduleTooFrequent", JobScheduleTooFrequent, "JOB_SCHEDULE_TOO_FREQUENT"},
		{"JobTimezoneInvalid", JobTimezoneInvalid, "JOB_TIMEZONE_INVALID"},
		{"JobKindInvalid", JobKindInvalid, "JOB_KIND_INVALID"},
		{"JobAppRequired", JobAppRequired, "JOB_APP_REQUIRED"},
		{"JobAppNotFound", JobAppNotFound, "JOB_APP_NOT_FOUND"},
		{"JobImageRequired", JobImageRequired, "JOB_IMAGE_REQUIRED"},
		{"JobConcurrencyUnsupported", JobConcurrencyUnsupported, "JOB_CONCURRENCY_UNSUPPORTED"},
		{"JobImageNotFound", JobImageNotFound, "JOB_IMAGE_NOT_FOUND"},
		{"JobImageUnauthorized", JobImageUnauthorized, "JOB_IMAGE_UNAUTHORIZED"},
		{"JobImageRegistryUnreachable", JobImageRegistryUnreachable, "JOB_IMAGE_REGISTRY_UNREACHABLE"},
		{"JobInvalidPricingSlug", JobInvalidPricingSlug, "JOB_INVALID_PRICING_SLUG"},
		{"JobValidationFailed", JobValidationFailed, "JOB_VALIDATION_FAILED"},
		{"JobUnauthorized", JobUnauthorized, "JOB_UNAUTHORIZED"},
		{"JobInvalidID", JobInvalidID, "JOB_INVALID_ID"},
		{"JobInvalidRequestBody", JobInvalidRequestBody, "JOB_INVALID_REQUEST_BODY"},
		{"JobInternalError", JobInternalError, "JOB_INTERNAL_ERROR"},
		{"JobMetricsInternalError", JobMetricsInternalError, "JOB_METRICS_INTERNAL_ERROR"},
		{"JobLogsInternalError", JobLogsInternalError, "JOB_LOGS_INTERNAL_ERROR"},
		{"JobOperationFailed", JobOperationFailed, "JOB_OPERATION_FAILED"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}
