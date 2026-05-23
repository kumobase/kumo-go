package types

import "time"

// CreateOrganizationRequest is the body for POST /api/v1/registry/organizations.
// Honors Idempotency-Key.
//
// Slug becomes the registry namespace under registry.kumo.run/<slug>/<repo>
// and is immutable once set. It is 2..64 chars, validated server-side
// against a reserved-word list; common reserved values return 400
// ORG_SLUG_RESERVED.
type CreateOrganizationRequest struct {
	Slug        string `json:"slug"`
	DisplayName string `json:"display_name"`
}

// UpdateOrganizationRequest is the body for PATCH /api/v1/registry/organizations/:slug.
// Slug is immutable — a body that includes the slug field is rejected with
// 400 ORG_SLUG_IMMUTABLE (the field is captured here only so the server can
// detect change attempts).
type UpdateOrganizationRequest struct {
	DisplayName             *string `json:"display_name,omitempty"`
	RegistryAutoCreateRepos *bool   `json:"registry_auto_create_repos,omitempty"`

	// Slug is captured only to surface a clear immutability error.
	Slug *string `json:"slug,omitempty"`
}

// OrganizationResponse is the detail shape returned by every org endpoint
// (create, list, get, update). RegistrySuspendedAt is non-nil when an
// admin has suspended the org's registry access (push/pull will fail).
type OrganizationResponse struct {
	ID                      uint       `json:"id"`
	Slug                    string     `json:"slug"`
	DisplayName             string     `json:"display_name"`
	OwnerUserID             uint       `json:"owner_user_id"`
	RegistryAutoCreateRepos bool       `json:"registry_auto_create_repos"`
	RegistrySuspendedAt     *time.Time `json:"registry_suspended_at,omitempty"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}
