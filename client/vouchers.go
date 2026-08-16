package client

import (
	"context"
	"fmt"

	"github.com/kumobase/kumo-go/types"
)

const vouchersPath = "/api/v1/vouchers"
const adminVoucherCampaignsPath = "/api/v1/admin/vouchers/campaigns"

// VouchersService provides voucher redemption history and automatic campaign
// discovery for the authenticated customer.
type VouchersService struct{ c *Client }

// Vouchers returns the customer voucher service.
func (c *Client) Vouchers() *VouchersService { return &VouchersService{c: c} }

// ListCampaigns returns currently relevant automatic voucher campaigns.
func (s *VouchersService) ListCampaigns(ctx context.Context) ([]types.PublicVoucherCampaignResponse, error) {
	var out []types.PublicVoucherCampaignResponse
	_, _, err := s.c.do(ctx, "GET", vouchersPath+"/campaigns", nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ListApplications returns the caller's automatic voucher benefit history.
func (s *VouchersService) ListApplications(ctx context.Context, opts ...ListOption) ([]types.VoucherApplicationResponse, *types.Meta, error) {
	var out []types.VoucherApplicationResponse
	meta, err := s.c.doList(ctx, "GET", withQuery(vouchersPath+"/applications", resolveListOpts(opts)), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, meta, nil
}

// AdminVouchersService manages automatic voucher campaigns. Existing code
// voucher CRUD remains on the server's legacy administrative surface.
type AdminVouchersService struct{ c *Client }

// AdminVouchers returns the protected voucher-administration service.
func (c *Client) AdminVouchers() *AdminVouchersService { return &AdminVouchersService{c: c} }

func adminVoucherCampaignPath(id uint, suffix string) string {
	path := fmt.Sprintf("%s/%d", adminVoucherCampaignsPath, id)
	if suffix != "" {
		path += "/" + suffix
	}
	return path
}

func (s *AdminVouchersService) CreateCampaign(ctx context.Context, req *types.CreateVoucherCampaignRequest, opts ...WriteOption) (*types.VoucherCampaignResponse, error) {
	return s.writeCampaign(ctx, "POST", adminVoucherCampaignsPath, req, opts...)
}

func (s *AdminVouchersService) ListCampaigns(ctx context.Context, opts ...ListOption) ([]types.VoucherCampaignResponse, *types.Meta, error) {
	var out []types.VoucherCampaignResponse
	meta, err := s.c.doList(ctx, "GET", withQuery(adminVoucherCampaignsPath, resolveListOpts(opts)), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, meta, nil
}

func (s *AdminVouchersService) GetCampaign(ctx context.Context, id uint) (*types.VoucherCampaignResponse, string, error) {
	var out types.VoucherCampaignResponse
	etag, _, err := s.c.do(ctx, "GET", adminVoucherCampaignPath(id, ""), nil, nil, &out)
	if err != nil {
		return nil, "", err
	}
	return &out, etag, nil
}

func (s *AdminVouchersService) UpdateCampaign(ctx context.Context, id uint, req *types.UpdateVoucherCampaignRequest, opts ...WriteOption) (*types.VoucherCampaignResponse, error) {
	return s.writeCampaign(ctx, "PATCH", adminVoucherCampaignPath(id, ""), req, opts...)
}

func (s *AdminVouchersService) ScheduleCampaign(ctx context.Context, id uint, req *types.VoucherCampaignActionRequest, opts ...WriteOption) (*types.VoucherCampaignResponse, error) {
	return s.writeCampaign(ctx, "POST", adminVoucherCampaignPath(id, "schedule"), req, opts...)
}

func (s *AdminVouchersService) PauseCampaign(ctx context.Context, id uint, req *types.VoucherCampaignActionRequest, opts ...WriteOption) (*types.VoucherCampaignResponse, error) {
	return s.writeCampaign(ctx, "POST", adminVoucherCampaignPath(id, "pause"), req, opts...)
}

func (s *AdminVouchersService) ResumeCampaign(ctx context.Context, id uint, req *types.VoucherCampaignActionRequest, opts ...WriteOption) (*types.VoucherCampaignResponse, error) {
	return s.writeCampaign(ctx, "POST", adminVoucherCampaignPath(id, "resume"), req, opts...)
}

func (s *AdminVouchersService) EndCampaign(ctx context.Context, id uint, req *types.VoucherCampaignActionRequest, opts ...WriteOption) (*types.VoucherCampaignResponse, error) {
	return s.writeCampaign(ctx, "POST", adminVoucherCampaignPath(id, "end"), req, opts...)
}

func (s *AdminVouchersService) ListCampaignApplications(ctx context.Context, id uint, opts ...ListOption) ([]types.VoucherApplicationResponse, *types.Meta, error) {
	var out []types.VoucherApplicationResponse
	meta, err := s.c.doList(ctx, "GET", withQuery(adminVoucherCampaignPath(id, "applications"), resolveListOpts(opts)), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, meta, nil
}

// The three application actions return the APPLICATION they mutated, not the
// campaign: it is the application's status, resolution_action and timestamps
// that change, and there is no other way to observe the result of your own
// write short of re-listing and scanning.

func (s *AdminVouchersService) ReverseApplication(ctx context.Context, campaignID, applicationID uint, req *types.VoucherApplicationActionRequest, opts ...WriteOption) (*types.VoucherApplicationResponse, error) {
	return s.writeApplicationAction(ctx, campaignID, applicationID, "reverse", req, opts...)
}

func (s *AdminVouchersService) RetryApplication(ctx context.Context, campaignID, applicationID uint, req *types.VoucherApplicationActionRequest, opts ...WriteOption) (*types.VoucherApplicationResponse, error) {
	return s.writeApplicationAction(ctx, campaignID, applicationID, "retry", req, opts...)
}

func (s *AdminVouchersService) WaiveApplication(ctx context.Context, campaignID, applicationID uint, req *types.VoucherApplicationActionRequest, opts ...WriteOption) (*types.VoucherApplicationResponse, error) {
	return s.writeApplicationAction(ctx, campaignID, applicationID, "waive", req, opts...)
}

// GetApplication fetches a single application, so a caller can poll one without
// paging the whole list.
func (s *AdminVouchersService) GetApplication(ctx context.Context, campaignID, applicationID uint) (*types.VoucherApplicationResponse, error) {
	path := fmt.Sprintf("%s/%d/applications/%d", adminVoucherCampaignsPath, campaignID, applicationID)
	var out types.VoucherApplicationResponse
	if _, _, err := s.c.do(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *AdminVouchersService) writeApplicationAction(ctx context.Context, campaignID, applicationID uint, action string, req *types.VoucherApplicationActionRequest, opts ...WriteOption) (*types.VoucherApplicationResponse, error) {
	path := fmt.Sprintf("%s/%d/applications/%d/%s", adminVoucherCampaignsPath, campaignID, applicationID, action)
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.VoucherApplicationResponse
	if _, _, err = s.c.do(ctx, "POST", path, req, &wopts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *AdminVouchersService) writeCampaign(ctx context.Context, method, path string, body any, opts ...WriteOption) (*types.VoucherCampaignResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.VoucherCampaignResponse
	_, _, err = s.c.do(ctx, method, path, body, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
