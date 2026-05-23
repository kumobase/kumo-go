package types

import "time"

// API-key management routes are sessionOnly — they cannot be called with
// another API key (leaked key must not mint replacements). A CLI client
// typically calls these with a JWT obtained from the email/password login
// flow, then issues an API key for subsequent resource calls.

// CreateAPIKeyRequest is the body for POST /api/v1/api-keys.
//
// Scopes is one or more of {"read", "write"}; an empty list defaults to
// both. ExpiresInDays is optional (1..365); when omitted the key does not
// expire automatically.
//
// RegistryScope makes the key a registry credential (Harbor / GHCR-style
// robot account) instead of a control-plane key. Registry keys are
// rejected on every /api/v1/* route (403 REGISTRY_KEY_HTTP_FORBIDDEN) and
// can only authenticate against the OCI /v2/token endpoint.
type CreateAPIKeyRequest struct {
	Name          string              `json:"name"` // 1..100 chars
	ExpiresInDays *int                `json:"expires_in_days,omitempty"`
	Scopes        []string            `json:"scopes,omitempty"`
	RegistryScope *RegistryScopeInput `json:"registry_scope,omitempty"`
}

// RegistryScopeInput turns the new key into a registry credential.
// OrgSlug is required; RepoName pins the key to a single repository
// (org-wide when omitted). Permissions is a non-empty subset of
// {"pull", "push", "delete"}.
type RegistryScopeInput struct {
	OrgSlug     string   `json:"org_slug"`
	RepoName    *string  `json:"repo_name,omitempty"`
	Permissions []string `json:"permissions"`
}

// UpdateAPIKeyRequest is the body for PATCH /api/v1/api-keys/:id. Only the
// display name and expiry can be rotated; scopes and registry binding are
// immutable post-create (replace by delete + create).
type UpdateAPIKeyRequest struct {
	Name          *string `json:"name,omitempty"`
	ExpiresInDays *int    `json:"expires_in_days,omitempty"`
}

// APIKeyResponse is the metadata shape returned by list/get/update.
// Note: the full key value is NEVER returned here — only at creation time
// via APIKeyCreateResponse. Use KeyPrefix to identify a key in audit logs.
type APIKeyResponse struct {
	ID         uint       `json:"id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"key_prefix"`
	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	Scopes     []string   `json:"scopes"`

	// Registry-key summary; nil on personal keys.
	RegistryOrgSlug     *string  `json:"registry_org_slug,omitempty"`
	RegistryRepoName    *string  `json:"registry_repo_name,omitempty"`
	RegistryPermissions []string `json:"registry_permissions,omitempty"`
}

// APIKeyCreateResponse is returned by POST /api/v1/api-keys. The full
// Key value (kumo_sk_…) is shown exactly once — store it immediately,
// it cannot be retrieved later.
type APIKeyCreateResponse struct {
	APIKeyResponse
	Key string `json:"key"`
}
