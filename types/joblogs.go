package types

import "time"

// JobLogsRetentionHours is how long a job execution's logs are retained in Loki
// (and, equivalently, how far back the per-execution logs endpoint will query
// and how far back executions remain accessible via the list/get endpoints).
// Jobs get a longer window than apps (LogsMaxLookbackHours) because each
// execution is a discrete artifact users inspect historically, and job log
// volume is tiny/bursty.
//
// This MUST match the Loki per-stream retention override for {container="job"}
// in the infra values (grafana-loki/prod-values-baremetal.yaml). Changing the
// horizon means changing both.
const JobLogsRetentionHours = 168 // 7 days

// Job execution logs DTOs returned by
// GET /api/v1/jobs/:id/executions/:execution_id/logs. Returns a page of log
// lines for one job execution's pod over a time window that defaults to the
// execution's lifetime (overridable via start/end), newest-first by default,
// with a cursor for paging further back. Reuses LogEntry and the generic log
// constants (LogDirection*, LogLevels, LogsDefault*/Max*) from applogs.go.
//
// Timestamps and the Next cursor are unix-NANOSECOND strings (see applogs.go
// for the precision rationale).

// JobExecutionLogsResponse is the Data payload of the per-execution logs
// endpoint.
//   - Start/End echo the unix-nanosecond window the server actually queried.
//   - Entries is never null; it is an empty array when the execution produced
//     no logs in the window (e.g. never started, or already aged out of
//     retention).
//   - Next is the cursor for the following page in the chosen Direction, or
//     null when the current page is the last one.
//   - LogsExpired is true when the execution finished longer ago than
//     JobLogsRetentionHours, so its logs have aged out of Loki: Entries is
//     empty by design and the server did not query the backend. Clients should
//     render an explicit "logs expired" state rather than "no logs".
//   - LogsExpiresAt is when this execution's logs will (or did) age out of
//     Loki — pod_finished_at + JobLogsRetentionHours. It is null while the
//     execution is still pending/running (no finish time yet).
type JobExecutionLogsResponse struct {
	JobID         uint       `json:"job_id"`
	ExecutionID   uint       `json:"execution_id"`
	Start         string     `json:"start"`
	End           string     `json:"end"`
	Direction     string     `json:"direction"`
	Limit         int        `json:"limit"`
	Entries       []LogEntry `json:"entries"`
	Next          *string    `json:"next,omitempty"`
	LogsExpired   bool       `json:"logs_expired"`
	LogsExpiresAt *time.Time `json:"logs_expires_at,omitempty"`
}
