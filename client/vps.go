package client

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/kumobase/kumo-go/types"
)

// VPSService backs /api/v1/vps/*. Reads (regions, plans, servers) are
// synchronous. Mutations (rent, reboot, reinstall, power on/off, renew,
// cancel, reset password) are async via the server's ActionStatus field —
// each *AndWait helper composes the action with PollResource until the
// ActionStatus flips back to "".
type VPSService struct {
	c *Client
}

// VPS returns the VPS service.
func (c *Client) VPS() *VPSService { return &VPSService{c: c} }

// ── Reads ──────────────────────────────────────────────────────────

// ListRegions returns the regions Kumo currently supports for VPS rental.
func (s *VPSService) ListRegions(ctx context.Context) ([]types.VPSRegionResponse, error) {
	var out []types.VPSRegionResponse
	_, _, err := s.c.do(ctx, "GET", "/api/v1/vps/regions", nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ListPlans returns the public VPS plan catalogue. Filter by region with
// WithExtraQuery("region", "sgp"). Server returns 400 MISSING_REGION when
// the catalogue requires a region selector and none is supplied.
func (s *VPSService) ListPlans(ctx context.Context, opts ...ListOption) ([]types.PublicVPSPlanResponse, error) {
	q := resolveListOpts(opts)
	var out []types.PublicVPSPlanResponse
	_, _, err := s.c.do(ctx, "GET", withQuery("/api/v1/vps/plans", q), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ListServers returns the user's rented VPS instances, paginated. Filters
// via WithExtraQuery: "status", "region", "provider_name", "expires_before"
// (RFC3339 timestamp).
func (s *VPSService) ListServers(ctx context.Context, opts ...ListOption) ([]types.VPSServerResponse, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.VPSServerResponse
	meta, err := s.c.doList(ctx, "GET", withQuery("/api/v1/vps/servers", q), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, meta, nil
}

// GetServer fetches a single VPS instance. Use the returned ActionStatus
// to determine whether an async action is still in flight ("" = idle).
func (s *VPSService) GetServer(ctx context.Context, id uint) (*types.VPSServerResponse, error) {
	var out types.VPSServerResponse
	_, _, err := s.c.do(ctx, "GET", fmt.Sprintf("/api/v1/vps/servers/%d", id), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetServerByName fetches a single VPS instance by its display_name. The
// server resolves a non-numeric path segment as a name; servers without a
// user-supplied display_name remain reachable only by id. Returns 404 if
// no server in the caller's scope matches the name.
func (s *VPSService) GetServerByName(ctx context.Context, name string) (*types.VPSServerResponse, error) {
	var out types.VPSServerResponse
	_, _, err := s.c.do(ctx, "GET", "/api/v1/vps/servers/"+url.PathEscape(name), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ── Lifecycle (writes) ─────────────────────────────────────────────

// RentServer provisions a new VPS instance. Honors Idempotency-Key — the
// server caches the response so duplicate creates are impossible (a real
// concern given billing involvement). Returns the freshly created server
// in its initial "provisioning" state; poll until Status="running".
func (s *VPSService) RentServer(ctx context.Context, req *types.RentServerRequest, opts ...WriteOption) (*types.VPSServerResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.VPSServerResponse
	_, _, err = s.c.do(ctx, "POST", "/api/v1/vps/servers", req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateServerName renames a VPS instance (display label only — the
// provider-side identifier is immutable).
func (s *VPSService) UpdateServerName(ctx context.Context, id uint, name string, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "PATCH", fmt.Sprintf("/api/v1/vps/servers/%d/name", id),
		&types.UpdateServerNameRequest{Name: name}, &wopts, nil)
	return err
}

// CancelSubscription disables auto-renew on a server (it will expire at
// ExpiresAt and not be billed for the next period). Returns 409
// AUTO_RENEW_ALREADY_CANCELLED if already cancelled.
func (s *VPSService) CancelSubscription(ctx context.Context, id uint, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/vps/servers/%d/cancel", id), nil, &wopts, nil)
	return err
}

// Renew bills the next period immediately (manual top-up of subscription
// life). Useful when the user has disabled auto-renew but wants to extend.
func (s *VPSService) Renew(ctx context.Context, id uint, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/vps/servers/%d/renew", id), nil, &wopts, nil)
	return err
}

// RevealInitialPassword returns the auto-generated root password set at
// provisioning time. Server logs the call for audit. Only meaningful on
// fresh instances — once the user changes the password via SSH, this
// returns the stale initial value.
func (s *VPSService) RevealInitialPassword(ctx context.Context, id uint, opts ...WriteOption) (string, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return "", err
	}
	var out struct {
		Password string `json:"password"`
	}
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/vps/servers/%d/password", id), nil, &wopts, &out)
	if err != nil {
		return "", err
	}
	return out.Password, nil
}

// ── Async actions ──────────────────────────────────────────────────

// Reboot triggers a provider-level reboot. ActionStatus transitions to
// "rebooting" and back to "" once the provider confirms completion.
// Use RebootAndWait to block.
func (s *VPSService) Reboot(ctx context.Context, id uint, opts ...WriteOption) error {
	return s.action(ctx, id, "reboot", opts...)
}

// Reinstall asks the provider to wipe and re-image the instance. Same
// async semantics as Reboot but takes much longer (minutes).
func (s *VPSService) Reinstall(ctx context.Context, id uint, opts ...WriteOption) error {
	return s.action(ctx, id, "reinstall", opts...)
}

// PowerOn boots a stopped instance.
func (s *VPSService) PowerOn(ctx context.Context, id uint, opts ...WriteOption) error {
	return s.action(ctx, id, "poweron", opts...)
}

// PowerOff cleanly halts the instance (not destructive — the disk is
// preserved and billing continues per the provider's policy).
func (s *VPSService) PowerOff(ctx context.Context, id uint, opts ...WriteOption) error {
	return s.action(ctx, id, "poweroff", opts...)
}

// ResetPassword regenerates the root password via the provider and
// returns the new value. Async — the new password is only valid after
// ActionStatus returns to "".
func (s *VPSService) ResetPassword(ctx context.Context, id uint, opts ...WriteOption) (string, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return "", err
	}
	var out struct {
		Password string `json:"password"`
	}
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/vps/servers/%d/reset-password", id), nil, &wopts, &out)
	if err != nil {
		return "", err
	}
	return out.Password, nil
}

// action centralises the "POST /vps/servers/:id/<verb>" pattern shared by
// reboot / reinstall / poweron / poweroff. Honors Idempotency-Key (the
// server rejects concurrent actions with 409 ACTION_IN_PROGRESS).
func (s *VPSService) action(ctx context.Context, id uint, verb string, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/vps/servers/%d/%s", id, verb), nil, &wopts, nil)
	return err
}

// ── Polling helpers ────────────────────────────────────────────────

// WaitForActionComplete blocks until the instance's ActionStatus returns
// to "" (idle), or until the operation surfaces a terminal error via the
// ActionError field. Use after any async action when the caller wants to
// confirm completion before moving on.
func (s *VPSService) WaitForActionComplete(ctx context.Context, id uint, opts ...PollOption) (*types.VPSServerResponse, error) {
	return PollResource(ctx,
		func(ctx context.Context) (*types.VPSServerResponse, error) {
			return s.GetServer(ctx, id)
		},
		func(v *types.VPSServerResponse) (bool, error) {
			if v == nil {
				return false, nil
			}
			if v.ActionError != "" {
				return true, errors.New("kumo: vps action failed: " + v.ActionError)
			}
			return v.ActionStatus == "", nil
		},
		opts...,
	)
}

// RebootAndWait composes Reboot + WaitForActionComplete.
func (s *VPSService) RebootAndWait(ctx context.Context, id uint, opts ...PollOption) (*types.VPSServerResponse, error) {
	if err := s.Reboot(ctx, id); err != nil {
		return nil, err
	}
	return s.WaitForActionComplete(ctx, id, opts...)
}

// ReinstallAndWait composes Reinstall + WaitForActionComplete. Note this
// can take 5-10 minutes on some providers — pass WithPollMaxWait if you
// want a tighter SLA than the default 10 minutes.
func (s *VPSService) ReinstallAndWait(ctx context.Context, id uint, opts ...PollOption) (*types.VPSServerResponse, error) {
	if err := s.Reinstall(ctx, id); err != nil {
		return nil, err
	}
	return s.WaitForActionComplete(ctx, id, opts...)
}

// PowerOnAndWait composes PowerOn + WaitForActionComplete.
func (s *VPSService) PowerOnAndWait(ctx context.Context, id uint, opts ...PollOption) (*types.VPSServerResponse, error) {
	if err := s.PowerOn(ctx, id); err != nil {
		return nil, err
	}
	return s.WaitForActionComplete(ctx, id, opts...)
}

// PowerOffAndWait composes PowerOff + WaitForActionComplete.
func (s *VPSService) PowerOffAndWait(ctx context.Context, id uint, opts ...PollOption) (*types.VPSServerResponse, error) {
	if err := s.PowerOff(ctx, id); err != nil {
		return nil, err
	}
	return s.WaitForActionComplete(ctx, id, opts...)
}
