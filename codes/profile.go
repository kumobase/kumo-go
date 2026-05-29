package codes

// Profile-module wire codes returned by /api/v1/profile/* read endpoints.
// Mutation endpoints (PATCH /profile, password management) are sessionOnly and
// not part of the api-key surface, so their codes are not exposed here.
const (
	// ProfileInternalError is the catch-all 500 code for profile read endpoints
	// (e.g. GET /profile/balance).
	ProfileInternalError = "PROFILE_INTERNAL_ERROR"
)
