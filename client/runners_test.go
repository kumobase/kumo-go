package client_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/kumobase/kumo-go/client"
	"github.com/kumobase/kumo-go/types"
)

func TestRunners_Smoke(t *testing.T) {
	var lastQuery string
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		lastQuery = r.URL.RawQuery
		switch {
		case r.Method == "GET" && r.URL.Path == "/api/v1/runner-specs":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message":"ok","data":[{"label":"kumo-2c-4g","display_name":"2 vCPU / 4 GB","cpu":2,"memory_mb":4096}],"meta":{"page":1,"page_size":20,"total_items":1,"total_pages":1}}`))
		case r.Method == "GET" && r.URL.Path == "/api/v1/runner-jobs":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message":"ok","data":[{"id":11,"spec_label":"kumo-2c-4g","github_job_id":99,"run_id":7,"repo_full_name":"acme/api","state":"running","queued_at":"2026-06-24T00:00:00Z"}],"meta":{"page":1,"page_size":20,"total_items":1,"total_pages":1}}`))
		case r.Method == "GET" && r.URL.Path == "/api/v1/runner-jobs/11":
			writeStruct(w, 200, "", "ok", &types.RunnerJobResponse{ID: 11, SpecLabel: "kumo-2c-4g", State: types.RunnerJobStateCompleted})
		case r.Method == "GET" && r.URL.Path == "/api/v1/runner-usage":
			writeStruct(w, 200, "", "ok", &types.RunnerUsageResponse{JobCount: 12, BilledMinutes: 87.5, EstimatedCents: 420})
		}
	})
	ctx := context.Background()

	specs, _, err := c.Runners().ListSpecs(ctx)
	if err != nil || len(specs) != 1 || specs[0].Label != "kumo-2c-4g" {
		t.Fatalf("ListSpecs: %v (%+v)", err, specs)
	}

	jobs, meta, err := c.Runners().ListJobs(ctx, client.WithExtraQuery("state", "running"))
	if err != nil || len(jobs) != 1 || jobs[0].State != types.RunnerJobStateRunning || meta == nil {
		t.Fatalf("ListJobs: %v (%+v)", err, jobs)
	}
	if !strings.Contains(lastQuery, "state=running") {
		t.Fatalf("state filter not forwarded: %q", lastQuery)
	}

	job, err := c.Runners().GetJob(ctx, 11)
	if err != nil || job.SpecLabel != "kumo-2c-4g" {
		t.Fatalf("GetJob: %v (%+v)", err, job)
	}

	usage, err := c.Runners().Usage(ctx)
	if err != nil || usage.JobCount != 12 || usage.BilledMinutes != 87.5 {
		t.Fatalf("Usage: %v (%+v)", err, usage)
	}
}
