package codes

// RDS-logs wire codes returned by GET /api/v1/rds/:id/logs. Mirrors the
// per-sentinel branches in modules/rdslogs/errors.go::handleRDSLogsError on the
// server. RDS_INSTANCE_NOT_FOUND (codes/rds.go), LogsBackendUnavailable and
// InvalidLogFilter (codes/applogs.go) and InvalidTimeRange (codes/appmetrics.go)
// are reused for ownership lookup, backend failures and window/filter validation,
// and are intentionally not redeclared here.
const (
	// RDSLogsInternalError — an unexpected server error while assembling the logs
	// response. Returned as 500.
	RDSLogsInternalError = "RDS_LOGS_INTERNAL_ERROR"
)
