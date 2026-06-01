package client_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/kumobase/kumo-go/types"
)

func TestJobs_Smoke(t *testing.T) {
	type call struct{ method, path string }
	var seen call
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seen = call{r.Method, r.URL.Path}
		switch {
		case r.Method == "POST" && r.URL.Path == "/api/v1/jobs":
			writeStruct(w, 202, "", "queued", &types.ResponseJobAsync{
				ID: 7, Name: "nightly", DeploymentStatus: types.JobDeploymentStatusDeploying, OperationID: "op_1",
			})
		case r.Method == "GET" && r.URL.Path == "/api/v1/jobs/7":
			writeStruct(w, 200, "", "ok", &types.JobResponse{
				ID: 7, Name: "nightly", Kind: types.JobKindStandalone, Image: "alpine:3.20",
				Schedule: "0 2 * * *", Timezone: "Asia/Jakarta",
				ConcurrencyPolicy: types.JobConcurrencyForbid,
				DeploymentStatus:  types.JobDeploymentStatusActive,
			})
		case r.Method == "GET" && r.URL.Path == "/api/v1/jobs/nightly":
			writeStruct(w, 200, "", "ok", &types.JobResponse{ID: 7, Name: "nightly"})
		case r.Method == "PATCH" && r.URL.Path == "/api/v1/jobs/7":
			writeStruct(w, 202, "", "queued", &types.ResponseJobAsync{ID: 7, OperationID: "op_2"})
		case r.Method == "POST" && r.URL.Path == "/api/v1/jobs/7/run":
			writeStruct(w, 202, "", "queued", &types.RunJobResponse{
				ExecutionID: 99, Status: types.JobExecutionStatusPending, OperationID: "op_3",
			})
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/suspend"):
			writeStruct(w, 200, "", "ok", &types.JobResponse{ID: 7, Suspended: true})
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/resume"):
			writeStruct(w, 200, "", "ok", &types.JobResponse{ID: 7, Suspended: false})
		case r.Method == "GET" && r.URL.Path == "/api/v1/jobs/7/executions":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"message":"ok","data":[{"id":11,"job_id":7,"trigger":"schedule","k8s_job_name":"nightly-1","status":"succeeded","created_at":"2026-05-31T00:00:00Z"}],"meta":{"page":1,"page_size":20,"total_items":1,"total_pages":1}}`))
		case r.Method == "GET" && r.URL.Path == "/api/v1/jobs/7/executions/11":
			writeStruct(w, 200, "", "ok", &types.JobExecution{
				ID: 11, JobID: 7, Trigger: types.JobExecutionTriggerSchedule, Status: types.JobExecutionStatusSucceeded,
			})
		case r.Method == "DELETE":
			writeStruct(w, 202, "", "queued", &types.ResponseJobAsync{ID: 7, OperationID: "op_4"})
		}
	})
	ctx := context.Background()

	created, err := c.Jobs().Create(ctx, &types.CreateJobRequest{
		Name: "nightly", Kind: types.JobKindStandalone, PricingSlug: "small",
		Image: "alpine:3.20", Schedule: "0 2 * * *",
	})
	if err != nil || created.OperationID != "op_1" {
		t.Fatalf("Create: %v (%+v)", err, created)
	}

	got, _, err := c.Jobs().Get(ctx, 7)
	if err != nil || got.Name != "nightly" || got.Kind != types.JobKindStandalone {
		t.Fatalf("Get: %v (%+v)", err, got)
	}

	byName, _, err := c.Jobs().GetByName(ctx, "nightly")
	if err != nil || byName.ID != 7 {
		t.Fatalf("GetByName: %v (%+v)", err, byName)
	}

	upd, err := c.Jobs().Update(ctx, 7, &types.UpdateJobRequest{ConcurrencyPolicy: types.JobConcurrencyForbid})
	if err != nil || upd.OperationID != "op_2" {
		t.Fatalf("Update: %v (%+v)", err, upd)
	}

	run, err := c.Jobs().RunNow(ctx, 7)
	if err != nil || run.ExecutionID != 99 {
		t.Fatalf("RunNow: %v (%+v)", err, run)
	}

	susp, err := c.Jobs().Suspend(ctx, 7)
	if err != nil || !susp.Suspended {
		t.Fatalf("Suspend: %v (%+v)", err, susp)
	}

	res, err := c.Jobs().Resume(ctx, 7)
	if err != nil || res.Suspended {
		t.Fatalf("Resume: %v (%+v)", err, res)
	}

	execs, meta, err := c.Jobs().ListExecutions(ctx, 7)
	if err != nil || len(execs) != 1 || meta.TotalItems != 1 {
		t.Fatalf("ListExecutions: %v (%v / %+v)", err, execs, meta)
	}

	exec, err := c.Jobs().GetExecution(ctx, 7, 11)
	if err != nil || exec.Status != types.JobExecutionStatusSucceeded {
		t.Fatalf("GetExecution: %v (%+v)", err, exec)
	}

	del, err := c.Jobs().Delete(ctx, 7)
	if err != nil || del.OperationID != "op_4" {
		t.Fatalf("Delete: %v (%+v)", err, del)
	}
	if seen.method != "DELETE" || seen.path != "/api/v1/jobs/7" {
		t.Errorf("last call: %s %s, want DELETE /api/v1/jobs/7", seen.method, seen.path)
	}
}
