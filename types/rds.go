package types

import "time"

// RDSStatus enumerates the lifecycle states of a managed database instance.
// Transient states (provisioning, scaling, resizing, deleting) resolve
// asynchronously — clients poll GET /api/v1/rds/:id (or the operation) and
// watch Status until it settles on ready / suspended / failed.
type RDSStatus string

const (
	RDSStatusProvisioning RDSStatus = "provisioning"
	RDSStatusReady        RDSStatus = "ready"
	RDSStatusScaling      RDSStatus = "scaling"
	RDSStatusResizing     RDSStatus = "resizing"
	// RDSStatusReconfiguring is the transient state while an engine-config change
	// is applied — a parameter-template attach/edit or a tls_mode (TLS
	// enforcement) change. The instance keeps serving; it returns to "ready" when
	// the reconfigure settles.
	RDSStatusReconfiguring RDSStatus = "reconfiguring"
	// RDSStatusSwitchingOver is the transient state during a planned HA
	// switchover (sync standby being promoted to primary). The instance stays
	// fully serving; it returns to "ready" when the switch settles.
	RDSStatusSwitchingOver RDSStatus = "switching_over"
	RDSStatusSuspended     RDSStatus = "suspended"
	RDSStatusFailed        RDSStatus = "failed"
	RDSStatusDeleting      RDSStatus = "deleting"
	// RDSStatusDeleted is filtered server-side from list/get responses; included
	// for completeness of the wire enum.
	RDSStatusDeleted RDSStatus = "deleted"
)

// RDSMode is the topology a database runs in. "standalone" is a single primary;
// "ha" adds one synchronous standby for no-data-loss failover. Read replicas
// (asynchronous standbys) are a separate count layered on either mode. Chosen at
// create and changeable live via PATCH (mode / read_replicas).
type RDSMode string

const (
	RDSModeStandalone RDSMode = "standalone"
	RDSModeHA         RDSMode = "ha"
)

// RDSTLSMode controls TLS for client connections to the database:
//
//   - "disabled" — no server certificate; plaintext only.
//   - "optional" — TLS is available (server cert issued) but the server still
//     accepts plaintext; the client chooses via its libpq sslmode. This is the
//     default and matches the historical "ssl on, not enforced" behaviour.
//   - "required" — the server rejects non-TLS connections (plaintext refused),
//     like AWS rds.force_ssl=1 / Alibaba "force SSL". Gated by a platform flag.
//
// Chosen at create. TLS availability (the certificate) is fixed at create, but
// "optional" <-> "required" can be changed later via PATCH /rds/:id/tls — a
// server-side pg_hba reload, no restart. A transition to/from "disabled" on a
// live instance is not supported (it would mint/remove the cert and restart).
type RDSTLSMode string

const (
	RDSTLSModeDisabled RDSTLSMode = "disabled"
	RDSTLSModeOptional RDSTLSMode = "optional"
	RDSTLSModeRequired RDSTLSMode = "required"
)

// RDS engine identifiers. PostgreSQL is the only engine at launch.
const (
	RDSEnginePostgreSQL = "postgresql"
)

// RDSOperationStatus tracks an async lifecycle operation returned by the
// operations polling endpoint.
type RDSOperationStatus string

const (
	RDSOperationQueued     RDSOperationStatus = "queued"
	RDSOperationInProgress RDSOperationStatus = "in_progress"
	RDSOperationSucceeded  RDSOperationStatus = "succeeded"
	RDSOperationFailed     RDSOperationStatus = "failed"
)

