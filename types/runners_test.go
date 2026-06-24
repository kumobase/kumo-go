package types

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// The job DTO is the user-facing LOGICAL view: it must never leak cloud
// internals (provider, region, instance type, backend) — those are admin-only.
func TestRunnerJobResponse_NoCloudInternals(t *testing.T) {
	in := RunnerJobResponse{
		ID: 5, SpecLabel: "kumo-2c-4g", GithubJobID: 99, RunID: 7,
		RepoFullName: "acme/api", State: RunnerJobStateRunning,
		QueuedAt: time.Unix(1_700_000_000, 0).UTC(),
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	js := string(b)
	for _, leak := range []string{"provider", "region", "instance_type", "\"az\"", "backend"} {
		if strings.Contains(js, leak) {
			t.Fatalf("job DTO leaks internal field %q: %s", leak, js)
		}
	}
	var out RunnerJobResponse
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.SpecLabel != "kumo-2c-4g" || out.State != RunnerJobStateRunning {
		t.Fatalf("round-trip mismatch: %+v", out)
	}
}

func TestRunnerSpecResponse_RoundTrip(t *testing.T) {
	in := RunnerSpecResponse{Label: "kumo-2c-4g", DisplayName: "2 vCPU / 4 GB", CPU: 2, MemoryMB: 4096}
	b, _ := json.Marshal(in)
	var out RunnerSpecResponse
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out != in {
		t.Fatalf("round-trip mismatch: %+v", out)
	}
}
