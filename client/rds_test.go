package client_test

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kumobase/kumo-go/client"
	"github.com/kumobase/kumo-go/codes"
	"github.com/kumobase/kumo-go/types"
)

func TestRDS_Smoke(t *testing.T) {
	type call struct{ method, path string }
	var seen call
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seen = call{r.Method, r.URL.Path}
		switch {
		case r.Method == "GET" && r.URL.Path == "/api/v1/rds/plans":
			writeStruct(w, 200, "", "ok", []types.PublicRDSPlanResponse{
				{Slug: "kumo.pg.small", Engine: types.RDSEnginePostgreSQL, Name: "Small", CPUvCPU: "1", MemoryMB: 2048, PriceHour: "12.5000"},
			})
		case r.Method == "POST" && r.URL.Path == "/api/v1/rds":
			writeStruct(w, 202, "", "queued", &types.RDSMutationResponse{ID: 7, OperationID: "op-1", Status: string(types.RDSStatusProvisioning)})
		case r.Method == "POST" && r.URL.Path == "/api/v1/rds/7/start":
			writeStruct(w, 202, "", "queued", &types.RDSMutationResponse{ID: 7, OperationID: "op-4", Status: string(types.RDSStatusProvisioning)})
		case r.Method == "POST" && r.URL.Path == "/api/v1/rds/7/switchover":
			writeStruct(w, 202, "", "queued", &types.RDSMutationResponse{ID: 7, OperationID: "op-5", Status: string(types.RDSStatusSwitchingOver)})
		case r.Method == "GET" && r.URL.Path == "/api/v1/rds/7":
			writeStruct(w, 200, "", "ok", &types.RDSInstanceResponse{ID: 7, Name: "my-pg", Engine: types.RDSEnginePostgreSQL, Status: string(types.RDSStatusReady)})
		case r.Method == "GET" && r.URL.Path == "/api/v1/rds/7/connection":
			writeStruct(w, 200, "", "ok", &types.RDSConnectionResponse{Host: "h", Port: 5432, Username: "kumo", Database: "kumo", Password: "s", SSLMode: "require"})
		case r.Method == "PATCH" && r.URL.Path == "/api/v1/rds/7/tls":
			writeStruct(w, 202, "", "queued", &types.RDSMutationResponse{ID: 7, OperationID: "op-6", Status: string(types.RDSStatusReconfiguring)})
		case r.Method == "PATCH" && r.URL.Path == "/api/v1/rds/7":
			writeStruct(w, 202, "", "queued", &types.RDSMutationResponse{ID: 7, OperationID: "op-2", Status: string(types.RDSStatusScaling)})
		case r.Method == "DELETE" && r.URL.Path == "/api/v1/rds/7":
			writeStruct(w, 202, "", "queued", &types.RDSMutationResponse{ID: 7, OperationID: "op-3", Status: string(types.RDSStatusDeleting)})
		}
	})
	ctx := context.Background()

	plans, err := c.RDS().ListPlans(ctx)
	if err != nil || len(plans) != 1 || plans[0].Slug != "kumo.pg.small" {
		t.Fatalf("ListPlans: %v (%+v)", err, plans)
	}
	created, err := c.RDS().Create(ctx, &types.CreateRDSInstanceRequest{
		Name: "my-pg", Engine: types.RDSEnginePostgreSQL, EngineVersion: "16", Plan: "kumo.pg.small", StorageGB: 20,
	})
	if err != nil || created.ID != 7 || created.OperationID != "op-1" {
		t.Fatalf("Create: %v (%+v)", err, created)
	}
	got, err := c.RDS().Get(ctx, 7)
	if err != nil || got.Name != "my-pg" {
		t.Fatalf("Get: %v (%+v)", err, got)
	}
	conn, err := c.RDS().GetConnection(ctx, 7)
	if err != nil || conn.Port != 5432 {
		t.Fatalf("GetConnection: %v (%+v)", err, conn)
	}
	scaled, err := c.RDS().Scale(ctx, 7, "kumo.pg.medium")
	if err != nil || scaled.OperationID != "op-2" {
		t.Fatalf("Scale: %v (%+v)", err, scaled)
	}
	resized, err := c.RDS().Resize(ctx, 7, 50)
	if err != nil || resized.OperationID != "op-2" {
		t.Fatalf("Resize: %v (%+v)", err, resized)
	}
	started, err := c.RDS().Start(ctx, 7)
	if err != nil || started.OperationID != "op-4" {
		t.Fatalf("Start: %v (%+v)", err, started)
	}
	switched, err := c.RDS().Switchover(ctx, 7)
	if err != nil || switched.OperationID != "op-5" || switched.Status != string(types.RDSStatusSwitchingOver) {
		t.Fatalf("Switchover: %v (%+v)", err, switched)
	}
	tlsd, err := c.RDS().UpdateTLS(ctx, 7, string(types.RDSTLSModeRequired))
	if err != nil || tlsd.OperationID != "op-6" || tlsd.Status != string(types.RDSStatusReconfiguring) {
		t.Fatalf("UpdateTLS: %v (%+v)", err, tlsd)
	}
	del, err := c.RDS().Delete(ctx, 7)
	if err != nil || del.OperationID != "op-3" {
		t.Fatalf("Delete: %v (%+v)", err, del)
	}
	if seen.method != "DELETE" || seen.path != "/api/v1/rds/7" {
		t.Errorf("last call: %s %s, want DELETE /api/v1/rds/7", seen.method, seen.path)
	}
}

func TestRDS_NotReadyError(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeStruct(w, 409, codes.RDSInstanceNotReady, "database is not ready", nil)
	})
	_, err := c.RDS().GetConnection(context.Background(), 7)
	if !client.IsCode(err, codes.RDSInstanceNotReady) {
		t.Errorf("expected IsCode(RDSInstanceNotReady), got %v", err)
	}
}

func TestRDS_CreateAndWait(t *testing.T) {
	var status atomic.Value
	status.Store(string(types.RDSStatusProvisioning))
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			time.AfterFunc(20*time.Millisecond, func() { status.Store(string(types.RDSStatusReady)) })
			writeStruct(w, 202, "", "queued", &types.RDSMutationResponse{ID: 7, OperationID: "op-1", Status: string(types.RDSStatusProvisioning)})
			return
		}
		writeStruct(w, 200, "", "ok", &types.RDSInstanceResponse{ID: 7, Status: status.Load().(string)})
	})
	got, err := c.RDS().CreateAndWait(context.Background(),
		&types.CreateRDSInstanceRequest{Name: "my-pg", Engine: types.RDSEnginePostgreSQL, EngineVersion: "16", Plan: "kumo.pg.small", StorageGB: 20},
		client.WithPollInterval(10*time.Millisecond),
	)
	if err != nil {
		t.Fatalf("CreateAndWait: %v", err)
	}
	if got.Status != string(types.RDSStatusReady) {
		t.Errorf("final status: %s, want ready", got.Status)
	}
}
