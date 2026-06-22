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

	// RDSEngineNotSupported — the requested engine is not offered (e.g. a
	// non-postgresql engine at launch).
	RDSEngineNotSupported = "RDS_ENGINE_NOT_SUPPORTED"

	// RDSEngineVersionNotSupported — the requested engine_version is not in the
	// active version catalogue (unknown, disabled, or end-of-life). Call
	// GET /api/v1/rds/engine-versions for the offered set.
	RDSEngineVersionNotSupported = "RDS_ENGINE_VERSION_NOT_SUPPORTED"

	// RDSEngineVersionUnavailable — an admin tried to register/update an engine
	// version whose kb_service_version is not offered by the installed KubeBlocks
	// PostgreSQL addon (ComponentVersion). Provisioning on it would hang, so the
	// catalogue edit is rejected up-front. Pick a serviceVersion the addon ships.
	RDSEngineVersionUnavailable = "RDS_ENGINE_VERSION_UNAVAILABLE"

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

	// RDSNamespaceTerminating — the tenant's Kubernetes namespace is still being
	// torn down from a previous (last) database deletion, so a new instance
	// cannot be created yet. Transient; retry after a few seconds.
	RDSNamespaceTerminating = "RDS_NAMESPACE_TERMINATING"

	// ── Parameter templates (DB-engine configuration groups) ──────────────────

	// RDSParameterTemplateNotFound — no parameter template with the given
	// id/slug exists in the caller's scope.
	RDSParameterTemplateNotFound = "RDS_PARAMETER_TEMPLATE_NOT_FOUND"

	// RDSParameterTemplateReadOnly — the template is a system/default template
	// and cannot be edited or deleted. Clone it, then edit the copy.
	RDSParameterTemplateReadOnly = "RDS_PARAMETER_TEMPLATE_READ_ONLY"

	// RDSParameterTemplateInUse — the template is attached to one or more
	// instances and cannot be deleted until they detach.
	RDSParameterTemplateInUse = "RDS_PARAMETER_TEMPLATE_IN_USE"

	// RDSParameterNotAllowed — a parameter in the request is not in the
	// engine-version's editable allowlist (managed/unsafe params are blocked).
	RDSParameterNotAllowed = "RDS_PARAMETER_NOT_ALLOWED"

	// RDSParameterInvalidValue — a parameter value failed type/range/enum
	// validation against the catalogue.
	RDSParameterInvalidValue = "RDS_PARAMETER_INVALID_VALUE"

	// RDSParameterTemplateVersionMismatch — the template's engine version does
	// not match the instance's (a template applies only within its version, like
	// an AWS parameter-group family).
	RDSParameterTemplateVersionMismatch = "RDS_PARAMETER_TEMPLATE_VERSION_MISMATCH"

	// RDSParameterTemplateDefaultProtected — the operation would delete or demote
	// the engine version's only default template; every version must always keep a
	// default. Promote another template to default first, or pick a new default in
	// the same request.
	RDSParameterTemplateDefaultProtected = "RDS_PARAMETER_TEMPLATE_DEFAULT_PROTECTED"

	// ── Read replicas / HA ────────────────────────────────────────────────────

	// RDSReadReplicaLimitExceeded — the requested read-replica count exceeds the
	// plan's per-flavor cap or the platform hard ceiling.
	RDSReadReplicaLimitExceeded = "RDS_READ_REPLICA_LIMIT_EXCEEDED"

	// RDSInvalidMode — the requested topology mode is not standalone|ha.
	RDSInvalidMode = "RDS_INVALID_MODE"

	// RDSReadReplicaSpecMismatch — read_replica_specs was supplied but is
	// inconsistent with the request: its length does not match read_replicas, a
	// replica plan slug is unknown, or a replica's resolved spec violates the
	// replica-must-be->=-primary-storage rule.
	RDSReadReplicaSpecMismatch = "RDS_READ_REPLICA_SPEC_MISMATCH"

	// RDSModeChangeUnsupported — switching topology mode (standalone<->ha) on a
	// live instance is not supported; HA must be selected at create time. Returned
	// only when the platform cannot apply synchronous replication to a running
	// cluster (see the RDS HA runbook).
	RDSModeChangeUnsupported = "RDS_MODE_CHANGE_UNSUPPORTED"

	// ── Manual switchover (planned HA role swap) ──────────────────────────────

	// RDSSwitchoverDisabled — a manual switchover was requested but the platform
	// RDS_SWITCHOVER_ENABLED flag is off. Planned switchover is not offered yet.
	RDSSwitchoverDisabled = "RDS_SWITCHOVER_DISABLED"

	// RDSSwitchoverNotHA — a manual switchover was requested on a standalone
	// instance. Switchover promotes the synchronous standby, so it requires HA
	// mode (primary + sync standby). Async read replicas are never candidates.
	RDSSwitchoverNotHA = "RDS_SWITCHOVER_NOT_HA"

	// RDSSwitchoverNotReady — a manual switchover was requested but there is no
	// healthy synchronous standby to promote (the standby is down, not yet
	// caught up, or otherwise ineligible). Retry once the standby is healthy.
	RDSSwitchoverNotReady = "RDS_SWITCHOVER_NOT_READY"

	// ── TLS mode / enforcement ────────────────────────────────────────────────

	// RDSInvalidTLSMode — the requested tls_mode is not disabled|optional|required.
	RDSInvalidTLSMode = "RDS_INVALID_TLS_MODE"

	// RDSTLSEnforcementDisabled — tls_mode "required" was requested but the
	// platform RDS_TLS_ENFORCE_ENABLED flag is off (server-side TLS enforcement is
	// not offered yet). Use "optional", or ask the operator to enable enforcement.
	RDSTLSEnforcementDisabled = "RDS_TLS_ENFORCEMENT_DISABLED"

	// RDSTLSModeChangeUnsupported — a day-2 tls_mode change to or from "disabled"
	// was requested. Only "optional" <-> "required" is changeable on a live
	// instance (a pg_hba reload); changing TLS availability (the cert) requires
	// recreating the instance.
	RDSTLSModeChangeUnsupported = "RDS_TLS_MODE_CHANGE_UNSUPPORTED"

	// ── Backups (to object storage) ───────────────────────────────────────────

	// RDSBackupDisabled — a backup operation (on-demand backup, schedule config,
	// or restore) was requested but the platform RDS_BACKUP_ENABLED flag is off.
	// Managed backups are not offered yet on this deployment.
	RDSBackupDisabled = "RDS_BACKUP_DISABLED"

	// RDSBackupInProgress — a backup is already running for this instance and the
	// per-instance concurrency cap was reached. Wait for the in-flight backup to
	// settle (poll its operation_id), then retry.
	RDSBackupInProgress = "RDS_BACKUP_IN_PROGRESS"

	// RDSBackupNotFound — no backup with the given id exists for this instance in
	// the caller's scope (also covers cross-tenant access attempts and backups
	// already deleted/expired).
	RDSBackupNotFound = "RDS_BACKUP_NOT_FOUND"

	// RDSBackupNotReady — the referenced backup is not in a usable state for the
	// requested operation (e.g. a restore was requested from a backup that is
	// still running or has failed). Only a completed backup can be restored.
	RDSBackupNotReady = "RDS_BACKUP_NOT_READY"

	// RDSBackupTierNotFound — the requested backup tier slug does not exist or is
	// not active in the catalogue. Call the backup-tier listing for the offered
	// set, or omit to use the default tier.
	RDSBackupTierNotFound = "RDS_BACKUP_TIER_NOT_FOUND"

	// RDSRestoreStorageTooSmall — a restore requested a target data-disk size
	// smaller than the source backup's storage. A restored database must be at
	// least as large as the database the backup was taken from.
	RDSRestoreStorageTooSmall = "RDS_RESTORE_STORAGE_TOO_SMALL"

	// RDSPITRNotEnabled — a point-in-time restore (restore_to_time) was requested
	// for a source database that does not have PITR (continuous WAL archiving)
	// enabled, so there is no WAL stream to replay. Enable PITR on the source
	// first (backup config), or omit restore_to_time to restore to a full backup.
	RDSPITRNotEnabled = "RDS_PITR_NOT_ENABLED"

	// RDSInvalidRestoreTime — restore_to_time is malformed (not RFC-3339) or falls
	// outside the source's restorable WAL window [earliest, latest]. Pick a
	// timestamp within the reported window.
	RDSInvalidRestoreTime = "RDS_INVALID_RESTORE_TIME"
)
