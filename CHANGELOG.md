# Changelog

All notable changes to kumo-go are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
