package client

import (
	"context"

	"github.com/kumobase/kumo-go/types"
)

// ProfileService backs the read-only side of /api/v1/profile/*. Profile
// mutations (name, email, password) are sessionOnly on the server and
// not exposed in this SDK — call them via the dashboard or a JWT session.
type ProfileService struct {
	c *Client
}

// Profile returns the profile service.
func (c *Client) Profile() *ProfileService { return &ProfileService{c: c} }

// Get returns the authenticated user's profile.
func (s *ProfileService) Get(ctx context.Context) (*types.GetProfileResponse, error) {
	var out types.GetProfileResponse
	_, _, err := s.c.do(ctx, "GET", "/api/v1/profile", nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetBalance returns the authenticated user's current balance in IDR.
// Balance is a decimal string ("1000.50") — parse with a decimal library
// if arithmetic is needed.
func (s *ProfileService) GetBalance(ctx context.Context) (*types.GetBalanceResponse, error) {
	var out types.GetBalanceResponse
	_, _, err := s.c.do(ctx, "GET", "/api/v1/profile/balance", nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// HasPassword reports whether the user has a password set (vs only Google
// OAuth login). Lets dashboards decide between "set password" and
// "change password" CTAs without leaking the rest of the profile.
func (s *ProfileService) HasPassword(ctx context.Context) (bool, error) {
	var out types.HasPasswordResponse
	_, _, err := s.c.do(ctx, "GET", "/api/v1/profile/has-password", nil, nil, &out)
	if err != nil {
		return false, err
	}
	return out.HasPassword, nil
}
