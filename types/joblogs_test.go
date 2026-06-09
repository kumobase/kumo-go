package types

import "testing"

func TestJobExecutionLogsRoundTrip(t *testing.T) {
	next := "1780000000000000000"

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
		Next: &next,
	})

	// Empty page (never started / aged out), last page (Next nil).
	roundTrip(t, "JobExecutionLogsResponse/empty", JobExecutionLogsResponse{
		JobID:       9,
		ExecutionID: 1,
		Start:       "1780000000000000000",
		End:         "1780000090000000000",
		Direction:   LogDirectionBackward,
		Limit:       LogsDefaultLimit,
		Entries:     []LogEntry{},
	})
}
