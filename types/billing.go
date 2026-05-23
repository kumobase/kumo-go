package types

import "time"

// Billing DTOs cover the read-only customer surface (charges, summary,
// breakdown). Write paths (/billing/topup, /vouchers/redeem) are
// sessionOnly today, so their request shapes are not exposed here.
//
// All monetary amounts are decimal strings (e.g. "4.99") to avoid float
// rounding bugs. Convert with a decimal library when arithmetic is needed.

// GroupByMode selects how list charges are aggregated server-side. Pass
// either value as the group_by query parameter; the response shape changes
// accordingly (PublicGroupedChargeResponse instead of PublicChargeResponse).
type GroupByMode string

const (
	GroupByDate         GroupByMode = "date"
	GroupBySubscription GroupByMode = "subscription"
)

// BreakdownGranularity selects the time-bucket width for the breakdown
// endpoint.
type BreakdownGranularity string

const (
	BreakdownGranularityDaily   BreakdownGranularity = "daily"
	BreakdownGranularityMonthly BreakdownGranularity = "monthly"
)

// BreakdownGroupBy selects the grouping dimension for the breakdown.
type BreakdownGroupBy string

const (
	BreakdownGroupByProductType  BreakdownGroupBy = "product_type"
	BreakdownGroupByChargeModel  BreakdownGroupBy = "charge_model"
	BreakdownGroupBySubscription BreakdownGroupBy = "subscription"
	BreakdownGroupByNone         BreakdownGroupBy = "none"
)

// PublicPlanResponse is the customer-facing plan DTO returned by billing
// plan list/get endpoints. Internal cost structure (base_cost, margin) is
// intentionally absent.
type PublicPlanResponse struct {
	ID            uint      `json:"id"`
	ProductType   string    `json:"product_type"`
	ProductRefID  string    `json:"product_ref_id"`
	Name          string    `json:"name"`
	BillingPeriod string    `json:"billing_period"`
	ChargeModel   string    `json:"charge_model"`
	Currency      string    `json:"currency"`
	Price         string    `json:"price"` // decimal string
	Enabled       bool      `json:"enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PublicChargeItemResponse is one line on a charge. Amount and Quantity are
// decimal strings.
type PublicChargeItemResponse struct {
	ID          uint      `json:"id"`
	ChargeID    uint      `json:"charge_id"`
	ItemType    string    `json:"item_type"`
	Description string    `json:"description"`
	Amount      string    `json:"amount"`
	Currency    string    `json:"currency"`
	Quantity    string    `json:"quantity"`
	SourceType  string    `json:"source_type,omitempty"`
	SourceID    string    `json:"source_id,omitempty"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
}

// PublicChargeResponse is one charge in the user's billing history.
// Items is populated when the list endpoint is called with include=items.
type PublicChargeResponse struct {
	ID             uint                       `json:"id"`
	SubscriptionID uint                       `json:"subscription_id"`
	ProductType    string                     `json:"product_type"`
	PlanName       string                     `json:"plan_name"`
	Amount         string                     `json:"amount"`
	Currency       string                     `json:"currency"`
	PeriodStart    time.Time                  `json:"period_start"`
	PeriodEnd      time.Time                  `json:"period_end"`
	ChargeType     string                     `json:"charge_type"`
	Status         string                     `json:"status"`
	FailureReason  *string                    `json:"failure_reason,omitempty"`
	ReferenceID    string                     `json:"reference_id"`
	CreatedAt      time.Time                  `json:"created_at"`
	Items          []PublicChargeItemResponse `json:"items,omitempty"`
}

// PublicGroupedChargeResponse is the summarised charge view returned when
// group_by=date or group_by=subscription. The list endpoint returns these
// instead of PublicChargeResponse in that mode.
type PublicGroupedChargeResponse struct {
	GroupKey        string           `json:"group_key"`
	TotalAmount     string           `json:"total_amount"`
	Currency        string           `json:"currency"`
	ChargeCount     int64            `json:"charge_count"`
	Date            *string          `json:"date,omitempty"`            // group_by=date
	SubscriptionID  *uint            `json:"subscription_id,omitempty"` // group_by=subscription
	ProductType     *string          `json:"product_type,omitempty"`
	PlanName        *string          `json:"plan_name,omitempty"`
	StatusBreakdown map[string]int64 `json:"status_breakdown"`
}

// ProductBreakdown is per-product spending totals.
type ProductBreakdown struct {
	VPS     string `json:"vps"`
	App     string `json:"app"`
	Storage string `json:"storage"`
}

// PeriodSummary is spend totals over a time period.
type PeriodSummary struct {
	Start        time.Time        `json:"start"`
	End          time.Time        `json:"end"`
	TotalCharged string           `json:"total_charged"`
	ByProduct    ProductBreakdown `json:"by_product"`
}

// BillingSummaryResponse is the spending overview returned by
// GET /api/v1/billing/v2/summary. Only prepaid charges are included;
// postpaid (pay-as-you-go) charges live in the breakdown endpoint.
type BillingSummaryResponse struct {
	Currency            string        `json:"currency"`
	CurrentPeriod       PeriodSummary `json:"current_period"`
	PreviousPeriodTotal string        `json:"previous_period_total"`
}

// BreakdownGroup is one slice of the pie within a bucket (or in totals).
type BreakdownGroup struct {
	Key         string  `json:"key"`
	Amount      string  `json:"amount"`
	Quantity    string  `json:"quantity,omitempty"`
	Unit        string  `json:"unit,omitempty"`
	ProductType *string `json:"product_type,omitempty"`
	PlanName    *string `json:"plan_name,omitempty"`
}

// BreakdownBucket is one time slice of the breakdown response.
type BreakdownBucket struct {
	PeriodStart time.Time        `json:"period_start"`
	PeriodEnd   time.Time        `json:"period_end"`
	Amount      string           `json:"amount"`
	Groups      []BreakdownGroup `json:"groups"`
}

// BreakdownTotals aggregates groups across the entire queried range.
type BreakdownTotals struct {
	Amount string           `json:"amount"`
	Groups []BreakdownGroup `json:"groups"`
}

// BillingBreakdownResponse is returned by GET /api/v1/billing/v2/breakdown
// (AWS Cost Explorer-style timeline). From and To are RFC3339 date strings
// at midnight Asia/Jakarta.
type BillingBreakdownResponse struct {
	Currency    string               `json:"currency"`
	Granularity BreakdownGranularity `json:"granularity"`
	GroupBy     BreakdownGroupBy     `json:"group_by"`
	From        string               `json:"from"`
	To          string               `json:"to"`
	Totals      BreakdownTotals      `json:"totals"`
	Buckets     []BreakdownBucket    `json:"buckets"`
}
