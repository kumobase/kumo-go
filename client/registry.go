package client

import (
	"context"
	"fmt"

	"github.com/kumobase/kumo-go/types"
)

// RegistryService is the entry point for /api/v1/registry/*. It exposes
// two sub-services:
//
//   - Orgs() — CRUD on registry organizations (namespaces).
//   - Repos(slug) — CRUD on repositories within a specific organization.
//
// The two are separate types so the slug-binding is explicit in the call
// site (caller can't accidentally call a repo method without specifying
// which org).
type RegistryService struct {
	c *Client
}

// Registry returns the registry service.
func (c *Client) Registry() *RegistryService { return &RegistryService{c: c} }

// Orgs returns the organizations sub-service.
func (r *RegistryService) Orgs() *OrganizationsService {
	return &OrganizationsService{c: r.c}
}

// Repos returns the repositories sub-service bound to a specific org
// slug. All Repos methods operate inside the supplied org.
func (r *RegistryService) Repos(orgSlug string) *RepositoriesService {
	return &RepositoriesService{c: r.c, orgSlug: orgSlug}
}

// ListPlans returns the public container-registry plan catalogue (storage /
// transfer billing tiers) from GET /api/v1/registry/plans. The server wraps
// the catalogue in {"plans":[…]} — this method flattens it to the inner
// slice for ergonomics. Internal cost structure (base cost, margin) is never
// exposed; see types.RegistryPlanOption.
func (r *RegistryService) ListPlans(ctx context.Context, opts ...ListOption) ([]types.RegistryPlanOption, error) {
	q := resolveListOpts(opts)
	var out types.RegistryPricingResponse
	_, _, err := r.c.do(ctx, "GET", withQuery("/api/v1/registry/plans", q), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return out.Plans, nil
}

// ── Organizations ──────────────────────────────────────────────────

// OrganizationsService backs
// /api/v1/registry/organizations/*. All operations are synchronous.
type OrganizationsService struct {
	c *Client
}

// Get fetches an org by its slug. Returns the ETag for use with IfMatch
// on a subsequent Update.
func (s *OrganizationsService) Get(ctx context.Context, slug string) (*types.OrganizationResponse, string, error) {
	var out types.OrganizationResponse
	etag, _, err := s.c.do(ctx, "GET",
		fmt.Sprintf("/api/v1/registry/organizations/%s", slug), nil, nil, &out)
	if err != nil {
		return nil, "", err
	}
	return &out, etag, nil
}

// List returns all organizations the authenticated user is a member of.
// Not paginated server-side today (low cardinality per user); returns the
// full list.
func (s *OrganizationsService) List(ctx context.Context) ([]types.OrganizationResponse, error) {
	var out []types.OrganizationResponse
	_, err := s.c.doList(ctx, "GET", "/api/v1/registry/organizations", &out)
	return out, err
}

// Create provisions a new organization. Honors Idempotency-Key. Slug is
// immutable once set and is validated against a reserved-word list — see
// codes.OrgSlugInvalid / codes.OrgSlugReserved / codes.OrgSlugTaken.
func (s *OrganizationsService) Create(ctx context.Context, req *types.CreateOrganizationRequest, opts ...WriteOption) (*types.OrganizationResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.OrganizationResponse
	_, _, err = s.c.do(ctx, "POST", "/api/v1/registry/organizations", req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Update mutates DisplayName / RegistryAutoCreateRepos on an org. Passing
// a Slug field on the request body returns 400 ORG_SLUG_IMMUTABLE — the
// SDK doesn't strip it for you so caller bugs surface clearly.
func (s *OrganizationsService) Update(ctx context.Context, slug string, req *types.UpdateOrganizationRequest, opts ...WriteOption) (*types.OrganizationResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.OrganizationResponse
	_, _, err = s.c.do(ctx, "PATCH",
		fmt.Sprintf("/api/v1/registry/organizations/%s", slug), req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete removes an organization. Returns 409 ORG_HAS_REPOS if any
// repositories still exist under the org — delete them first.
func (s *OrganizationsService) Delete(ctx context.Context, slug string, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "DELETE",
		fmt.Sprintf("/api/v1/registry/organizations/%s", slug), nil, &wopts, nil)
	return err
}

// ── Repositories ───────────────────────────────────────────────────

// RepositoriesService backs
// /api/v1/registry/organizations/:slug/repositories/*. Bound to a
// specific organization slug at construction.
type RepositoriesService struct {
	c       *Client
	orgSlug string
}

// reposURL builds a path under the bound org.
func (s *RepositoriesService) reposURL(suffix string) string {
	if suffix == "" {
		return fmt.Sprintf("/api/v1/registry/organizations/%s/repositories", s.orgSlug)
	}
	return fmt.Sprintf("/api/v1/registry/organizations/%s/repositories/%s", s.orgSlug, suffix)
}

// Get fetches a repository by name within the bound org.
func (s *RepositoriesService) Get(ctx context.Context, name string) (*types.RepositoryResponse, error) {
	var out types.RepositoryResponse
	_, _, err := s.c.do(ctx, "GET", s.reposURL(name), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// List returns all repositories in the bound org, paginated.
func (s *RepositoriesService) List(ctx context.Context, opts ...ListOption) ([]types.RepositoryResponse, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.RepositoryResponse
	meta, err := s.c.doList(ctx, "GET", withQuery(s.reposURL(""), q), &out)
	return out, meta, err
}

// Create provisions a new repository in the bound org. Name must match
// the OCI distribution name-component grammar (lowercase, dot/dash/underscore
// separators); see codes.RegistryInvalidRepositoryName.
func (s *RepositoriesService) Create(ctx context.Context, req *types.CreateRepositoryRequest, opts ...WriteOption) (*types.RepositoryResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.RepositoryResponse
	_, _, err = s.c.do(ctx, "POST", s.reposURL(""), req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Update changes tag mutability / soft-delete settings on a repo. Pass
// only the fields you want to change (pointer fields).
func (s *RepositoriesService) Update(ctx context.Context, name string, req *types.UpdateRepositoryRequest, opts ...WriteOption) (*types.RepositoryResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.RepositoryResponse
	_, _, err = s.c.do(ctx, "PATCH", s.reposURL(name), req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete removes a repository. Server enforces a soft-delete window
// configured per-org/per-repo; the repo is purged after the window
// expires.
func (s *RepositoriesService) Delete(ctx context.Context, name string, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "DELETE", s.reposURL(name), nil, &wopts, nil)
	return err
}

// ── Manifests ──────────────────────────────────────────────────────

// ListManifests returns pushed manifests for a repository, paginated.
// Hydration runs asynchronously after push — newly-pushed manifests may
// have HydratedAt == nil and minimal metadata until the hydrator
// completes.
func (s *RepositoriesService) ListManifests(ctx context.Context, repoName string, opts ...ListOption) ([]types.ManifestResponse, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.ManifestResponse
	meta, err := s.c.doList(ctx, "GET",
		withQuery(s.reposURL(repoName+"/manifests"), q), &out)
	return out, meta, err
}

// GetManifest fetches a single manifest by digest within a repository.
func (s *RepositoriesService) GetManifest(ctx context.Context, repoName, digest string) (*types.ManifestResponse, error) {
	var out types.ManifestResponse
	_, _, err := s.c.do(ctx, "GET",
		s.reposURL(repoName+"/manifests/"+digest), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
