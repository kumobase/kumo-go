package codes

import "testing"

// Wire codes are a public contract (terraform-provider-kumo, kumo-cli). These
// assert the exact string values so an accidental rename is caught here before
// release.
func TestAppMetricsCodeValues(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"MetricsBackendUnavailable", MetricsBackendUnavailable, "METRICS_BACKEND_UNAVAILABLE"},
		{"InvalidTimeRange", InvalidTimeRange, "INVALID_TIME_RANGE"},
		{"AppMetricsInternalError", AppMetricsInternalError, "APP_METRICS_INTERNAL_ERROR"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}
