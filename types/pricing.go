package types

// TemplateWithPricing is the customer-facing app plan DTO returned by
// GET /api/v1/apps/plans — one purchasable app resource template (CPU/memory
// envelope) with its hourly/daily/monthly price.
//
// PriceHour / PriceDay / PriceMonth are decimal strings (e.g. "1500.00", all
// prices in IDR) — parse with strconv.ParseFloat or a decimal library if you
// need arithmetic. CPURequestvCPU / CPULimitvCPU are likewise decimal strings
// ("0.25", "1").
type TemplateWithPricing struct {
	Slug            string        `json:"slug"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	CPURequestvCPU  string        `json:"cpu_request_vcpu"`
	CPULimitvCPU    string        `json:"cpu_limit_vcpu"`
	MemoryRequestMB int           `json:"memory_request_mb"`
	MemoryLimitMB   int           `json:"memory_limit_mb"`
	PriceHour       string        `json:"price_hour"`
	PriceDay        string        `json:"price_day"`
	PriceMonth      string        `json:"price_month"`
	Availability    *Availability `json:"availability,omitempty"`
}

// PricingResponse wraps the app plan catalogue returned by
// GET /api/v1/apps/plans. Templates is a list so new tiers can ship without
// breaking older clients. The SDK's AppsService.ListPlans flattens this to
// the inner slice for ergonomics.
type PricingResponse struct {
	Templates []TemplateWithPricing `json:"templates"`
}
