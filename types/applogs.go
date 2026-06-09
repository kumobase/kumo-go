package types

// App logs DTOs returned by GET /api/v1/apps/:id/logs. The endpoint returns a
// page of log lines for one app over a caller-chosen time window, newest-first
// by default, with a cursor for paging further back.
//
// Timestamps and the Next cursor are unix-NANOSECOND strings, not numbers:
// Loki's native resolution is nanoseconds and those values exceed 2^53, so a
// JSON number would lose precision in JavaScript clients. The cursor must also
// feed back verbatim into the next request's `start`/`end`, which Loki expects
// in nanoseconds — so entry timestamp and cursor share one representation.

// Log line direction. Backward is newest-first (the default); forward is
// oldest-first. The direction also fixes how the Next cursor is consumed.
const (
	LogDirectionForward  = "forward"  // oldest first
	LogDirectionBackward = "backward" // newest first (default)
)

// LogLevels is the allowlist for the optional `level` query parameter, matched
// against Loki's auto-detected `detected_level` structured metadata. Unknown
// values are rejected server-side with code INVALID_LOG_FILTER.
var LogLevels = []string{"error", "warn", "info", "debug"}

// Log paging/window bounds. The server enforces these; they are exported so
// clients (CLI/terraform) can validate before calling.
const (
	// LogsDefaultLimit is applied when `limit` is omitted.
	LogsDefaultLimit = 100
	// LogsMaxLimit caps `limit`; larger values are clamped down.
	LogsMaxLimit = 1000
	// LogsDefaultLookback is the window applied when `start` is omitted
	// (relative to `end`, which defaults to now).
	LogsDefaultLookback = "1h"
	// LogsMaxLookbackHours bounds how far back `start` may reach. It matches
	// Loki's free-tier retention; older data has been deleted, so the server
	// rejects such requests with INVALID_TIME_RANGE rather than returning a
	// silently-narrowed window.
	LogsMaxLookbackHours = 72
)

// LogEntry is one merged, globally-ordered log line.
//   - Timestamp is a unix-nanosecond string.
//   - Stream is "stdout" or "stderr" (from Loki structured metadata), omitted
//     when unknown.
//   - Level is the auto-detected level (error/warn/info/debug/…), omitted when
//     Loki did not detect one.
//   - Pod is the source pod name, useful to disambiguate lines when an app runs
//     multiple replicas; omitted when unavailable.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Line      string `json:"line"`
	Stream    string `json:"stream,omitempty"`
	Level     string `json:"level,omitempty"`
	Pod       string `json:"pod,omitempty"`
}

// AppLogsResponse is the Data payload of GET /api/v1/apps/:id/logs.
//   - Start/End echo the unix-nanosecond window the server actually queried.
//   - Entries is never null; it is an empty array when the app produced no logs
//     in the window (e.g. never deployed, scaled to zero, or quiet).
//   - Next is the cursor for the following page in the chosen Direction, or null
//     when the current page is the last one. To page, repeat the request with
//     `end=<next>` (backward) or `start=<next>` (forward).
type AppLogsResponse struct {
	AppID     uint       `json:"app_id"`
	Start     string     `json:"start"`
	End       string     `json:"end"`
	Direction string     `json:"direction"`
	Limit     int        `json:"limit"`
	Entries   []LogEntry `json:"entries"`
	Next      *string    `json:"next,omitempty"`
}
