package types

import "time"

// TagMutability controls whether a pushed tag can be overwritten by a
// subsequent push. Mirrors AWS ECR's imageTagMutability setting.
type TagMutability string

const (
	TagMutabilityMutable   TagMutability = "MUTABLE"
	TagMutabilityImmutable TagMutability = "IMMUTABLE"
)

// CreateRepositoryRequest is the body for POST
// /api/v1/registry/organizations/:slug/repositories.
//
// Name must match the OCI distribution name-component grammar (lowercase
// letters/digits with single '.', '_' or '-' as internal separators), max
// 255 chars. TagMutability defaults to MUTABLE; SoftDeleteDays defaults to
// the org's setting when nil.
type CreateRepositoryRequest struct {
	Name           string        `json:"name"`
	TagMutability  TagMutability `json:"tag_mutability,omitempty"`
	SoftDeleteDays *int          `json:"soft_delete_days,omitempty"`
}

// UpdateRepositoryRequest is the body for PATCH
// /api/v1/registry/organizations/:slug/repositories/:repo. Pointer fields
// signal "only update if provided".
type UpdateRepositoryRequest struct {
	TagMutability  *TagMutability `json:"tag_mutability,omitempty"`
	SoftDeleteDays *int           `json:"soft_delete_days,omitempty"`
}

// UpdateSettingsRequest is the body for PATCH /api/v1/registry/settings
// (per-org auto-create-repos toggle).
type UpdateSettingsRequest struct {
	RegistryAutoCreateRepos *bool `json:"registry_auto_create_repos,omitempty"`
}

// RepositoryResponse is the detail shape returned by every repo endpoint.
type RepositoryResponse struct {
	ID             uint          `json:"id"`
	Name           string        `json:"name"`
	TagMutability  TagMutability `json:"tag_mutability"`
	SoftDeleteDays int           `json:"soft_delete_days"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// SettingsResponse is the per-org registry settings DTO.
type SettingsResponse struct {
	RegistryAutoCreateRepos bool `json:"registry_auto_create_repos"`
}

// ManifestResponse is the wire shape for a single pushed manifest. Fields
// populated during hydration are pointer / omitempty so a row still queued
// for hydration renders compactly (digest + media_type + pushed_at).
//
// Exactly one of (Architecture + OS) or Platforms will be populated after
// hydration: image manifest → the former, manifest index → the latter.
type ManifestResponse struct {
	ID        uint      `json:"id"`
	Digest    string    `json:"digest"`
	Tag       *string   `json:"tag,omitempty"`
	MediaType string    `json:"media_type"`
	SizeBytes int64     `json:"size_bytes"` // size of the manifest document itself
	PushedAt  time.Time `json:"pushed_at"`

	ConfigDigest   *string            `json:"config_digest,omitempty"`
	Architecture   *string            `json:"architecture,omitempty"`
	OS             *string            `json:"os,omitempty"`
	OSVersion      *string            `json:"os_version,omitempty"`
	Variant        *string            `json:"variant,omitempty"`
	Platform       *string            `json:"platform,omitempty"` // display string "linux/arm64" or "linux/arm/v7"
	ImageCreatedAt *time.Time         `json:"image_created_at,omitempty"`
	Labels         map[string]string  `json:"labels,omitempty"`
	LayerCount     *int               `json:"layer_count,omitempty"`
	// ImageSizeBytes is the total compressed image size — the sum of the
	// manifest's layer blob sizes (config excluded), the figure registries
	// like Docker Hub display. Populated at hydration for image manifests;
	// 0/omitted for manifest indexes and not-yet-hydrated rows.
	ImageSizeBytes int64              `json:"image_size_bytes,omitempty"`
	Platforms      []ManifestPlatform `json:"platforms,omitempty"`
	HydratedAt     *time.Time         `json:"hydrated_at,omitempty"`
	HydrationError *string            `json:"hydration_error,omitempty"`
}

// ManifestPlatform is one entry of a manifest index's child list. Platform
// is the canonical "os/arch[/variant]" string used by docker buildx and k8s
// node selectors.
type ManifestPlatform struct {
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	Variant      string `json:"variant,omitempty"`
	OSVersion    string `json:"os_version,omitempty"`
	Digest       string `json:"digest"`
	Size         int64  `json:"size"`
	Platform     string `json:"platform"`
}

// RegistryPricingResponse is the public-facing pricing surface returned by
// GET /api/v1/registry/pricing. Plans is a list so new tiers can ship
// without breaking older clients.
type RegistryPricingResponse struct {
	Plans []RegistryPlanOption `json:"plans"`
}

// RegistryPlanOption is the sanitized projection of a billing plan — it
// intentionally omits BaseCost and Margin so internal cost structure never
// leaks. PricePerUnit is a decimal string.
type RegistryPlanOption struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Unit          string `json:"unit"`
	PricePerUnit  string `json:"price_per_unit"`
	Currency      string `json:"currency"`
	ChargeModel   string `json:"charge_model"`
	BillingPeriod string `json:"billing_period"`
}
