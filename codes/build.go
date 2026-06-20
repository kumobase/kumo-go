package codes

// Build / git-build-app wire codes returned by the
// /api/v1/source-connections/:id/apps and /api/v1/apps/:id/builds endpoints.
// Mirrors the per-sentinel branches in modules/build/errors.go on the server.
const (
	// BuildNotFound — no build with the given id exists for the app, or the
	// app isn't owned by the caller.
	BuildNotFound = "BUILD_NOT_FOUND"

	// BuildAppImageImmutable — attempted to change the image of a git-build
	// app via PATCH /api/v1/apps/:id. The image is system-owned; trigger a
	// rebuild instead.
	BuildAppImageImmutable = "BUILD_APP_IMAGE_IMMUTABLE"

	// BuildConnectionRequired — the operation needs a usable source
	// connection that the app does not have (e.g. rebuild on an app whose
	// connection is gone).
	BuildConnectionRequired = "BUILD_CONNECTION_REQUIRED"

	// BuildConnectionInUse — a source connection cannot be disconnected
	// while git-build apps are still bound to it. Delete the apps first.
	BuildConnectionInUse = "BUILD_CONNECTION_IN_USE"

	// BuildSourceUnavailable — the build cannot proceed because the source
	// is unreachable: the GitHub installation was removed/suspended, or the
	// org's registry is suspended (settle the balance). Reconnect/settle and
	// retry.
	BuildSourceUnavailable = "BUILD_SOURCE_UNAVAILABLE"

	// BuildAlreadyRunning — a build for this app is already pending/running.
	// A new push supersedes it automatically; manual rebuilds are rejected
	// while one is active.
	BuildAlreadyRunning = "BUILD_ALREADY_RUNNING"

	// BuildProviderError — a call to the upstream git provider (GitHub)
	// failed while setting up the build. Usually transient; safe to retry.
	BuildProviderError = "BUILD_PROVIDER_ERROR"

	// BuildInternalError — unexpected server-side failure.
	BuildInternalError = "BUILD_INTERNAL_ERROR"

	// BuildTriggerRequired — Create or Update would leave the app with no
	// build trigger at all (both Branch and TagPattern empty). At least one
	// must be set.
	BuildTriggerRequired = "BUILD_TRIGGER_REQUIRED"

	// BuildInvalidTagPattern — supplied TagPattern is not a valid glob
	// (path.Match syntax). Common causes: unbalanced bracket like "v[".
	BuildInvalidTagPattern = "BUILD_INVALID_TAG_PATTERN"

	// BuildNeedsBranch — manual rebuild (POST /api/v1/apps/:id/builds) was
	// attempted on a tag-only app (Branch empty). Move/re-push a matching
	// tag instead, or PATCH the build-config to add a branch trigger.
	BuildNeedsBranch = "BUILD_NEEDS_BRANCH"

	// BuildLogNotAvailable — GET /api/v1/apps/:id/builds/:buildId/log-url was
	// requested for a build that has no persisted log (e.g. still
	// pending/running, never started, or the log upload failed).
	BuildLogNotAvailable = "BUILD_LOG_NOT_AVAILABLE"

	// BuildInvalidDockerfilePath — the supplied dockerfile_path is not a clean
	// relative path (it is absolute, or contains ".." traversal). Returned 400
	// at create/update time, before any build runs.
	BuildInvalidDockerfilePath = "BUILD_INVALID_DOCKERFILE_PATH"

	// BuildNoDockerfile — a "dockerfile" build ran but no file exists at the
	// configured dockerfile_path in the cloned repo. Surfaced as a failed build
	// (not a create-time error), since the path is only known at build time.
	BuildNoDockerfile = "BUILD_NO_DOCKERFILE"

	// BuildNoRailpackPlan — an "auto"/"railpack" build ran but Railpack could not
	// detect a buildable project (no Dockerfile and no recognized language).
	// Surfaced as a failed build (the repo is only inspected at build time).
	BuildNoRailpackPlan = "BUILD_NO_RAILPACK_PLAN"
)
