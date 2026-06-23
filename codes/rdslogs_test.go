package codes

import "testing"

// Wire codes are a public contract (terraform-provider-kumo, kumo-cli). Assert
// the exact string value so an accidental rename is caught here before release.
func TestRDSLogsCodeValues(t *testing.T) {
	if RDSLogsInternalError != "RDS_LOGS_INTERNAL_ERROR" {
		t.Errorf("RDSLogsInternalError = %q, want %q", RDSLogsInternalError, "RDS_LOGS_INTERNAL_ERROR")
	}
}