// CreateRDSInstanceRequest is the body for POST /api/v1/rds. Honors
// Idempotency-Key — retrying the same key + body replays the cached response
// rather than provisioning a second database.
//
// Engine/EngineVersion select the database (e.g. "postgresql"/"16"). Plan is a
// catalogue slug (see PublicRDSPlanResponse.Slug). StorageGB is the data disk
// size; it can be grown later (never shrunk) via PATCH.
type CreateRDSInstanceRequest struct {
	Name          string `json:"name"`           // 1..63 chars, RFC-1035 label
	Engine        string `json:"engine"`         // "postgresql"
	EngineVersion string `json:"engine_version"` // e.g. "16"
	Plan          string `json:"plan"`           // plan slug, e.g. "kumo.pg.small"
	StorageGB     int    `json:"storage_gb"`     // >= tier minimum
	// TLSMode controls TLS for client connections: "disabled" (plaintext only),
	// "optional" (TLS available, client may still connect plaintext — the secure
	// default), or "required" (server rejects non-TLS connections). Omit to take
	// the default ("optional"). See RDSTLSMode. "required" is gated by a platform
	// flag (400 RDS_TLS_ENFORCEMENT_DISABLED when off).
	TLSMode string `json:"tls_mode,omitempty"`
	// ParameterTemplate is the slug of a parameter template to attach (must match
	// the instance's engine version). Omit to use the version's default template.
	ParameterTemplate string `json:"parameter_template,omitempty"`
	// Mode selects the topology: "standalone" (default; primary only) or "ha"
	// (primary + 1 synchronous standby for no-data-loss failover).
	Mode string `json:"mode,omitempty"`
	// ReadReplicas is the number of asynchronous read-only standbys (0..plan cap).
	// Layered on either mode; pointer so omitted (nil) means 0.
	ReadReplicas *int `json:"read_replicas,omitempty"`
	// ReadReplicaSpecs optionally sizes each read replica independently of the
	// primary (e.g. smaller, cheaper read replicas). When present, its length must
	// equal ReadReplicas (or ReadReplicas may be omitted and inferred from the
	// list); when omitted, all replicas use the primary's plan. A read replica's
	// storage is always >= the primary's (a streaming standby holds a full copy).
	ReadReplicaSpecs []ReadReplicaSpec `json:"read_replica_specs,omitempty"`
}

// ReadReplicaSpec requests one asynchronous read replica at a specific plan
// (instance class), so read replicas can be sized independently of the primary.
type ReadReplicaSpec struct {
	// Plan is the catalogue slug for this replica's instance class (see
	// PublicRDSPlanResponse.Slug). The replica is excluded from failover
	// (Patroni nofailover), so an undersized replica can never become primary.
	Plan string `json:"plan"`
}

// UpdateRDSInstanceRequest is the body for PATCH /api/v1/rds/:id. Exactly one
// dimension changes per call: supply Plan to vertically scale compute, or
// StorageGB to grow the data disk (shrink is rejected). Both async (202 +
// operation_id). Send If-Match with the instance's ETag to guard against
// concurrent writes (stale → 412 ETAG_MISMATCH).
type UpdateRDSInstanceRequest struct {
	Plan      string `json:"plan,omitempty"`
	StorageGB *int   `json:"storage_gb,omitempty"`
	// ParameterTemplate, when set, attaches a different parameter template and
	// live-reconfigures the instance (KubeBlocks reload or rolling restart). A
	// static-parameter change leaves PendingRestart=true until the restart lands.
	ParameterTemplate string `json:"parameter_template,omitempty"`
	// Mode and ReadReplicas live-scale the topology (KubeBlocks HorizontalScaling
	// + sync-mode reconfigure). Supply at most one mutation dimension per call
	// (plan, storage_gb, parameter_template, or mode/read_replicas).
	Mode         string `json:"mode,omitempty"`
	ReadReplicas *int   `json:"read_replicas,omitempty"`
	// ReadReplicaSpecs sizes the read replicas independently of the primary; see
	// CreateRDSInstanceRequest.ReadReplicaSpecs. When present, its length must
	// equal ReadReplicas. Replacing a replica's plan is a delete+add of that
	// replica node.
	ReadReplicaSpecs []ReadReplicaSpec `json:"read_replica_specs,omitempty"`
}

// UpdateRDSTLSRequest is the body for PATCH /api/v1/rds/:id/tls. It changes the
// connection TLS enforcement of a running instance between "optional" and
// "required" (a server-side pg_hba reload — no restart). Transitions to/from
// "disabled" are rejected (409 RDS_TLS_MODE_CHANGE_UNSUPPORTED); requesting
// "required" while the platform enforcement flag is off is rejected (400
// RDS_TLS_ENFORCEMENT_DISABLED). Async (202 + operation_id). Send If-Match with
// the instance ETag to guard against concurrent writes.
type UpdateRDSTLSRequest struct {
	TLSMode string `json:"tls_mode"`
}

