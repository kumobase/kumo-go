package types

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/kumobase/kumo-go/codes"
)

// TestCreateGitBuildAppRequest_TagPatternRoundTrip pins the new optional
// tag_pattern field. Server contract: omitempty so an empty value is absent.
func TestCreateGitBuildAppRequest_TagPatternRoundTrip(t *testing.T) {
	req := CreateGitBuildAppRequest{
		Name: "rel-app", Port: 8080, Replicas: 1,
		RepoFullName: "acme/web",
		Branch:       "",
		TagPattern:   "v*",
		PricingSlug:  "kumo.nano",
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(b), `"tag_pattern":"v*"`) {
		t.Fatalf("tag_pattern key missing or wrong: %s", b)
	}
	// Empty branch should be omitted, not "branch":"".
	if strings.Contains(string(b), `"branch":""`) {
		t.Fatalf("empty branch should be omitted, got %s", b)
	}
	var back CreateGitBuildAppRequest
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if back.TagPattern != "v*" {
		t.Fatalf("round-trip lost tag_pattern: %q", back.TagPattern)
	}
}

// TestUpdateBuildConfigRequest_TagPatternPointer mirrors the v0.7.2 PATCH
// pattern: nil = no change, non-nil "" = clear.
func TestUpdateBuildConfigRequest_TagPatternPointer(t *testing.T) {
	b, _ := json.Marshal(UpdateBuildConfigRequest{Language: "static"})
	if strings.Contains(string(b), `"tag_pattern"`) {
		t.Fatalf("nil TagPattern must be absent from JSON, got %s", b)
	}

	empty := ""
	b, _ = json.Marshal(UpdateBuildConfigRequest{TagPattern: &empty})
	if !strings.Contains(string(b), `"tag_pattern":""`) {
		t.Fatalf("non-nil empty TagPattern must serialise as \"\", got %s", b)
	}

	v := "v[0-9]*"
	b, _ = json.Marshal(UpdateBuildConfigRequest{TagPattern: &v})
	if !strings.Contains(string(b), `"tag_pattern":"v[0-9]*"`) {
		t.Fatalf("TagPattern value not preserved, got %s", b)
	}
}

// TestGitBuildInfo_TagPatternRoundTrip pins the read-side surface on the app
// detail.
func TestGitBuildInfo_TagPatternRoundTrip(t *testing.T) {
	gb := GitBuildInfo{RepoFullName: "acme/web", Branch: "main", TagPattern: "v*"}
	b, _ := json.Marshal(gb)
	if !strings.Contains(string(b), `"tag_pattern":"v*"`) {
		t.Fatalf("missing tag_pattern: %s", b)
	}
	gb.TagPattern = ""
	b, _ = json.Marshal(gb)
	if strings.Contains(string(b), `"tag_pattern"`) {
		t.Fatalf("empty TagPattern must be omitted, got %s", b)
	}
}

// TestBuildCodes_NewSentinels guards the wire strings against accidental
// renames. These codes are consumed by Terraform/CLI users via switch on
// APIError.Code.
func TestBuildCodes_NewSentinels(t *testing.T) {
	if codes.BuildTriggerRequired != "BUILD_TRIGGER_REQUIRED" {
		t.Fatalf("BuildTriggerRequired drifted: %q", codes.BuildTriggerRequired)
	}
	if codes.BuildInvalidTagPattern != "BUILD_INVALID_TAG_PATTERN" {
		t.Fatalf("BuildInvalidTagPattern drifted: %q", codes.BuildInvalidTagPattern)
	}
	if codes.BuildNeedsBranch != "BUILD_NEEDS_BRANCH" {
		t.Fatalf("BuildNeedsBranch drifted: %q", codes.BuildNeedsBranch)
	}
}
