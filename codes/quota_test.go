package codes

import "testing"

// Wire codes are a public contract. These assert the exact string values
// so an accidental rename is caught here before release.
func TestQuotaCodeValues(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"QuotaBelowUsage", QuotaBelowUsage, "QUOTA_BELOW_USAGE"},
		{"QuotaInternalError", QuotaInternalError, "QUOTA_INTERNAL_ERROR"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}
