package types

import "time"

// AppDeploymentStatus is the async deployment-operation state surfaced on
// app responses. Empty string is the unset/idle state.
type AppDeploymentStatus string

const (
	AppDeploymentStatusNone      AppDeploymentStatus = ""
	AppDeploymentStatusDeploying AppDeploymentStatus = "deploying"
	AppDeploymentStatusDeleting  AppDeploymentStatus = "deleting"
	AppDeploymentStatusFailed    AppDeploymentStatus = "failed"
)

// DomainVerificationStatus is the verification lifecycle of a custom domain
// attached to an app.
type DomainVerificationStatus string

const (
	DomainVerificationStatusPending  DomainVerificationStatus = "pending"
	DomainVerificationStatusVerified DomainVerificationStatus = "verified"
	DomainVerificationStatusFailed   DomainVerificationStatus = "failed"
)

// AutoscalingConfig is the optional HPA configuration on an app.
//
// When Enabled is true:
//   - MaxReplicas MUST be greater than MinReplicas
//   - At least one of CPUTargetPercentage / MemoryTargetPercentage MUST be set
//   - Replicas (in the parent request) MUST sit within [MinReplicas, MaxReplicas]
//
// The server returns 400 VALIDATION_FAILED on rule violations.
type AutoscalingConfig struct {
	Enabled                bool `json:"enabled"`
	MinReplicas            int  `json:"min_replicas,omitempty"`
	MaxReplicas            int  `json:"max_replicas,omitempty"`
	CPUTargetPercentage    *int `json:"cpu_target_percentage,omitempty"`
	MemoryTargetPercentage *int `json:"memory_target_percentage,omitempty"`
}

// HealthCheck declares an HTTP or TCP probe Kumo uses for liveness/readiness.
type HealthCheck struct {
	Type string `json:"type"`           // "http" or "tcp"
	Path string `json:"path,omitempty"` // required when Type == "http"
	Port uint16 `json:"port,omitempty"` // defaults to the app's Port when unset
}

// SecretFileMountType enumerates the kinds of secret a SecretFileMount can
// reference. Today only secret_file is supported.
type SecretFileMountType string

const (
	SecretFileMountTypeSecretFile SecretFileMountType = "secret_file"
)

// SecretVar mounts a Kumo Secret as an environment variable inside the app.
type SecretVar struct {
	SecretId           uint `json:"secret_id"`
	RestartWhenUpdated bool `json:"restart_when_updated"`
}

// SecretFileMount projects the contents of a Kumo Secret as a file at
// MountTo inside the app's container.
type SecretFileMount struct {
	Type    SecretFileMountType `json:"type"`
	MountTo string              `json:"mount_to"`

	// SecretId is required when Type == "secret_file".
	SecretId           uint `json:"secret_id,omitempty"`
	RestartWhenUpdated bool `json:"restart_when_updated"`
}

// BaseCreateApp is the shared shape for app create/update. Server validation:
//   - Name: 6..100 chars
//   - Image: must be a parseable Docker image reference
//   - Port: 1..65535
//   - Replicas: >= 1 (and within Autoscaling bounds when Autoscaling.Enabled)
type BaseCreateApp struct {
	Name        string             `json:"name"`
	Image       string             `json:"image"`
	Port        uint16             `json:"port"`
	IsExposed   bool               `json:"is_exposed"`
	Replicas    int                `json:"replicas"`
	Autoscaling *AutoscalingConfig `json:"autoscaling,omitempty"`
}

// CreateAppRequest is the body for POST /api/v1/apps. Supports Idempotency-Key.
type CreateAppRequest struct {
	BaseCreateApp
	EnvironmentVariables []EnvironmentVariable `json:"environment_variables,omitempty"`

	PricingSlug string `json:"pricing_slug"`

	RegistryCredentialId uint              `json:"registry_credential_id,omitempty"`
	TLSSecretId          *uint             `json:"tls_secret_id,omitempty"`
	SecretVars           []SecretVar       `json:"secret_vars,omitempty"`
	SecretFileMounts     []SecretFileMount `json:"secret_file_mounts,omitempty"`
	HealthCheck          *HealthCheck      `json:"healthcheck,omitempty"`
}

