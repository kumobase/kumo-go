package types

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestProductBreakdownRoundTrip pins the by_product wire shape, including the
// container_registry and database fields added so no charged product type is
// silently dropped from the summary.
func TestProductBreakdownRoundTrip(t *testing.T) {
	pb := ProductBreakdown{
		VPS:               "44750",
		App:               "2971.6227",
		Storage:           "9.2167",
		ContainerRegistry: "44.2241",
		Database:          "0",
		Jobs:              "12.5",
		VMRunners:         "37.5",
		Packages:          "5.75",
	}
	roundTrip(t, "ProductBreakdown", pb)

	b, err := json.Marshal(pb)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, key := range []string{"vps", "app", "storage", "container_registry", "database", "jobs", "vm_runners", "packages"} {
		if !strings.Contains(string(b), `"`+key+`"`) {
			t.Errorf("ProductBreakdown JSON missing %q key: %s", key, b)
		}
	}
}

// TestChargeFiltersResponseRoundTrip pins the split-out filter-vocabulary DTO.
func TestChargeFiltersResponseRoundTrip(t *testing.T) {
	roundTrip(t, "ChargeFiltersResponse", ChargeFiltersResponse{
		AvailableProductTypes: []string{"vps", "app", "storage", "container_registry"},
		AvailableStatuses:     []string{"charged", "pending", "failed", "refunded"},
	})
}

// TestBillingSummaryResponseRoundTrip exercises the full summary envelope.
func TestBillingSummaryResponseRoundTrip(t *testing.T) {
	roundTrip(t, "BillingSummaryResponse", BillingSummaryResponse{
		Currency:            "IDR",
		PreviousPeriodTotal: "100",
		CurrentPeriod: PeriodSummary{
			TotalCharged: "47730.8394",
			ByProduct: ProductBreakdown{
				VPS: "44750", App: "2971.6227", Storage: "9.2167",
				ContainerRegistry: "0", Database: "0",
			},
			// container_registry is monthly-postpaid: 0 charged so far, but
			// accruing surfaces the in-flight usage estimate for the period.
			Accruing: ProductBreakdown{
				VPS: "0", App: "1.2500", Storage: "0.4000",
				ContainerRegistry: "44.2241", Database: "0",
			},
			AccruingTotal: "45.8741",
		},
	})
}
