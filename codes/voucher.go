package codes

// Voucher-module wire codes. Only the read surface (GET /vouchers/history) is
// api-key-callable; redemption is sessionOnly, so its codes are not exposed.
const (
	// VoucherInternalError is the catch-all 500 code for voucher read endpoints
	// (e.g. GET /vouchers/history).
	VoucherInternalError = "VOUCHER_INTERNAL_ERROR"
)
