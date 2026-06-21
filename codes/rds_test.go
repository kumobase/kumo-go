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
		{"RDSEngineVersionUnavailable", RDSEngineVersionUnavailable, "RDS_ENGINE_VERSION_UNAVAILABLE"},
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
		{"RDSParameterTemplateDefaultProtected", RDSParameterTemplateDefaultProtected, "RDS_PARAMETER_TEMPLATE_DEFAULT_PROTECTED"},
		{"RDSReadReplicaLimitExceeded", RDSReadReplicaLimitExceeded, "RDS_READ_REPLICA_LIMIT_EXCEEDED"},
		{"RDSInvalidMode", RDSInvalidMode, "RDS_INVALID_MODE"},
		{"RDSReadReplicaSpecMismatch", RDSReadReplicaSpecMismatch, "RDS_READ_REPLICA_SPEC_MISMATCH"},
		{"RDSModeChangeUnsupported", RDSModeChangeUnsupported, "RDS_MODE_CHANGE_UNSUPPORTED"},
		{"RDSSwitchoverDisabled", RDSSwitchoverDisabled, "RDS_SWITCHOVER_DISABLED"},
		{"RDSSwitchoverNotHA", RDSSwitchoverNotHA, "RDS_SWITCHOVER_NOT_HA"},
		{"RDSSwitchoverNotReady", RDSSwitchoverNotReady, "RDS_SWITCHOVER_NOT_READY"},
		{"RDSInvalidTLSMode", RDSInvalidTLSMode, "RDS_INVALID_TLS_MODE"},
		{"RDSTLSEnforcementDisabled", RDSTLSEnforcementDisabled, "RDS_TLS_ENFORCEMENT_DISABLED"},
		{"RDSTLSModeChangeUnsupported", RDSTLSModeChangeUnsupported, "RDS_TLS_MODE_CHANGE_UNSUPPORTED"},
		{"RDSBackupDisabled", RDSBackupDisabled, "RDS_BACKUP_DISABLED"},
		{"RDSBackupInProgress", RDSBackupInProgress, "RDS_BACKUP_IN_PROGRESS"},
		{"RDSBackupNotFound", RDSBackupNotFound, "RDS_BACKUP_NOT_FOUND"},
		{"RDSBackupNotReady", RDSBackupNotReady, "RDS_BACKUP_NOT_READY"},
		{"RDSBackupTierNotFound", RDSBackupTierNotFound, "RDS_BACKUP_TIER_NOT_FOUND"},
		{"RDSRestoreStorageTooSmall", RDSRestoreStorageTooSmall, "RDS_RESTORE_STORAGE_TOO_SMALL"},
		{"RDSPITRNotEnabled", RDSPITRNotEnabled, "RDS_PITR_NOT_ENABLED"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}
