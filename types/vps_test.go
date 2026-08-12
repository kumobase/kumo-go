package types

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestVPSProviderResponse_RoundTrip(t *testing.T) {
	in := VPSProviderResponse{Name: "Tencent Cloud"}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out VPSProviderResponse
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(out, in) {
		t.Fatalf("round-trip mismatch: got %+v want %+v", out, in)
	}
}
