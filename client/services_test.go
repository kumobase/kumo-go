package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kumobase/kumo-go/client"
	"github.com/kumobase/kumo-go/codes"
	"github.com/kumobase/kumo-go/types"
)

// One happy-path round-trip per service confirms the URL, method, body
// shape, and decoding contract. Per-service edge cases (validation, code
// branching, etc.) are exercised by the shared infrastructure tests in
// client_test.go.

// ── Secrets ──────────────────────────────────────────────────────────

func TestSecrets_Smoke(t *testing.T) {
	type call struct{ method, path string }
	var seen call
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seen = call{r.Method, r.URL.Path}
		switch r.Method {
		case "POST":
			writeStruct(w, 201, "", "ok", &types.GetSecretAllResponse{ID: 5, Name: "db", Type: types.SecretTypeEnvVar})
		case "GET":
			writeStruct(w, 200, "", "ok", &types.ResponseGetSecret{ID: 5, Name: "db", Type: types.SecretTypeEnvVar})
		case "PATCH":
			writeStruct(w, 200, "", "ok", &types.ResponseGetSecret{ID: 5, Name: "db-2", Type: types.SecretTypeEnvVar})
		case "DELETE":
			writeStruct(w, 200, "", "ok", nil)
		}
	})
	ctx := context.Background()

	created, err := c.Secrets().Create(ctx, &types.CreateSecretRequest{
		RequestSecretBase: types.RequestSecretBase{Name: "db", Type: types.SecretTypeEnvVar},
		EnvironmentVariables: []types.EnvironmentVariable{{Key: "URL", Value: "x"}},
	})
	if err != nil || created.ID != 5 {
		t.Fatalf("Create: %v (id=%d)", err, created.ID)
	}
	got, _, err := c.Secrets().Get(ctx, 5)
	if err != nil || got.Name != "db" {
		t.Fatalf("Get: %v (%+v)", err, got)
	}
	upd, err := c.Secrets().Update(ctx, 5, &types.UpdateSecretRequest{})
	if err != nil || upd.Name != "db-2" {
		t.Fatalf("Update: %v (%+v)", err, upd)
	}
	if err := c.Secrets().Delete(ctx, 5); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if seen.method != "DELETE" || seen.path != "/api/v1/secrets/5" {
		t.Errorf("last call: %s %s, want DELETE /api/v1/secrets/5", seen.method, seen.path)
	}
}

// ── Volumes ──────────────────────────────────────────────────────────

func TestVolumes_Smoke(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			writeStruct(w, 201, "", "ok", &types.VolumeResponse{ID: 3, Name: "data", SizeGB: 10, Status: "ready"})
		case "GET":
			writeStruct(w, 200, "", "ok", &types.VolumeResponse{ID: 3, SizeGB: 10, Status: "ready"})
		case "DELETE":
			writeStruct(w, 200, "", "ok", nil)
		}
	})
	ctx := context.Background()
	v, err := c.Volumes().Create(ctx, &types.CreateVolumeRequest{Name: "data", StorageTier: "ssd-std", SizeGB: 10})
	if err != nil || v.ID != 3 {
		t.Fatalf("Create: %v (%+v)", err, v)
	}
	got, _, err := c.Volumes().Get(ctx, 3)
	if err != nil || got.SizeGB != 10 {
		t.Fatalf("Get: %v (%+v)", err, got)
	}
	if err := c.Volumes().Delete(ctx, 3); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestVolumes_ResizeAndWait(t *testing.T) {
	var status atomic.Value
	status.Store("resizing")
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {
			// kick off the resize — status flips to "ready" after a moment
			time.AfterFunc(20*time.Millisecond, func() { status.Store("ready") })
			writeStruct(w, 202, "", "queued", &types.VolumeResponse{ID: 3, Status: "resizing"})
			return
		}
		writeStruct(w, 200, "", "ok", &types.VolumeResponse{ID: 3, Status: status.Load().(string)})
	})
	got, err := c.Volumes().ResizeAndWait(context.Background(), 3,
		&types.ResizeVolumeRequest{SizeGB: 20},
		client.WithPollInterval(10*time.Millisecond),
	)
	if err != nil {
		t.Fatalf("ResizeAndWait: %v", err)
	}
	if got.Status != "ready" {
		t.Errorf("final status: %s, want ready", got.Status)
	}
}

// ── VPS ──────────────────────────────────────────────────────────────

