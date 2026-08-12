package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/kumobase/kumo-go/types"
)

func writePublicEnvelope(w http.ResponseWriter, status int, data string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, `{"message":"ok","data":%s}`, data)
}

func TestNewPublicValidation(t *testing.T) {
	if _, err := NewPublic(""); err == nil {
		t.Fatal("expected empty base URL to fail")
	}
	for name, opt := range map[string]Option{
		"api key": WithAPIKey("kumo_sk_test"),
		"jwt":     WithJWT("jwt"),
	} {
		t.Run(name, func(t *testing.T) {
			_, err := NewPublic("https://api.example.com", opt)
			if !errors.Is(err, ErrPublicClientRestricted) {
				t.Fatalf("expected ErrPublicClientRestricted, got %v", err)
			}
		})
	}
	_, err := NewPublic("https://api.example.com", WithAPIKey("key"), WithJWT("jwt"))
	if !errors.Is(err, ErrPublicClientRestricted) {
		t.Fatalf("expected both auth options to be rejected, got %v", err)
	}
}

func TestPublicClient_AllCatalogues(t *testing.T) {
	var requests atomic.Int32
	var logged atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		if got := r.Header.Get("Authorization"); got != "" {
			t.Errorf("public request sent Authorization %q", got)
		}
		if got := r.Header.Get("User-Agent"); !strings.Contains(got, "catalogue-test/1.0") {
			t.Errorf("User-Agent %q does not include caller suffix", got)
		}
		switch r.URL.Path {
		case "/api/v1/apps/plans":
			writePublicEnvelope(w, http.StatusOK, `{"templates":[]}`)
		case "/api/v1/vps/regions":
			writePublicEnvelope(w, http.StatusOK, `[]`)
		case "/api/v1/vps/providers":
			writePublicEnvelope(w, http.StatusOK, `[{"name":"Tencent Cloud"}]`)
		case "/api/v1/vps/plans":
			if got := r.URL.Query().Get("region"); got != "jakarta" {
				t.Errorf("region query = %q, want jakarta", got)
			}
			writePublicEnvelope(w, http.StatusOK, `[]`)
		case "/api/v1/volumes/plans":
			writePublicEnvelope(w, http.StatusOK, `[]`)
		case "/api/v1/registry/plans":
			writePublicEnvelope(w, http.StatusOK, `{"plans":[]}`)
		case "/api/v1/packages/plans":
			writePublicEnvelope(w, http.StatusOK, `{"plans":[]}`)
		case "/api/v1/runners/plans":
			writePublicEnvelope(w, http.StatusOK, `[{"label":"kumo-2c-4g","display_name":"2 vCPU / 4 GB","cpu":2,"memory_mb":4096,"price_per_minute":"12.5","currency":"IDR"}]`)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	c, err := NewPublic(srv.URL,
		WithRetries(0),
		WithUserAgent("catalogue-test/1.0"),
		WithLogger(func(string, ...any) { logged.Add(1) }),
	)
	if err != nil {
		t.Fatalf("NewPublic: %v", err)
	}
	ctx := context.Background()

	apps, err := c.Apps().ListPlans(ctx)
	if err != nil || apps == nil || len(apps) != 0 {
		t.Fatalf("Apps.ListPlans: plans=%v err=%v", apps, err)
	}
	regions, err := c.VPS().ListRegions(ctx)
	if err != nil || regions == nil || len(regions) != 0 {
		t.Fatalf("VPS.ListRegions: regions=%v err=%v", regions, err)
	}
	providers, err := c.VPS().ListProviders(ctx)
	if err != nil || len(providers) != 1 || providers[0].Name != "Tencent Cloud" {
		t.Fatalf("VPS.ListProviders: providers=%v err=%v", providers, err)
	}
	vpsPlans, err := c.VPS().ListPlans(ctx, WithExtraQuery("region", "jakarta"))
	if err != nil || vpsPlans == nil || len(vpsPlans) != 0 {
		t.Fatalf("VPS.ListPlans: plans=%v err=%v", vpsPlans, err)
	}
	volumePlans, _, err := c.Volumes().ListPlans(ctx)
	if err != nil || volumePlans == nil || len(volumePlans) != 0 {
		t.Fatalf("Volumes.ListPlans: plans=%v err=%v", volumePlans, err)
	}
	registryPlans, err := c.Registry().ListPlans(ctx)
	if err != nil || registryPlans == nil || len(registryPlans) != 0 {
		t.Fatalf("Registry.ListPlans: plans=%v err=%v", registryPlans, err)
	}
	packagePlans, err := c.Packages().ListPlans(ctx)
	if err != nil || packagePlans == nil || len(packagePlans) != 0 {
		t.Fatalf("Packages.ListPlans: plans=%v err=%v", packagePlans, err)
	}
	runnerPlans, err := c.Runners().ListPlans(ctx)
	if err != nil || len(runnerPlans) != 1 || runnerPlans[0].Label != "kumo-2c-4g" {
		t.Fatalf("Runners.ListPlans: plans=%v err=%v", runnerPlans, err)
	}
	if requests.Load() != 8 || logged.Load() != 8 {
		t.Fatalf("requests=%d logs=%d, want 8 each", requests.Load(), logged.Load())
	}
}

func TestPublicClient_RejectsDisallowedRequestsWithoutHTTP(t *testing.T) {
	var requests atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		writePublicEnvelope(w, http.StatusOK, `[]`)
	}))
	t.Cleanup(srv.Close)
	c, err := NewPublic(srv.URL)
	if err != nil {
		t.Fatalf("NewPublic: %v", err)
	}

	checks := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/apps"},
		{http.MethodGet, "/api/v1/apps/plans/extra"},
		{http.MethodGet, "/api/v1/apps/plans-extra"},
		{http.MethodGet, "/api/v1/apps/%70lans"},
		{http.MethodGet, "/api/v1/vps/regions/"},
		{http.MethodPost, "/api/v1/apps/plans"},
		{http.MethodPut, "/api/v1/apps/plans"},
		{http.MethodPatch, "/api/v1/apps/plans"},
		{http.MethodDelete, "/api/v1/apps/plans"},
	}
	for _, tc := range checks {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			_, _, err := c.do(context.Background(), tc.method, tc.path, nil, nil, nil)
			if !errors.Is(err, ErrPublicClientRestricted) {
				t.Fatalf("expected ErrPublicClientRestricted, got %v", err)
			}
		})
	}
	if requests.Load() != 0 {
		t.Fatalf("disallowed calls contacted server %d times", requests.Load())
	}
}

