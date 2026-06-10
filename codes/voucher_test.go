package codes

import "testing"

// Wire codes are a public contract. These assert the exact string values
// so an accidental rename is caught here before release.
func TestVoucherCodeValues(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"VoucherInternalError", VoucherInternalError, "VOUCHER_INTERNAL_ERROR"},
		{"VoucherNotFound", VoucherNotFound, "VOUCHER_NOT_FOUND"},
		{"VoucherExpired", VoucherExpired, "VOUCHER_EXPIRED"},
		{"VoucherExhausted", VoucherExhausted, "VOUCHER_EXHAUSTED"},
		{"VoucherNotActive", VoucherNotActive, "VOUCHER_NOT_ACTIVE"},
		{"VoucherAlreadyRedeemed", VoucherAlreadyRedeemed, "VOUCHER_ALREADY_REDEEMED"},
		{"VoucherWhitelistDenied", VoucherWhitelistDenied, "VOUCHER_WHITELIST_DENIED"},
		{"VoucherNewUsersOnly", VoucherNewUsersOnly, "VOUCHER_NEW_USERS_ONLY"},
		{"VoucherPoolInsufficient", VoucherPoolInsufficient, "VOUCHER_POOL_INSUFFICIENT"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}
