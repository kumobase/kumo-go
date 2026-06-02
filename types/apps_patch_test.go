package types

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestUpdateAppRequest_PatchSemantics_OmitsAbsent asserts that a sparse
// UpdateAppRequest only marshals the fields the client actually set. The
// PATCH semantics contract is "absent JSON key = no change", so any
// zero-value pollution would silently zero those fields server-side.
func TestUpdateAppRequest_PatchSemantics_OmitsAbsent(t *testing.T) {
	req := UpdateAppRequest{Name: strPtr("my-app")}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(b)
	want := `{"name":"my-app"}`
	if got != want {
		t.Fatalf("UpdateAppRequest{Name: ptr} marshalled to %s, want %s", got, want)
	}
}

// TestUpdateAppRequest_PatchSemantics_ReportedPayload mirrors the exact
// bug-report payload (no image) and asserts UpdateAppRequest unmarshals it
// with Image == nil, so the server can distinguish "no image change" from
// "set image to empty string".
func TestUpdateAppRequest_PatchSemantics_ReportedPayload(t *testing.T) {
	body := `{"name":"testtokodjaringa","port":8080,"is_exposed":true,"replicas":3,` +
		`"environment_variables":[],"pricing_slug":"kumo.nano",` +
		`"secret_vars":[],"secret_file_mounts":[]}`
	var req UpdateAppRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if req.Image != nil {
		t.Fatalf("Image should be nil when key is absent, got %q", *req.Image)
	}
	if req.Name == nil || *req.Name != "testtokodjaringa" {
		t.Fatalf("Name not parsed: %v", req.Name)
	}
	if req.EnvironmentVariables == nil {
		t.Fatal("EnvironmentVariables should be non-nil empty slice (client said: clear them)")
	}
	if len(req.EnvironmentVariables) != 0 {
		t.Fatalf("EnvironmentVariables should be empty, got %d", len(req.EnvironmentVariables))
	}
}

// TestUpdateAppRequest_PatchSemantics_EmptyBody asserts that an entirely
// empty body unmarshals to a zero-value request with every pointer nil.
func TestUpdateAppRequest_PatchSemantics_EmptyBody(t *testing.T) {
	var req UpdateAppRequest
	if err := json.Unmarshal([]byte(`{}`), &req); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if req.Name != nil || req.Image != nil || req.Port != nil || req.Replicas != nil ||
		req.IsExposed != nil || req.PricingSlug != nil || req.RegistryCredentialId != nil ||
		req.EnvironmentVariables != nil || req.SecretVars != nil || req.SecretFileMounts != nil {
		t.Fatalf("empty body should leave every pointer/slice nil, got %+v", req)
	}
}

// TestUpdateAppRequest_PatchSemantics_ImageOmittedNotInJSON guards against a
// regression where Image stops being a pointer and silently re-introduces
// the zero-value-pollution failure mode.
func TestUpdateAppRequest_PatchSemantics_ImageOmittedNotInJSON(t *testing.T) {
	req := UpdateAppRequest{Replicas: intPtr(3)}
	b, _ := json.Marshal(req)
	if strings.Contains(string(b), `"image"`) {
		t.Fatalf("Image key should be absent from marshalled output, got %s", b)
	}
}

// TestCreateAppRequest_VolumeOmittedWhenNil asserts the optional volume block
// is absent from the wire when not requested, so existing no-volume creates
// are byte-identical to before the field was added.
func TestCreateAppRequest_VolumeOmittedWhenNil(t *testing.T) {
	req := CreateAppRequest{BaseCreateApp: BaseCreateApp{Name: "my-app", Image: "nginx", Port: 80, Replicas: 1}}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(b), `"volume"`) {
		t.Fatalf("volume key should be absent when Volume is nil, got %s", b)
	}
}

// TestCreateAppVolume_RoundTrip covers both attach-by-id and attach-by-name
// shapes, asserting the volume block survives a marshal/unmarshal cycle and
// that VolumeID stays a pointer (so omitempty distinguishes "unset" from 0).
func TestCreateAppVolume_RoundTrip(t *testing.T) {
	roundTrip(t, "CreateAppVolume/by-id", CreateAppRequest{
		BaseCreateApp: BaseCreateApp{Name: "my-app", Image: "nginx", Port: 80, Replicas: 1},
		Volume:        &CreateAppVolume{VolumeID: uintPtr(7), MountPath: "/data"},
	})
	roundTrip(t, "CreateAppVolume/by-name", CreateAppRequest{
		BaseCreateApp: BaseCreateApp{Name: "my-app", Image: "nginx", Port: 80, Replicas: 1},
		Volume:        &CreateAppVolume{VolumeName: "pgdata"},
	})
}
