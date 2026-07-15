package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/kumobase/kumo-go/types"
)

// PackagesService is the entry point for /api/v1/packages/*. It exposes the
// public plan catalogue plus an org-bound sub-service:
//
//   - ListPlans() — the unauthenticated pricing catalogue.
//   - Org(slug)  — packages within a specific organization.
//
// The ecosystem protocol endpoints (npm/maven/pypi/nuget/rubygems, served at
// /<format>/:org/*) are deliberately NOT modelled here: their wire shapes are
// dictated by each package manager, so callers should point the native tool
// (npm, mvn, pip, …) at Kumo rather than drive them through this SDK.
type PackagesService struct {
	c *Client
}

// Packages returns the packages service.
func (c *Client) Packages() *PackagesService { return &PackagesService{c: c} }

// Org returns the package sub-service bound to a specific org slug. All of its
// methods operate inside the supplied org.
func (p *PackagesService) Org(orgSlug string) *OrgPackagesService {
	return &OrgPackagesService{c: p.c, orgSlug: orgSlug}
}

// ListPlans returns the public Kumo Packages plan catalogue from
// GET /api/v1/packages/plans. The server wraps the catalogue in {"plans":[…]} —
// this method flattens it to the inner slice for ergonomics. Internal cost
// structure (base cost, margin, the per-GB-hour rate) is never exposed; see
// types.PackagesPlanOption.
//
// The route is unauthenticated server-side; the SDK sends credentials anyway,
// which the server ignores.
func (p *PackagesService) ListPlans(ctx context.Context, opts ...ListOption) ([]types.PackagesPlanOption, error) {
	q := resolveListOpts(opts)
	var out types.PackagesPricingResponse
	_, _, err := p.c.do(ctx, "GET", withQuery("/api/v1/packages/plans", q), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return out.Plans, nil
}

// ── Org-scoped packages ────────────────────────────────────────────

// OrgPackagesService backs
// /api/v1/packages/organizations/:slug/packages/*. Bound to a specific
// organization slug at construction.
type OrgPackagesService struct {
	c       *Client
	orgSlug string
}

// packagesURL builds a path under the bound org. The list route is registered
// server-side as "/", hence the trailing slash when suffix is empty.
func (s *OrgPackagesService) packagesURL(suffix string) string {
	base := fmt.Sprintf("/api/v1/packages/organizations/%s/packages/", s.orgSlug)
	return base + suffix
}

// pkgPath builds the {format}/{name} portion addressing a single package.
//
// The name is escaped with url.PathEscape, which turns "/" into "%2F" while
// leaving "@" alone — so a scoped name like "@acme/utils" arrives as one path
// segment. This matters: the server routes these as a greedy wildcard and
// URL-decodes the NAME segment only (packages.parseManagementPath). Do not
// substitute url.QueryEscape, which also escapes "@" and encodes spaces as "+".
//
// The VERSION is deliberately NOT escaped — see versionPath.
func (s *OrgPackagesService) pkgPath(format types.PackageFormat, name string) string {
	return string(format) + "/" + url.PathEscape(name)
}

// versionPath addresses a single version of a package.
//
// The version is appended raw, on purpose. The server unescapes only the name
// (packages.parseManagementPath), so an escaped version would arrive still
// encoded and never string-match the stored value. Loose semver is path-safe
// ([0-9A-Za-z.+-]) so escaping would be a no-op at best — and a bug at worst,
// e.g. "1.0.0+build" would become "1.0.0%2Bbuild" and 404. Leave it raw.
func (s *OrgPackagesService) versionPath(format types.PackageFormat, name, version string) string {
	return s.pkgPath(format, name) + "/versions/" + version
}

// List returns the packages in the bound org, paginated. Results span every
// format; PackageResponse.Format says which.
//
// Only "name" and "created_at" are honoured by WithSort — any other value
// falls back to "updated_at" server-side. There is no format or search filter
// today; WithExtraQuery is the escape hatch if one is added.
func (s *OrgPackagesService) List(ctx context.Context, opts ...ListOption) ([]types.PackageResponse, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.PackageResponse
	meta, err := s.c.doList(ctx, "GET", withQuery(s.packagesURL(""), q), &out)
	return out, meta, err
}

// Get fetches a package and all its versions. Scoped names ("@acme/utils") and
// Maven coordinates ("com.acme:lib") are supported.
//
// format is part of the package's identity, not a filter: the same name may
// exist under several formats in one org. Returns the resource ETag alongside
// the detail.
func (s *OrgPackagesService) Get(ctx context.Context, format types.PackageFormat, name string) (*types.PackageDetailResponse, string, error) {
	var out types.PackageDetailResponse
	etag, _, err := s.c.do(ctx, "GET", s.packagesURL(s.pkgPath(format, name)), nil, nil, &out)
	if err != nil {
		return nil, "", err
	}
	return &out, etag, nil
}

// GetVersion fetches a single published version. Returns an APIError carrying
// codes.PackageVersionNotFound when the package exists but the version does
// not.
//
// No ETag is returned: the server emits the parent PACKAGE's ETag on this
// route, which would be misleading attached to a version resource.
func (s *OrgPackagesService) GetVersion(ctx context.Context, format types.PackageFormat, name, version string) (*types.PackageVersionResponse, error) {
	var out types.PackageVersionResponse
	_, _, err := s.c.do(ctx, "GET", s.packagesURL(s.versionPath(format, name, version)), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete unpublishes an entire package. The server soft-deletes and schedules a
// GC purge, so this returns once the deletion is SCHEDULED — not once the blobs
// are reclaimed.
//
// Note: the server does not currently dedupe deletes by Idempotency-Key. The
// header is sent (see WithIdempotencyKey) and harmlessly ignored.
func (s *OrgPackagesService) Delete(ctx context.Context, format types.PackageFormat, name string, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "DELETE", s.packagesURL(s.pkgPath(format, name)), nil, &wopts, nil)
	return err
}

// DeleteVersion unpublishes a single version. Same scheduling semantics as
// Delete.
func (s *OrgPackagesService) DeleteVersion(ctx context.Context, format types.PackageFormat, name, version string, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "DELETE", s.packagesURL(s.versionPath(format, name, version)), nil, &wopts, nil)
	return err
}
