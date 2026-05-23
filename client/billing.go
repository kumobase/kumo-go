package client

import (
	"context"

	"github.com/kumobase/kumo-go/types"
)

// BillingService backs the read-only customer surface of /api/v1/billing/*.
// Write paths (topup, voucher redeem) are sessionOnly on the server and
// not exposed in the SDK — call them via the dashboard or a CLI JWT
// session.
type BillingService struct {
	c *Client
}

// Billing returns the billing service.
func (c *Client) Billing() *BillingService { return &BillingService{c: c} }

// ListCharges returns the user's charge history, paginated. When
// WithExtraQuery("group_by", "date") or "subscription" is supplied the
// response shape changes to PublicGroupedChargeResponse — call
// ListGroupedCharges instead for the typed variant.
func (s *BillingService) ListCharges(ctx context.Context, opts ...ListOption) ([]types.PublicChargeResponse, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.PublicChargeResponse
	meta, err := s.c.doList(ctx, "GET", withQuery("/api/v1/billing/v2/charges", q), &out)
	return out, meta, err
}

// ListGroupedCharges returns the summarised view returned by the same
// endpoint when group_by=date or group_by=subscription is set. Provide
// the group_by via opts (WithExtraQuery("group_by", "date")).
func (s *BillingService) ListGroupedCharges(ctx context.Context, opts ...ListOption) ([]types.PublicGroupedChargeResponse, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.PublicGroupedChargeResponse
	meta, err := s.c.doList(ctx, "GET", withQuery("/api/v1/billing/v2/charges", q), &out)
	return out, meta, err
}

// GetSummary returns a spending overview (current period + previous
// period total). Prepaid only — postpaid pay-as-you-go charges live in
// GetBreakdown.
func (s *BillingService) GetSummary(ctx context.Context) (*types.BillingSummaryResponse, error) {
	var out types.BillingSummaryResponse
	_, _, err := s.c.do(ctx, "GET", "/api/v1/billing/v2/summary", nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetBreakdown returns an AWS-Cost-Explorer-style timeline of charges.
// Required query params via opts: WithExtraQuery("from", "2026-05-01"),
// WithExtraQuery("to", "2026-05-23"), WithExtraQuery("granularity", "daily"),
// WithExtraQuery("group_by", "product_type"). Server returns 400 with
// codes.BillingInvalidDateRange / BillingInvalidEnumValue on bad input.
func (s *BillingService) GetBreakdown(ctx context.Context, opts ...ListOption) (*types.BillingBreakdownResponse, error) {
	q := resolveListOpts(opts)
	var out types.BillingBreakdownResponse
	_, _, err := s.c.do(ctx, "GET", withQuery("/api/v1/billing/v2/breakdown", q), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
