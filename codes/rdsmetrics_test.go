package codes

import "testing"

// Wire codes are a public contract (terraform-provider-kumo, kumo-cli). Assert
// the exact string value so an accidental rename is caught here before release.
func TestRDSMetricsCodeValues(t *testing.T) {
	if RDSMetricsInternalError != "RDS_METRICS_INTERNAL_ERROR" {
		t.Errorf("RDSMetricsInternalError = %q, want %q", RDSMetricsInternalError, "RDS_METRICS_INTERNAL_ERROR")
	}
}
