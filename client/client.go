package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kumobase/kumo-go/types"
)

// Client is the top-level entry point. Construct once at program start with
// New and share across goroutines — it is safe for concurrent use.
//
// All resource access goes through the per-service accessors (Apps(), VPS(),
// Secrets(), …) so the SDK can evolve internals without touching the
// public method surface.
type Client struct {
	cfg *config
}

// New constructs a Client. baseURL is required; opts configure auth and
// behaviour. Auth is mandatory — pass exactly one of WithAPIKey or
// WithJWT. Calling neither, or both, returns an error here rather than
// failing later at request time.
//
// The returned *Client is safe for concurrent use by multiple goroutines.
func New(baseURL string, opts ...Option) (*Client, error) {
	cfg := defaultConfig()
	cfg.baseURL = strings.TrimRight(baseURL, "/")
	for _, opt := range opts {
		opt(cfg)
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &Client{cfg: cfg}, nil
}

// do is the single request entry point used by every service method. It
// owns the full per-request lifecycle:
//
//  1. Marshal body to JSON (skipped if body == nil)
//  2. Build the *http.Request with base URL, auth header, User-Agent, and
//     any per-call headers (Idempotency-Key, If-Match) from writeOpts
//  3. Run the retry loop (per cfg.retry):
//     - send request
//     - read response body fully
//     - parse StructureResponse to extract the wire Code
//     - classify (retry? wait? give up?)
//     - on retryable, sleep ctx-aware, repeat
//  4. On 2xx: return raw body bytes + ETag for the caller to decode
//  5. On 4xx/5xx: build *APIError and return it
//
// out is decoded from the Data field of the StructureResponse on success;
// pass nil if you don't need the body decoded (e.g. 204 No Content).
//
// wopts is non-nil only on POST/PATCH/DELETE — readers (GET) pass nil and
// skip Idempotency-Key + If-Match generation.
func (c *Client) do(
	ctx context.Context,
	method, path string,
	body any,
	wopts *writeOpts,
	out any,
) (etag string, resp *http.Response, err error) {
	url := c.cfg.baseURL + path

	// Marshal body once; we re-use the bytes across retries.
	var bodyBytes []byte
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return "", nil, fmt.Errorf("kumo: marshal request: %w", err)
		}
	}

	attempts := c.cfg.retry.maxAttempts
	if attempts < 1 {
		attempts = 1 // even with WithRetries(0), make one attempt
	}

	var (
		lastResp *http.Response
		lastBody []byte
		lastErr  error
	)

	for attempt := 0; attempt < attempts; attempt++ {
		// ctx may be cancelled between attempts — check before each.
		if err := ctx.Err(); err != nil {
			return "", nil, err
		}

		req, rerr := buildRequest(ctx, method, url, bodyBytes, c.cfg.userAgent, wopts)
		if rerr != nil {
			return "", nil, rerr
		}
		c.cfg.auth.apply(req)

		start := time.Now()
		resp, lastErr = c.cfg.httpClient.Do(req)
		dur := time.Since(start)

		// Drain body so the connection can be reused. Body may be nil on
		// transport errors.
		if resp != nil {
			lastBody, _ = io.ReadAll(resp.Body)
			_ = resp.Body.Close()
		} else {
			lastBody = nil
		}
		lastResp = resp

		// Extract Code early so classify() can branch on it (the
		// IDEMPOTENCY_IN_PROGRESS / IDEMPOTENCY_KEY_CONFLICT distinction
		// hinges on the body Code, not the status alone).
		code := extractCode(lastBody)

		c.cfg.logger("kumo: %s %s -> %s in %s (attempt %d/%d, code=%q)",
			method, path, statusOf(resp), dur, attempt+1, attempts, code)

		decision := classify(resp, lastErr, code)
		if !decision.retry || attempt == attempts-1 {
			break
		}

		// Choose sleep: explicit Retry-After wins; otherwise exponential.
		sleep := decision.retryAfter
		if sleep <= 0 {
			sleep = c.cfg.retry.backoffFor(attempt)
		}
		if sleepErr := sleepCtx(ctx, sleep); sleepErr != nil {
			return "", nil, sleepErr
		}
	}

	// Network failure with no response — propagate.
	if lastErr != nil && lastResp == nil {
		return "", nil, fmt.Errorf("kumo: %s %s: %w", method, path, lastErr)
	}

	// HTTP-level failure — build APIError.
	if lastResp.StatusCode >= 400 {
		return "", lastResp, buildAPIError(lastResp.StatusCode, lastBody)
	}

	// 2xx: optionally decode Data into out.
	tag := readETag(lastResp)
	if out != nil && len(lastBody) > 0 {
		var env types.StructureResponse
		if jerr := json.Unmarshal(lastBody, &env); jerr != nil {
			return tag, lastResp, fmt.Errorf("kumo: decode response envelope: %w (body=%s)", jerr, truncForErr(lastBody))
		}
		if len(env.Data) > 0 {
			if jerr := json.Unmarshal(env.Data, out); jerr != nil {
				return tag, lastResp, fmt.Errorf("kumo: decode response Data: %w (body=%s)", jerr, truncForErr(env.Data))
			}
		}
	}
	return tag, lastResp, nil
}

