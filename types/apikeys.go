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
//
// Grants is the unified permission model that supersedes the Scopes /
// RegistryScope split: one key carries a list of per-product grants, each
// naming a Domain, its Actions, and (for org-scoped domains) the Orgs it
// applies to. When Grants is set, Scopes and RegistryScope are ignored.
// Scopes / RegistryScope remain supported for backward compatibility and
// are translated into equivalent grants server-side. Conditions is a
// forward-compatible, token-level constraint block (e.g. IP allowlist);
// its fields are validated but not all are enforced yet.
type CreateAPIKeyRequest struct {
	Name          string              `json:"name"` // 1..100 chars
	ExpiresInDays *int                `json:"expires_in_days,omitempty"`
	Scopes        []string            `json:"scopes,omitempty"`
	RegistryScope *RegistryScopeInput `json:"registry_scope,omitempty"`
	Grants        []Grant             `json:"grants,omitempty"`
	Conditions    *TokenConditions    `json:"conditions,omitempty"`
}

// Grant is one entry in the unified permission model: it authorizes a set
// of Actions within a single product Domain, optionally restricted to
// specific Orgs.
//
//   - Domain is a product area, e.g. "control_plane" or "registry".
//   - Actions are domain-specific verbs: control_plane → {read, write};
//     registry → {pull, push, delete} (push implies pull).
//   - Orgs are organization slugs. It is only meaningful for org-scoped
//     domains (e.g. registry); an empty list means "every organization the
//     owning user belongs to". Listing several orgs grants the actions in
//     each of them (membership is still enforced per request).
type Grant struct {
	Domain  string   `json:"domain"`
	Actions []string `json:"actions"`
	Orgs    []string `json:"orgs,omitempty"`
}

// TokenConditions holds token-level constraints applied to every request
// made with the key. It is intentionally forward-compatible: new fields may
// be added over time. IPAllowlist is reserved — accepted and validated, but
// not yet enforced.
type TokenConditions struct {
	IPAllowlist []string `json:"ip_allowlist,omitempty"`
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

	// Grants is the unified permission view of the key (org slugs resolved
	// for display). Present on keys created via the unified model; omitted
	// for legacy keys that still carry only Scopes / registry_* fields.
	// Conditions echoes any token-level constraints attached at creation.
	Grants     []Grant          `json:"grants,omitempty"`
	Conditions *TokenConditions `json:"conditions,omitempty"`
}

// APIKeyCreateResponse is returned by POST /api/v1/api-keys. The full
// Key value (kumo_sk_…) is shown exactly once — store it immediately,
// it cannot be retrieved later.
type APIKeyCreateResponse struct {
	APIKeyResponse
	Key string `json:"key"`
}
