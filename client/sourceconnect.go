package client

import (
	"context"
	"fmt"

	"github.com/kumobase/kumo-go/types"
)

// SourceConnectionsService backs /api/v1/source-connections/*. A source
// connection links a customer's git account (via the Kumo Build GitHub App)
// so repositories can later feed automated builds.
//
// Connecting an account is an interactive browser flow — the App install and
// the OAuth/setup callback happen on the provider and are deliberately not
// part of this SDK. These methods cover the machine-usable surface: listing
// connections, listing the repositories a connection can access (the repo
// picker), and disconnecting.
type SourceConnectionsService struct {
	c *Client
}

// SourceConnections returns the source-connections service.
func (c *Client) SourceConnections() *SourceConnectionsService {
	return &SourceConnectionsService{c: c}
}

// List returns every source connection owned by the authenticated user. Not
// paginated server-side (low cardinality per user); returns the full list.
func (s *SourceConnectionsService) List(ctx context.Context) ([]types.SourceConnectionResponse, error) {
	var out []types.SourceConnectionResponse
	_, err := s.c.doList(ctx, "GET", "/api/v1/source-connections", &out)
	return out, err
}

// ListRepos returns the repositories the given connection has been granted
// access to, fetched live from the provider so the result always reflects the
// account owner's current grant. Powers the repo picker.
//
// Returns codes.SourceConnectionNotFound when the connection doesn't exist or
// isn't owned by the caller, codes.SourceConnectionSuspended when the install
// is suspended, or codes.SourceProviderError when the provider call fails.
func (s *SourceConnectionsService) ListRepos(ctx context.Context, id uint) ([]types.SourceRepoResponse, error) {
	var out []types.SourceRepoResponse
	_, err := s.c.doList(ctx, "GET", fmt.Sprintf("/api/v1/source-connections/%d/repos", id), &out)
	return out, err
}

// Disconnect removes a source connection and uninstalls the Kumo Build app
// from the provider account, fully revoking access. Returns the disconnected
// connection snapshot.
func (s *SourceConnectionsService) Disconnect(ctx context.Context, id uint, opts ...WriteOption) (*types.SourceConnectionResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.SourceConnectionResponse
	_, _, err = s.c.do(ctx, "DELETE", fmt.Sprintf("/api/v1/source-connections/%d", id), nil, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
