package client

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/kumobase/kumo-go/codes"
)

// retryPolicy controls the do() pipeline's retry behaviour. Configured via
// WithRetries. Defaults set in defaultConfig.
type retryPolicy struct {
	maxAttempts int           // total attempts including the first try; 0 disables retries
	baseDelay   time.Duration // backoff = baseDelay * 2^attempt + jitter, capped at maxDelay
	maxDelay    time.Duration
}

// classifyResult decides whether to retry a request based on the HTTP
// response (if any), the network error (if any), and the response's parsed
// Code field. Returns:
//   - retryAfter > 0  → sleep retryAfter then retry
//   - retry = true && retryAfter = 0 → backoff per policy then retry
//   - retry = false   → surface the result to the caller
//
// Called once per attempt by do().
type retryDecision struct {
	retry      bool
	retryAfter time.Duration // explicit server hint (Retry-After header)
}

// classify is pure — no side effects. Easy to unit-test.
func classify(resp *http.Response, netErr error, code string) retryDecision {
	// Network-level failure → retry (caller's ctx still wins).
	if netErr != nil {
		// context.Canceled and DeadlineExceeded are caller-intent — don't retry.
		if errors.Is(netErr, context.Canceled) || errors.Is(netErr, context.DeadlineExceeded) {
			return retryDecision{}
		}
		// net.OpError, EOF, DNS failures, refused connections — all retryable.
		return retryDecision{retry: true}
	}
	if resp == nil {
		// shouldn't happen, but defensive: treat as non-retryable
		return retryDecision{}
	}
	switch resp.StatusCode {
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		// 502/503/504 — transient upstream issue, retry with backoff.
		return retryDecision{retry: true}
	case http.StatusTooManyRequests:
		// 429 — honor Retry-After exactly when present; fall back to backoff.
		return retryDecision{retry: true, retryAfter: parseRetryAfter(resp.Header.Get("Retry-After"))}
	case http.StatusConflict:
		// 409 has two flavours we care about:
		//   IDEMPOTENCY_IN_PROGRESS — the same key is still being processed,
		//     retry with backoff (server will replay the cached response once
		//     the original request completes).
		//   IDEMPOTENCY_KEY_CONFLICT — caller bug: same key, different body.
		//     Don't retry; surface to caller so they can pick a fresh key.
		if code == codes.IdempotencyInProgress {
			return retryDecision{retry: true}
		}
		return retryDecision{}
	}
	// Everything else (2xx success, other 4xx, unknown 5xx) is terminal.
	return retryDecision{}
}

// parseRetryAfter accepts the two formats RFC 9110 defines: integer
// seconds or HTTP-date. Returns 0 on parse failure (caller falls back to
// exponential backoff).
func parseRetryAfter(v string) time.Duration {
	if v == "" {
		return 0
	}
	if secs, err := strconv.Atoi(v); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}

// backoffFor returns the sleep duration for the given retry attempt (0-indexed:
// attempt 0 = first retry, after the first send failed). Adds 0-25% jitter to
// avoid thundering-herd behaviour when many clients retry in lockstep.
func (p retryPolicy) backoffFor(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	// 2^attempt can overflow uint64 around attempt=63; cap conservatively.
	if attempt > 30 {
		attempt = 30
	}
	d := p.baseDelay * (1 << uint(attempt))
	if d > p.maxDelay || d <= 0 {
		d = p.maxDelay
	}
	// Jitter: add up to 25% of d to spread retries.
	jitter := time.Duration(rand.Int63n(int64(d / 4))) // #nosec G404 — non-crypto jitter
	return d + jitter
}

// sleepCtx sleeps for d or until ctx is cancelled, whichever comes first.
// Returns ctx.Err() on cancellation so the caller can short-circuit the
// retry loop without firing another attempt.
func sleepCtx(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// isNetworkTimeout helps tests assert that a synthetic net.Error is treated
// as retryable. Not used by classify directly — kept here so the property
// stays visible alongside the policy.
func isNetworkTimeout(err error) bool {
	var ne net.Error
	return errors.As(err, &ne) && ne.Timeout()
}
