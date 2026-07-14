package types

// PackagesPricingResponse is the public-facing pricing surface returned by
// GET /api/v1/public/packages/plans. Plans is a list so new tiers can ship
// without breaking older clients. Mirrors RegistryPricingResponse.
type PackagesPricingResponse struct {
	Plans []PackagesPlanOption `json:"plans"`
}

// PackagesPlanOption is the sanitized projection of a billing plan for the
// Kumo Packages product — it intentionally omits base cost, margin, and the
// internal per-GB-hour rate so cost structure never leaks. PricePerUnit is a
// decimal string quoting the per-GB-month price to match the registry/ECR
// convention.
type PackagesPlanOption struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Unit          string `json:"unit"`
	PricePerUnit  string `json:"price_per_unit"`
	Currency      string `json:"currency"`
	ChargeModel   string `json:"charge_model"`
	BillingPeriod string `json:"billing_period"`
}
