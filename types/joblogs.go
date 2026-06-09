package types

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
type JobExecutionLogsResponse struct {
	JobID       uint       `json:"job_id"`
	ExecutionID uint       `json:"execution_id"`
	Start       string     `json:"start"`
	End         string     `json:"end"`
	Direction   string     `json:"direction"`
	Limit       int        `json:"limit"`
	Entries     []LogEntry `json:"entries"`
	Next        *string    `json:"next,omitempty"`
}
