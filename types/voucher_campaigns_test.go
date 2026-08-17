package types

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func campaignString(value string) *string { return &value }
func campaignInt(value int) *int          { return &value }

func TestVoucherCampaignRequestsRoundTrip(t *testing.T) {
	start := time.Date(2026, 8, 16, 17, 0, 0, 0, time.UTC)
	end := time.Date(2026, 9, 1, 17, 0, 0, 0, time.UTC)
	planID := uint(17)
	req := CreateVoucherCampaignRequest{
		Slug: "merdeka-17", Name: "Merdeka 17", Trigger: VoucherCampaignTriggerBillingCharge,
		BenefitType: VoucherCampaignBenefitPercentageDiscount, Audience: VoucherCampaignAudienceAll,
		Currency: "IDR", Value: "17.0000", Budget: campaignString("1000000.0000"),
		Priority: 100, StartsAt: start, EndsAt: end,
		Scopes: []VoucherCampaignScope{
			{ProductType: VoucherCampaignProductApp},
			{ProductType: VoucherCampaignProductJobs, PlanID: &planID},
			{ProductType: VoucherCampaignProductStorage},
			{ProductType: VoucherCampaignProductContainerRegistry},
			{ProductType: VoucherCampaignProductVPS},
			{ProductType: VoucherCampaignProductDatabase},
			{ProductType: VoucherCampaignProductVMRunners},
			{ProductType: VoucherCampaignProductPackages},
		},
	}
	roundTrip(t, "CreateVoucherCampaignRequest", req)
	roundTrip(t, "UpdateVoucherCampaignRequest", UpdateVoucherCampaignRequest{
		Value: campaignString("20.0000"), PerUserLimit: campaignInt(1), ClearBudget: true,
	})
	roundTrip(t, "VoucherApplicationActionRequest", VoucherApplicationActionRequest{Reason: "operator action"})
}

func TestVoucherCampaignProductTypesAreComplete(t *testing.T) {
	all := AllVoucherCampaignProductTypes()
	if len(all) != 8 {
		t.Fatalf("expected 8 campaign product types, got %d: %v", len(all), all)
	}
	seen := map[VoucherCampaignProductType]bool{}
	for _, product := range all {
		if seen[product] {
			t.Errorf("duplicate product type %q", product)
		}
		seen[product] = true
		if !product.Valid() {
			t.Errorf("product type %q returned from AllVoucherCampaignProductTypes but is not Valid", product)
		}
	}
	// The string values are the wire contract and are compared against the
	// server's billing product types as plain strings. Changing one silently
	// stops a campaign from ever matching that product.
	for _, want := range []VoucherCampaignProductType{
		"vps", "app", "database", "storage", "container_registry", "jobs", "vm_runners", "packages",
	} {
		if !seen[want] {
			t.Errorf("missing product type %q", want)
		}
	}
	if VoucherCampaignProductType("bogus").Valid() {
		t.Error("unknown product type reported as valid")
	}
	if VoucherCampaignProductType("").Valid() {
		t.Error("empty product type reported as valid")
	}
}

func TestVoucherCampaignResponsesRoundTrip(t *testing.T) {
	now := time.Date(2026, 8, 17, 0, 0, 0, 0, time.UTC)
	ledgerID := "3f891951-0785-4cc8-a054-4a8f013430ab"
	roundTrip(t, "VoucherCampaignResponse", VoucherCampaignResponse{
		ID: 1, Slug: "first-topup", Name: "First Top-Up 2x", Trigger: VoucherCampaignTriggerTopupSuccess,
		BenefitType: VoucherCampaignBenefitPercentageBonus, Audience: VoucherCampaignAudienceNeverSuccessfulTopup,
		Currency: "IDR", Value: "100.0000", PerEventCap: campaignString("81000.0000"),
		PerUserLimit: campaignInt(1), Priority: 100, StartsAt: now,
		EndsAt: now.Add(24 * time.Hour), Status: VoucherCampaignStatusActive,
		ReservedAmount: "81000.0000", AppliedAmount: "81000.0000", AvailableAmount: campaignString("838000.0000"),
		CreatedAt: now, UpdatedAt: now,
	})
	roundTrip(t, "PublicVoucherCampaignResponse", PublicVoucherCampaignResponse{
		ID: 1, Slug: "first-topup", Name: "First Top-Up 2x", Trigger: VoucherCampaignTriggerTopupSuccess,
		BenefitType: VoucherCampaignBenefitPercentageBonus, Currency: "IDR", Value: "100.0000",
		StartsAt: now, EndsAt: now.Add(24 * time.Hour),
		Eligible: true, Availability: VoucherCampaignAvailabilityAvailable,
	})
	roundTrip(t, "VoucherApplicationResponse", VoucherApplicationResponse{
		ID: 2, CampaignID: 1, CampaignSlug: "first-topup", CampaignName: "First Top-Up 2x",
		Status: VoucherApplicationApplied, SourceType: "topup", SourceReference: "topup-123",
		GrossAmount: "100000.0000", BenefitAmount: "81000.0000", NetAmount: "181000.0000",
		Currency: "IDR", LedgerTransactionID: &ledgerID, AppliedAt: &now, CreatedAt: now, UpdatedAt: now,
	})
}

func TestVoucherCampaignMonetaryFieldsRemainJSONStrings(t *testing.T) {
	body, err := json.Marshal(CreateVoucherCampaignRequest{
		Value: "17.0000", PerEventCap: campaignString("81000.0000"), Budget: campaignString("1000000.0000"),
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{`"value":"17.0000"`, `"per_event_cap":"81000.0000"`, `"budget":"1000000.0000"`} {
		if !strings.Contains(string(body), fragment) {
			t.Errorf("JSON missing %s: %s", fragment, body)
		}
	}
}

func TestVoucherCampaignEnumWireValues(t *testing.T) {
	cases := map[string]string{
		string(VoucherCampaignTriggerTopupSuccess):          "topup_success",
		string(VoucherCampaignBenefitPercentageBonus):       "percentage_bonus",
		string(VoucherCampaignAudienceNeverSuccessfulTopup): "never_successful_topup",
		string(VoucherCampaignStatusPaused):                 "paused",
		string(VoucherApplicationManualReview):              "manual_review",
		string(VoucherApplicationNotAwarded):                "not_awarded",
		string(VoucherApplicationWaived):                    "waived",
		string(VoucherCampaignAvailabilityPoolInsufficient): "pool_insufficient",
		string(VoucherApplicationResolutionRetry):           "retry",
	}
	for got, want := range cases {
		if got != want {
			t.Errorf("enum = %q, want %q", got, want)
		}
	}
}
