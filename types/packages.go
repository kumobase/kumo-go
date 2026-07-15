package types

import "time"

// PackageFormat identifies the ecosystem a package belongs to. It is part of a
// package's identity, not a filter: the server enforces uniqueness on
// (organization, format, normalized name), so the same name can exist as both
// an npm and a PyPI package in one org. Every management call that addresses a
// single package therefore takes a format.
type PackageFormat string

const (
	PackageFormatNPM      PackageFormat = "npm"
	PackageFormatMaven    PackageFormat = "maven"
	PackageFormatPyPI     PackageFormat = "pypi"
	PackageFormatNuGet    PackageFormat = "nuget"
	PackageFormatRubyGems PackageFormat = "rubygems"
)

// PackagesPricingResponse is the public-facing pricing surface returned by
// GET /api/v1/packages/plans. Plans is a list so new tiers can ship
// without breaking older clients. Mirrors RegistryPricingResponse.
type PackagesPricingResponse struct {
	Plans []PackagesPlanOption `json:"plans"`
}

// PackagesPlanOption is the sanitized projection of a billing plan for the
// Kumo Packages product — it intentionally omits base cost, margin, and the
// internal per-GB-hour rate so cost structure never leaks. PricePerUnit is a
// decimal string quoting the per-GB-month price to match the registry/ECR
// convention.
type PackagesPlanOption struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Unit          string `json:"unit"`
	PricePerUnit  string `json:"price_per_unit"`
	Currency      string `json:"currency"`
	ChargeModel   string `json:"charge_model"`
	BillingPeriod string `json:"billing_period"`
}

// PackageResponse is the management-API summary projection of a package,
// returned by GET /api/v1/packages/organizations/{slug}/packages/ and embedded
// in PackageDetailResponse. LatestVersion is the "latest" dist-tag when the
// format has one (npm), else the newest published version; it is empty for a
// package whose versions have all been unpublished.
type PackageResponse struct {
	ID             uint      `json:"id"`
	OrganizationID uint      `json:"organization_id"`
	Format         string    `json:"format"`
	Name           string    `json:"name"`
	LatestVersion  string    `json:"latest_version,omitempty"`
	VersionCount   int       `json:"version_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// PackageVersionResponse is a single published version. Deprecated is nil for
// live versions and carries the deprecation message when set. Shasum (SHA-1)
// and Integrity (SRI) are npm dist metadata and are empty for other formats.
type PackageVersionResponse struct {
	Version     string    `json:"version"`
	SizeBytes   int64     `json:"size_bytes"`
	Deprecated  *string   `json:"deprecated,omitempty"`
	Shasum      string    `json:"shasum,omitempty"`
	Integrity   string    `json:"integrity,omitempty"`
	PublishedAt time.Time `json:"published_at"`
}

// PackageDetailResponse is the full package projection returned by
// GET /api/v1/packages/organizations/{slug}/packages/{format}/{name}. Versions
// excludes the internal Maven metadata pseudo-version. DistTags is npm-only —
// an empty map for Maven/PyPI/NuGet/RubyGems.
//
// Mirrors packages.PackageDetail server-side.
type PackageDetailResponse struct {
	Package  PackageResponse          `json:"package"`
	Versions []PackageVersionResponse `json:"versions"`
	DistTags map[string]string        `json:"dist_tags"`
}
