package codes

// API-key access-policy codes. These can be returned by any /api/v1/* route
// depending on the caller's auth and the route's guard, not just by
// /api-keys endpoints.
const (
	// APIKeyAdminForbidden — an API key tried to call an /admin/* route.
	// API keys are resource credentials and never gain admin access, even
	// when the owning user is an admin. Use the dashboard or a CLI JWT
	// bearer instead.
	APIKeyAdminForbidden = "API_KEY_ADMIN_FORBIDDEN"

	// APIKeySessionRequired — an API key tried to call an account-takeover
	// or financial endpoint (password, /billing/topup, /vouchers/redeem,
	// /api-keys/*, /profile mutations, /tickets/*, /referral/*). Use
	// dashboard cookie or CLI JWT bearer.
	APIKeySessionRequired = "API_KEY_SESSION_REQUIRED"

	// APIKeyForbiddenScope — the key lacks the required scope (read/write)
	// for the requested action.
	APIKeyForbiddenScope = "API_KEY_FORBIDDEN_SCOPE"

	// APIKeyUnauthorized — the supplied key is invalid, expired, revoked,
	// or otherwise not authenticatable.
	APIKeyUnauthorized = "API_KEY_UNAUTHORIZED"

	// APIKeyNotFound — referenced API key id does not exist (CRUD endpoints).
	APIKeyNotFound = "API_KEY_NOT_FOUND"

	// APIKeyValidation — request body for API-key CRUD failed validation.
	APIKeyValidation = "API_KEY_VALIDATION_FAILED"

	// APIKeyBadRequest — generic 400 from the API-key CRUD surface (invalid
	// id, malformed body, etc.).
	APIKeyBadRequest = "API_KEY_BAD_REQUEST"

	// APIKeyInternalError — unmapped server error during API-key operations.
	APIKeyInternalError = "API_KEY_INTERNAL_ERROR"

	// Registry-key codes. Registry keys (those with a registry_scope) are
	// rejected on /api/v1/* (REGISTRY_KEY_HTTP_FORBIDDEN) and only valid at
	// /v2/token.

	// RegistryKeyHTTPForbidden — a registry-scoped key tried to call the
	// control-plane API.
	RegistryKeyHTTPForbidden = "REGISTRY_KEY_HTTP_FORBIDDEN"

	// RegistryKeyInvalidScope — registry_scope.permissions contained a
	// value outside {pull, push, delete}.
	RegistryKeyInvalidScope = "REGISTRY_KEY_INVALID_SCOPE"

	// RegistryKeyRepoPinDisabled — repo-level pinning was requested but is
	// disabled on this deployment.
	RegistryKeyRepoPinDisabled = "REGISTRY_KEY_REPO_PIN_DISABLED"
)
