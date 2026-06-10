package codes

// Voucher-module wire codes. Only the read surface (GET /vouchers/history) is
// api-key-callable; redemption is sessionOnly, so its codes are not exposed.
const (
	// VoucherInternalError is the catch-all 500 code for voucher read endpoints
	// (e.g. GET /vouchers/history).
	VoucherInternalError = "VOUCHER_INTERNAL_ERROR"
)

// Voucher redemption wire codes — one per typed sentinel branched in the
// server's POST /vouchers/redeem switch (modules/voucher/handler.go). The
// redeem endpoint is sessionOnly (not api-key-callable), but the codes are a
// stable contract for the web/CLI client so it can branch on the precise
// failure rather than the human-readable Message.
const (
	// VoucherNotFound — no voucher matched the supplied code (HTTP 404).
	VoucherNotFound = "VOUCHER_NOT_FOUND"

	// VoucherExpired — the voucher's expiry date has passed (HTTP 410).
	VoucherExpired = "VOUCHER_EXPIRED"

	// VoucherExhausted — the voucher has reached its redemption limit (HTTP 410).
	VoucherExhausted = "VOUCHER_EXHAUSTED"

	// VoucherNotActive — the voucher exists but is not currently active, e.g.
	// disabled or outside its active window (HTTP 409).
	VoucherNotActive = "VOUCHER_NOT_ACTIVE"

	// VoucherAlreadyRedeemed — the caller has already redeemed this voucher
	// (HTTP 403).
	VoucherAlreadyRedeemed = "VOUCHER_ALREADY_REDEEMED"

	// VoucherWhitelistDenied — the caller is not on the voucher's redemption
	// whitelist (HTTP 403).
	VoucherWhitelistDenied = "VOUCHER_WHITELIST_DENIED"

	// VoucherNewUsersOnly — the voucher is restricted to new users and the
	// caller does not qualify (HTTP 403).
	VoucherNewUsersOnly = "VOUCHER_NEW_USERS_ONLY"

	// VoucherPoolInsufficient — the shared voucher pool has insufficient balance
	// to fund this redemption; surfaced as a temporary outage (HTTP 503).
	VoucherPoolInsufficient = "VOUCHER_POOL_INSUFFICIENT"
)
