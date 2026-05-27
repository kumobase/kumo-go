package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/kumobase/kumo-go/types"
)

// SecretsService backs /api/v1/secrets/*. All operations are synchronous;
// there's no async polling pattern for secrets (k8s Secret resources
// reconcile internally without exposing operation IDs).
type SecretsService struct {
	c *Client
}

// Secrets returns the secrets service.
func (c *Client) Secrets() *SecretsService { return &SecretsService{c: c} }

// Get fetches a single secret (with its sensitive payload — env vars,
// registry password, file content, certificate). Returns the response and
// the ETag for use with IfMatch on a subsequent Update.
//
// The detail endpoint is the only way to read a secret's contents back —
// the List shape excludes the payload for log/audit safety.
func (s *SecretsService) Get(ctx context.Context, id uint) (*types.ResponseGetSecret, string, error) {
	var out types.ResponseGetSecret
	etag, _, err := s.c.do(ctx, "GET", fmt.Sprintf("/api/v1/secrets/%d", id), nil, nil, &out)
	if err != nil {
		return nil, "", err
	}
	return &out, etag, nil
}

// GetByName fetches a secret by its name instead of its numeric id, hitting
// the same detail endpoint as Get (the server resolves a non-numeric path
// segment as a name). Returns 409 AMBIGUOUS_NAME if more than one secret
// shares the name — fall back to Get(ctx, id) to disambiguate. Returns the
// ETag like Get.
func (s *SecretsService) GetByName(ctx context.Context, name string) (*types.ResponseGetSecret, string, error) {
	var out types.ResponseGetSecret
	etag, _, err := s.c.do(ctx, "GET", "/api/v1/secrets/"+url.PathEscape(name), nil, nil, &out)
	if err != nil {
		return nil, "", err
	}
	return &out, etag, nil
}

// List returns paginated secrets (payload omitted from list items).
// Accepts the canonical ListOption set plus endpoint-specific filters
// via WithExtraQuery("type", "registry") or WithExtraQuery("search", "foo").
func (s *SecretsService) List(ctx context.Context, opts ...ListOption) ([]types.GetSecretAllResponse, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.GetSecretAllResponse
	meta, err := s.c.doList(ctx, "GET", withQuery("/api/v1/secrets", q), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, meta, nil
}

// Create persists a new secret. Honors Idempotency-Key. Exactly one of
// SecretRegistry / EnvironmentVariables / FileContent / CertificateContent
// must be populated matching req.Type — the server returns 400
// SECRET_UNSUPPORTED_TYPE or SECRET_ENV_VARS_EMPTY on shape violations.
//
// The returned response is the list-item shape (id, name, type,
// timestamps); fetch the secret by id for the full payload.
func (s *SecretsService) Create(ctx context.Context, req *types.CreateSecretRequest, opts ...WriteOption) (*types.GetSecretAllResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.GetSecretAllResponse
	_, _, err = s.c.do(ctx, "POST", "/api/v1/secrets", req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Update mutates an existing secret. Type is immutable — passing a
// different type than the stored secret returns 400 SECRET_TYPE_IMMUTABLE.
//
// Pass IfMatch(etag) for optimistic concurrency. The returned response is
// the freshly re-fetched secret (server re-reads after the write per
// CLAUDE.md rule 10).
func (s *SecretsService) Update(ctx context.Context, id uint, req *types.UpdateSecretRequest, opts ...WriteOption) (*types.ResponseGetSecret, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.ResponseGetSecret
	_, _, err = s.c.do(ctx, "PATCH", fmt.Sprintf("/api/v1/secrets/%d", id), req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete removes a secret. Returns 409 SECRET_IN_USE (errors.Is matches
// nothing useful — branch via IsCode or check (*APIError).Code) if any app
// still references the secret. List the referencing apps via Get and the
// UsedBy field.
func (s *SecretsService) Delete(ctx context.Context, id uint, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "DELETE", fmt.Sprintf("/api/v1/secrets/%d", id), nil, &wopts, nil)
	return err
}
