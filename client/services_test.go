package client_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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

// ── Name addressing (id-or-name) ─────────────────────────────────────

// GetByName must hit the same detail endpoint as Get with the name in the
// path segment (the server resolves a non-numeric segment as a name).
func TestNameAddressing_GetByName(t *testing.T) {
	var seen string
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seen = r.URL.Path
		switch r.URL.Path {
		case "/api/v1/apps/my-api":
			writeStruct(w, 200, "", "ok", &types.AppByIdResponse{Id: 7})
		case "/api/v1/secrets/db-creds":
			writeStruct(w, 200, "", "ok", &types.ResponseGetSecret{ID: 8, Name: "db-creds", Type: types.SecretTypeEnvVar})
		case "/api/v1/volumes/data-vol":
			writeStruct(w, 200, "", "ok", &types.VolumeResponse{ID: 9, Name: "data-vol", Status: "ready"})
		default:
			writeStruct(w, 404, codes.AppNotFound, "not found", nil)
		}
	})
	ctx := context.Background()

	if app, _, err := c.Apps().GetByName(ctx, "my-api"); err != nil || app.Id != 7 {
		t.Fatalf("Apps.GetByName: %v (%+v)", err, app)
	}
	if sec, _, err := c.Secrets().GetByName(ctx, "db-creds"); err != nil || sec.ID != 8 {
		t.Fatalf("Secrets.GetByName: %v (%+v)", err, sec)
	}
	if vol, _, err := c.Volumes().GetByName(ctx, "data-vol"); err != nil || vol.ID != 9 {
		t.Fatalf("Volumes.GetByName: %v (%+v)", err, vol)
	}
	if seen != "/api/v1/volumes/data-vol" {
		t.Errorf("last path = %q, want /api/v1/volumes/data-vol", seen)
	}
}

// VPS GetServerByName + API keys Get/GetByName — same id-or-name mechanism,
// added in v0.6.0 alongside cross-resource attach-by-name DTOs.
func TestNameAddressing_VPSAndAPIKeys(t *testing.T) {
	var seen string
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seen = r.URL.Path
		switch r.URL.Path {
		case "/api/v1/vps/servers/edge-1":
			writeStruct(w, 200, "", "ok", &types.VPSServerResponse{ID: 11, DisplayName: "edge-1", Status: "running"})
		case "/api/v1/api-keys/ci-deployer":
			writeStruct(w, 200, "", "ok", &types.APIKeyResponse{ID: 12, Name: "ci-deployer"})
		case "/api/v1/api-keys/13":
			writeStruct(w, 200, "", "ok", &types.APIKeyResponse{ID: 13, Name: "byid"})
		default:
			writeStruct(w, 404, codes.AppNotFound, "not found", nil)
		}
	})
	ctx := context.Background()

	if srv, err := c.VPS().GetServerByName(ctx, "edge-1"); err != nil || srv.ID != 11 {
		t.Fatalf("VPS.GetServerByName: %v (%+v)", err, srv)
	}
	if k, err := c.APIKeys().GetByName(ctx, "ci-deployer"); err != nil || k.ID != 12 {
		t.Fatalf("APIKeys.GetByName: %v (%+v)", err, k)
	}
	if k, err := c.APIKeys().Get(ctx, 13); err != nil || k.ID != 13 {
		t.Fatalf("APIKeys.Get: %v (%+v)", err, k)
	}
	if seen != "/api/v1/api-keys/13" {
		t.Errorf("last path = %q", seen)
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
		case r.Method == "GET" && r.URL.Path == "/api/v1/source-connections/1/repos" && r.URL.RawQuery != "":
			// Paginated/filtered surface: echoes a Meta block. Captures the
			// query so the test can assert page/page_size/q were forwarded.
			seen.method, seen.path = r.Method, r.URL.Path+"?"+r.URL.RawQuery
			fmt.Fprint(w, `{"message":"ok","data":[{"id":99,"full_name":"acme/web","private":true,"default_branch":"main"}],"meta":{"page":1,"page_size":30,"total_items":1,"total_pages":1}}`)
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
	// Unpaginated surface: full list, nil Meta (back-compat).
	repos, meta, err := c.SourceConnections().ListRepos(ctx, 1)
	if err != nil || len(repos) != 1 || repos[0].FullName != "acme/web" {
		t.Fatalf("ListRepos: %v (%+v)", err, repos)
	}
	if meta != nil {
		t.Errorf("unpaginated ListRepos should return nil Meta, got %+v", meta)
	}
	// Paginated/filtered surface: page+size+q forwarded, Meta returned.
	repos, meta, err = c.SourceConnections().ListRepos(ctx, 1,
		client.WithPage(1), client.WithPageSize(30), client.WithExtraQuery("q", "web"))
	if err != nil || len(repos) != 1 || meta == nil || meta.TotalItems != 1 {
		t.Fatalf("ListRepos paginated: %v (repos=%+v meta=%+v)", err, repos, meta)
	}
	if !strings.Contains(seen.path, "page=1") || !strings.Contains(seen.path, "page_size=30") || !strings.Contains(seen.path, "q=web") {
		t.Errorf("ListRepos paginated query not forwarded: %s", seen.path)
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
			writeStruct(w, 200, "", "ok", &types.BuildResponse{ID: 3, AppID: 7, Status: types.BuildStatusSucceeded})
		case r.Method == "GET" && r.URL.Path == "/api/v1/apps/7/builds/3/log-url":
			writeStruct(w, 200, "", "ok", &types.BuildLogURLResponse{LogURL: "https://logs/3.txt?sig=x"})
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
	if err != nil || got.ID != 3 {
		t.Fatalf("Get: %v (%+v)", err, got)
	}
	logURL, err := c.Builds().GetLogURL(ctx, 7, 3)
	if err != nil || logURL != "https://logs/3.txt?sig=x" {
		t.Fatalf("GetLogURL: %v (%q)", err, logURL)
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

// GetLogURL surfaces BUILD_LOG_NOT_AVAILABLE (404) as a typed APIError when the
// build has no persisted log yet.
func TestBuilds_GetLogURL_NotAvailable(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeStruct(w, 404, codes.BuildLogNotAvailable, "build log not available", nil)
	})
	_, err := c.Builds().GetLogURL(context.Background(), 7, 3)
	var apiErr *client.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T (%v)", err, err)
	}
	if apiErr.Code != codes.BuildLogNotAvailable {
		t.Errorf("Code: got %q, want %q", apiErr.Code, codes.BuildLogNotAvailable)
	}
}

