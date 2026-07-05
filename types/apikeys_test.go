package types

import "testing"

func TestAPIKeyRoundTrip(t *testing.T) {
	// Legacy create shapes still round-trip unchanged.
	roundTrip(t, "CreateAPIKeyRequest/scopes", CreateAPIKeyRequest{
		Name:   "ci",
		Scopes: []string{"read", "write"},
	})
	roundTrip(t, "CreateAPIKeyRequest/registry_scope", CreateAPIKeyRequest{
		Name:          "docker-push",
		RegistryScope: &RegistryScopeInput{OrgSlug: "acme", Permissions: []string{"pull", "push"}},
	})

	// Unified grants shape, including a multi-org registry grant and the
	// forward-compatible conditions block.
	roundTrip(t, "CreateAPIKeyRequest/grants", CreateAPIKeyRequest{
		Name: "unified",
		Grants: []Grant{
			{Domain: "control_plane", Actions: []string{"read", "write"}},
			{Domain: "registry", Actions: []string{"pull", "push"}, Orgs: []string{"acme", "beta"}},
		},
		Conditions: &TokenConditions{IPAllowlist: []string{"203.0.113.0/24"}},
	})

	roundTrip(t, "APIKeyResponse/grants", APIKeyResponse{
		ID:        7,
		Name:      "unified",
		KeyPrefix: "kumo_sk_live_",
		Scopes:    []string{},
		Grants: []Grant{
			{Domain: "registry", Actions: []string{"push"}, Orgs: []string{"acme"}},
		},
		Conditions: &TokenConditions{IPAllowlist: []string{"10.0.0.0/8"}},
	})
}
