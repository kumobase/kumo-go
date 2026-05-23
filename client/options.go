package client

import (
	"errors"
	"net/http"
	"time"
)

// Option configures the Client. Apply via New(baseURL, opts...).
type Option func(*config)

// config is the internal Client configuration. Fields are unexported on
// purpose — consumers should never reach into them directly; everything
// they can tune is exposed via With… functions for forward compatibility.
type config struct {
	baseURL    string
	httpClient *http.Client
	auth       *authConfig
	userAgent  string
	logger     func(format string, args ...any)

	retry retryPolicy

	// authErr captures errors from WithAPIKey/WithJWT (e.g. both set)
	// so New can return them rather than silently picking one.
	authErr error
}

// defaults mirror the values used by other production Go SDKs (Stripe,
// GitHub, AWS) and the values the Kumo server actually expects.
func defaultConfig() *config {
	return &config{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		userAgent:  "kumo-go/" + sdkVersionForUA() + " (https://github.com/kumobase/kumo-go)",
		logger:     func(string, ...any) {}, // no-op default
		retry: retryPolicy{
			maxAttempts: 5,
			baseDelay:   50 * time.Millisecond,
			maxDelay:    5 * time.Second,
		},
	}
}

// WithBaseURL is implicit — the baseURL is the first arg to New. This
// option exists for symmetry with other With… options if a future
// constructor variant accepts only options.
func WithBaseURL(url string) Option {
	return func(c *config) { c.baseURL = url }
}

// WithHTTPClient lets the caller bring their own *http.Client. Use this to
// inject custom transports (telemetry, proxy, timeouts, mTLS). The Client
// composes around your http.Client; it does NOT mutate it.
//
// If you pass a client with Timeout=0 the SDK respects that (caller knows
// what they're doing). The SDK's own retry logic respects ctx deadlines
// independently of your http.Client.Timeout.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *config) {
		if hc != nil {
			c.httpClient = hc
		}
	}
}

// WithRetries sets the maximum total attempts for retryable failures.
// Default 5. Pass 0 to disable retries entirely (one attempt, fail fast).
// Negative values are clamped to 0.
//
// The SDK retries on: network errors, HTTP 502/503/504, HTTP 429 (honoring
// Retry-After), HTTP 409 IDEMPOTENCY_IN_PROGRESS. It never retries 4xx
// other than those two, or 409 IDEMPOTENCY_KEY_CONFLICT (which indicates a
// caller bug — different request body with the same Idempotency-Key).
func WithRetries(maxAttempts int) Option {
	return func(c *config) {
		if maxAttempts < 0 {
			maxAttempts = 0
		}
		c.retry.maxAttempts = maxAttempts
	}
}

// WithUserAgent appends an extra string to the default User-Agent so server
// logs can attribute traffic to your integration:
//
//	WithUserAgent("terraform-provider-kumo/0.3.1")  →
//	  User-Agent: kumo-go/0.2.0 (…) terraform-provider-kumo/0.3.1
func WithUserAgent(extra string) Option {
	return func(c *config) {
		if extra != "" {
			c.userAgent = c.userAgent + " " + extra
		}
	}
}

// WithLogger sets an optional logger called once per request with a
// compact summary ("POST /api/v1/apps -> 202 in 142ms"). The default is a
// no-op. SECURITY: do not log request/response bodies unless you trust the
// log sink — they contain auth tokens, passwords, certificates, etc.
func WithLogger(logger func(format string, args ...any)) Option {
	return func(c *config) {
		if logger != nil {
			c.logger = logger
		}
	}
}

// validate is run by New after applying every option. Catches caller bugs
// (no auth, both auth modes, missing base URL) before any request flies.
func (c *config) validate() error {
	if c.authErr != nil {
		return c.authErr
	}
	if c.baseURL == "" {
		return errors.New("kumo: WithBaseURL or base URL argument is required")
	}
	if c.auth == nil {
		return errors.New("kumo: exactly one of WithAPIKey or WithJWT is required")
	}
	return nil
}

// sdkVersionForUA returns the SDKVersion string without the leading 'v' so
// the User-Agent reads cleanly. Kept here (vs importing version/) to avoid
// a circular-feeling dep between client/ and version/ — the value flows
// the other way (cmd-line tooling reads version.SDKVersion).
//
// On every SDK release, bump this constant in lockstep with
// version.SDKVersion. The pinning test in the server cross-checks the wire
// codes; the User-Agent is purely cosmetic so no test gates it.
func sdkVersionForUA() string { return "0.2.0" }
