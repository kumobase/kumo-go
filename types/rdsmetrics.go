package types

// RDS instance metrics DTOs returned by GET /api/v1/rds/:id/metrics. The
// endpoint returns CPU, memory and disk time series for one managed-PostgreSQL
// instance over a server-chosen resolution. CPU is in cores and memory/disk in
// bytes — reusing the generic MetricSeries/MetricPoint shapes from appmetrics.go
// (and the UnitCores/UnitBytes units and AppMetricsRanges/AppMetricsDefaultRange
// presets, which are shared and intentionally not redeclared here).
//
// Disk is the worst-case fill across the instance's data volumes (primary plus
// any read replicas): Points carry used bytes, and Limit is the provisioned
// volume capacity in bytes (authoritative from the instance's storage_gb, so a
// "% full" can be rendered even before the volume-usage metric is scraped).

// RDSMetrics groups the per-instance series the endpoint returns.
type RDSMetrics struct {
	CPU    MetricSeries `json:"cpu"`    // Unit "cores"
	Memory MetricSeries `json:"memory"` // Unit "bytes"
	Disk   MetricSeries `json:"disk"`   // Unit "bytes"; Limit = volume capacity bytes
	// Database holds PostgreSQL-internal series from the instance's exporter
	// (Tier 2). Each series is empty (non-null) when the exporter isn't scraped.
	Database RDSDatabaseMetrics `json:"database"`
}

// RDSDatabaseMetrics groups the PostgreSQL-internal series scraped from the
// instance's postgres_exporter. Logical metrics (connections/commits/rollbacks/
// cache-hit/up) reflect the PRIMARY node; ReplicationLag is the worst-case lag
// across standbys (0 on a standalone instance).
type RDSDatabaseMetrics struct {
	Up             MetricSeries `json:"up"`              // count, 0/1 (primary reachable)
	Connections    MetricSeries `json:"connections"`    // count; Limit = max_connections
	ReplicationLag MetricSeries `json:"replication_lag"` // seconds (worst standby)
	Commits        MetricSeries `json:"commits"`         // per_second
	Rollbacks      MetricSeries `json:"rollbacks"`       // per_second
	CacheHitRatio  MetricSeries `json:"cache_hit_ratio"` // percent (0..100)
}

// RDSMetricsResponse is the Data payload of GET /api/v1/rds/:id/metrics.
// Start/End/StepSeconds describe the window the server actually queried (after
// aligning to the step). RateWindowSeconds is the rate() window used for the CPU
// counter, exposed for transparency.
type RDSMetricsResponse struct {
	InstanceID        uint       `json:"instance_id"`
	Range             string     `json:"range"`
	Start             int64      `json:"start"`
	End               int64      `json:"end"`
	StepSeconds       int        `json:"step_seconds"`
	RateWindowSeconds int        `json:"rate_window_seconds"`
	Metrics           RDSMetrics `json:"metrics"`
}
