package codes

// Billing-module wire codes returned by /api/v1/billing/v2/* endpoints
// (read-only customer surface). Mutation codes (topup, voucher redeem) live
// elsewhere; those endpoints are sessionOnly and not callable with an API
// key.
const (
	BillingInvalidFilterCombination = "INVALID_FILTER_COMBINATION"
	BillingInvalidDateRange         = "INVALID_DATE_RANGE"
	BillingInvalidEnumValue         = "INVALID_ENUM_VALUE"
	BillingBreakdownFailed          = "BREAKDOWN_FAILED"
	// BillingInternalError is the catch-all 500 code for billing read endpoints
	// (charges list, charge filters, summary). Stable so consumers can branch on
	// it instead of the human-readable message.
	BillingInternalError = "BILLING_INTERNAL_ERROR"
)
