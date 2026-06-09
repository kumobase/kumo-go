package types

import (
	"encoding/json"
	"testing"
)

func TestAppLogsRoundTrip(t *testing.T) {
	next := "1780000000000000000"

	roundTrip(t, "AppLogsResponse/populated", AppLogsResponse{
		AppID:     42,
		Start:     "1779999999000000000",
		End:       "1780003600000000000",
		Direction: LogDirectionBackward,
		Limit:     100,
		Entries: []LogEntry{
			{Timestamp: "1780003599000000000", Line: "listening on :8080", Stream: "stdout", Level: "info", Pod: "app-x-1"},
			{Timestamp: "1780003500000000000", Line: "connection refused", Stream: "stderr", Level: "error", Pod: "app-x-2"},
		},
		Next: &next,
	})

	// Last page (fewer than Limit) + minimal entry (no metadata) — Next omitted.
	roundTrip(t, "AppLogsResponse/lastpage", AppLogsResponse{
		AppID:     7,
		Start:     "1780000000000000000",
		End:       "1780003600000000000",
		Direction: LogDirectionForward,
		Limit:     300,
		Entries:   []LogEntry{{Timestamp: "1780000001000000000", Line: "hello"}},
	})

	// Empty window — Entries marshals to [] (never null), Next omitted.
	roundTrip(t, "AppLogsResponse/empty", AppLogsResponse{
		AppID:     9,
		Start:     "1780000000000000000",
		End:       "1780003600000000000",
		Direction: LogDirectionBackward,
		Limit:     100,
		Entries:   []LogEntry{},
	})
}

// Next must be omitted from the wire when nil; entry metadata omitted when empty.
func TestAppLogsOmitempty(t *testing.T) {
	noNext, _ := json.Marshal(AppLogsResponse{
		AppID: 1, Start: "1", End: "2", Direction: LogDirectionBackward, Limit: 100,
		Entries: []LogEntry{},
	})
	if got := string(noNext); got != `{"app_id":1,"start":"1","end":"2","direction":"backward","limit":100,"entries":[]}` {
		t.Errorf("nil next / empty entries: got %s", got)
	}

	bareEntry, _ := json.Marshal(LogEntry{Timestamp: "1", Line: "x"})
	if got := string(bareEntry); got != `{"timestamp":"1","line":"x"}` {
		t.Errorf("empty stream/level/pod should be omitted, got %s", got)
	}
}
