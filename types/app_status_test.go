package types

import "testing"

// TestAppStatusValues pins the wire string values of the computed app_status
// field. These are the contract Terraform/CLI consumers branch on; changing a
// value silently breaks them, so treat a failure here as a breaking change.
func TestAppStatusValues(t *testing.T) {
	cases := map[string]string{
		AppStatusRunning:    "running",
		AppStatusStopped:    "stopped",
		AppStatusBuilding:   "building",
		AppStatusDeploying:  "deploying",
		AppStatusDegraded:   "degraded",
		AppStatusCrashing:   "crashing",
		AppStatusImageError: "image_error",
		AppStatusFailed:     "failed",
		AppStatusSuspended:  "suspended",
		AppStatusUnknown:    "unknown",
	}
	for got, want := range cases {
		if got != want {
			t.Errorf("app status constant = %q, want %q", got, want)
		}
	}
}
