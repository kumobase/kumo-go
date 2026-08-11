package types

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestJobsRoundTrip(t *testing.T) {
	now := time.Date(2026, 5, 31, 14, 0, 0, 0, time.UTC)
	billed := "0.0123"
	appID := uint(42)
	appName := "my-app"
	exit := 0
	dur := int64(5000)

	roundTrip(t, "CreateJobRequest/standalone", CreateJobRequest{
		Name:        "nightly-cleanup",
		Kind:        JobKindStandalone,
		PricingSlug: "small",
		Image:       "alpine:3.20",
		Command:     []string{"sh", "-c", "echo hi"},
		Env:         []EnvironmentVariable{{Key: "FOO", Value: "bar"}},
		SecretRefs: []JobSecretRef{
			{SecretName: "db-creds", SourceFrom: "PASSWORD", MountTo: "DB_PASSWORD"},
		},
		Schedule:              "0 2 * * *",
		Timezone:              "Asia/Jakarta",
		ConcurrencyPolicy:     JobConcurrencyForbid,
		ActiveDeadlineSeconds: 900,
		BackoffLimit:          1,
	})

	roundTrip(t, "CreateJobRequest/app_attached", CreateJobRequest{
		Name:        "app-cron",
		Kind:        JobKindAppAttached,
		PricingSlug: "small",
		AppID:       &appID,
		Command:     []string{"./worker", "send-emails"},
		Schedule:    "*/10 * * * *",
	})

	roundTrip(t, "JobResponse", JobResponse{
		ID:                    1,
		Name:                  "nightly-cleanup",
		Kind:                  JobKindStandalone,
		Image:                 "alpine:3.20",
		Command:               []string{"sh", "-c", "echo hi"},
		Schedule:              "0 2 * * *",
		Timezone:              "Asia/Jakarta",
		ConcurrencyPolicy:     JobConcurrencyForbid,
		ActiveDeadlineSeconds: 900,
		BackoffLimit:          1,
		ResourceTemplate: JobResourceTemplate{
			Slug: "small", Name: "Small",
			CPUvCPU: "0.25", MemoryMB: 256,
			CPURequestvCPU: "0.25", CPULimitvCPU: "0.50",
			MemoryRequestMB: 256, MemoryLimitMB: 256,
			PricePerHour: "0.005",
		},
		Suspended:        false,
		DeploymentStatus: JobDeploymentStatusActive,
		LastExecutionAt:  &now,
		NextRunTimes:     []time.Time{now.Add(time.Hour)},
		CreatedAt:        now,
		UpdatedAt:        now,
	})

	roundTrip(t, "JobResponse/app_attached", JobResponse{
		ID:                    2,
		Name:                  "app-cron",
		Kind:                  JobKindAppAttached,
		AppID:                 &appID,
		AppName:               &appName,
		Schedule:              "*/10 * * * *",
		Timezone:              "Asia/Jakarta",
		ConcurrencyPolicy:     JobConcurrencyForbid,
		ActiveDeadlineSeconds: 900,
		BackoffLimit:          0,
		ResourceTemplate: JobResourceTemplate{
			Slug: "small", Name: "Small",
			CPUvCPU: "0.25", MemoryMB: 256,
			CPURequestvCPU: "0.25", CPULimitvCPU: "0.50",
			MemoryRequestMB: 256, MemoryLimitMB: 256,
			PricePerHour: "0.005",
		},
		DeploymentStatus: JobDeploymentStatusActive,
		CreatedAt:        now,
		UpdatedAt:        now,
	})

	roundTrip(t, "JobListItem", JobListItem{
		ID:               1,
		Name:             "nightly-cleanup",
		Kind:             JobKindStandalone,
		Schedule:         "0 2 * * *",
		Timezone:         "Asia/Jakarta",
		DeploymentStatus: JobDeploymentStatusActive,
		CreatedAt:        now,
		UpdatedAt:        now,
	})

	roundTrip(t, "ResponseJobAsync", ResponseJobAsync{
		ID:               1,
		Name:             "nightly-cleanup",
		DeploymentStatus: JobDeploymentStatusDeploying,
		OperationID:      "op_abc123",
		UpdatedAt:        now,
	})

	roundTrip(t, "JobExecution", JobExecution{
		ID:            10,
		JobID:         1,
		Trigger:       JobExecutionTriggerSchedule,
		K8sJobName:    "nightly-cleanup-28401300",
		Status:        JobExecutionStatusSucceeded,
		ExitCode:      &exit,
		PodStartedAt:  &now,
		PodFinishedAt: &now,
		DurationMS:    &dur,
		CPUvCPU:       "0.25",
		MemoryMB:      256,

		CPURequestvCPU:  "0.25",
		CPULimitvCPU:    "0.50",
		MemoryRequestMB: 256,
		MemoryLimitMB:   256,

		BilledAmount: &billed,
		CreatedAt:    now,
	})

	// Executions that predate the limit snapshot carry only the request. The
	// *Limit fields must drop out of the wire entirely rather than serialize
	// as "0" / "", which a client would read as a real zero ceiling.
	roundTrip(t, "JobExecution/legacy_no_limit_snapshot", JobExecution{
		ID:         12,
		JobID:      1,
		Trigger:    JobExecutionTriggerManual,
		K8sJobName: "nightly-cleanup-manual-a1b2c3d4",
		Status:     JobExecutionStatusRunning,
		CPUvCPU:    "0.25",
		MemoryMB:   256,
		CreatedAt:  now,
	})

	roundTrip(t, "RunJobResponse", RunJobResponse{
		ExecutionID: 11,
		Status:      JobExecutionStatusPending,
		OperationID: "op_run_xyz",
	})

	sched := "0 5 * * *"
	bo := 3
	roundTrip(t, "UpdateJobRequest", UpdateJobRequest{
		PricingSlug:           "medium",
		Schedule:              &sched,
		ConcurrencyPolicy:     JobConcurrencyForbid,
		ActiveDeadlineSeconds: 600,
		BackoffLimit:          &bo,
	})
}

