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
)
