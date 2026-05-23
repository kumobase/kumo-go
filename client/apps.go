package client

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/kumobase/kumo-go/types"
)

// AppsService is the access point for /api/v1/apps/* endpoints. Get the
// instance via Client.Apps(); the service is stateless and safe for
// concurrent use.
type AppsService struct {
	c *Client
}

// Apps returns the apps service.
func (c *Client) Apps() *AppsService { return &AppsService{c: c} }

// ListOptions controls pagination + sorting for List endpoints. The
// canonical query parameters (page, page_size, sort, sort_order) match
// the server's common.ParseListParams.
type ListOptions struct {
	Page      int    // 1-based; 0 = server default (1)
	PageSize  int    // capped at 100 server-side
	Sort      string // column name; server whitelists
	SortOrder string // "asc" or "desc"; anything else coerces to desc
	// Extra is for ad-hoc filters specific to a list endpoint (e.g.
	// "status=running" on VPS). Keys are appended verbatim; values are
	// URL-escaped.
	Extra map[string]string
}

// ListOption configures a list call via the functional-option pattern.
type ListOption func(*ListOptions)

// WithPage sets the page number on a list call (1-based).
func WithPage(page int) ListOption {
	return func(o *ListOptions) { o.Page = page }
}

// WithPageSize sets items per page (capped at 100 server-side).
func WithPageSize(size int) ListOption {
	return func(o *ListOptions) { o.PageSize = size }
}

// WithSort sets the sort column. The server whitelists allowed values
// per endpoint; invalid columns return 400.
func WithSort(col, order string) ListOption {
	return func(o *ListOptions) {
		o.Sort = col
		o.SortOrder = order
	}
}

// WithExtraQuery appends an ad-hoc query parameter to a list call. Use
// this for endpoint-specific filters not modelled here yet (e.g. VPS
// status filter, secret type filter).
func WithExtraQuery(key, value string) ListOption {
	return func(o *ListOptions) {
		if o.Extra == nil {
			o.Extra = map[string]string{}
		}
		o.Extra[key] = value
	}
}

// resolveListOpts applies opts and returns the URL-encoded query string
// (without the leading "?").
func resolveListOpts(opts []ListOption) string {
	o := ListOptions{}
	for _, opt := range opts {
		opt(&o)
	}
	v := url.Values{}
	if o.Page > 0 {
		v.Set("page", strconv.Itoa(o.Page))
	}
	if o.PageSize > 0 {
		v.Set("page_size", strconv.Itoa(o.PageSize))
	}
	if o.Sort != "" {
		v.Set("sort", o.Sort)
	}
	if o.SortOrder != "" {
		v.Set("sort_order", o.SortOrder)
	}
	for k, val := range o.Extra {
		v.Set(k, val)
	}
	return v.Encode()
}

// withQuery appends the encoded query string to path if non-empty.
func withQuery(path, query string) string {
	if query == "" {
		return path
	}
	return path + "?" + query
}

// ─── Synchronous reads ─────────────────────────────────────────────────

// Get fetches the full app detail. Returns the response and the ETag
// header for use with IfMatch on a subsequent Update.
func (s *AppsService) Get(ctx context.Context, id uint) (*types.AppByIdResponse, string, error) {
	var out types.AppByIdResponse
	etag, _, err := s.c.do(ctx, "GET", fmt.Sprintf("/api/v1/apps/%d", id), nil, nil, &out)
	if err != nil {
		return nil, "", err
	}
	return &out, etag, nil
}

