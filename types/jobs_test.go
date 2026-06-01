package types

import (
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
			Slug: "small", Name: "Small", CPUvCPU: "0.25", MemoryMB: 256, PricePerHour: "0.005",
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
			Slug: "small", Name: "Small", CPUvCPU: "0.25", MemoryMB: 256, PricePerHour: "0.005",
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
		BilledAmount:  &billed,
		CreatedAt:     now,
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
