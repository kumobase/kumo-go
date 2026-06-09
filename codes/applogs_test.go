package codes

import "testing"

// Wire codes are a public contract (terraform-provider-kumo, kumo-cli). These
// assert the exact string values so an accidental rename is caught here before
// release.
func TestAppLogsCodeValues(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"LogsBackendUnavailable", LogsBackendUnavailable, "LOGS_BACKEND_UNAVAILABLE"},
		{"InvalidLogFilter", InvalidLogFilter, "INVALID_LOG_FILTER"},
		{"AppLogsInternalError", AppLogsInternalError, "APP_LOGS_INTERNAL_ERROR"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}
