# Changelog

All notable changes to kumo-go are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.26.0]

### Added
- `types/build.go` — `DockerfilePath` on `CreateGitBuildAppRequest` and
  `UpdateBuildConfigRequest` (the "dockerfile" build preset). New discovery
  types `BuildersResponse`, `BuilderOption`, `LanguageOption`.
- `types/apps.go` — `DockerfilePath` on `AppByIdResponse` (read-back of the
  configured Dockerfile path).
- `client/build.go` — `Builds().ListBuilders()` for `GET /api/v1/builders`
  (the selectable builder kinds + CNB language presets; lets clients stop
  hardcoding the list).
- `codes/build.go` — `BUILD_INVALID_DOCKERFILE_PATH`, `BUILD_NO_DOCKERFILE`
  (the "dockerfile" preset), and `BUILD_NO_RAILPACK_PLAN` (the "auto"/"railpack"
  preset when Railpack can't detect a buildable project).

### Changed
- `Language` now also accepts `"auto"` (Dockerfile-if-present-else-Railpack, the
  new zero-config default), `"railpack"`, `"dockerfile"`, and `"cnb"` (the
  legacy buildpack auto-detect). Doc comments updated; no wire-string changes to
  existing values.

## [v0.24.0]

### Added
- `codes/auth.go` — six new refresh-flow wire codes returned by the new
  `POST /api/v1/auth/refresh` endpoint: `REFRESH_TOKEN_MISSING`,
  `REFRESH_TOKEN_INVALID`, `REFRESH_TOKEN_EXPIRED`, `REFRESH_TOKEN_REVOKED`,
  `REFRESH_TOKEN_REUSED`, and `REFRESH_ACCOUNT_INACTIVE`.
- `types/auth.go` — `RefreshRequest` and `RefreshResponse` DTOs for the
  refresh endpoint (rotating refresh tokens).
- `client/auth_service.go` — `Client.Auth()` with `Refresh`, `Logout`, and
  `LogoutAll` methods for session-lifecycle management from a CLI/SDK client.
  Automatic single-flight "refresh on 401" remains a calling-application
  concern (browser cookie clients), not built into this stateless SDK.

## [v0.17.0]

### Changed
- **BREAKING:** `types/registry.go` — removed the `SoftDeleteDays` field from
  `CreateRepositoryRequest`, `UpdateRepositoryRequest`, and `RepositoryResponse`.
  The registry purge window is now a fixed internal policy (7 days) and is no
  longer user-configurable or surfaced on the wire. Clients that previously set
  or read `soft_delete_days` must drop it; the server silently ignores the field
  on requests and never emits it on responses.

### Deprecated
- `codes/registry.go` — `RegistryInvalidSoftDeleteDays` is retained for wire
  stability but is no longer returned by the server (the validation it guarded
  was removed alongside the field).

## [v0.12.1]

### Added
- `types/apps.go` — `CreateAppVolume` and a new optional `Volume *CreateAppVolume`
  field on `CreateAppRequest`. Lets `POST /api/v1/apps` bind an existing,
  unattached volume to the app in the same request (mounted into the app's
  first deployment), identified by exactly one of `VolumeID` / `VolumeName`.
  Reuses existing wire codes: `VALIDATION_FAILED` (shape / git-build rejection),
  `APP_VOLUME_CONFLICT` (replicas>1 or autoscaling, or volume already attached),
  `VOLUME_NOT_FOUND`. Additive and backward compatible — omitting `volume`
  preserves the prior wire shape byte-for-byte.

## [v0.12.0]

### Added
- `types/jobs.go` — Jobs product surface: `Job`, `JobKind` (`app_attached` /
  `standalone`), `JobDeploymentStatus`, `JobConcurrencyPolicy`,
  `JobExecutionTrigger`, `JobExecutionStatus`, `JobSecretRef`,
  `JobResourceTemplate`, `CreateJobRequest`, `UpdateJobRequest`, `JobResponse`,
  `JobListItem`, `ResponseJobAsync`, `JobExecution`, `RunJobResponse`.
- `codes/jobs.go` — 22 stable wire codes for `/api/v1/jobs/*`:
  `JOB_NOT_FOUND`, `JOB_EXECUTION_NOT_FOUND`, `JOB_OPERATION_NOT_FOUND`,
  `JOB_DEPLOYMENT_IN_PROGRESS`, `JOB_ALREADY_SUSPENDED`, `JOB_NOT_SUSPENDED`,
  `JOB_QUOTA_EXCEEDED`, `JOB_INSUFFICIENT_BALANCE`, `JOB_SCHEDULE_INVALID`,
  `JOB_SCHEDULE_TOO_FREQUENT`, `JOB_TIMEZONE_INVALID`, `JOB_KIND_INVALID`,
  `JOB_APP_REQUIRED`, `JOB_APP_NOT_FOUND`, `JOB_IMAGE_REQUIRED`,
  `JOB_CONCURRENCY_UNSUPPORTED`, `JOB_INVALID_PRICING_SLUG`,
  `JOB_VALIDATION_FAILED`, `JOB_UNAUTHORIZED`, `JOB_INVALID_ID`,
  `JOB_INVALID_REQUEST_BODY`, `JOB_INTERNAL_ERROR`.
- `client/jobs.go` — `JobsService` with `Create`, `List`, `Get`, `GetByName`,
  `Update`, `Delete`, `RunNow`, `Suspend`, `Resume`, `ListExecutions`,
  `GetExecution`. Async mutations return `*ResponseJobAsync` (202 +
  `operation_id`); reads return `(*JobResponse, etag, error)` for use with
  `IfMatch(etag)` on subsequent `Update`.

## [v0.7.0]

### Added
- `types/apps.go` — neutral instance/autoscaling fields on `AppByIdResponse`:
  `TotalInstances`, `PendingInstances`, `RunningInstances`, `FailedInstances`,
  `HasFailure`, and `AutoscalingStatus` (plus a new `AutoscalingStatus` type).
  The server populates both these and the old pod-named fields during the
  deprecation window so existing callers keep working.
- `types/apps.go` — `GitBuildInfo` (`repo_full_name`, `branch`) and
  `BuildSummary` (latest build snapshot). Embedded on `AppByIdResponse` as
  `GitBuild` and `LatestBuild`; both nil for `registry-image` apps. Lets
  callers render a git-build app's source + last-build state from one
  `GET /apps/:id` instead of also calling `/apps/:id/builds`.

### Deprecated
- `AppByIdResponse.TotalPods` / `PendingPods` / `RunningPods` / `FailedPods` —
  use the matching `*Instances` fields.
- `AppByIdResponse.HasReplicaFailure` — use `HasFailure`.
- `AppByIdResponse.HPAStatus` — use `AutoscalingStatus`. The `HPAStatusInfo`
  type is now an alias for `AutoscalingStatus` so existing code compiles.

Deprecated fields are still populated by the server and round-trip identically
to today's response; they will be removed in a future minor.

### Changed
- `version.SDKVersion` bumped to `v0.7.0`.
- Comments on `AppByIdResponse.InternalDNS` and the autoscaling/runtime fields
  no longer name internal infrastructure terms.

## [v0.4.1]

### Added
- `codes/registry.go` — `REGISTRY_REPOSITORY_SYSTEM_OWNED`, returned when a
  user tries to delete a platform-owned (auto-provisioned) repository, such as
  the registry repo backing a git-build app.

### Changed
- `version.SDKVersion` bumped to `v0.4.1`.

## [v0.4.0]

### Added
- `types/build.go` — git-build app surface: `BuildResponse`, `BuildStatus`
  enum, `CreateGitBuildAppRequest`, and the `AppSource` enum
  (`registry-image` | `git-build`).
- `types/apps.go` — `Source` field on `AppByIdResponse` and
  `AppImageResponse` (additive; defaults to `registry-image`).
- `codes/build.go` — build wire codes: `BUILD_NOT_FOUND`,
  `BUILD_APP_IMAGE_IMMUTABLE`, `BUILD_CONNECTION_REQUIRED`,
  `BUILD_CONNECTION_IN_USE`, `BUILD_SOURCE_UNAVAILABLE`,
  `BUILD_ALREADY_RUNNING`, `BUILD_PROVIDER_ERROR`, `BUILD_INTERNAL_ERROR`.
- `client/build.go` — `client.Builds()` service: `CreateGitBuildApp(connID, …)`
  plus `List` / `Get` / `Rebuild` / `Cancel` for an app's builds.

### Changed
- `version.SDKVersion` bumped to `v0.4.0`.

## [v0.2.0]

### Added
- `client/` — typed HTTP client for the Kumo API. Wraps `net/http` (zero
  transitive deps), composes with `types/` and `codes/`.
  - `client.New(baseURL, …)` with functional options: `WithAPIKey`,
    `WithJWT`, `WithHTTPClient`, `WithRetries`, `WithUserAgent`,
    `WithLogger`.
  - Resource-grouped services: `client.Apps()`, `client.VPS()`,
    `client.Secrets()`, `client.Volumes()`, `client.Registry().Orgs()`,
    `client.Registry().Repos(slug)`, `client.Billing()`, `client.Profile()`,
    `client.APIKeys()`.
  - Auto `Idempotency-Key` (UUIDv4 per write call, held constant across
    retries). Override via `WithIdempotencyKey("…")`.
  - Optional `IfMatch("W/\"…\"")` write option for `ETag`-gated PATCHes.
    Returns `ErrETagMismatch` (`errors.Is`-matchable) on 412.
  - Smart retry defaults (5 attempts): network errors, 5xx, 429 honoring
    `Retry-After`, 409 `IDEMPOTENCY_IN_PROGRESS`. Never retries 4xx
    (other than the two above) or 409 `IDEMPOTENCY_KEY_CONFLICT`. Disable
    via `WithRetries(0)`. `ctx` deadline always wins.
  - `APIError` typed error wrapping `StructureResponse`: `IsCode(err, …)`,
    `IsNotFound(err)`, `IsConflict(err)` predicates.
  - Async polling: `client.Apps().PollOperation(…)` for operation_id
    polling and generic `PollResource[T](…)` for resources where status
    lives on the resource itself (VPS `ActionStatus`, Volume `Status`).

### Changed
- `version.SDKVersion` bumped to `v0.2.0`.

### Unchanged (still v0.1.0 wire contract)
- All `types/*` shapes — server contract unchanged.
- All `codes/*` strings — server contract unchanged.

## [v0.1.0]

### Added
- Initial public SDK surface: types and stable wire codes for every
  user-facing `/api/v1/*` endpoint a customer API key can reach.
- `types/common.go` — shared response envelope (`StructureResponse`),
  pagination metadata (`Meta`), and request/response shapes shared across
  modules.
- `codes/` — stable error codes per module plus cross-cutting codes
  (`IDEMPOTENCY_KEY_CONFLICT`, `ETAG_MISMATCH`, `API_KEY_SESSION_REQUIRED`,
  …).
- `version/version.go` — `SDKVersion` constant for `User-Agent`/compat
  headers.
