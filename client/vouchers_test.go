package client_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kumobase/kumo-go/client"
	"github.com/kumobase/kumo-go/codes"
	"github.com/kumobase/kumo-go/types"
)

func campaignRequest() *types.CreateVoucherCampaignRequest {
	return &types.CreateVoucherCampaignRequest{
		Slug: "merdeka-17", Name: "Merdeka 17", Trigger: types.VoucherCampaignTriggerBillingCharge,
		BenefitType: types.VoucherCampaignBenefitPercentageDiscount, Audience: types.VoucherCampaignAudienceAll,
		Currency: "IDR", Value: "17.0000", Priority: 100, Timezone: "Asia/Jakarta",
		StartsAt: time.Date(2026, 8, 16, 17, 0, 0, 0, time.UTC),
		EndsAt:   time.Date(2026, 9, 1, 17, 0, 0, 0, time.UTC),
		Scopes:   []types.VoucherCampaignScope{{ProductType: types.VoucherCampaignProductApp}},
	}
}

func writeVoucherCampaignPage(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "ok", "data": data,
		"meta": map[string]any{"page": 2, "page_size": 25, "total_items": 51, "total_pages": 3},
	})
}

func TestVoucherCampaignsCustomerMethods(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/vouchers/campaigns":
			writeStruct(w, http.StatusOK, "", "ok", []types.PublicVoucherCampaignResponse{{ID: 1, Eligible: true}})
		case "/api/v1/vouchers/applications":
			if r.URL.Query().Get("page") != "2" {
				t.Errorf("page = %q", r.URL.Query().Get("page"))
			}
			writeVoucherCampaignPage(w, []types.VoucherApplicationResponse{{ID: 9, Status: types.VoucherApplicationApplied}})
		default:
			http.NotFound(w, r)
		}
	})
	campaigns, err := c.Vouchers().ListCampaigns(context.Background())
	if err != nil || len(campaigns) != 1 || !campaigns[0].Eligible {
		t.Fatalf("List: err=%v campaigns=%+v", err, campaigns)
	}
	applications, meta, err := c.Vouchers().ListApplications(context.Background(), client.WithPage(2))
	if err != nil || len(applications) != 1 || meta.TotalItems != 51 {
		t.Fatalf("ListApplications: err=%v applications=%+v meta=%+v", err, applications, meta)
	}
}

func TestAdminVoucherCampaignsCRUD(t *testing.T) {
	var createKey, patchMatch string
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/admin/vouchers/campaigns":
			createKey = r.Header.Get("Idempotency-Key")
			var req types.CreateVoucherCampaignRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Slug != "merdeka-17" {
				t.Errorf("create body: err=%v req=%+v", err, req)
			}
			writeStruct(w, http.StatusCreated, "", "created", &types.VoucherCampaignResponse{ID: 1, Slug: req.Slug})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/admin/vouchers/campaigns":
			writeVoucherCampaignPage(w, []types.VoucherCampaignResponse{{ID: 1}})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/admin/vouchers/campaigns/1":
			w.Header().Set("ETag", `W/"abc"`)
			writeStruct(w, http.StatusOK, "", "ok", &types.VoucherCampaignResponse{ID: 1})
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v1/admin/vouchers/campaigns/1":
			patchMatch = r.Header.Get("If-Match")
			writeStruct(w, http.StatusOK, "", "updated", &types.VoucherCampaignResponse{ID: 1, Name: "Updated"})
		default:
			http.NotFound(w, r)
		}
	})
	ctx := context.Background()
	created, err := c.AdminVouchers().CreateCampaign(ctx, campaignRequest(), client.WithIdempotencyKey("create-1"))
	if err != nil || created.ID != 1 || createKey != "create-1" {
		t.Fatalf("Create: err=%v created=%+v key=%q", err, created, createKey)
	}
	listed, meta, err := c.AdminVouchers().ListCampaigns(ctx, client.WithPage(2))
	if err != nil || len(listed) != 1 || meta.Page != 2 {
		t.Fatalf("List: err=%v listed=%+v meta=%+v", err, listed, meta)
	}
	_, etag, err := c.AdminVouchers().GetCampaign(ctx, 1)
	if err != nil || etag != `W/"abc"` {
		t.Fatalf("Get: err=%v etag=%q", err, etag)
	}
	name := "Updated"
	updated, err := c.AdminVouchers().UpdateCampaign(ctx, 1, &types.UpdateVoucherCampaignRequest{Name: &name}, client.IfMatch(etag))
	if err != nil || updated.Name != name || patchMatch != etag {
		t.Fatalf("Update: err=%v updated=%+v If-Match=%q", err, updated, patchMatch)
	}
}