// ── Public plan catalogues (pricing) ─────────────────────────────────

// Apps().ListPlans hits GET /api/v1/apps/plans and flattens the
// {"templates":[…]} wrapper to the inner slice.
func TestApps_ListPlans(t *testing.T) {
	var seen string
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seen = r.URL.Path
		writeStruct(w, 200, "", "ok", &types.PricingResponse{
			Templates: []types.TemplateWithPricing{
				{Slug: "kumo.nano", Name: "Nano", PriceMonth: "1080.00"},
				{Slug: "kumo.small", Name: "Small", PriceMonth: "2160.00"},
			},
		})
	})
	plans, err := c.Apps().ListPlans(context.Background())
	if err != nil || len(plans) != 2 || plans[0].Slug != "kumo.nano" || plans[0].PriceMonth != "1080.00" {
		t.Fatalf("ListPlans: %v (%+v)", err, plans)
	}
	if seen != "/api/v1/apps/plans" {
		t.Errorf("path: got %q, want /api/v1/apps/plans", seen)
	}
}

// Registry().ListPlans hits GET /api/v1/registry/plans and flattens the
// {"plans":[…]} wrapper.
func TestRegistry_ListPlans(t *testing.T) {
	var seen string
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seen = r.URL.Path
		writeStruct(w, 200, "", "ok", &types.RegistryPricingResponse{
			Plans: []types.RegistryPlanOption{
				{ID: 1, Name: "Storage", Unit: "GB-month", PricePerUnit: "0.50", Currency: "IDR", ChargeModel: "metered", BillingPeriod: "monthly"},
			},
		})
	})
	plans, err := c.Registry().ListPlans(context.Background())
	if err != nil || len(plans) != 1 || plans[0].PricePerUnit != "0.50" {
		t.Fatalf("ListPlans: %v (%+v)", err, plans)
	}
	if seen != "/api/v1/registry/plans" {
		t.Errorf("path: got %q, want /api/v1/registry/plans", seen)
	}
}

// Volumes().ListPlans hits GET /api/v1/volumes/plans — the one paginated
// catalogue, so it returns *Meta alongside the slice.
func TestVolumes_ListPlans(t *testing.T) {
	var seen string
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seen = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprint(w, `{"message":"ok","data":[{"id":1,"slug":"ssd-std","name":"SSD Standard","price_per_gb_hour":"0.0001234","min_size_gb":1,"max_size_gb":1000}],"meta":{"page":1,"page_size":20,"total_items":1,"total_pages":1}}`)
	})
	tiers, meta, err := c.Volumes().ListPlans(context.Background())
	if err != nil || len(tiers) != 1 || tiers[0].Slug != "ssd-std" || tiers[0].PricePerGBHour != "0.0001234" {
		t.Fatalf("ListPlans: %v (%+v)", err, tiers)
	}
	if meta == nil || meta.TotalItems != 1 {
		t.Fatalf("ListPlans meta: %+v", meta)
	}
	if seen != "/api/v1/volumes/plans" {
		t.Errorf("path: got %q, want /api/v1/volumes/plans", seen)
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
