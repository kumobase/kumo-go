package types

import "testing"

func TestPackagesPricingRoundTrip(t *testing.T) {
	roundTrip(t, "PackagesPlanOption", PackagesPlanOption{
		ID:            72,
		Name:          "Pay-as-you-go",
		Description:   "metered storage, billed monthly",
		Unit:          "GB-month",
		PricePerUnit:  "1060",
		Currency:      "IDR",
		ChargeModel:   "postpaid",
		BillingPeriod: "monthly",
	})
	// Description omitted — exercises the omitempty field.
	roundTrip(t, "PackagesPlanOption+NoDescription", PackagesPlanOption{
		ID:            73,
		Name:          "Free tier",
		Unit:          "GB-month",
		PricePerUnit:  "0",
		Currency:      "IDR",
		ChargeModel:   "postpaid",
		BillingPeriod: "monthly",
	})
	roundTrip(t, "PackagesPricingResponse", PackagesPricingResponse{
		Plans: []PackagesPlanOption{
			{ID: 72, Name: "Pay-as-you-go", Unit: "GB-month", PricePerUnit: "1060", Currency: "IDR", ChargeModel: "postpaid", BillingPeriod: "monthly"},
		},
	})
	// Empty list — the no-plan / all-zero-price case must serialise as {"plans":[]}.
	roundTrip(t, "PackagesPricingResponse+Empty", PackagesPricingResponse{Plans: []PackagesPlanOption{}})
}