func TestAdminVoucherCampaignsActions(t *testing.T) {
	seen := map[string]bool{}
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/api/v1/admin/vouchers/campaigns/1/applications" {
			writeVoucherCampaignPage(w, []types.VoucherApplicationResponse{{ID: 7}})
			return
		}
		if r.Method != http.MethodPost || !strings.HasPrefix(r.URL.Path, "/api/v1/admin/vouchers/campaigns/1/") {
			http.NotFound(w, r)
			return
		}
		suffix := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/vouchers/campaigns/1/")
		seen[suffix] = true
		if r.Header.Get("Idempotency-Key") == "" || r.Header.Get("If-Match") != `W/"abc"` {
			t.Errorf("%s headers idempotency=%q if-match=%q", suffix, r.Header.Get("Idempotency-Key"), r.Header.Get("If-Match"))
		}
		writeStruct(w, http.StatusOK, "", "ok", &types.VoucherCampaignResponse{ID: 1})
	})
	ctx := context.Background()
	opts := []client.WriteOption{client.IfMatch(`W/"abc"`)}
	action := &types.VoucherCampaignActionRequest{Reason: "operator request"}
	calls := map[string]func() (*types.VoucherCampaignResponse, error){
		"schedule": func() (*types.VoucherCampaignResponse, error) {
			return c.AdminVouchers().ScheduleCampaign(ctx, 1, action, opts...)
		},
		"pause": func() (*types.VoucherCampaignResponse, error) {
			return c.AdminVouchers().PauseCampaign(ctx, 1, action, opts...)
		},
		"resume": func() (*types.VoucherCampaignResponse, error) {
			return c.AdminVouchers().ResumeCampaign(ctx, 1, action, opts...)
		},
		"end": func() (*types.VoucherCampaignResponse, error) {
			return c.AdminVouchers().EndCampaign(ctx, 1, action, opts...)
		},
		"applications/7/reverse": func() (*types.VoucherCampaignResponse, error) {
			return c.AdminVouchers().ReverseApplication(ctx, 1, 7, &types.VoucherApplicationActionRequest{Reason: "refund"}, opts...)
		},
		"applications/7/retry": func() (*types.VoucherCampaignResponse, error) {
			return c.AdminVouchers().RetryApplication(ctx, 1, 7, &types.VoucherApplicationActionRequest{Reason: "retry"}, opts...)
		},
		"applications/7/waive": func() (*types.VoucherCampaignResponse, error) {
			return c.AdminVouchers().WaiveApplication(ctx, 1, 7, &types.VoucherApplicationActionRequest{Reason: "waive"}, opts...)
		},
	}
	for name, call := range calls {
		if _, err := call(); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if !seen[name] {
			t.Errorf("%s path not observed", name)
		}
	}
	applications, meta, err := c.AdminVouchers().ListCampaignApplications(ctx, 1, client.WithPage(2))
	if err != nil || len(applications) != 1 || meta.Page != 2 {
		t.Fatalf("ListApplications: err=%v applications=%+v meta=%+v", err, applications, meta)
	}
}

func TestVoucherCampaignErrorClassification(t *testing.T) {
	for _, code := range []string{codes.VoucherCampaignNotFound, codes.VoucherApplicationNotFound} {
		if !client.IsNotFound(&client.APIError{StatusCode: http.StatusNotFound, Code: code}) {
			t.Errorf("IsNotFound(%s) = false", code)
		}
	}
	for _, code := range []string{
		codes.VoucherCampaignAlreadyExists, codes.VoucherCampaignInvalidState, codes.VoucherCampaignRulesLocked,
		codes.VoucherCampaignPoolInsufficient, codes.VoucherCampaignReversalInsufficientBalance,
	} {
		if !client.IsConflict(&client.APIError{StatusCode: http.StatusConflict, Code: code}) {
			t.Errorf("IsConflict(%s) = false", code)
		}
	}
	opErr := &client.APIError{StatusCode: http.StatusBadRequest, Code: codes.VoucherCampaignInvalidRule}
	if !client.IsCode(errors.Join(errors.New("outer"), opErr), codes.VoucherCampaignInvalidRule) {
		t.Error("IsCode did not match wrapped campaign error")
	}
}
