package codes

import "testing"

// Wire codes are a public contract. These assert the exact string values so
// an accidental rename is caught here before release.
func TestAPIKeyCodeValues(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"APIKeyAdminForbidden", APIKeyAdminForbidden, "API_KEY_ADMIN_FORBIDDEN"},
		{"APIKeySessionRequired", APIKeySessionRequired, "API_KEY_SESSION_REQUIRED"},
		{"APIKeyForbiddenScope", APIKeyForbiddenScope, "API_KEY_FORBIDDEN_SCOPE"},
		{"APIKeyUnauthorized", APIKeyUnauthorized, "API_KEY_UNAUTHORIZED"},
		{"APIKeyNotFound", APIKeyNotFound, "API_KEY_NOT_FOUND"},
		{"APIKeyValidation", APIKeyValidation, "API_KEY_VALIDATION_FAILED"},
		{"APIKeyBadRequest", APIKeyBadRequest, "API_KEY_BAD_REQUEST"},
		{"APIKeyInternalError", APIKeyInternalError, "API_KEY_INTERNAL_ERROR"},
		{"RegistryKeyHTTPForbidden", RegistryKeyHTTPForbidden, "REGISTRY_KEY_HTTP_FORBIDDEN"},
		{"RegistryKeyInvalidScope", RegistryKeyInvalidScope, "REGISTRY_KEY_INVALID_SCOPE"},
		{"RegistryKeyRepoPinDisabled", RegistryKeyRepoPinDisabled, "REGISTRY_KEY_REPO_PIN_DISABLED"},
		// Unified-grants additions.
		{"APIKeyUnknownDomain", APIKeyUnknownDomain, "API_KEY_UNKNOWN_DOMAIN"},
		{"APIKeyInvalidGrant", APIKeyInvalidGrant, "API_KEY_INVALID_GRANT"},
		{"APIKeyInvalidCondition", APIKeyInvalidCondition, "API_KEY_INVALID_CONDITION"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}
