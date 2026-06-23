package codes

// RDS-metrics wire codes returned by GET /api/v1/rds/:id/metrics. Mirrors the
// per-sentinel branches in modules/rdsmetrics/errors.go::handleRDSMetricsError
// on the server. RDS_INSTANCE_NOT_FOUND (codes/rds.go), MetricsBackendUnavailable
// and InvalidTimeRange (codes/appmetrics.go) are reused for ownership lookup,
// backend failures and range validation, and are intentionally not redeclared
// here.
const (
	// RDSMetricsInternalError — an unexpected server error while assembling the
	// metrics response. Returned as 500.
	RDSMetricsInternalError = "RDS_METRICS_INTERNAL_ERROR"
)