func TestVPS_Smoke(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/vps/regions":
			writeStruct(w, 200, "", "ok", []types.VPSRegionResponse{{ID: "sgp", Name: "Singapore"}})
		case "/api/v1/vps/plans":
			writeStruct(w, 200, "", "ok", []types.PublicVPSPlanResponse{{PlanID: 1, Name: "1c1g", SellingPrice: "4.99"}})
		case "/api/v1/vps/servers":
			if r.Method == "POST" {
				writeStruct(w, 201, "", "ok", &types.VPSServerResponse{ID: 9, Status: "provisioning"})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			fmt.Fprint(w, `{"message":"ok","data":[{"id":9,"status":"running","display_provider":"vultr","region_id":"sgp","ssh_port":22,"auto_renew":true,"created_at":"2026-05-23T10:00:00Z","ssh_setup_completed":true,"action_status":""}],"meta":{"page":1,"page_size":20,"total_items":1,"total_pages":1}}`)
		default:
			writeStruct(w, 200, "", "ok", &types.VPSServerResponse{ID: 9, Status: "running"})
		}
	})
	ctx := context.Background()
	regions, err := c.VPS().ListRegions(ctx)
	if err != nil || len(regions) != 1 {
		t.Fatalf("ListRegions: %v (%v)", err, regions)
	}
	plans, err := c.VPS().ListPlans(ctx)
	if err != nil || len(plans) != 1 || plans[0].SellingPrice != "4.99" {
		t.Fatalf("ListPlans: %v (%v)", err, plans)
	}
	servers, _, err := c.VPS().ListServers(ctx)
	if err != nil || len(servers) != 1 {
		t.Fatalf("ListServers: %v (%v)", err, servers)
	}
	rented, err := c.VPS().RentServer(ctx, &types.RentServerRequest{Provider: "vultr", Region: "sgp", Plan: "1c1g", Name: "edge"})
	if err != nil || rented.ID != 9 {
		t.Fatalf("RentServer: %v (%+v)", err, rented)
	}
}

func TestVPS_RebootAndWait(t *testing.T) {
	var action atomic.Value
	action.Store("rebooting")
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			time.AfterFunc(20*time.Millisecond, func() { action.Store("") })
			writeStruct(w, 202, "", "queued", nil)
			return
		}
		writeStruct(w, 200, "", "ok", &types.VPSServerResponse{ID: 9, ActionStatus: action.Load().(string)})
	})
	got, err := c.VPS().RebootAndWait(context.Background(), 9,
		client.WithPollInterval(10*time.Millisecond),
	)
	if err != nil {
		t.Fatalf("RebootAndWait: %v", err)
	}
	if got.ActionStatus != "" {
		t.Errorf("final ActionStatus: %q, want \"\"", got.ActionStatus)
	}
}

// ── Registry ─────────────────────────────────────────────────────────

func TestRegistry_OrgsSmoke(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			writeStruct(w, 201, "", "ok", &types.OrganizationResponse{ID: 1, Slug: "acme", DisplayName: "Acme"})
		case "GET":
			if r.URL.Path == "/api/v1/registry/organizations" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				fmt.Fprint(w, `{"message":"ok","data":[{"id":1,"slug":"acme","display_name":"Acme","owner_user_id":1,"registry_auto_create_repos":true,"created_at":"2026-05-23T10:00:00Z","updated_at":"2026-05-23T10:00:00Z"}]}`)
				return
			}
			writeStruct(w, 200, "", "ok", &types.OrganizationResponse{ID: 1, Slug: "acme", DisplayName: "Acme"})
		}
	})
	ctx := context.Background()
	org, err := c.Registry().Orgs().Create(ctx, &types.CreateOrganizationRequest{Slug: "acme", DisplayName: "Acme"})
	if err != nil || org.Slug != "acme" {
		t.Fatalf("Create: %v (%+v)", err, org)
	}
	got, _, err := c.Registry().Orgs().Get(ctx, "acme")
	if err != nil || got.DisplayName != "Acme" {
		t.Fatalf("Get: %v (%+v)", err, got)
	}
	orgs, err := c.Registry().Orgs().List(ctx)
	if err != nil || len(orgs) != 1 {
		t.Fatalf("List: %v (%v)", err, orgs)
	}
}

