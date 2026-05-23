package client

import (
	"context"
	"fmt"

	"github.com/kumobase/kumo-go/types"
)

// APIKeysService backs /api/v1/api-keys/*. Every endpoint is sessionOnly
// on the server (a leaked API key must not be able to mint replacements),
// so calls made with a Client configured via WithAPIKey will return 403
// API_KEY_SESSION_REQUIRED.
//
// The CLI's `kumo auth keys create` flow is the canonical caller: log in
// with email+password → exchange for a JWT → construct the Client with
// WithJWT → issue an API key for subsequent resource work.
type APIKeysService struct {
	c *Client
}

// APIKeys returns the API-keys service.
func (c *Client) APIKeys() *APIKeysService { return &APIKeysService{c: c} }

// List returns all of the user's API keys (metadata only; the raw key
// values are only shown once at creation time via Create).
func (s *APIKeysService) List(ctx context.Context) ([]types.APIKeyResponse, error) {
	var out []types.APIKeyResponse
	_, err := s.c.doList(ctx, "GET", "/api/v1/api-keys", &out)
	return out, err
}

// Create issues a new API key. The full key value (kumo_sk_…) is
// returned in APIKeyCreateResponse.Key exactly ONCE — store it
// immediately, it cannot be retrieved later. KeyPrefix is preserved on
// subsequent List/Get for identification.
//
// req.RegistryScope makes the key a registry-only credential (Harbor-style
// robot account). Registry keys are rejected on /api/v1/* and can only
// authenticate against the OCI /v2/token endpoint.
func (s *APIKeysService) Create(ctx context.Context, req *types.CreateAPIKeyRequest, opts ...WriteOption) (*types.APIKeyCreateResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.APIKeyCreateResponse
	_, _, err = s.c.do(ctx, "POST", "/api/v1/api-keys", req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Update rotates the display name and/or expiry of an existing key. Scopes
// and registry binding are immutable post-create — replace via delete + create.
func (s *APIKeysService) Update(ctx context.Context, id uint, req *types.UpdateAPIKeyRequest, opts ...WriteOption) (*types.APIKeyResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.APIKeyResponse
	_, _, err = s.c.do(ctx, "PATCH", fmt.Sprintf("/api/v1/api-keys/%d", id), req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete revokes an API key immediately. Outstanding requests in flight
// with the revoked key complete; new requests get 401.
func (s *APIKeysService) Delete(ctx context.Context, id uint, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "DELETE", fmt.Sprintf("/api/v1/api-keys/%d", id), nil, &wopts, nil)
	return err
}
