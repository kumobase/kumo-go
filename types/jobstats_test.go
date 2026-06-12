package types

import (
	"testing"
	"time"
)

func TestJobExecutionStatsRoundTrip(t *testing.T) {
	from := time.Date(2026, 6, 5, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC)
	rate := 0.846
	p50, p95, p99 := int64(1200), int64(5300), int64(9000)
	avg, min, max := int64(1800), int64(400), int64(9000)

	roundTrip(t, "JobExecutionStatsResponse/populated", JobExecutionStatsResponse{
		JobID:       2,
		From:        from,
		To:          to,
		Granularity: "day",
		Counts:      JobExecutionStatusCounts{Total: 40, Pending: 0, Running: 1, Succeeded: 33, Failed: 5, Timeout: 1},
		SuccessRate: &rate,
		DurationMS:  JobExecutionDurationStats{P50: &p50, P95: &p95, P99: &p99, Avg: &avg, Min: &min, Max: &max},
		Timeline: []JobExecutionTimelineBucket{
			{Bucket: time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC), Total: 6, Succeeded: 5, Failed: 1, Timeout: 0},
		},
	})

	// Empty window: no terminal executions, no finished durations, sparse
	// timeline collapses to an empty slice. success_rate + every duration field
	// omitted.
	roundTrip(t, "JobExecutionStatsResponse/empty", JobExecutionStatsResponse{
		JobID:       9,
		From:        from,
		To:          to,
		Granularity: "hour",
		Counts:      JobExecutionStatusCounts{},
		DurationMS:  JobExecutionDurationStats{},
		Timeline:    []JobExecutionTimelineBucket{},
	})
}