// TestJobExecution_OmitsMissingLimitSnapshot pins the wire behaviour for
// executions recorded before job_executions grew its limit columns. Those rows
// read back NULL, and a client must be able to tell "no limit recorded" apart
// from "limit is zero" — a zero CPU limit means *unlimited* to Kubernetes, so
// serializing "cpu_limit_vcpu":"" or 0 here would be actively misleading.
// roundTrip alone cannot catch this: it only compares a re-marshal.
func TestJobExecution_OmitsMissingLimitSnapshot(t *testing.T) {
	body, err := json.Marshal(JobExecution{
		ID:       12,
		JobID:    1,
		CPUvCPU:  "0.12",
		MemoryMB: 256,
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, absent := range []string{"cpu_limit_vcpu", "memory_limit_mb", "cpu_request_vcpu", "memory_request_mb"} {
		if strings.Contains(string(body), absent) {
			t.Errorf("expected %q to be omitted for a legacy execution, got %s", absent, body)
		}
	}
	// The deprecated aliases stay on the wire for existing consumers.
	for _, present := range []string{`"cpu_vcpu":"0.12"`, `"memory_mb":256`} {
		if !strings.Contains(string(body), present) {
			t.Errorf("expected %s in payload, got %s", present, body)
		}
	}
}

// TestJobResourceTemplate_AlwaysEmitsBothNumbers is the counterpart for the
// plan projection: unlike the execution snapshot, these fields are not
// omitempty, because every pinned template version always has both numbers.
// A client rendering the picker can rely on them being present.
func TestJobResourceTemplate_AlwaysEmitsBothNumbers(t *testing.T) {
	body, err := json.Marshal(JobResourceTemplate{Slug: "nano", Name: "Nano"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, key := range []string{"cpu_request_vcpu", "cpu_limit_vcpu", "memory_request_mb", "memory_limit_mb", "cpu_vcpu", "memory_mb"} {
		if !strings.Contains(string(body), key) {
			t.Errorf("expected %q always present, got %s", key, body)
		}
	}
}