// PublicRDSPlanResponse is the customer-facing plan (instance class) DTO
// returned by GET /api/v1/rds/plans and embedded in RDSInstanceResponse.
// Internal pricing inputs (base cost, margin) are intentionally absent — only
// the final PriceHour is exposed.
//
// CPUvCPU and PriceHour are decimal strings (e.g. "1", "0.5", "12.5000") —
// parse with a decimal library if you need arithmetic.
type PublicRDSPlanResponse struct {
	Slug     string `json:"slug"`
	Engine   string `json:"engine"`
	Name     string `json:"name"`
	CPUvCPU  string `json:"cpu_vcpu"`
	MemoryMB int    `json:"memory_mb"`
	// MaxReadReplicas is the per-plan cap on asynchronous read replicas (bounded
	// by a platform hard ceiling).
	MaxReadReplicas int           `json:"max_read_replicas"`
	PriceHour       string        `json:"price_hour"`
	Availability    *Availability `json:"availability,omitempty"`
}

// PublicRDSEngineVersionResponse is a customer-facing entry of the engine
// version catalogue returned by GET /api/v1/rds/engine-versions. Version is the
// label the client passes as CreateRDSInstanceRequest.EngineVersion (e.g. "16");
// the server maps it to the exact image internally. Status is
// "supported" | "deprecated" | "eol" — only non-eol versions accept new
// instances. Exactly one entry per engine has IsDefault=true.
type PublicRDSEngineVersionResponse struct {
	Engine    string `json:"engine"`
	Version   string `json:"version"`
	IsDefault bool   `json:"is_default"`
	Status    string `json:"status"`
}

// RDSParameter is a single engine configuration key/value within a template.
type RDSParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// RDSPgParameterResponse describes an editable parameter from the engine
// version's catalogue (allowlist) returned by
// GET /api/v1/rds/engine-versions/:id/parameters. DataType is
// int|real|bool|string|enum. ApplyMethod is "dynamic" (hot reload) or "static"
// (needs a restart). Min/Max/EnumValues bound the allowed values.
type RDSPgParameterResponse struct {
	Name         string   `json:"name"`
	DataType     string   `json:"data_type"`
	ApplyMethod  string   `json:"apply_method"`
	DefaultValue string   `json:"default_value,omitempty"`
	MinValue     string   `json:"min_value,omitempty"`
	MaxValue     string   `json:"max_value,omitempty"`
	EnumValues   []string `json:"enum_values,omitempty"`
	Unit         string   `json:"unit,omitempty"`
	Description  string   `json:"description,omitempty"`
}

