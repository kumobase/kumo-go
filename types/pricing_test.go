package types

import "testing"

func TestPricingRoundTrip(t *testing.T) {
	roundTrip(t, "TemplateWithPricing", TemplateWithPricing{
		Slug:            "kumo.nano",
		Name:            "Nano",
		Description:     "smallest app tier",
		CPURequestvCPU:  "0.25",
		CPULimitvCPU:    "1",
		MemoryRequestMB: 256,
		MemoryLimitMB:   512,
		PriceHour:       "1.50",
		PriceDay:        "36.00",
		PriceMonth:      "1080.00",
	})
	// With Availability populated — exercises the omitempty pointer field.
	roundTrip(t, "TemplateWithPricing+Availability", TemplateWithPricing{
		Slug:         "kumo.large",
		Name:         "Large",
		PriceMonth:   "9000.00",
		Availability: &Availability{Available: false, Reason: AvailabilityReasonCPUFull},
	})
	roundTrip(t, "PricingResponse", PricingResponse{
		Templates: []TemplateWithPricing{
			{Slug: "kumo.nano", Name: "Nano", PriceMonth: "1080.00"},
			{Slug: "kumo.small", Name: "Small", PriceMonth: "2160.00"},
		},
	})
}
