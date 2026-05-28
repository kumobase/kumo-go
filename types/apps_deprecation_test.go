package types

import (
	"encoding/json"
	"testing"
	"time"
)

// TestAppByIdResponse_DeprecatedFieldsPinning pins the deprecation contract:
// the server populates BOTH the new instance/autoscaling fields AND the
// old pod-named fields, and the JSON keys for both appear on the wire. If
// this test fires, either:
//   - someone removed a deprecated field before the announced minor (breaks
//     existing callers pinned to today's wire); or
//   - someone added a new deprecation without also pinning it here.
//
// Remove the relevant assertions only as part of the same commit that drops
// the field, and bump SDKVersion to the planned removal minor.
func TestAppByIdResponse_DeprecatedFieldsPinning(t *testing.T) {
	resp := AppByIdResponse{
		Id:               1,
		CreateAppRequest: CreateAppRequest{BaseCreateApp: BaseCreateApp{Name: "x", Image: "nginx", Port: 80, Replicas: 2}, PricingSlug: "kumo.nano"},
		TotalInstances:   3, PendingInstances: 1, RunningInstances: 2, FailedInstances: 0,
		TotalPods: 3, PendingPods: 1, RunningPods: 2, FailedPods: 0,
		HasFailure: true, HasReplicaFailure: true,
		AutoscalingStatus: &AutoscalingStatus{CurrentReplicas: 2, DesiredReplicas: 3, MinReplicas: 1, MaxReplicas: 5},
		HPAStatus:         &AutoscalingStatus{CurrentReplicas: 2, DesiredReplicas: 3, MinReplicas: 1, MaxReplicas: 5},
		InternalDNS:       "x.ns",
		CreatedAt:         time.Unix(0, 0).UTC(), UpdatedAt: time.Unix(0, 0).UTC(),
	}

	body, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	// Pairs of (new, old) keys that MUST both appear during the deprecation window.
	pairs := [][2]string{
		{"total_instances", "total_pods"},
		{"pending_instances", "pending_pods"},
		{"running_instances", "running_pods"},
		{"failed_instances", "failed_pods"},
		{"has_failure", "has_replica_failure"},
		{"autoscaling_status", "hpa_status"},
	}
	for _, p := range pairs {
		newKey, oldKey := p[0], p[1]
		newVal, hasNew := raw[newKey]
		oldVal, hasOld := raw[oldKey]
		if !hasNew {
			t.Errorf("new key %q missing from JSON: %s", newKey, body)
		}
		if !hasOld {
			t.Errorf("deprecated key %q missing from JSON: %s", oldKey, body)
		}
		if hasNew && hasOld && string(newVal) != string(oldVal) {
			t.Errorf("%s vs %s diverged: %s vs %s", newKey, oldKey, newVal, oldVal)
		}
	}
}

// TestHPAStatusInfoAlias confirms HPAStatusInfo is still usable as a name for
// AutoscalingStatus, so callers pinned to the old type identifier keep
// compiling during the deprecation window.
func TestHPAStatusInfoAlias(t *testing.T) {
	var _ HPAStatusInfo = AutoscalingStatus{}
	var _ AutoscalingStatus = HPAStatusInfo{}
}
