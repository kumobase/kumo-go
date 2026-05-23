package client

import (
	"errors"
	"net/http"
)

// authMode discriminates between the two Authorization: Bearer flavours
// the server accepts. The wire header is identical for both; the server
// recognises API keys by the "kumo_sk_" prefix.
type authMode int

const (
	authAPIKey authMode = iota + 1
	authJWT
)

// authConfig captures the chosen auth strategy. Set exactly once via
// WithAPIKey or WithJWT; New rejects both-set / neither-set.
type authConfig struct {
	mode  authMode
	token string
}

// apply attaches the Authorization header to outgoing requests. Called by
// the do() pipeline on every request.
func (a *authConfig) apply(req *http.Request) {
	if a == nil || a.token == "" {
		return
	}
	req.Header.Set("Authorization", "Bearer "+a.token)
}

// WithAPIKey configures the client to authenticate with a Kumo API key
// (kumo_sk_…). API keys are resource credentials — they manage apps, vps,
// volumes, secrets, registry, and may read billing/quota. They cannot
// reach admin or account-takeover routes; the server returns 403
// API_KEY_SESSION_REQUIRED or API_KEY_ADMIN_FORBIDDEN.
//
// Cannot be combined with WithJWT.
func WithAPIKey(key string) Option {
	return func(c *config) {
		if c.auth != nil {
			c.auth = nil // mark for validation rejection
			c.authErr = errors.New("kumo: WithAPIKey and WithJWT are mutually exclusive")
			return
		}
		c.auth = &authConfig{mode: authAPIKey, token: key}
	}
}

// WithJWT configures the client to authenticate with a JWT bearer token —
// typically obtained from a `kumo login` flow that exchanges
// email+password for a session JWT. JWTs reach the full user surface
// including account/financial routes that API keys cannot.
//
// Cannot be combined with WithAPIKey.
func WithJWT(token string) Option {
	return func(c *config) {
		if c.auth != nil {
			c.auth = nil
			c.authErr = errors.New("kumo: WithAPIKey and WithJWT are mutually exclusive")
			return
		}
		c.auth = &authConfig{mode: authJWT, token: token}
	}
}
