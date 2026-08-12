package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var publicCataloguePaths = map[string]struct{}{
	"/api/v1/apps/plans":     {},
	"/api/v1/vps/regions":    {},
	"/api/v1/vps/providers":  {},
	"/api/v1/vps/plans":      {},
	"/api/v1/volumes/plans":  {},
	"/api/v1/registry/plans": {},
	"/api/v1/packages/plans": {},
	"/api/v1/runners/plans":  {},
}

// NewPublic constructs an unauthenticated Client restricted to the public
// catalogue GET endpoints documented in this package. Non-authentication
// options such as WithHTTPClient, WithRetries, WithLogger, and WithUserAgent
// are supported. WithAPIKey and WithJWT are rejected.
//
// The returned *Client never sends an Authorization header and is safe for
// concurrent use by multiple goroutines.
func NewPublic(baseURL string, opts ...Option) (*Client, error) {
	cfg := defaultConfig()
	cfg.baseURL = strings.TrimRight(baseURL, "/")
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.auth != nil || cfg.authErr != nil {
		return nil, fmt.Errorf("%w: authentication options are not allowed", ErrPublicClientRestricted)
	}
	if cfg.baseURL == "" {
		return nil, fmt.Errorf("kumo: WithBaseURL or base URL argument is required")
	}
	cfg.public = true
	return &Client{cfg: cfg}, nil
}

func (c *Client) checkPublicRequest(method, requestPath string) error {
	if !c.cfg.public {
		return nil
	}
	if method != http.MethodGet {
		return fmt.Errorf("%w: %s %s", ErrPublicClientRestricted, method, requestPath)
	}
	u, err := url.ParseRequestURI(requestPath)
	if err != nil || u.IsAbs() || u.Host != "" || u.Fragment != "" {
		return fmt.Errorf("%w: GET %s", ErrPublicClientRestricted, requestPath)
	}
	if _, ok := publicCataloguePaths[u.EscapedPath()]; !ok {
		return fmt.Errorf("%w: GET %s", ErrPublicClientRestricted, requestPath)
	}
	return nil
}
