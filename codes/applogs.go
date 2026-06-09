package codes

// App-logs wire codes returned by GET /api/v1/apps/:id/logs. Mirrors the
// per-sentinel branches in modules/applogs/errors.go::handleAppLogsError on the
// server. APP_NOT_FOUND (codes/apps.go), AMBIGUOUS_NAME (codes/common.go) and
// INVALID_TIME_RANGE (codes/appmetrics.go) are reused for lookup and window
// validation and are intentionally not redeclared here.
const (
	// LogsBackendUnavailable — the logs backend (Grafana Loki) could not be
	// reached or refused the server-built query. The request was valid; retry
	// after a short delay. Returned as 503 (unreachable/timeout) or 502
	// (backend rejected the query).
	LogsBackendUnavailable = "LOGS_BACKEND_UNAVAILABLE"

	// InvalidLogFilter — a log query parameter that is not a time window was
	// invalid: `level` outside the allowlist (error/warn/info/debug),
	// `direction` not forward/backward, or `filter` too long. Returned as 400.
	InvalidLogFilter = "INVALID_LOG_FILTER"

	// AppLogsInternalError — an unexpected server error while assembling the
	// logs response. Returned as 500.
	AppLogsInternalError = "APP_LOGS_INTERNAL_ERROR"
)
