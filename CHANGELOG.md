# Changelog

All notable changes to kumo-go are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
