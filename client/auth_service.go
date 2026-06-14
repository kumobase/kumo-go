package client

import (
	"context"

	"github.com/kumobase/kumo-go/types"
)

// AuthService backs the session-lifecycle endpoints under /api/v1/auth that a
// non-browser client can use. The dashboard drives login/refresh/logout over
// HttpOnly cookies; a CLI/SDK holding a session JWT uses Refresh to renew a
// short-lived access token from a stored refresh token without a full
// re-login, and Logout/LogoutAll to revoke server-side sessions.
//
// Note: automatic, single-flight "refresh on 401" belongs in the calling
// application (typically a browser client that owns a cookie jar), not in this
// stateless SDK — the Go Client holds a single static credential set via
// WithJWT/WithAPIKey and does not mutate it on your behalf. Call Refresh
// explicitly and reconstruct the Client with the new token.
type AuthService struct {
	c *Client
}

// Auth returns the auth/session service.
func (c *Client) Auth() *AuthService { return &AuthService{c: c} }

// Refresh exchanges a refresh token for a fresh access token, rotating the
// refresh token in the process. The returned RefreshResponse.RefreshToken
// supersedes the one passed in — persist it and discard the old value, which
// is now invalid (re-presenting it outside the server's grace window revokes
// the whole session).
//
// The endpoint does not require an Authorization header; the refresh token in
// the request body is the credential.
func (s *AuthService) Refresh(ctx context.Context, req *types.RefreshRequest) (*types.RefreshResponse, error) {
	var out types.RefreshResponse
	_, _, err := s.c.do(ctx, "POST", "/api/v1/auth/refresh", req, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Logout revokes the session family the supplied refresh token belongs to.
// It is idempotent — an unknown or empty token still returns nil so callers
// can log out without first checking token validity.
func (s *AuthService) Logout(ctx context.Context, req *types.RefreshRequest) error {
	_, _, err := s.c.do(ctx, "POST", "/api/v1/auth/logout", req, nil, nil)
	return err
}

// LogoutAll revokes every live refresh-token session for the authenticated
// user (all devices). It requires a session JWT (WithJWT); API-key clients are
// rejected with 403 API_KEY_SESSION_REQUIRED.
func (s *AuthService) LogoutAll(ctx context.Context) error {
	_, _, err := s.c.do(ctx, "POST", "/api/v1/auth/logout-all", nil, nil, nil)
	return err
}
