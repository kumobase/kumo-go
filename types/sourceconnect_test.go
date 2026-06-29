package types

import (
	"encoding/json"
	"reflect"
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
		AccountID:      9001,
		AvatarURL:      "https://avatars.githubusercontent.com/u/9001?v=4",
		ManageURL:      "https://github.com/organizations/acme/settings/installations/4242",
		RepoSelection:  "selected",
		RepoCount:      4,
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
	// GitHub repo summary must be on the wire for the picker UI.
	if !strings.Contains(string(b), `"repo_selection":"selected"`) || !strings.Contains(string(b), `"repo_count":4`) {
		t.Fatalf("missing repo summary on wire: %s", b)
	}
	// A github row carries no nested gitlab object.
	if strings.Contains(string(b), `"gitlab":`) {
		t.Fatalf("github row should not emit gitlab sub-object: %s", b)
	}
	var out SourceConnectionResponse
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	// GitLab is nil here, so == comparison is valid.
	if out != in {
		t.Fatalf("round-trip mismatch:\n got %+v\nwant %+v", out, in)
	}
}

func TestSourceConnectionResponse_NestedGitLab_RoundTrip(t *testing.T) {
	in := SourceConnectionResponse{
		ID:           12,
		Provider:     SourceProviderGitLab,
		AccountLogin: "rahadiangg",
		Status:       SourceConnectionStatusActive,
		AppKind:      "runner",
		GitLab: &GitLabConnectionResponse{
			ID:            12,
			Provider:      SourceProviderGitLab,
			InstanceID:    2,
			BaseURL:       "https://gitlab.com",
			InstanceKind:  "saas",
			Host:          "gitlab.com",
			Kind:          GitLabNamespaceGroup,
			NamespaceID:   123,
			NamespacePath: "myorg/team",
			DisplayName:   "My Team",
			AvatarURL:     "https://gitlab.com/uploads/-/system/group/avatar/123/a.png",
			Visibility:    "private",
			WebURL:        "https://gitlab.com/myorg/team",
			Status:        SourceConnectionStatusActive,
			CreatedAt:     time.Unix(1_700_000_000, 0).UTC(),
			UpdatedAt:     time.Unix(1_700_000_500, 0).UTC(),
		},
		CreatedAt: time.Unix(1_700_000_000, 0).UTC(),
		UpdatedAt: time.Unix(1_700_000_500, 0).UTC(),
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	// The nested object and its discriminating fields must reach the wire.
	for _, want := range []string{`"gitlab":`, `"instance_kind":"saas"`, `"host":"gitlab.com"`, `"display_name":"My Team"`, `"visibility":"private"`} {
		if !strings.Contains(string(b), want) {
			t.Fatalf("missing %s on wire: %s", want, b)
		}
	}
	var out SourceConnectionResponse
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	// Pointer field present → DeepEqual, not ==.
	if !reflect.DeepEqual(out, in) {
		t.Fatalf("round-trip mismatch:\n got %+v\nwant %+v", out, in)
	}
}

func TestGitLabConnectionResponse_RoundTrip(t *testing.T) {
	in := GitLabConnectionResponse{
		ID:            5,
		Provider:      SourceProviderGitLab,
		InstanceID:    3,
		BaseURL:       "https://gitlab.acme.com",
		InstanceKind:  "self_managed",
		Host:          "gitlab.acme.com",
		Kind:          GitLabNamespaceProject,
		NamespaceID:   77,
		NamespacePath: "platform/api",
		DisplayName:   "API",
		AvatarURL:     "",
		Visibility:    "internal",
		WebURL:        "https://gitlab.acme.com/platform/api",
		Status:        SourceConnectionStatusActive,
		CreatedAt:     time.Unix(1_700_000_000, 0).UTC(),
		UpdatedAt:     time.Unix(1_700_000_500, 0).UTC(),
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	// Self-managed instance kind + internal visibility must round-trip.
	if !strings.Contains(string(b), `"instance_kind":"self_managed"`) || !strings.Contains(string(b), `"visibility":"internal"`) {
		t.Fatalf("missing instance_kind/visibility on wire: %s", b)
	}
	// Empty avatar must be omitted.
	if strings.Contains(string(b), `"avatar_url"`) {
		t.Fatalf("empty avatar_url should be omitted: %s", b)
	}
	var out GitLabConnectionResponse
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out != in {
		t.Fatalf("round-trip mismatch:\n got %+v\nwant %+v", out, in)
	}
}
