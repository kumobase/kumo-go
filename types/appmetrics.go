package types

// App metrics DTOs returned by GET /api/v1/apps/:id/metrics. The endpoint
// returns CPU and memory time series for one app over a server-chosen
// resolution, plus the app's configured limit so consumers can render
// "% of limit". Values are plain numbers (not decimal strings): CPU in cores,
// memory in bytes — these are observability samples, not money.

// Metric units carried in MetricSeries.Unit.
const (
	UnitCores = "cores" // CPU, fractional cores (e.g. 0.25)
	UnitBytes = "bytes" // memory working set
)

// AppMetricsRanges are the allowed values of the `range` query parameter.
// Unknown values are rejected server-side with code INVALID_TIME_RANGE. The
// server picks the step/resolution per range; clients cannot override it.
var AppMetricsRanges = []string{"15m", "1h", "6h", "24h", "7d"}

// AppMetricsDefaultRange is applied when the `range` parameter is omitted.
const AppMetricsDefaultRange = "1h"

// MetricPoint is a single sample: T is unix seconds (step-aligned), V is the
// value in the series' Unit.
type MetricPoint struct {
	T int64   `json:"t"`
	V float64 `json:"v"`
}

// MetricSeries is one metric's time series plus its configured limit.
//   - Points is never null; it is an empty array when the app has no data in
//     the window (e.g. never deployed or scaled to zero).
//   - Limit is the app's configured limit in the same Unit, or null when no
//     limit is set. Consumers compute "% of limit" as point.V / Limit.
type MetricSeries struct {
	Unit   string        `json:"unit"`
	Points []MetricPoint `json:"points"`
	Limit  *float64      `json:"limit,omitempty"`
}

// AppMetrics groups the per-app series the endpoint returns.
type AppMetrics struct {
	CPU    MetricSeries `json:"cpu"`    // Unit "cores"
	Memory MetricSeries `json:"memory"` // Unit "bytes"
}

// AppMetricsResponse is the Data payload of GET /api/v1/apps/:id/metrics.
// Start/End/StepSeconds describe the window the server actually queried (after
// aligning to the step). RateWindowSeconds is the rate() window used for the
// CPU counter, exposed for transparency.
type AppMetricsResponse struct {
	AppID             uint       `json:"app_id"`
	Range             string     `json:"range"`
	Start             int64      `json:"start"`
	End               int64      `json:"end"`
	StepSeconds       int        `json:"step_seconds"`
	RateWindowSeconds int        `json:"rate_window_seconds"`
	Metrics           AppMetrics `json:"metrics"`
}
