package types

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestSourceConnectionResponse_RoundTrip(t *testing.T) {
	in := SourceConnectionResponse{
		ID:             7,
		Provider:       SourceProviderGitHub,
		InstallationID: 4242,
		AccountLogin:   "acme",
		AccountType:    "Organization",
		ManageURL:      "https://github.com/organizations/acme/settings/installations/4242",
		Status:         SourceConnectionStatusActive,
		AppKind:        "runner",
		CreatedAt:      time.Unix(1_700_000_000, 0).UTC(),
		UpdatedAt:      time.Unix(1_700_000_500, 0).UTC(),
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	// The discriminator clients branch on must be on the wire.
	if !strings.Contains(string(b), `"app_kind":"runner"`) {
		t.Fatalf("missing app_kind on wire: %s", b)
	}
	var out SourceConnectionResponse
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out != in {
		t.Fatalf("round-trip mismatch:\n got %+v\nwant %+v", out, in)
	}
}
