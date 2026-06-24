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
)
