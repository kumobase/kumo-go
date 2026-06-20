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

// TestUpdateBuildConfigRequest_BranchPointer mirrors the TagPattern PATCH
// pattern for the editable branch trigger: nil = no change, non-nil "" = clear,
// non-nil non-empty = set.
func TestUpdateBuildConfigRequest_BranchPointer(t *testing.T) {
	b, _ := json.Marshal(UpdateBuildConfigRequest{Language: "auto"})
	if strings.Contains(string(b), `"branch"`) {
		t.Fatalf("nil Branch must be absent from JSON, got %s", b)
	}

	empty := ""
	b, _ = json.Marshal(UpdateBuildConfigRequest{Branch: &empty})
	if !strings.Contains(string(b), `"branch":""`) {
		t.Fatalf("non-nil empty Branch must serialise as \"\", got %s", b)
	}

	v := "develop"
	b, _ = json.Marshal(UpdateBuildConfigRequest{Branch: &v})
	if !strings.Contains(string(b), `"branch":"develop"`) {
		t.Fatalf("Branch value not preserved, got %s", b)
	}
	var back UpdateBuildConfigRequest
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if back.Branch == nil || *back.Branch != "develop" {
		t.Fatalf("round-trip lost Branch: %v", back.Branch)
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
	if codes.BuildLogNotAvailable != "BUILD_LOG_NOT_AVAILABLE" {
		t.Fatalf("BuildLogNotAvailable drifted: %q", codes.BuildLogNotAvailable)
	}
}

// TestCreateGitBuildAppRequest_DockerfilePathRoundTrip pins the new optional
// dockerfile_path field. omitempty so an empty value is absent from the wire.
func TestCreateGitBuildAppRequest_DockerfilePathRoundTrip(t *testing.T) {
	req := CreateGitBuildAppRequest{
		Name: "df-app", Port: 8080, Replicas: 1,
		RepoFullName:   "acme/web",
		Branch:         "main",
		Language:       "dockerfile",
		DockerfilePath: "docker/prod.Dockerfile",
		PricingSlug:    "kumo.nano",
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(b), `"dockerfile_path":"docker/prod.Dockerfile"`) {
		t.Fatalf("dockerfile_path key missing or wrong: %s", b)
	}
	var back CreateGitBuildAppRequest
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if back.DockerfilePath != "docker/prod.Dockerfile" {
		t.Fatalf("round-trip lost dockerfile_path: %q", back.DockerfilePath)
	}

	// Empty DockerfilePath must be omitted (back-compat: existing non-dockerfile
	// callers never send the key).
	empty := CreateGitBuildAppRequest{Name: "x", Port: 8080, RepoFullName: "a/b", Branch: "main"}
	b, _ = json.Marshal(empty)
	if strings.Contains(string(b), `"dockerfile_path"`) {
		t.Fatalf("empty dockerfile_path must be omitted, got %s", b)
	}
}

// TestUpdateBuildConfigRequest_DockerfilePathRoundTrip pins the PATCH-side field.
func TestUpdateBuildConfigRequest_DockerfilePathRoundTrip(t *testing.T) {
	b, _ := json.Marshal(UpdateBuildConfigRequest{Language: "dockerfile", DockerfilePath: "Dockerfile"})
	if !strings.Contains(string(b), `"dockerfile_path":"Dockerfile"`) {
		t.Fatalf("dockerfile_path not serialised: %s", b)
	}
	// Absent when empty.
	b, _ = json.Marshal(UpdateBuildConfigRequest{Language: "static"})
	if strings.Contains(string(b), `"dockerfile_path"`) {
		t.Fatalf("empty dockerfile_path must be absent, got %s", b)
	}
}

// TestBuildCodes_DockerfileSentinels guards the new wire strings against
// accidental renames (consumed by Terraform/CLI via switch on APIError.Code).
func TestBuildCodes_DockerfileSentinels(t *testing.T) {
	if codes.BuildInvalidDockerfilePath != "BUILD_INVALID_DOCKERFILE_PATH" {
		t.Fatalf("BuildInvalidDockerfilePath drifted: %q", codes.BuildInvalidDockerfilePath)
	}
	if codes.BuildNoDockerfile != "BUILD_NO_DOCKERFILE" {
		t.Fatalf("BuildNoDockerfile drifted: %q", codes.BuildNoDockerfile)
	}
	if codes.BuildNoRailpackPlan != "BUILD_NO_RAILPACK_PLAN" {
		t.Fatalf("BuildNoRailpackPlan drifted: %q", codes.BuildNoRailpackPlan)
	}
}

// TestBuildLogURLResponse_RoundTrip pins the dedicated log-url endpoint's body.
func TestBuildLogURLResponse_RoundTrip(t *testing.T) {
	b, err := json.Marshal(BuildLogURLResponse{LogURL: "https://logs/3.txt?sig=x"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(b), `"log_url":"https://logs/3.txt?sig=x"`) {
		t.Fatalf("log_url key missing or wrong: %s", b)
	}
	var back BuildLogURLResponse
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if back.LogURL != "https://logs/3.txt?sig=x" {
		t.Fatalf("round-trip lost log_url: %q", back.LogURL)
	}
}

// TestBuildersResponse_RoundTrip pins the discovery endpoint body shape.
func TestBuildersResponse_RoundTrip(t *testing.T) {
	in := BuildersResponse{
		Builders: []BuilderOption{
			{Kind: "auto", Label: "Auto (Dockerfile or Railpack)", Default: true},
			{Kind: "railpack", Label: "Railpack"},
			{Kind: "cnb", Label: "Buildpacks"},
		},
		Languages: []LanguageOption{{Value: "nodejs"}, {Value: "go"}},
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(b), `"kind":"auto"`) || !strings.Contains(string(b), `"default":true`) {
		t.Fatalf("builders shape wrong: %s", b)
	}
	// Non-default builders omit the default flag (omitempty).
	if strings.Contains(string(b), `"kind":"railpack","label":"Railpack","default"`) {
		t.Fatalf("non-default builder must omit default flag: %s", b)
	}
	var back BuildersResponse
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(back.Builders) != 3 || !back.Builders[0].Default || back.Languages[1].Value != "go" {
		t.Fatalf("round-trip lost data: %+v", back)
	}
}

// TestAppByIdResponse_DockerfilePathRoundTrip pins the read-back of the
// configured Dockerfile path on app detail.
func TestAppByIdResponse_DockerfilePathRoundTrip(t *testing.T) {
	b, _ := json.Marshal(AppByIdResponse{Language: "dockerfile", DockerfilePath: "docker/prod.Dockerfile"})
	if !strings.Contains(string(b), `"dockerfile_path":"docker/prod.Dockerfile"`) {
		t.Fatalf("dockerfile_path missing from app detail: %s", b)
	}
	// Omitted when empty (registry-image / non-dockerfile apps).
	b, _ = json.Marshal(AppByIdResponse{Language: "static"})
	if strings.Contains(string(b), `"dockerfile_path"`) {
		t.Fatalf("empty dockerfile_path must be omitted: %s", b)
	}
}
