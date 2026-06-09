package types

// Job execution metrics DTOs returned by
// GET /api/v1/jobs/:id/executions/:execution_id/metrics. Unlike apps (one
// long-running deployment queried over a preset range), a job runs as discrete
// short-lived executions; the server queries the window of THIS execution's
// lifetime (pod start..finish, or ..now while running). CPU is in cores, memory
// in bytes — reusing the generic MetricSeries/MetricPoint shapes. The series
// Limit is the execution's pinned resource snapshot.

// JobExecutionMetrics groups the per-execution series the endpoint returns.
type JobExecutionMetrics struct {
	CPU    MetricSeries `json:"cpu"`    // Unit "cores"
	Memory MetricSeries `json:"memory"` // Unit "bytes"
}

// JobExecutionMetricsResponse is the Data payload of the per-execution metrics
// endpoint. Start/End/StepSeconds describe the window the server actually
// queried (the execution lifetime, step-aligned). RateWindowSeconds is the
// rate() window used for the CPU counter, exposed for transparency.
type JobExecutionMetricsResponse struct {
	JobID             uint                `json:"job_id"`
	ExecutionID       uint                `json:"execution_id"`
	Start             int64               `json:"start"`
	End               int64               `json:"end"`
	StepSeconds       int                 `json:"step_seconds"`
	RateWindowSeconds int                 `json:"rate_window_seconds"`
	Metrics           JobExecutionMetrics `json:"metrics"`
}
