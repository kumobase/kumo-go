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
	// LogURL is deprecated and no longer populated by the list or detail
	// endpoints (it always presigned a short-lived URL the caller usually
	// never used). Fetch a fresh URL on demand via Builds().GetLogURL.
	LogURL string `json:"log_url,omitempty"`
	Error  string `json:"error,omitempty"`
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
// RepoFullName ("owner/repo") must be a repository the connection can access.
// At least one of Branch or TagPattern must be set:
//   - Branch: exact-match branch name; pushes trigger a build. HEAD is also
//     built once on create so the app deploys without a dummy push.
//   - TagPattern: glob pattern (path.Match syntax: *, ?, [abc]) matched
//     against the bare tag name (no refs/tags/ prefix). Tag pushes whose
//     name matches trigger a build. No build runs on create for tag-only
//     apps — push or move a matching tag to deploy.
//
// Both may be set on the same app (build on main pushes AND on vX.Y.Z tags).
// Server returns 400 BUILD_TRIGGER_REQUIRED if both are empty, and 400
// BUILD_INVALID_TAG_PATTERN if TagPattern fails glob validation.
//
// Supports Idempotency-Key.
type CreateGitBuildAppRequest struct {
	Name        string             `json:"name"`
	Port        uint16             `json:"port"`
	IsExposed   bool               `json:"is_exposed"`
	Replicas    int                `json:"replicas"`
	Autoscaling *AutoscalingConfig `json:"autoscaling,omitempty"`

	RepoFullName string `json:"repo_full_name"` // "owner/repo"
	Branch       string `json:"branch,omitempty"`
	TagPattern   string `json:"tag_pattern,omitempty"` // glob, e.g. "v*", "release/*"

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

// UpdateBuildConfigRequest is the body for PATCH /api/v1/apps/:id/build-config.
// It updates the build preset + trigger config of an existing git-build app.
//
// PATCH semantics for Branch and TagPattern: nil pointer / absent key = no
// change; non-nil empty string = clear that trigger; non-nil non-empty = set
// it. The Language / OutputDir / BuildCommand strings stay wholesale-set
// (empty clears) for back-compat with the v0.7.x shape; pass the full
// intended value.
//
// Changes apply on the NEXT build (this does not trigger one). Switching
// Language to "static" forces the app's port to 8080 (nginx). Clearing both
// the branch and tag triggers (evaluated against the resulting state) returns
// 400 BUILD_TRIGGER_REQUIRED. Setting a branch on a previously tag-only app
// re-enables manual rebuild. Supports optional If-Match for optimistic
// concurrency.
type UpdateBuildConfigRequest struct {
	Language     string  `json:"language,omitempty"`      // "auto" (default) | a language | "static"
	OutputDir    string  `json:"output_dir,omitempty"`    // static only → BP_WEB_SERVER_ROOT
	BuildCommand string  `json:"build_command,omitempty"` // static only → BP_NODE_RUN_SCRIPTS (npm script)
	// Branch is exact-match (not a glob). nil = no change, "" = clear the
	// branch trigger, non-empty = set. Cross-field (at-least-one trigger)
	// checked server-side.
	Branch     *string `json:"branch,omitempty"`
	TagPattern *string `json:"tag_pattern,omitempty"` // glob; nil = no change, "" = clear
}

// BuildLogURLResponse is the body of GET /api/v1/apps/:id/builds/:buildId/log-url.
// LogURL is a freshly-minted, short-lived presigned link to the build's
// plain-text log. Fetch it only when you're about to open the log — it expires
// quickly and is not cacheable. Returns 404 BUILD_LOG_NOT_AVAILABLE when the
// build has no persisted log (still pending/running, or the upload failed).
type BuildLogURLResponse struct {
	LogURL string `json:"log_url"`
}
