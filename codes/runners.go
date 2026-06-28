package codes

// Runner-module wire codes returned by /api/v1/runner-specs and
// /api/v1/runner-jobs (the VM-backed CI-runner product). Mirrors the
// per-sentinel branches in modules/runners/errors.go::handleRunnerError.
//
// The runner product is opt-in by connecting the GitHub App and using a known
// `kumo-*` label in runs-on — there are no user-managed pools. Capacity is
// governed by admin-managed system specs + a global admission gate, so codes
// for spec/backend administration are server-only (admin endpoints), and
// capacity outcomes (no_capacity_timeout, spot_capacity_unavailable) are job
// *conclusions*, not API error codes.
const (
	// RunnerSpecNotFound — no runner spec with the given id/label, or it is
	// disabled.
	RunnerSpecNotFound = "RUNNER_SPEC_NOT_FOUND"

	// RunnerJobNotFound — no runner job with the given id owned by the caller.
	RunnerJobNotFound = "RUNNER_JOB_NOT_FOUND"

	// RunnerValidationFailed — runner-specific request validation rejection.
	RunnerValidationFailed = "RUNNER_VALIDATION_FAILED"

	// RunnerUnauthorized — the caller may not access this job.
	RunnerUnauthorized = "RUNNER_UNAUTHORIZED"

	// RunnerInvalidID — a path id segment was not a valid identifier.
	RunnerInvalidID = "RUNNER_INVALID_ID"

	// RunnerInternalError — unexpected server-side failure.
	RunnerInternalError = "RUNNER_INTERNAL_ERROR"

	// --- GitLab CI provider ---------------------------------------------------
	// The runner product also fulfills GitLab CI jobs (jobs tagged with a known
	// `kumo-*` tag on a connected group/project). These codes surface on the
	// GitLab connect + webhook surfaces; they mirror the GitLab-specific
	// branches in modules/runners/errors.go and the GitLab connect handlers.

	// RunnerGitLabOAuthFailed — the OAuth authorization-code exchange or a
	// later token refresh failed; the connection cannot call GitLab. For an
	// expired/revoked grant the connection is suspended and must be reconnected.
	RunnerGitLabOAuthFailed = "RUNNER_GITLAB_OAUTH_FAILED"

	// RunnerGitLabTokenInvalid — the bootstrap-presented or stored GitLab
	// runner token was rejected (expired, revoked, or already consumed).
	RunnerGitLabTokenInvalid = "RUNNER_GITLAB_TOKEN_INVALID"

	// RunnerGitLabInstanceUnknown — the referenced GitLab instance is not
	// registered, or its base URL failed validation/connectivity.
	RunnerGitLabInstanceUnknown = "RUNNER_GITLAB_INSTANCE_UNKNOWN"

	// RunnerGitLabWebhookFailed — creating or deleting the group/project
	// webhook on GitLab failed (insufficient scope, or a provider error).
	RunnerGitLabWebhookFailed = "RUNNER_GITLAB_WEBHOOK_FAILED"

	// RunnerGitLabScopeDenied — the OAuth grant lacks a scope the action needs
	// (e.g. create_runner or api), or the user lost access to the namespace.
	RunnerGitLabScopeDenied = "RUNNER_GITLAB_SCOPE_DENIED"
)