// doList is a list-endpoint variant of do that also returns the Meta block
// alongside the decoded Data items. Caller passes a pointer to a slice as
// out.
func (c *Client) doList(
	ctx context.Context,
	method, path string,
	out any,
) (*types.Meta, error) {
	url := c.cfg.baseURL + path

	attempts := c.cfg.retry.maxAttempts
	if attempts < 1 {
		attempts = 1
	}

	var (
		lastResp *http.Response
		lastBody []byte
		lastErr  error
	)
	for attempt := 0; attempt < attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		req, rerr := buildRequest(ctx, method, url, nil, c.cfg.userAgent, nil)
		if rerr != nil {
			return nil, rerr
		}
		c.cfg.auth.apply(req)

		start := time.Now()
		lastResp, lastErr = c.cfg.httpClient.Do(req)
		dur := time.Since(start)
		if lastResp != nil {
			lastBody, _ = io.ReadAll(lastResp.Body)
			_ = lastResp.Body.Close()
		} else {
			lastBody = nil
		}
		code := extractCode(lastBody)
		c.cfg.logger("kumo: %s %s -> %s in %s (attempt %d/%d, code=%q)",
			method, path, statusOf(lastResp), dur, attempt+1, attempts, code)

		decision := classify(lastResp, lastErr, code)
		if !decision.retry || attempt == attempts-1 {
			break
		}
		sleep := decision.retryAfter
		if sleep <= 0 {
			sleep = c.cfg.retry.backoffFor(attempt)
		}
		if serr := sleepCtx(ctx, sleep); serr != nil {
			return nil, serr
		}
	}

	if lastErr != nil && lastResp == nil {
		return nil, fmt.Errorf("kumo: %s %s: %w", method, path, lastErr)
	}
	if lastResp.StatusCode >= 400 {
		return nil, buildAPIError(lastResp.StatusCode, lastBody)
	}

	var env types.StructureResponse
	if err := json.Unmarshal(lastBody, &env); err != nil {
		return nil, fmt.Errorf("kumo: decode response envelope: %w (body=%s)", err, truncForErr(lastBody))
	}
	if out != nil && len(env.Data) > 0 {
		if err := json.Unmarshal(env.Data, out); err != nil {
			return nil, fmt.Errorf("kumo: decode response Data: %w (body=%s)", err, truncForErr(env.Data))
		}
	}
	return env.Meta, nil
}

// buildRequest is extracted so do/doList share the header/idempotency/etag
// wiring without duplicating it. Returns a fresh request on each call so
// retries don't accidentally re-use a consumed body.
func buildRequest(
	ctx context.Context,
	method, url string,
	bodyBytes []byte,
	userAgent string,
	wopts *writeOpts,
) (*http.Request, error) {
	var bodyReader io.Reader
	if bodyBytes != nil {
		bodyReader = bytes.NewReader(bodyBytes)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("kumo: build request: %w", err)
	}
	if bodyBytes != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent)
	if wopts != nil {
		if wopts.idempotencyKey != "" {
			req.Header.Set("Idempotency-Key", wopts.idempotencyKey)
		}
		if wopts.ifMatch != "" {
			req.Header.Set("If-Match", wopts.ifMatch)
		}
	}
	return req, nil
}

// extractCode pulls the StructureResponse.Code from a response body
// without committing to a full decode. Returns "" if the body isn't a JSON
// object or has no code field. Used to discriminate idempotency-related
// 409s before deciding whether to retry.
func extractCode(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	var env struct {
		Code string `json:"code"`
	}
	if json.Unmarshal(body, &env) != nil {
		return ""
	}
	return env.Code
}

// buildAPIError parses the response body into a StructureResponse and
// wraps the relevant fields into *APIError. Falls back gracefully when the
// body is missing, empty, or malformed (proxy stripped it, gateway error
// page, etc.) so callers always get a usable error message.
func buildAPIError(status int, body []byte) *APIError {
	out := &APIError{StatusCode: status, Body: append([]byte(nil), body...)}
	if len(body) == 0 {
		out.Message = fmt.Sprintf("http %d", status)
		return out
	}
	var env types.StructureResponse
	if err := json.Unmarshal(body, &env); err != nil {
		// Not JSON — surface a snippet of the raw body as the message.
		out.Message = fmt.Sprintf("http %d: %s", status, string(truncForErr(body)))
		return out
	}
	out.Code = env.Code
	out.Message = env.Message
	if out.Message == "" {
		out.Message = fmt.Sprintf("http %d", status)
	}
	return out
}

// truncForErr trims a body to a manageable size for inclusion in error
// messages (avoid dumping multi-MB payloads into logs).
func truncForErr(b []byte) []byte {
	const max = 512
	if len(b) <= max {
		return b
	}
	return append(b[:max:max], []byte("…(truncated)")...)
}

// statusOf returns a printable status string for the logger.
func statusOf(resp *http.Response) string {
	if resp == nil {
		return "<no response>"
	}
	return resp.Status
}