func TestRegistry_ReposScopedBySlug(t *testing.T) {
	var seenPath string
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		writeStruct(w, 201, "", "ok", &types.RepositoryResponse{ID: 1, Name: "demo", TagMutability: types.TagMutabilityMutable})
	})
	repos := c.Registry().Repos("acme")
	_, err := repos.Create(context.Background(), &types.CreateRepositoryRequest{Name: "demo"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	want := "/api/v1/registry/organizations/acme/repositories"
	if seenPath != want {
		t.Errorf("path: got %q, want %q", seenPath, want)
	}
}

// ── Billing / Profile ────────────────────────────────────────────────

func TestBilling_GetSummary(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeStruct(w, 200, "", "ok", &types.BillingSummaryResponse{
			Currency: "IDR", PreviousPeriodTotal: "10000.00",
		})
	})
	got, err := c.Billing().GetSummary(context.Background())
	if err != nil || got.Currency != "IDR" {
		t.Fatalf("GetSummary: %v (%+v)", err, got)
	}
}

func TestProfile_Smoke(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/profile":
			writeStruct(w, 200, "", "ok", &types.GetProfileResponse{FullName: "Alice", Email: "a@example.com"})
		case "/api/v1/profile/balance":
			writeStruct(w, 200, "", "ok", &types.GetBalanceResponse{Balance: "100.50"})
		case "/api/v1/profile/has-password":
			writeStruct(w, 200, "", "ok", &types.HasPasswordResponse{HasPassword: true})
		}
	})
	ctx := context.Background()
	prof, err := c.Profile().Get(ctx)
	if err != nil || prof.FullName != "Alice" {
		t.Fatalf("Get: %v (%+v)", err, prof)
	}
	bal, err := c.Profile().GetBalance(ctx)
	if err != nil || bal.Balance != "100.50" {
		t.Fatalf("GetBalance: %v (%+v)", err, bal)
	}
	hp, err := c.Profile().HasPassword(ctx)
	if err != nil || !hp {
		t.Fatalf("HasPassword: %v (%v)", err, hp)
	}
}

// ── SourceConnections ────────────────────────────────────────────────

func TestSourceConnections_Smoke(t *testing.T) {
	type call struct{ method, path string }
	var seen call
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seen = call{r.Method, r.URL.Path}
		switch {
		case r.Method == "GET" && r.URL.Path == "/api/v1/source-connections":
			fmt.Fprint(w, `{"message":"ok","data":[{"id":1,"provider":"github","installation_id":12345,"account_login":"acme","account_type":"Organization","status":"active","created_at":"2026-05-23T10:00:00Z","updated_at":"2026-05-23T10:00:00Z"}]}`)
		case r.Method == "GET" && r.URL.Path == "/api/v1/source-connections/1/repos":
			fmt.Fprint(w, `{"message":"ok","data":[{"id":99,"full_name":"acme/web","private":true,"default_branch":"main"}]}`)
		case r.Method == "DELETE":
			writeStruct(w, 200, "", "ok", &types.SourceConnectionResponse{ID: 1, Provider: types.SourceProviderGitHub, Status: types.SourceConnectionStatusActive})
		}
	})
	ctx := context.Background()

	conns, err := c.SourceConnections().List(ctx)
	if err != nil || len(conns) != 1 || conns[0].AccountLogin != "acme" {
		t.Fatalf("List: %v (%+v)", err, conns)
	}
	repos, err := c.SourceConnections().ListRepos(ctx, 1)
	if err != nil || len(repos) != 1 || repos[0].FullName != "acme/web" {
		t.Fatalf("ListRepos: %v (%+v)", err, repos)
	}
	disc, err := c.SourceConnections().Disconnect(ctx, 1)
	if err != nil || disc.ID != 1 {
		t.Fatalf("Disconnect: %v (%+v)", err, disc)
	}
	if seen.method != "DELETE" || seen.path != "/api/v1/source-connections/1" {
		t.Errorf("last call: %s %s, want DELETE /api/v1/source-connections/1", seen.method, seen.path)
	}
}

// ── Builds (git-build apps) ──────────────────────────────────────────

