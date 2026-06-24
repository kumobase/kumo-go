package types

import "time"

// RunnerJobState is the lifecycle state of one CI job's runner VM. Mirrors the
// chk_runner_jobs_state CHECK constraint on the server.
//
//	queued       — matched a spec, waiting for capacity in the admission queue
//	provisioning — admitted; a spot VM launch has been requested
//	booting       — VM up; JIT runner config minted, runner starting
//	registered   — the ephemeral runner registered with GitHub
//	running       — GitHub reports the job in_progress on our runner
//	completed    — job finished; VM terminated
//	failed       — provisioning failed terminally, or it waited past the queue
//	               deadline without capacity (conclusion no_capacity_timeout)
//	interrupted  — the spot VM was reclaimed mid-job (re-queued by GitHub)
//	cancelled    — cancelled before the VM became usable
//	orphaned     — reconciled by the reaper
type RunnerJobState string

const (
	RunnerJobStateQueued       RunnerJobState = "queued"
	RunnerJobStateProvisioning RunnerJobState = "provisioning"
	RunnerJobStateBooting      RunnerJobState = "booting"
	RunnerJobStateRegistered   RunnerJobState = "registered"
	RunnerJobStateRunning      RunnerJobState = "running"
	RunnerJobStateCompleted    RunnerJobState = "completed"
	RunnerJobStateFailed       RunnerJobState = "failed"
	RunnerJobStateInterrupted  RunnerJobState = "interrupted"
	RunnerJobStateCancelled    RunnerJobState = "cancelled"
	RunnerJobStateOrphaned     RunnerJobState = "orphaned"
)

// RunnerSpecResponse is a logical runner size, returned by GET
// /api/v1/runner-specs so users can discover which `kumo-*` labels they may put
// in `runs-on`. It is intentionally sanitized: the cloud backends (provider,
// region, instance types, capacity) that fulfill a spec are an internal,
// admin-only concern and are never exposed here.
type RunnerSpecResponse struct {
	Label       string `json:"label"`        // e.g. "kumo-2c-4g" — use in runs-on
	DisplayName string `json:"display_name"`
	CPU         int    `json:"cpu"`
	MemoryMB    int    `json:"memory_mb"`
}

// RunnerJobResponse is one CI job that ran (or is queued/running) on Kumo
// runners, returned by GET /api/v1/runner-jobs. This is the LOGICAL view: it
// reports the size requested and the job's status, but never which cloud,
// region, or instance handled it (that's an internal scheduling detail).
//
// State `queued` means the job is waiting for capacity in the admission queue.
// Conclusion mirrors GitHub's terminal conclusion (or a Kumo reason like
// no_capacity_timeout) and is empty until the job finishes.
type RunnerJobResponse struct {
	ID           uint           `json:"id"`
	SpecLabel    string         `json:"spec_label"` // the size requested (runs-on label)
	GithubJobID  int64          `json:"github_job_id"`
	RunID        int64          `json:"run_id"`
	RepoFullName string         `json:"repo_full_name"`
	State        RunnerJobState `json:"state"`
	Conclusion   string         `json:"conclusion,omitempty"`

	QueuedAt   time.Time  `json:"queued_at"`
	StartedAt  *time.Time `json:"started_at,omitempty"`  // when the runner registered
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

// RunnerUsageResponse is the aggregated usage summary returned by GET
// /api/v1/runner-usage over an optional from/to window (default: current
// billing period). Minutes is metered VM wall-clock; interrupted time is
// excluded. EstimatedCents is indicative, not invoiced.
type RunnerUsageResponse struct {
	WindowStart    time.Time `json:"window_start"`
	WindowEnd      time.Time `json:"window_end"`
	JobCount       int       `json:"job_count"`
	BilledMinutes  float64   `json:"billed_minutes"`
	EstimatedCents int64     `json:"estimated_cents"`
}
