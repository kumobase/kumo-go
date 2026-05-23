package codes

// Organization-module wire codes returned by
// /api/v1/registry/organizations/* endpoints.
const (
	OrgNotFound             = "ORG_NOT_FOUND"
	OrgSlugTaken            = "ORG_SLUG_TAKEN"
	OrgSlugInvalid          = "ORG_SLUG_INVALID"
	OrgSlugReserved         = "ORG_SLUG_RESERVED"
	OrgSlugImmutable        = "ORG_SLUG_IMMUTABLE"
	OrgNotMember            = "ORG_NOT_MEMBER"
	OrgMaxOrganizationsReached = "ORG_MAX_REACHED"
	OrgHasRepos             = "ORG_HAS_REPOS"

	OrgInvalidRequestBody = "INVALID_REQUEST_BODY"
	OrgValidationError    = "VALIDATION_ERROR"
	OrgUnauthorized       = "UNAUTHORIZED"
	OrgInternalError      = "INTERNAL_ERROR"
)