// PublicRDSParameterTemplateResponse is a parameter template (a reusable named
// set of engine config, scoped to one engine version — like an AWS parameter
// group / Huawei parameter template). IsSystem/IsDefault templates are
// read-only; clone to customise.
type PublicRDSParameterTemplateResponse struct {
	ID            uint           `json:"id"`
	Slug          string         `json:"slug"`
	Name          string         `json:"name"`
	Description   string         `json:"description,omitempty"`
	Engine        string         `json:"engine"`
	EngineVersion string         `json:"engine_version"`
	IsSystem      bool           `json:"is_system"`
	IsDefault     bool           `json:"is_default"`
	Parameters    []RDSParameter `json:"parameters"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// CreateRDSParameterTemplateRequest is the body for
// POST /api/v1/rds/parameter-templates. EngineVersion scopes the template; every
// parameter must be in that version's editable allowlist.
type CreateRDSParameterTemplateRequest struct {
	Name          string         `json:"name"`
	Description   string         `json:"description,omitempty"`
	Engine        string         `json:"engine"`
	EngineVersion string         `json:"engine_version"`
	Parameters    []RDSParameter `json:"parameters,omitempty"`
}

// UpdateRDSParameterTemplateRequest is the body for
// PATCH /api/v1/rds/parameter-templates/:id. Parameters replaces the full set.
// Send If-Match with the template ETag to guard concurrent edits.
type UpdateRDSParameterTemplateRequest struct {
	Name        *string        `json:"name,omitempty"`
	Description *string        `json:"description,omitempty"`
	Parameters  []RDSParameter `json:"parameters,omitempty"`
}

// RDSInstanceResponse is the detail returned by GET /api/v1/rds/:id and the
// items of GET /api/v1/rds. The server sets ETag from UpdatedAt; echo it back
// in If-Match on PATCH to detect concurrent writes.
//
// Endpoint{Host,Port} are populated once the instance is ready. Credentials are
// not included here — fetch them via GET /api/v1/rds/:id/connection.
type RDSInstanceResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Engine        string `json:"engine"`
	EngineVersion string `json:"engine_version"`
	Mode          string `json:"mode"`
	Replicas      int    `json:"replicas"`
	// ReadReplicas is the number of asynchronous read-only standbys.
	ReadReplicas int `json:"read_replicas"`
	// ReadReplicaDetails lists each read replica's resolved plan and state, so a
	// client can see heterogeneous (independently-sized) replicas. Empty when
	// there are no replicas or when all replicas share the primary's plan.
	ReadReplicaDetails []ReadReplicaDetail    `json:"read_replica_details,omitempty"`
	Plan               *PublicRDSPlanResponse `json:"plan,omitempty"`
	StorageGB          int                    `json:"storage_gb"`
	Status             string                 `json:"status"`
	EndpointHost       string                 `json:"endpoint_host,omitempty"`
	// ReadEndpointHost is the read-only Service host that load-balances across
	// standbys; populated once at least one replica exists.
	ReadEndpointHost string `json:"read_endpoint_host,omitempty"`
	EndpointPort     int    `json:"endpoint_port,omitempty"`
	StatusMessage    string `json:"status_message,omitempty"`
	// IsSuspended is true when the database was stopped for non-payment; the
	// user must top up and POST /rds/:id/start to resume. SuspendReason explains
	// why ("insufficient balance"); SuspendedAt is when the retention window
	// (auto-deletion) started.
	IsSuspended   bool       `json:"is_suspended"`
	SuspendReason string     `json:"suspend_reason,omitempty"`
	SuspendedAt   *time.Time `json:"suspended_at,omitempty"`
	// TLSMode is the connection TLS policy: "disabled" | "optional" | "required".
	// "required" means the server rejects non-TLS connections; "optional" means
	// TLS is available but plaintext is still accepted. See RDSTLSMode.
	TLSMode string `json:"tls_mode"`
	// ParameterTemplate is the attached configuration template, when set.
	ParameterTemplate *PublicRDSParameterTemplateResponse `json:"parameter_template,omitempty"`
	// PendingRestart is true when a static parameter change has been applied but
	// the rolling restart that activates it has not completed yet.
	PendingRestart bool      `json:"pending_restart,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ReadReplicaDetail is one read replica's resolved spec and lifecycle state
// within RDSInstanceResponse. Ordinal is the stable index used to name the
// replica node; Plan is its instance class (may differ from the primary's).
type ReadReplicaDetail struct {
	Ordinal int                    `json:"ordinal"`
	Plan    *PublicRDSPlanResponse `json:"plan,omitempty"`
	Status  string                 `json:"status"`
}

// RDSMutationResponse is the 202 Accepted body for async POST/PATCH/DELETE on
// an RDS instance. Poll GET /api/v1/rds/:id/operations/:operation_id (or the
// instance's Status) until the operation settles.
type RDSMutationResponse struct {
	ID          uint   `json:"id"`
	OperationID string `json:"operation_id"`
	Status      string `json:"status"`
}

// RDSOperationResponse is returned by
// GET /api/v1/rds/:id/operations/:operation_id. Status is one of the
// RDSOperation* constants; ErrorCode/ErrorMessage are set on terminal failure.
type RDSOperationResponse struct {
	OperationID  string     `json:"operation_id"`
	ActionType   string     `json:"action_type"` // create | scale | resize | delete
	Status       string     `json:"status"`
	ErrorCode    string     `json:"error_code,omitempty"`
	ErrorMessage string     `json:"error_message,omitempty"`
	QueuedAt     time.Time  `json:"queued_at"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
}

// ── Backups (to object storage) ──────────────────────────────────────

// RDSBackupStatus enumerates the lifecycle states of a database backup stored
// in object storage. Transient states (pending, running) resolve asynchronously
// — clients poll GET /api/v1/rds/:id/backups/:bid until Status settles on
// completed / failed. expired and deleting are terminal cleanup states.
type RDSBackupStatus string

const (
	RDSBackupStatusPending   RDSBackupStatus = "pending"
	RDSBackupStatusRunning   RDSBackupStatus = "running"
	RDSBackupStatusCompleted RDSBackupStatus = "completed"
	RDSBackupStatusFailed    RDSBackupStatus = "failed"
	RDSBackupStatusDeleting  RDSBackupStatus = "deleting"
	RDSBackupStatusExpired   RDSBackupStatus = "expired"
)

// RDSBackupMethod is how a backup was taken. "full" is a complete base backup
// (pg_basebackup); "continuous" is WAL archiving for point-in-time recovery
// (PITR). Only "full" is offered at launch.
type RDSBackupMethod string

const (
	RDSBackupMethodFull       RDSBackupMethod = "full"
	RDSBackupMethodContinuous RDSBackupMethod = "continuous"
)

// CreateRDSBackupRequest is the body for POST /api/v1/rds/:id/backups — an
// on-demand backup of the instance to object storage. Honors Idempotency-Key
// (retrying the same key returns the same backup rather than starting a second).
// Async (202 + operation_id); poll the backup until Status="completed".
type CreateRDSBackupRequest struct {
	// RetentionDays optionally overrides the backup tier's default retention for
	// this one backup (bounded by the tier's min/max). 0 / omitted uses the tier
	// default. After this many days the backup is automatically deleted from
	// object storage and billing for it stops.
	RetentionDays int `json:"retention_days,omitempty"`
}

// RDSBackupResponse is the detail returned by GET /api/v1/rds/:id/backups/:bid
// and the items of GET /api/v1/rds/:id/backups. The server sets ETag from
// UpdatedAt. SizeBytes is the measured size in object storage (0 until the
// backup completes).
type RDSBackupResponse struct {
	ID            uint `json:"id"`
	RDSInstanceID uint `json:"rds_instance_id"`
	// SourceInstanceName / SourceEngineVersion describe the database the backup
	// was taken from — retained for the global backups view even after that
	// database is deleted (the backup itself survives and stays restorable).
	SourceInstanceName  string `json:"source_instance_name,omitempty"`
	SourceEngineVersion string `json:"source_engine_version,omitempty"`
	Method              string `json:"method"` // full | continuous
	Status              string `json:"status"`
	SizeBytes           int64  `json:"size_bytes"`
	// TierSlug is the backup tier this backup is billed on (price + storage
	// destination). RetentionDays / ExpiresAt describe its automatic cleanup.
	TierSlug      string     `json:"tier_slug,omitempty"`
	RetentionDays int        `json:"retention_days"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	ErrorMessage  string     `json:"error_message,omitempty"`
	StartedAt     *time.Time `json:"started_at,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// UpdateRDSBackupConfigRequest is the body for PUT /api/v1/rds/:id/backup-config.
// It enables/disables automatic scheduled backups for an instance and sets the
// schedule, retention, and backup tier. When Enabled is false the schedule is
// removed (existing backups are retained until their own expiry).
type UpdateRDSBackupConfigRequest struct {
	Enabled bool `json:"enabled"`
	// ScheduleCron is a standard 5-field cron expression for automatic full
	// backups (e.g. "0 2 * * *" = daily 02:00 UTC). Required when Enabled.
	ScheduleCron string `json:"schedule_cron,omitempty"`
	// RetentionDays is how long scheduled backups are kept (bounded by the tier).
	RetentionDays int `json:"retention_days,omitempty"`
	// TierSlug selects the backup tier (price + object-storage destination). Omit
	// to use the platform default tier.
	TierSlug string `json:"tier_slug,omitempty"`
	// PITREnabled turns on continuous WAL archiving (point-in-time recovery). With
	// it on, the database can be restored to any second within the retention
	// window (not just to a discrete full backup) — see Restore.RestoreToTime.
	// Requires Enabled (scheduled full backups provide the base for WAL replay).
	// The archived WAL is billed on the same backup tier as full backups.
	PITREnabled bool `json:"pitr_enabled,omitempty"`
}

// RDSBackupConfigResponse is the current automatic-backup configuration of an
// instance, returned by PUT /api/v1/rds/:id/backup-config.
type RDSBackupConfigResponse struct {
	Enabled       bool   `json:"enabled"`
	ScheduleCron  string `json:"schedule_cron,omitempty"`
	RetentionDays int    `json:"retention_days"`
	TierSlug      string `json:"tier_slug,omitempty"`
	PITREnabled   bool   `json:"pitr_enabled"`
}

// RestoreRDSBackupRequest is the body for POST /api/v1/rds/:id/restore. It
// provisions a NEW database instance from a completed backup (the source
// instance is untouched). The new instance is billed independently. Async
// (202 + operation_id on the new instance); poll the new instance until ready.
type RestoreRDSBackupRequest struct {
	// BackupID is the completed backup to restore from (must belong to the source
	// instance referenced in the path).
	BackupID uint `json:"backup_id"`
	// Name is the new instance's name (1..40 chars, RFC-1035 label).
	Name string `json:"name"`
	// Plan is the new instance's compute plan slug. Omit to inherit the source's.
	Plan string `json:"plan,omitempty"`
	// StorageGB is the new instance's data-disk size; must be >= the source
	// backup's storage (400 RDS_RESTORE_STORAGE_TOO_SMALL otherwise). Omit to
	// inherit the source's size.
	StorageGB int `json:"storage_gb,omitempty"`
	// RestoreToTime, when set (RFC-3339, e.g. "2026-06-21T05:43:33Z"), performs a
	// point-in-time recovery: the new instance is restored to the database state
	// as of that instant by replaying archived WAL. Requires the source to have
	// had PITR (continuous WAL archiving) enabled, and the time must fall within
	// the available WAL window. When empty, restores to the backup's own point.
	RestoreToTime string `json:"restore_to_time,omitempty"`
}

// RDSConnectionResponse is returned by GET /api/v1/rds/:id/connection. The
// Password is read live from the instance's credentials secret and is only
// available once the instance is ready (otherwise 409 RDS_INSTANCE_NOT_READY).
type RDSConnectionResponse struct {
	Host string `json:"host"`
	// ReadHost is the load-balanced read-only endpoint (routes to any current
	// standby); present when replicas exist. Use it for read-only / reporting
	// traffic to offload the primary.
	ReadHost string `json:"read_host,omitempty"`
	// ReadReplicaHosts are the per-replica direct endpoints — one stable host per
	// current standby (vs ReadHost which load-balances across them). Present when
	// replicas exist. NOTE: these address individual nodes; on failover a node's
	// role can change, so prefer ReadHost (or the primary Host) for role-stable
	// routing. Empty for a single-node instance.
	ReadReplicaHosts []string `json:"read_replica_hosts,omitempty"`
	Port             int      `json:"port"`
	Username         string   `json:"username"`
	Database         string   `json:"database"`
	Password         string   `json:"password"`
	// SSLMode is the recommended libpq sslmode for this instance's TLSMode:
	// "require" when TLS is enforced (tls_mode=required), "prefer" when TLS is
	// available but not enforced (tls_mode=optional), "disable" when there is no
	// server cert (tls_mode=disabled). It is a client hint; the server only
	// rejects plaintext when tls_mode=required.
	SSLMode string `json:"ssl_mode"`
	// CACert is the PEM-encoded CA bundle that signed the server certificate.
	// Present whenever a cert exists (tls_mode optional or required); supply it as
	// sslrootcert to connect with sslmode=verify-full. Omitted when tls_mode is
	// disabled (or the cert is momentarily unavailable).
	CACert string `json:"ca_cert,omitempty"`
}