// UpdateAppRequest is the body for PATCH /api/v1/apps/:id. Same shape as
// create, kept as a distinct type so future fields can diverge.
type UpdateAppRequest struct {
	BaseCreateApp
	EnvironmentVariables []EnvironmentVariable `json:"environment_variables,omitempty"`

	PricingSlug          string            `json:"pricing_slug"`
	RegistryCredentialId uint              `json:"registry_credential_id,omitempty"`
	TLSSecretId          *uint             `json:"tls_secret_id,omitempty"`
	SecretVars           []SecretVar       `json:"secret_vars,omitempty"`
	SecretFileMounts     []SecretFileMount `json:"secret_file_mounts,omitempty"`
	HealthCheck          *HealthCheck      `json:"healthcheck,omitempty"`
}

// CreateAppResponse is the 202 Accepted payload returned by POST /api/v1/apps.
// OperationID is the polling handle for GET /api/v1/apps/:id/operations/:opid.
// UpdatedAt seeds the ETag (W/"<unix-nano-hex>") that PATCH should echo via
// If-Match to detect concurrent writes.
type CreateAppResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	GenerateAppName  string    `json:"generate_app_name"`
	DeploymentStatus string    `json:"deployment_status"`
	OperationID      string    `json:"operation_id,omitempty"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// AddCustomDomainRequest is the body for POST /api/v1/apps/:id/custom-domain.
// Server validates Domain as an FQDN and rejects platform-owned zones.
type AddCustomDomainRequest struct {
	Domain string `json:"domain"`
}

// CustomDomainInfo is the custom-domain summary embedded in app detail
// responses.
type CustomDomainInfo struct {
	Domain             string `json:"domain"`
	VerificationStatus string `json:"verification_status"`
}

// VolumeInfo is the attached-volume summary embedded in AppByIdResponse.
type VolumeInfo struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	SizeGB    int    `json:"size_gb"`
	MountPath string `json:"mount_path"`
	Status    string `json:"status"`
}

// HPAStatusInfo is the runtime HorizontalPodAutoscaler snapshot included
// when an app has autoscaling enabled.
type HPAStatusInfo struct {
	CurrentReplicas int32   `json:"current_replicas"`
	DesiredReplicas int32   `json:"desired_replicas"`
	MinReplicas     int32   `json:"min_replicas"`
	MaxReplicas     int32   `json:"max_replicas"`
	CurrentCPUUsage *int32  `json:"current_cpu_usage,omitempty"`
	CurrentMemUsage *int32  `json:"current_mem_usage,omitempty"`
	LastScaleTime   *string `json:"last_scale_time,omitempty"`
}

// AppByIdResponse is the full app detail returned by GET /api/v1/apps/:id.
// The server sets ETag from UpdatedAt; echo it back in If-Match on PATCH.
type AppByIdResponse struct {
	Id uint `json:"id"`
	CreateAppRequest
	GeneratedSubDomain string `json:"generated_sub_domain"`
	TLSSecretId        *uint  `json:"tls_secret_id,omitempty"`

	// Source distinguishes a normal "registry-image" app from a "git-build"
	// app whose image is produced by the platform. For git-build apps the
	// Image field is system-owned and cannot be changed via PATCH.
	Source AppSource `json:"source"`

	// Language is the build language preset for git-build apps ("auto" or a
	// specific language: nodejs/python/go/java/ruby/php/dotnet, or "static").
	// Empty/omitted for registry-image apps.
	Language string `json:"language,omitempty"`

	// OutputDir and BuildCommand are the static-site preset config for git-build
	// apps (Language == "static"): the directory nginx serves and the npm build
	// script run before serving. Empty/omitted for non-static and registry-image apps.
	OutputDir    string `json:"output_dir,omitempty"`
	BuildCommand string `json:"build_command,omitempty"`

	// Suspension state
	IsSuspended   bool   `json:"is_suspended"`
	SuspendReason string `json:"suspend_reason,omitempty"`

	// Runtime status
	AppStatus           string `json:"app_status"`
	StatusMessage       string `json:"status_message"`
	DesiredReplicas     int    `json:"desired_replicas"`
	ReadyReplicas       int    `json:"ready_replicas"`
	AvailableReplicas   int    `json:"available_replicas"`
	UpdatedReplicas     int    `json:"updated_replicas"`
	UnavailableReplicas int    `json:"unavailable_replicas"`
	TotalPods           int    `json:"total_pods"`
	PendingPods         int    `json:"pending_pods"`
	RunningPods         int    `json:"running_pods"`
	FailedPods          int    `json:"failed_pods"`
	IsDeploying         bool   `json:"is_deploying"`
	HasReplicaFailure   bool   `json:"has_replica_failure"`

	CustomDomain *CustomDomainInfo `json:"custom_domain,omitempty"`
	InternalDNS  string            `json:"internal_dns"`
	HPAStatus    *HPAStatusInfo    `json:"hpa_status,omitempty"`
	Volume       *VolumeInfo       `json:"volume,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AppImageResponse is the list-item shape returned by GET /api/v1/apps (both
