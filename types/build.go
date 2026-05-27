package types

import "time"

// AppSource identifies where an app's container image comes from. The default
// "registry-image" app deploys a user-supplied image; a "git-build" app has
// its image produced by the platform from a connected git repository and is
// therefore system-owned (read-only over the API).
type AppSource string

const (
	AppSourceRegistryImage AppSource = "registry-image"
	AppSourceGitBuild      AppSource = "git-build"
)

// BuildStatus is the lifecycle of a single git-build run.
//
//   - pending    — queued, no builder pod yet
//   - running    — builder pod is executing
//   - succeeded  — image built + pushed; the app's image was set to the digest
//   - failed     — clone/detect/build/push failed (see Error / log)
//   - canceled   — explicitly canceled by the user
//   - superseded — a newer push for the same app replaced this build
//
// succeeded, failed, canceled, and superseded are terminal.
type BuildStatus string

const (
	BuildStatusPending    BuildStatus = "pending"
	BuildStatusRunning    BuildStatus = "running"
	BuildStatusSucceeded  BuildStatus = "succeeded"
	BuildStatusFailed     BuildStatus = "failed"
	BuildStatusCanceled   BuildStatus = "canceled"
	BuildStatusSuperseded BuildStatus = "superseded"
)

// BuildResponse is the wire shape for one build of a git-build app, returned
// by the /api/v1/apps/:id/builds endpoints.
//
// ImageDigest is the pushed image's content digest (sha256:…), set only on a
// successful build. LogURL is a short-lived presigned link to the plain-text
// build log on object storage; it may be empty when no log was persisted (e.g.
// the build never started or the upload failed) and should be re-fetched via
// Get rather than cached, as it expires.
type BuildResponse struct {
	ID          uint        `json:"id"`
	AppID       uint        `json:"app_id"`
	CommitSHA   string      `json:"commit_sha"`
	Ref         string      `json:"ref"` // e.g. "refs/heads/main"
	Status      BuildStatus `json:"status"`
	ImageDigest string      `json:"image_digest,omitempty"`
	LogURL      string      `json:"log_url,omitempty"`
	Error       string      `json:"error,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	StartedAt   *time.Time  `json:"started_at,omitempty"`
	FinishedAt  *time.Time  `json:"finished_at,omitempty"`
}

// CreateGitBuildAppRequest is the body for
// POST /api/v1/source-connections/:id/apps. It mirrors CreateAppRequest's app
// configuration but omits Image and RegistryCredentialId: a git-build app's
// image is produced by the platform and pushed to an auto-provisioned, system-
// owned registry repository, so neither is user-supplied.
//
// RepoFullName ("owner/repo") must be a repository the connection can access;
// Branch is the branch whose pushes trigger a build (HEAD is also built once
// on create so the app deploys without a dummy push). Supports Idempotency-Key.
type CreateGitBuildAppRequest struct {
	Name        string             `json:"name"`
	Port        uint16             `json:"port"`
	IsExposed   bool               `json:"is_exposed"`
	Replicas    int                `json:"replicas"`
	Autoscaling *AutoscalingConfig `json:"autoscaling,omitempty"`

	RepoFullName string `json:"repo_full_name"` // "owner/repo"
	Branch       string `json:"branch"`

	// Language is the build language preset. Empty or "auto" (default) lets the
	// platform auto-detect the language; a specific value (e.g. "nodejs",
	// "python", "go") pins the build to that language's buildpack. The special
	// value "static" builds a static site served by nginx (see OutputDir /
	// BuildCommand below).
	Language string `json:"language,omitempty"`

	// OutputDir and BuildCommand apply only to the "static" preset and are
	// otherwise ignored. OutputDir is the directory nginx serves (default ".",
	// the repo root for plain HTML; e.g. "dist"/"build" for a framework that
	// compiles to static output). BuildCommand is the npm script name run before
	// serving (e.g. "build") for framework→static; leave empty for pure static.
	// A static app is always served on port 8080 (the port field is forced).
	OutputDir    string `json:"output_dir,omitempty"`
	BuildCommand string `json:"build_command,omitempty"`

	EnvironmentVariables []EnvironmentVariable `json:"environment_variables,omitempty"`
	PricingSlug          string                `json:"pricing_slug"`
	TLSSecretId          *uint                 `json:"tls_secret_id,omitempty"`
	SecretVars           []SecretVar           `json:"secret_vars,omitempty"`
	SecretFileMounts     []SecretFileMount     `json:"secret_file_mounts,omitempty"`
	HealthCheck          *HealthCheck          `json:"healthcheck,omitempty"`
}