// List returns one paginated page of apps and the pagination metadata.
// To iterate all pages: bump opts page until *Meta.TotalPages.
func (s *AppsService) List(ctx context.Context, opts ...ListOption) ([]types.AppImageResponse, *types.Meta, error) {
	q := resolveListOpts(opts)
	// Force pagination by always sending at least page_size; the server's
	// "opt-in" pagination triggers on any of page/page_size/name/sort_by/
	// sort_order so this ensures consistent envelopes for SDK callers.
	if q == "" {
		q = "page=1&page_size=20"
	}
	var out []types.AppImageResponse
	meta, err := s.c.doList(ctx, "GET", withQuery("/api/v1/apps", q), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, meta, nil
}

// ListOperations returns the async-operation history for an app, paginated.
func (s *AppsService) ListOperations(ctx context.Context, appID uint, opts ...ListOption) ([]types.AppOperation, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.AppOperation
	meta, err := s.c.doList(ctx, "GET",
		withQuery(fmt.Sprintf("/api/v1/apps/%d/operations", appID), q), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, meta, nil
}

// GetOperation fetches a single async operation by its UUID. Used as the
// poll target for the create/update/delete/start/stop flows.
func (s *AppsService) GetOperation(ctx context.Context, appID uint, opID string) (*types.AppOperation, error) {
	var out types.AppOperation
	_, _, err := s.c.do(ctx, "GET",
		fmt.Sprintf("/api/v1/apps/%d/operations/%s", appID, opID), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ─── Asynchronous mutations ────────────────────────────────────────────

// Create starts an app deployment. Returns the slim 202 response with the
// operation_id for polling. Use CreateAndWait for the common "block until
// deployed" flow.
//
// Honors Idempotency-Key — the SDK auto-generates a fresh UUID per call
// and replays it across retries so duplicate creates are impossible.
func (s *AppsService) Create(ctx context.Context, req *types.CreateAppRequest, opts ...WriteOption) (*types.CreateAppResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.CreateAppResponse
	_, _, err = s.c.do(ctx, "POST", "/api/v1/apps", req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Update queues an app update. Pass IfMatch(etag) for optimistic
// concurrency — server returns ErrETagMismatch if the resource has
// changed since the caller captured the tag.
func (s *AppsService) Update(ctx context.Context, id uint, req *types.UpdateAppRequest, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "PATCH", fmt.Sprintf("/api/v1/apps/%d", id), req, &wopts, nil)
	return err
}

// Delete queues an app deletion. Returns once the request is accepted —
// the actual k8s teardown happens asynchronously. Poll Get until the app
// disappears (404) to confirm deletion.
func (s *AppsService) Delete(ctx context.Context, id uint, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "DELETE", fmt.Sprintf("/api/v1/apps/%d", id), nil, &wopts, nil)
	return err
}

// Start un-suspends an app (scales replicas back from 0).
func (s *AppsService) Start(ctx context.Context, id uint, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/apps/%d/start", id), nil, &wopts, nil)
	return err
}

// Stop suspends an app (scales replicas to 0). Idempotent — calling Stop
// on an already-stopped app returns 409 APP_ALREADY_STOPPED.
func (s *AppsService) Stop(ctx context.Context, id uint, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/apps/%d/stop", id), nil, &wopts, nil)
	return err
}

// ─── Custom domains ────────────────────────────────────────────────────

// AddCustomDomain attaches a custom FQDN to an exposed app.
func (s *AppsService) AddCustomDomain(ctx context.Context, appID uint, domain string, opts ...WriteOption) (*types.CustomDomainInfo, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.CustomDomainInfo
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/apps/%d/custom-domain", appID),
		&types.AddCustomDomainRequest{Domain: domain}, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCustomDomain fetches the custom domain currently bound to an app
// (or returns IsNotFound() == true if no custom domain is attached).
func (s *AppsService) GetCustomDomain(ctx context.Context, appID uint) (*types.CustomDomainInfo, error) {
	var out types.CustomDomainInfo
	_, _, err := s.c.do(ctx, "GET", fmt.Sprintf("/api/v1/apps/%d/custom-domain", appID), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCustomDomain detaches the custom domain from an app.
func (s *AppsService) DeleteCustomDomain(ctx context.Context, appID uint, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "DELETE", fmt.Sprintf("/api/v1/apps/%d/custom-domain", appID), nil, &wopts, nil)
	return err
}

// VerifyCustomDomain triggers a DNS check on the attached custom domain
// and returns the updated verification status.
func (s *AppsService) VerifyCustomDomain(ctx context.Context, appID uint, opts ...WriteOption) (*types.CustomDomainInfo, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.CustomDomainInfo
	_, _, err = s.c.do(ctx, "POST",
		fmt.Sprintf("/api/v1/apps/%d/custom-domain/verify", appID), nil, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ValidateImagePullable sanity-checks that a Docker image reference (and
// optional registry credential) can be pulled. Useful for pre-flighting a
// CreateApp call so the user sees the failure before the deploy queue.
func (s *AppsService) ValidateImagePullable(ctx context.Context, req *types.ValidateImagePullableRequest) (*types.ValidateImagePullableResponse, error) {
	var out types.ValidateImagePullableResponse
	_, _, err := s.c.do(ctx, "POST", "/api/v1/apps/validate-image", req, &writeOpts{}, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ─── Polling helpers ───────────────────────────────────────────────────

// PollOperation blocks until the given operation reaches a terminal state
// (succeeded, failed, or cancelled) or until ctx / PollMaxWait expires.
//
// On Status == "failed", returns the operation alongside *APIError whose
// Code is the operation's own ErrorCode — so callers branch on Code the
// same way they would on a synchronous failure.
//
// On Status == "succeeded", returns (operation, nil).
// On Status == "cancelled", returns (operation, *APIError{Code: "OPERATION_CANCELLED"}).
func (s *AppsService) PollOperation(ctx context.Context, appID uint, opID string, opts ...PollOption) (*types.AppOperation, error) {
	return PollResource(ctx,
		func(ctx context.Context) (*types.AppOperation, error) {
			return s.GetOperation(ctx, appID, opID)
		},
		func(op *types.AppOperation) (bool, error) {
			if op == nil {
				return false, nil
			}
			switch op.Status {
			case types.AppOperationStatusSucceeded:
				return true, nil
			case types.AppOperationStatusFailed:
				code := ""
				msg := "app operation failed"
				if op.ErrorCode != nil {
					code = *op.ErrorCode
				}
				if op.ErrorMsg != nil {
					msg = *op.ErrorMsg
				}
				return true, &APIError{StatusCode: 0, Code: code, Message: msg}
			case types.AppOperationStatusCancelled:
				return true, &APIError{StatusCode: 0, Code: "OPERATION_CANCELLED", Message: "operation was cancelled"}
			default:
				return false, nil
			}
		},
		opts...,
	)
}

// CreateAndWait composes Create + PollOperation + Get — the common "kick
// off a deploy and block until ready" flow. Returns the freshly-fetched
// AppByIdResponse on success so callers don't need a second Get round-trip.
//
// On operation failure (status=failed), returns nil and the *APIError
// from PollOperation. On ctx cancellation, returns ctx.Err().
//
// Pass PollOption (WithPollInterval, WithPollMaxWait, etc.) to tune the
// polling cadence; WriteOption is consumed by the underlying Create.
//
// Note: WriteOption and PollOption are different types — pass them via
// CreateAndWaitOpts if you need both, or split into Create + PollOperation
// + Get manually.
func (s *AppsService) CreateAndWait(ctx context.Context, req *types.CreateAppRequest, pollOpts ...PollOption) (*types.AppByIdResponse, error) {
	created, err := s.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	if created.OperationID == "" {
		// Server returned 202 without an operation_id — shouldn't happen,
		// but fall back to fetching the app directly.
		app, _, gerr := s.Get(ctx, created.ID)
		return app, gerr
	}
	if _, perr := s.PollOperation(ctx, created.ID, created.OperationID, pollOpts...); perr != nil {
		return nil, perr
	}
	app, _, gerr := s.Get(ctx, created.ID)
	return app, gerr
}
