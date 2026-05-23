package types

import (
	"encoding/json"
	"reflect"
	"testing"
)

// roundTrip serialises v, decodes the result back into a fresh value of the
// same concrete type, then serialises *that* and compares the two JSON
// payloads byte-for-byte. Catches accidental tag changes, dropped fields,
// or omitempty misconfiguration.
func roundTrip[T any](t *testing.T, name string, original T) {
	t.Helper()
	first, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("%s: marshal: %v", name, err)
	}
	var decoded T
	if err := json.Unmarshal(first, &decoded); err != nil {
		t.Fatalf("%s: unmarshal: %v (json=%s)", name, err, first)
	}
	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("%s: not deep-equal after round-trip\n  orig: %+v\n  back: %+v", name, original, decoded)
	}
	second, err := json.Marshal(decoded)
	if err != nil {
		t.Fatalf("%s: re-marshal: %v", name, err)
	}
	if string(first) != string(second) {
		t.Errorf("%s: re-marshal not byte-identical\n  first:  %s\n  second: %s", name, first, second)
	}
}

func TestCommonRoundTrip(t *testing.T) {
	roundTrip(t, "Meta", Meta{Page: 2, PageSize: 50, TotalItems: 137, TotalPages: 3})
	roundTrip(t, "EnvironmentVariable", EnvironmentVariable{Key: "DB_URL", Value: "postgres://x"})
	roundTrip(t, "ValidationError", ValidationError{Field: "name", Message: "is required"})
	roundTrip(t, "ValidationErrorsResponse", ValidationErrorsResponse{
		Errors: []*ValidationError{{Field: "image", Message: "invalid image reference"}},
	})
	roundTrip(t, "Availability", Availability{Available: false, Reason: AvailabilityReasonCPUFull})
	roundTrip(t, "StructureResponse", StructureResponse{
		Code:    "APP_NOT_FOUND",
		Message: "app not found",
		Data:    json.RawMessage(`{"id":1}`),
		Meta:    &Meta{Page: 1, PageSize: 20, TotalItems: 0, TotalPages: 1},
	})
}

// TestStructureResponse_OmitsEmptyOptional verifies that an envelope with
// only Message set does not emit "code", "data", or "meta" — important
// because dashboards rely on the absence of those fields to decide layout.
func TestStructureResponse_OmitsEmptyOptional(t *testing.T) {
	body, _ := json.Marshal(StructureResponse{Message: "ok"})
	if got, want := string(body), `{"message":"ok"}`; got != want {
		t.Errorf("envelope: got %s, want %s", got, want)
	}
}
