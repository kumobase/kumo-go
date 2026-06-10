package codes

// Quota-module wire codes returned by the quota and admin-quota endpoints.
// Branch on these constants rather than the human-readable Message.
const (
	// QuotaBelowUsage — an admin attempted to set a user's quota below their
	// current usage for that resource. Returned with HTTP 409.
	QuotaBelowUsage = "QUOTA_BELOW_USAGE"

	// QuotaInternalError — the catch-all 500 code for quota endpoints when no
	// typed sentinel matched (e.g. failed to update or read quota/usage).
	QuotaInternalError = "QUOTA_INTERNAL_ERROR"
)
