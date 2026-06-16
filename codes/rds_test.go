package codes

import "testing"

// Wire codes are a public contract (terraform-provider-kumo, kumo-cli). These
// assert the exact string values so an accidental rename is caught here before
// release.
func TestRDSCodeValues(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"RDSInstanceNotFound", RDSInstanceNotFound, "RDS_INSTANCE_NOT_FOUND"},
		{"RDSFlavorNotFound", RDSFlavorNotFound, "RDS_FLAVOR_NOT_FOUND"},
		{"RDSFlavorDisabled", RDSFlavorDisabled, "RDS_FLAVOR_DISABLED"},
		{"RDSEngineNotSupported", RDSEngineNotSupported, "RDS_ENGINE_NOT_SUPPORTED"},
		{"RDSEngineVersionNotSupported", RDSEngineVersionNotSupported, "RDS_ENGINE_VERSION_NOT_SUPPORTED"},
		{"RDSActionInProgress", RDSActionInProgress, "RDS_ACTION_IN_PROGRESS"},
		{"RDSInstanceNotReady", RDSInstanceNotReady, "RDS_INSTANCE_NOT_READY"},
		{"RDSInstanceNotSuspended", RDSInstanceNotSuspended, "RDS_INSTANCE_NOT_SUSPENDED"},
		{"RDSInvalidStorageSize", RDSInvalidStorageSize, "RDS_INVALID_STORAGE_SIZE"},
		{"RDSOperationNotFound", RDSOperationNotFound, "RDS_OPERATION_NOT_FOUND"},
		{"RDSInsufficientBalance", RDSInsufficientBalance, "RDS_INSUFFICIENT_BALANCE"},
		{"RDSUnauthorized", RDSUnauthorized, "RDS_UNAUTHORIZED"},
		{"RDSInvalidRequestBody", RDSInvalidRequestBody, "RDS_INVALID_REQUEST_BODY"},
		{"RDSValidationError", RDSValidationError, "RDS_VALIDATION_ERROR"},
		{"RDSInvalidInstanceID", RDSInvalidInstanceID, "RDS_INVALID_INSTANCE_ID"},
		{"RDSInvalidPagination", RDSInvalidPagination, "RDS_INVALID_PAGINATION"},
		{"RDSInvalidStatusFilter", RDSInvalidStatusFilter, "RDS_INVALID_STATUS_FILTER"},
		{"RDSInternalError", RDSInternalError, "RDS_INTERNAL_ERROR"},
		{"RDSNamespaceTerminating", RDSNamespaceTerminating, "RDS_NAMESPACE_TERMINATING"},
		{"RDSParameterTemplateNotFound", RDSParameterTemplateNotFound, "RDS_PARAMETER_TEMPLATE_NOT_FOUND"},
		{"RDSParameterTemplateReadOnly", RDSParameterTemplateReadOnly, "RDS_PARAMETER_TEMPLATE_READ_ONLY"},
		{"RDSParameterTemplateInUse", RDSParameterTemplateInUse, "RDS_PARAMETER_TEMPLATE_IN_USE"},
		{"RDSParameterNotAllowed", RDSParameterNotAllowed, "RDS_PARAMETER_NOT_ALLOWED"},
		{"RDSParameterInvalidValue", RDSParameterInvalidValue, "RDS_PARAMETER_INVALID_VALUE"},
		{"RDSParameterTemplateVersionMismatch", RDSParameterTemplateVersionMismatch, "RDS_PARAMETER_TEMPLATE_VERSION_MISMATCH"},
		{"RDSReadReplicaLimitExceeded", RDSReadReplicaLimitExceeded, "RDS_READ_REPLICA_LIMIT_EXCEEDED"},
		{"RDSInvalidMode", RDSInvalidMode, "RDS_INVALID_MODE"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}
