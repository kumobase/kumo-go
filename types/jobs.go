package types

import "time"

// JobKind is the create-shape discriminator. "app_attached" jobs piggyback
// on an existing App (image, env, secret mounts are inherited from the App
// at run time). "standalone" jobs supply their own image and config.
type JobKind string

const (
	JobKindAppAttached JobKind = "app_attached"
	JobKindStandalone  JobKind = "standalone"
)

// JobDeploymentStatus tracks the asynchronous lifecycle of the underlying
// k8s CronJob (or, for scheduleless jobs, the readiness of the job
// definition). Mirrors AppDeploymentStatus.
type JobDeploymentStatus string

const (
	JobDeploymentStatusPending   JobDeploymentStatus = "pending"
	JobDeploymentStatusDeploying JobDeploymentStatus = "deploying"
	JobDeploymentStatusActive    JobDeploymentStatus = "active"
	JobDeploymentStatusFailed    JobDeploymentStatus = "failed"
	JobDeploymentStatusDeleting  JobDeploymentStatus = "deleting"
	JobDeploymentStatusDeleted   JobDeploymentStatus = "deleted"
)

// JobConcurrencyPolicy mirrors k8s CronJob.spec.concurrencyPolicy. The MVP
// only accepts Forbid; the field exists on the wire so adding "Allow" /
// "Replace" later is additive (no major bump).
type JobConcurrencyPolicy string

const (
	JobConcurrencyForbid  JobConcurrencyPolicy = "Forbid"
	JobConcurrencyAllow   JobConcurrencyPolicy = "Allow"
	JobConcurrencyReplace JobConcurrencyPolicy = "Replace"
)

// JobExecutionTrigger says how this execution row came into being.
type JobExecutionTrigger string

const (
	JobExecutionTriggerSchedule JobExecutionTrigger = "schedule"
	JobExecutionTriggerManual   JobExecutionTrigger = "manual"
	JobExecutionTriggerRetry    JobExecutionTrigger = "retry"
)

// JobExecutionStatus is the terminal-or-pending state of one execution.
// "pending" means the row exists but no pod container has started yet
// (image pulling, k8s scheduling). "running" means a container is up and
// metering against billing. The three terminal states distinguish exit
// reasons for billing audits and UI badges.
type JobExecutionStatus string

const (
	JobExecutionStatusPending   JobExecutionStatus = "pending"
	JobExecutionStatusRunning   JobExecutionStatus = "running"
	JobExecutionStatusSucceeded JobExecutionStatus = "succeeded"
	JobExecutionStatusFailed    JobExecutionStatus = "failed"
	JobExecutionStatusTimeout   JobExecutionStatus = "timeout"
)

// JobSecretRef references a Kumo Secret to mount into a standalone job's
// container. Mirrors the App secret-mount shape (modules/app's secret_var_apps
// + secret_file_apps), narrowed to the fields the job runtime needs.
//
// For env-var secrets, SourceFrom is the key inside the secret; MountTo is
// the env variable name the container sees.
// For file secrets, SourceFrom is left empty and MountTo is the absolute
// path the file is mounted at.
type JobSecretRef struct {
	SecretID   uint   `json:"secret_id,omitempty"`
	SecretName string `json:"secret_name,omitempty"`
	SourceFrom string `json:"source_from,omitempty"`
	MountTo    string `json:"mount_to"`
}

// CreateJobRequest is the body for POST /api/v1/jobs. Honors Idempotency-Key.
// Returns 202 + ResponseJobAsync (poll the operation_id).
//
// Shape rules enforced by the server:
//   - Kind = "app_attached": AppID or AppName required (exactly one);
//     Image, Env, SecretRefs must be empty.
//   - Kind = "standalone":   Image required; AppID/AppName must be empty.
//
// Schedule is optional. When empty the job is on-demand only (POST /run).
// Timezone defaults server-side to Asia/Jakarta.
//
// PricingSlug picks a ResourceTemplate from the Apps catalog (the same
// catalog GET /api/v1/apps/plans returns); jobs share the App templates in
// MVP. The resolved template version id is pinned on the job at create.
type CreateJobRequest struct {
	Name        string  `json:"name"`
	Kind        JobKind `json:"kind"`
	PricingSlug string  `json:"pricing_slug"`

	AppID   *uint  `json:"app_id,omitempty"`
	AppName string `json:"app_name,omitempty"`

	Image       string                `json:"image,omitempty"`
	Command     []string              `json:"command,omitempty"`
	Args        []string              `json:"args,omitempty"`
	Env         []EnvironmentVariable `json:"env,omitempty"`
	SecretRefs  []JobSecretRef        `json:"secret_refs,omitempty"`

	Schedule              string               `json:"schedule,omitempty"`
	Timezone              string               `json:"timezone,omitempty"`
	ConcurrencyPolicy     JobConcurrencyPolicy `json:"concurrency_policy,omitempty"`
	ActiveDeadlineSeconds int                  `json:"active_deadline_seconds,omitempty"`
	BackoffLimit          int                  `json:"backoff_limit,omitempty"`
}

