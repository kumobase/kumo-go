package types

import (
	"testing"
	"time"
)

func TestJobExecutionLogsRoundTrip(t *testing.T) {
	next := "1780000000000000000"
	expiresAt := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)

	roundTrip(t, "JobExecutionLogsResponse/populated", JobExecutionLogsResponse{
		JobID:       42,
		ExecutionID: 7,
		Start:       "1780000000000000000",
		End:         "1780000090000000000",
		Direction:   LogDirectionBackward,
		Limit:       LogsDefaultLimit,
		Entries: []LogEntry{
			{Timestamp: "1780000090000000000", Line: "done", Stream: "stdout", Level: "info", Pod: "kumo-job-x-ab12-manual-cd34-9xz"},
			{Timestamp: "1780000000000000000", Line: "starting", Stream: "stderr"},
		},
		Next:          &next,
		LogsExpiresAt: &expiresAt,
	})

	// Empty page (never started), last page (Next nil), logs still live.
	roundTrip(t, "JobExecutionLogsResponse/empty", JobExecutionLogsResponse{
		JobID:       9,
		ExecutionID: 1,
		Start:       "1780000000000000000",
		End:         "1780000090000000000",
		Direction:   LogDirectionBackward,
		Limit:       LogsDefaultLimit,
		Entries:     []LogEntry{},
	})

	// Aged-out execution: LogsExpired true, no entries, no backend query.
	roundTrip(t, "JobExecutionLogsResponse/expired", JobExecutionLogsResponse{
		JobID:         3,
		ExecutionID:   2,
		Start:         "1780000000000000000",
		End:           "1780000090000000000",
		Direction:     LogDirectionBackward,
		Limit:         LogsDefaultLimit,
		Entries:       []LogEntry{},
		LogsExpired:   true,
		LogsExpiresAt: &expiresAt,
	})
}

func TestJobLogsRetentionHours(t *testing.T) {
	if JobLogsRetentionHours != 168 {
		t.Errorf("JobLogsRetentionHours = %d, want 168 (7 days)", JobLogsRetentionHours)
	}
}
