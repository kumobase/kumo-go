# Changelog

All notable changes to kumo-go are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