// paginated and legacy unbounded forms). Same shape as AppByIdResponse minus
// the per-pod runtime stats.
type AppImageResponse struct {
	Id uint `json:"id"`
	CreateAppRequest
	GeneratedSubDomain string `json:"generated_sub_domain"`

	// Source distinguishes "registry-image" from "git-build" apps. See
	// AppByIdResponse.Source.
	Source AppSource `json:"source"`

	IsSuspended   bool   `json:"is_suspended"`
	SuspendReason string `json:"suspend_reason,omitempty"`

	AppStatus       string `json:"app_status"`
	StatusMessage   string `json:"status_message"`
	DesiredReplicas int    `json:"desired_replicas"`
	ReadyReplicas   int    `json:"ready_replicas"`

	CustomDomain *CustomDomainInfo `json:"custom_domain,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ValidateImagePullableRequest is the body for POST /api/v1/apps/validate-image.
type ValidateImagePullableRequest struct {
	Image                  string  `json:"image"`
	RegistryCredentialId   uint    `json:"registry_credential_id,omitempty"`
	RegistryCredentialName *string `json:"registry_credential_name,omitempty"`
}

// ValidateImagePullableResponse reports per-architecture pull success for
// the requested image + credential combination.
type ValidateImagePullableResponse struct {
	LinuxAmd64 bool `json:"linux_amd64"`
	LinuxArm64 bool `json:"linux_arm64"`
}

// AppOperationActionType discriminates the kind of async operation a
// row in /apps/:id/operations represents.
type AppOperationActionType string

const (
	AppOperationActionCreate  AppOperationActionType = "create"
	AppOperationActionUpdate  AppOperationActionType = "update"
	AppOperationActionDelete  AppOperationActionType = "delete"
	AppOperationActionRestart AppOperationActionType = "restart"
	AppOperationActionStart   AppOperationActionType = "start"
	AppOperationActionStop    AppOperationActionType = "stop"
)

// AppOperationStatus is the lifecycle state of an async app operation.
// "succeeded", "failed", and "cancelled" are terminal.
type AppOperationStatus string

const (
	AppOperationStatusQueued     AppOperationStatus = "queued"
	AppOperationStatusInProgress AppOperationStatus = "in_progress"
	AppOperationStatusSucceeded  AppOperationStatus = "succeeded"
	AppOperationStatusFailed     AppOperationStatus = "failed"
	AppOperationStatusCancelled  AppOperationStatus = "cancelled"
)

// AppOperation is one row of /api/v1/apps/:id/operations. Returned in
// CreateAppResponse.OperationID (UUID) and as the polling response from
// GET /api/v1/apps/:id/operations/:operation_id.
//
// ErrorCode is the wire code (see github.com/kumobase/kumo-go/codes/apps)
// that failed the operation, when Status == "failed".
type AppOperation struct {
	OperationID string                 `json:"operation_id"` // UUID
	AppID       uint                   `json:"app_id"`
	ActionType  AppOperationActionType `json:"action_type"`
	Status      AppOperationStatus     `json:"status"`
	ErrorCode   *string                `json:"error_code,omitempty"`
	ErrorMsg    *string                `json:"error_message,omitempty"`
	RequestedBy *string                `json:"requested_by,omitempty"` // "api_key" / "bearer_jwt" / "cookie"
	QueuedAt    time.Time              `json:"queued_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
}
