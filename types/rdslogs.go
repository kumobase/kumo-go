package types

// RDS instance logs DTOs returned by GET /api/v1/rds/:id/logs. Returns a page
// of PostgreSQL log lines for one managed-database instance (primary plus any
// read replicas) over a caller-chosen time window, newest-first by default, with
// a cursor for paging further back. Reuses LogEntry and the generic log
// constants (LogDirection*, LogLevels, LogsDefault*/Max*) from applogs.go.
//
// Timestamps and the Next cursor are unix-NANOSECOND strings (see applogs.go for
// the precision rationale).

// RDSLogsRetentionHours bounds how far back `start` may reach for RDS logs. It
// matches the Loki per-stream retention override for {container="postgresql"};
// older data has been deleted, so the server rejects such requests with
// INVALID_TIME_RANGE rather than returning a silently-narrowed window.
const RDSLogsRetentionHours = 72

// RDSLogsResponse is the Data payload of GET /api/v1/rds/:id/logs.
//   - Start/End echo the unix-nanosecond window the server actually queried.
//   - Entries is never null; it is an empty array when the instance produced no
//     logs in the window (e.g. just provisioned, suspended, or quiet).
//   - Next is the cursor for the following page in the chosen Direction, or null
//     when the current page is the last one. To page, repeat the request with
//     `end=<next>` (backward) or `start=<next>` (forward).
type RDSLogsResponse struct {
	InstanceID uint       `json:"instance_id"`
	Start      string     `json:"start"`
	End        string     `json:"end"`
	Direction  string     `json:"direction"`
	Limit      int        `json:"limit"`
	Entries    []LogEntry `json:"entries"`
	Next       *string    `json:"next,omitempty"`
}