func TestPublicClient_AuthenticatedClientStillUsesCatalogues(t *testing.T) {
	var authHeaders atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "Bearer kumo_sk_test" {
			authHeaders.Add(1)
		}
		writePublicEnvelope(w, http.StatusOK, `[]`)
	}))
	t.Cleanup(srv.Close)
	c, err := New(srv.URL, WithAPIKey("kumo_sk_test"), WithRetries(0))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := c.VPS().ListProviders(context.Background()); err != nil {
		t.Fatalf("VPS.ListProviders: %v", err)
	}
	if _, err := c.Runners().ListPlans(context.Background()); err != nil {
		t.Fatalf("Runners.ListPlans: %v", err)
	}
	if authHeaders.Load() != 2 {
		t.Fatalf("authenticated catalogue requests with auth header=%d, want 2", authHeaders.Load())
	}
}

func TestPublicClient_ResponseErrorsAndCancellation(t *testing.T) {
	t.Run("malformed", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(`{`)) }))
		defer srv.Close()
		c, _ := NewPublic(srv.URL, WithRetries(0))
		if _, err := c.Runners().ListPlans(context.Background()); err == nil || !strings.Contains(err.Error(), "decode response envelope") {
			t.Fatalf("expected envelope decode error, got %v", err)
		}
	})

	for _, status := range []int{http.StatusNotFound, http.StatusInternalServerError} {
		t.Run(fmt.Sprintf("http_%d", status), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				writePublicEnvelope(w, status, `null`)
			}))
			defer srv.Close()
			c, _ := NewPublic(srv.URL, WithRetries(0))
			_, err := c.VPS().ListProviders(context.Background())
			var apiErr *APIError
			if !errors.As(err, &apiErr) || apiErr.StatusCode != status {
				t.Fatalf("expected APIError status %d, got %v", status, err)
			}
		})
	}

	t.Run("cancelled", func(t *testing.T) {
		var requests atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			requests.Add(1)
			writePublicEnvelope(w, http.StatusOK, `[]`)
		}))
		defer srv.Close()
		c, _ := NewPublic(srv.URL)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := c.VPS().ListProviders(ctx)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
		if requests.Load() != 0 {
			t.Fatalf("cancelled request contacted server %d times", requests.Load())
		}
	})
}

func TestPublicClient_Retries(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if attempts.Add(1) < 3 {
			writePublicEnvelope(w, http.StatusServiceUnavailable, `null`)
			return
		}
		writePublicEnvelope(w, http.StatusOK, `[{"name":"Tencent Cloud"}]`)
	}))
	defer srv.Close()
	c, _ := NewPublic(srv.URL, WithRetries(3))
	providers, err := c.VPS().ListProviders(context.Background())
	if err != nil || len(providers) != 1 {
		t.Fatalf("providers=%v err=%v", providers, err)
	}
	if attempts.Load() != 3 {
		t.Fatalf("attempts=%d, want 3", attempts.Load())
	}
}

func TestPublicCatalogueAllowlist_ExactPathsWithQueriesOnly(t *testing.T) {
	c := &Client{cfg: &config{public: true}}
	for path := range publicCataloguePaths {
		if err := c.checkPublicRequest(http.MethodGet, path+"?page=2&sort=name"); err != nil {
			t.Errorf("allowlisted path %s rejected: %v", path, err)
		}
	}
}

func TestPublicClient_ProviderJSONShape(t *testing.T) {
	b, err := json.Marshal(types.VPSProviderResponse{Name: "Tencent Cloud"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(b) != `{"name":"Tencent Cloud"}` {
		t.Fatalf("provider JSON = %s", b)
	}
}