func TestBuilds_Smoke(t *testing.T) {
	type call struct{ method, path string }
	var seen call
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seen = call{r.Method, r.URL.Path}
		switch {
		case r.Method == "POST" && r.URL.Path == "/api/v1/source-connections/1/apps":
			writeStruct(w, 202, "", "queued", &types.CreateAppResponse{ID: 7, Name: "my-app", GenerateAppName: "my-appx4z8a1b2"})
		case r.Method == "GET" && r.URL.Path == "/api/v1/apps/7/builds":
			fmt.Fprint(w, `{"message":"ok","data":[{"id":3,"app_id":7,"commit_sha":"abc","ref":"refs/heads/main","status":"succeeded","image_digest":"sha256:x","created_at":"2026-05-23T10:00:00Z"}],"meta":{"page":1,"page_size":20,"total_items":1,"total_pages":1}}`)
		case r.Method == "GET" && r.URL.Path == "/api/v1/apps/7/builds/3":
			writeStruct(w, 200, "", "ok", &types.BuildResponse{ID: 3, AppID: 7, Status: types.BuildStatusSucceeded, LogURL: "https://logs/3.txt?sig=x"})
		case r.Method == "POST" && r.URL.Path == "/api/v1/apps/7/builds":
			writeStruct(w, 202, "", "queued", &types.BuildResponse{ID: 4, AppID: 7, Status: types.BuildStatusPending})
		case r.Method == "POST" && r.URL.Path == "/api/v1/apps/7/builds/4/cancel":
			writeStruct(w, 200, "", "ok", &types.BuildResponse{ID: 4, AppID: 7, Status: types.BuildStatusCanceled})
		}
	})
	ctx := context.Background()

	created, err := c.Builds().CreateGitBuildApp(ctx, 1, &types.CreateGitBuildAppRequest{
		Name: "my-app", Port: 8080, IsExposed: true, Replicas: 1,
		RepoFullName: "acme/web", Branch: "main", PricingSlug: "kumo.nano",
	})
	if err != nil || created.ID != 7 {
		t.Fatalf("CreateGitBuildApp: %v (%+v)", err, created)
	}
	builds, _, err := c.Builds().List(ctx, 7)
	if err != nil || len(builds) != 1 || builds[0].Status != types.BuildStatusSucceeded {
		t.Fatalf("List: %v (%+v)", err, builds)
	}
	got, err := c.Builds().Get(ctx, 7, 3)
	if err != nil || got.LogURL == "" {
		t.Fatalf("Get: %v (%+v)", err, got)
	}
	rebuilt, err := c.Builds().Rebuild(ctx, 7)
	if err != nil || rebuilt.ID != 4 {
		t.Fatalf("Rebuild: %v (%+v)", err, rebuilt)
	}
	canceled, err := c.Builds().Cancel(ctx, 7, 4)
	if err != nil || canceled.Status != types.BuildStatusCanceled {
		t.Fatalf("Cancel: %v (%+v)", err, canceled)
	}
	if seen.method != "POST" || seen.path != "/api/v1/apps/7/builds/4/cancel" {
		t.Errorf("last call: %s %s, want POST /api/v1/apps/7/builds/4/cancel", seen.method, seen.path)
	}
}

// ── APIKeys (sessionOnly) ────────────────────────────────────────────

func TestAPIKeys_SessionOnlyErrorSurfaces(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeStruct(w, 403, codes.APIKeySessionRequired, "api keys cannot manage api keys", nil)
	})
	_, err := c.APIKeys().Create(context.Background(), &types.CreateAPIKeyRequest{Name: "ci"})
	if !client.IsCode(err, codes.APIKeySessionRequired) {
		t.Errorf("expected IsCode(APIKeySessionRequired), got %v", err)
	}
}

func TestAPIKeys_Smoke(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(201)
			// APIKeyCreateResponse embeds APIKeyResponse — write directly.
			payload := map[string]any{
				"id": 1, "name": "ci", "key_prefix": "kumo_sk_abc",
				"created_at": time.Now().UTC().Format(time.RFC3339),
				"scopes":     []string{"read", "write"},
				"key":        "kumo_sk_full_secret",
			}
			b, _ := json.Marshal(payload)
			fmt.Fprintf(w, `{"message":"ok","data":%s}`, string(b))
		case "GET":
			fmt.Fprint(w, `{"message":"ok","data":[{"id":1,"name":"ci","key_prefix":"kumo_sk_abc","created_at":"2026-05-23T10:00:00Z","scopes":["read","write"]}]}`)
		case "DELETE":
			writeStruct(w, 200, "", "ok", nil)
		}
	})
	ctx := context.Background()
	created, err := c.APIKeys().Create(ctx, &types.CreateAPIKeyRequest{Name: "ci"})
	if err != nil || created.Key != "kumo_sk_full_secret" {
		t.Fatalf("Create: %v (%+v)", err, created)
	}
	keys, err := c.APIKeys().List(ctx)
	if err != nil || len(keys) != 1 {
		t.Fatalf("List: %v (%v)", err, keys)
	}
	if err := c.APIKeys().Delete(ctx, 1); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}
