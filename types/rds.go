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
	RDSStatusSuspended    RDSStatus = "suspended"
	RDSStatusFailed       RDSStatus = "failed"
	RDSStatusDeleting     RDSStatus = "deleting"
	// RDSStatusDeleted is filtered server-side from list/get responses; included
	// for completeness of the wire enum.
	RDSStatusDeleted RDSStatus = "deleted"
)

// RDSMode is the topology a database is running in. It is derived from the
// replica count (standalone == 1 replica, ha == 3) — never set directly by the
// client; switching modes is a future scale operation.
type RDSMode string

const (
	RDSModeStandalone RDSMode = "standalone"
	RDSModeHA         RDSMode = "ha"
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
}

// UpdateRDSInstanceRequest is the body for PATCH /api/v1/rds/:id. Exactly one
// dimension changes per call: supply Plan to vertically scale compute, or
// StorageGB to grow the data disk (shrink is rejected). Both async (202 +
// operation_id). Send If-Match with the instance's ETag to guard against
// concurrent writes (stale → 412 ETAG_MISMATCH).
type UpdateRDSInstanceRequest struct {
	Plan      string `json:"plan,omitempty"`
	StorageGB *int   `json:"storage_gb,omitempty"`
}

// PublicRDSPlanResponse is the customer-facing plan (instance class) DTO
// returned by GET /api/v1/rds/plans and embedded in RDSInstanceResponse.
// Internal pricing inputs (base cost, margin) are intentionally absent — only
// the final PriceHour is exposed.
//
// CPUvCPU and PriceHour are decimal strings (e.g. "1", "0.5", "12.5000") —
// parse with a decimal library if you need arithmetic.
type PublicRDSPlanResponse struct {
	Slug         string        `json:"slug"`
	Engine       string        `json:"engine"`
	Name         string        `json:"name"`
	CPUvCPU      string        `json:"cpu_vcpu"`
	MemoryMB     int           `json:"memory_mb"`
	PriceHour    string        `json:"price_hour"`
	Availability *Availability `json:"availability,omitempty"`
}

// RDSInstanceResponse is the detail returned by GET /api/v1/rds/:id and the
// items of GET /api/v1/rds. The server sets ETag from UpdatedAt; echo it back
// in If-Match on PATCH to detect concurrent writes.
//
// Endpoint{Host,Port} are populated once the instance is ready. Credentials are
// not included here — fetch them via GET /api/v1/rds/:id/connection.
type RDSInstanceResponse struct {
	ID            uint                     `json:"id"`
	Name          string                   `json:"name"`
	Engine        string                   `json:"engine"`
	EngineVersion string                   `json:"engine_version"`
	Mode          string                 `json:"mode"`
	Replicas      int                    `json:"replicas"`
	Plan          *PublicRDSPlanResponse `json:"plan,omitempty"`
	StorageGB     int                    `json:"storage_gb"`
	Status        string                   `json:"status"`
	EndpointHost  string                   `json:"endpoint_host,omitempty"`
	EndpointPort  int                      `json:"endpoint_port,omitempty"`
	StatusMessage string                   `json:"status_message,omitempty"`
	// IsSuspended is true when the database was stopped for non-payment; the
	// user must top up and POST /rds/:id/start to resume. SuspendReason explains
	// why ("insufficient balance"); SuspendedAt is when the retention window
	// (auto-deletion) started.
	IsSuspended   bool       `json:"is_suspended"`
	SuspendReason string     `json:"suspend_reason,omitempty"`
	SuspendedAt   *time.Time `json:"suspended_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
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
	OperationID string     `json:"operation_id"`
	ActionType  string     `json:"action_type"` // create | scale | resize | delete
	Status      string     `json:"status"`
	ErrorCode   string     `json:"error_code,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
	QueuedAt    time.Time  `json:"queued_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// RDSConnectionResponse is returned by GET /api/v1/rds/:id/connection. The
// Password is read live from the instance's credentials secret and is only
// available once the instance is ready (otherwise 409 RDS_INSTANCE_NOT_READY).
type RDSConnectionResponse struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Database string `json:"database"`
	Password string `json:"password"`
	SSLMode  string `json:"ssl_mode"`
}
