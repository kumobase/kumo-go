package types

import "time"

// Job execution stats DTOs returned by GET /api/v1/jobs/:id/stats. Unlike the
// per-execution metrics/logs endpoints (one execution's pod lifetime, sourced
// from cAdvisor/Loki), this is an AGGREGATE "job health" view computed entirely
// from the job_executions table in Postgres — so it is reliable regardless of
// how short an execution ran (no scrape-interval sampling gap). It surfaces the
// metrics batch/cron products converge on: counts by result, success rate, and
// duration percentiles, plus a timeline of buckets over the queried window.

// JobExecutionStatusCounts is the breakdown of executions by terminal/non-terminal
// status within the queried window. Total is the sum of all five.
type JobExecutionStatusCounts struct {
	Total     int64 `json:"total"`
	Pending   int64 `json:"pending"`
	Running   int64 `json:"running"`
	Succeeded int64 `json:"succeeded"`
	Failed    int64 `json:"failed"`
	Timeout   int64 `json:"timeout"`
}

// JobExecutionDurationStats summarises execution durations (milliseconds) over
// the window. Only finished executions carry a duration, so every field is a
// pointer and is null when no finished execution falls in the window. Percentiles
// are continuous (interpolated); Avg is rounded to whole milliseconds.
type JobExecutionDurationStats struct {
	P50 *int64 `json:"p50,omitempty"`
	P95 *int64 `json:"p95,omitempty"`
	P99 *int64 `json:"p99,omitempty"`
	Avg *int64 `json:"avg,omitempty"`
	Min *int64 `json:"min,omitempty"`
	Max *int64 `json:"max,omitempty"`
}

// JobExecutionTimelineBucket is one time bucket in the executions timeline. Bucket
// is the bucket's start instant (truncated in Asia/Jakarta). Only buckets with at
// least one execution are returned (the series is sparse); clients fill gaps.
type JobExecutionTimelineBucket struct {
	Bucket    time.Time `json:"bucket"`
	Total     int64     `json:"total"`
	Succeeded int64     `json:"succeeded"`
	Failed    int64     `json:"failed"`
	Timeout   int64     `json:"timeout"`
}

// JobExecutionStatsResponse is the Data payload of GET /api/v1/jobs/:id/stats.
// From/To echo the resolved window (defaulting to the last 7 days when the caller
// omits them); Granularity is the bucket size of Timeline ("hour" or "day").
// SuccessRate is succeeded / (succeeded+failed+timeout) in [0,1], null when there
// are no terminal executions in the window (avoids a meaningless 0).
type JobExecutionStatsResponse struct {
	JobID       uint                         `json:"job_id"`
	From        time.Time                    `json:"from"`
	To          time.Time                    `json:"to"`
	Granularity string                       `json:"granularity"`
	Counts      JobExecutionStatusCounts     `json:"counts"`
	SuccessRate *float64                     `json:"success_rate,omitempty"`
	DurationMS  JobExecutionDurationStats    `json:"duration_ms"`
	Timeline    []JobExecutionTimelineBucket `json:"timeline"`
}
