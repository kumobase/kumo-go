package types

import "testing"

func TestRDSLogsRoundTrip(t *testing.T) {
	next := "1780000000000000000"

	roundTrip(t, "RDSLogsResponse/populated", RDSLogsResponse{
		InstanceID: 42,
		Start:      "1780000000000000000",
		End:        "1780003600000000000",
		Direction:  LogDirectionBackward,
		Limit:      100,
		Entries: []LogEntry{
			{Timestamp: "1780000000000000000", Line: "LOG: database system is ready to accept connections", Stream: "stderr", Level: "info", Pod: "db-postgresql-0"},
			{Timestamp: "1780000001000000000", Line: "FATAL: too many connections", Stream: "stderr", Level: "error", Pod: "db-postgresql-0"},
		},
		Next: &next,
	})

	// Empty page (just provisioned / quiet) — Entries is a non-nil empty array.
	roundTrip(t, "RDSLogsResponse/empty", RDSLogsResponse{
		InstanceID: 7,
		Start:      "1780000000000000000",
		End:        "1780000900000000000",
		Direction:  LogDirectionForward,
		Limit:      100,
		Entries:    []LogEntry{},
	})
}

func TestRDSLogsRetentionHours(t *testing.T) {
	if RDSLogsRetentionHours != 72 {
		t.Errorf("RDSLogsRetentionHours = %d, want 72", RDSLogsRetentionHours)
	}
}