// UpdateJobRequest is the body for PATCH /api/v1/jobs/:id. Honors If-Match.
// Field semantics match CreateJobRequest. Fields left at zero value are
// untouched. Kind is immutable post-create (server returns 400
// JOB_VALIDATION_FAILED).
type UpdateJobRequest struct {
	PricingSlug string `json:"pricing_slug,omitempty"`

	Command    []string              `json:"command,omitempty"`
	Args       []string              `json:"args,omitempty"`
	Env        []EnvironmentVariable `json:"env,omitempty"`
	SecretRefs []JobSecretRef        `json:"secret_refs,omitempty"`

	Schedule              *string               `json:"schedule,omitempty"`
	Timezone              string                `json:"timezone,omitempty"`
	ConcurrencyPolicy     JobConcurrencyPolicy  `json:"concurrency_policy,omitempty"`
	ActiveDeadlineSeconds int                   `json:"active_deadline_seconds,omitempty"`
	BackoffLimit          *int                  `json:"backoff_limit,omitempty"`
}

// JobResourceTemplate is the slimmed projection of the pinned
// ResourceTemplateVersion exposed on the job. Same fields the apps surface
// publishes; admin-only fields (base_price, margin_*) are intentionally
// stripped.
type JobResourceTemplate struct {
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	CPUvCPU      string `json:"cpu_vcpu"`        // decimal string
	MemoryMB     int    `json:"memory_mb"`
	PricePerHour string `json:"price_per_hour"`  // decimal string
}

// JobResponse is the detail shape returned by GET /api/v1/jobs/:id and
// PATCH /api/v1/jobs/:id. Server sets ETag from UpdatedAt; echo it back in
// If-Match on PATCH for optimistic concurrency.
//
// NextRunTimes is computed server-side from Schedule + Timezone (up to 3
// upcoming runs, in the job's timezone). Empty when the job is on-demand
// only or suspended.
type JobResponse struct {
	ID   uint    `json:"id"`
	Name string  `json:"name"`
	Kind JobKind `json:"kind"`

	AppID   *uint   `json:"app_id,omitempty"`
	AppName *string `json:"app_name,omitempty"`

	Image      string                `json:"image,omitempty"`
	Command    []string              `json:"command,omitempty"`
	Args       []string              `json:"args,omitempty"`
	Env        []EnvironmentVariable `json:"env,omitempty"`
	SecretRefs []JobSecretRef        `json:"secret_refs,omitempty"`

	Schedule              string               `json:"schedule,omitempty"`
	Timezone              string               `json:"timezone"`
	ConcurrencyPolicy     JobConcurrencyPolicy `json:"concurrency_policy"`
	ActiveDeadlineSeconds int                  `json:"active_deadline_seconds"`
	BackoffLimit          int                  `json:"backoff_limit"`

	ResourceTemplate     JobResourceTemplate `json:"resource_template"`
	Suspended            bool                `json:"suspended"`
	DeploymentStatus     JobDeploymentStatus `json:"deployment_status"`
	LastExecutionAt      *time.Time          `json:"last_execution_at,omitempty"`
	NextRunTimes         []time.Time         `json:"next_run_times,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// JobListItem is the slimmed list-item shape returned by GET /api/v1/jobs.
// Heavy fields (Env, SecretRefs, NextRunTimes) are omitted to keep list
// responses cheap; fetch by id-or-name for the full JobResponse.
type JobListItem struct {
	ID                uint                `json:"id"`
	Name              string              `json:"name"`
	Kind              JobKind             `json:"kind"`
	AppID             *uint               `json:"app_id,omitempty"`
	Schedule          string              `json:"schedule,omitempty"`
	Timezone          string              `json:"timezone"`
	Suspended         bool                `json:"suspended"`
	DeploymentStatus  JobDeploymentStatus `json:"deployment_status"`
	LastExecutionAt   *time.Time          `json:"last_execution_at,omitempty"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
}

// ResponseJobAsync is the 202 payload for POST/PATCH/DELETE on /api/v1/jobs.
// Clients poll GET /api/v1/operations/:operation_id for terminal state.
type ResponseJobAsync struct {
	ID               uint                `json:"id"`
	Name             string              `json:"name"`
	DeploymentStatus JobDeploymentStatus `json:"deployment_status"`
	OperationID      string              `json:"operation_id"`
	UpdatedAt        time.Time           `json:"updated_at"`
}

// JobExecution is one row in the executions history, the source of truth
// for billing and audit. BilledAmount is a decimal string; null until the
// jobs_charger worker has settled the execution.
type JobExecution struct {
	ID            uint                `json:"id"`
	JobID         uint                `json:"job_id"`
	Trigger       JobExecutionTrigger `json:"trigger"`
	K8sJobName    string              `json:"k8s_job_name"`
	Status        JobExecutionStatus  `json:"status"`
	ExitCode      *int                `json:"exit_code,omitempty"`
	PodStartedAt  *time.Time          `json:"pod_started_at,omitempty"`
	PodFinishedAt *time.Time          `json:"pod_finished_at,omitempty"`
	DurationMS    *int64              `json:"duration_ms,omitempty"`
	CPUvCPU       string              `json:"cpu_vcpu,omitempty"`  // decimal string snapshot
	MemoryMB      int                 `json:"memory_mb,omitempty"`
	BilledAmount  *string             `json:"billed_amount,omitempty"`
	CreatedAt     time.Time           `json:"created_at"`
}

// RunJobResponse is the 202 payload for POST /api/v1/jobs/:id/run. Clients
// poll GET /api/v1/jobs/:id/executions/:execution_id for terminal state.
type RunJobResponse struct {
	ExecutionID uint               `json:"execution_id"`
	Status      JobExecutionStatus `json:"status"`
	OperationID string             `json:"operation_id,omitempty"`
}
