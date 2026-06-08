package codes

// App-metrics wire codes returned by GET /api/v1/apps/:id/metrics. Mirrors the
// per-sentinel branches in modules/appmetrics/errors.go::handleAppMetricsError
// on the server. APP_NOT_FOUND (codes/apps.go) and AMBIGUOUS_NAME
// (codes/common.go) are reused for the id-or-name lookup and are intentionally
// not redeclared here.
const (
	// MetricsBackendUnavailable — the metrics backend (Grafana Mimir) could
	// not be reached or refused the server-built query. The request was valid;
	// retry after a short delay. Returned as 503 (unreachable/timeout) or 502
	// (backend rejected the query).
	MetricsBackendUnavailable = "METRICS_BACKEND_UNAVAILABLE"

	// InvalidTimeRange — the `range` query parameter was not one of the allowed
	// presets (15m, 1h, 6h, 24h, 7d). Returned as 400.
	InvalidTimeRange = "INVALID_TIME_RANGE"

	// AppMetricsInternalError — an unexpected server error while assembling the
	// metrics response. Returned as 500.
	AppMetricsInternalError = "APP_METRICS_INTERNAL_ERROR"
)
