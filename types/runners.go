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
// /api/v1/runners/plans (the public price catalog) so users can discover which
// `kumo-*` labels they may put in `runs-on` and what each costs. It is
// intentionally sanitized: the cloud
// backends (provider, region, instance types, capacity) that fulfill a spec are
// an internal, admin-only concern and are never exposed here.
//
// PricePerMinute is the customer rate in IDR per VM-minute, a decimal string
// (e.g. "12.5000") — mirrors the jobs PricePerHour convention; never float.
// Currency is the ISO code for that rate (always "IDR" today).
type RunnerSpecResponse struct {
	Label          string `json:"label"`        // e.g. "kumo-2c-4g" — use in runs-on
	DisplayName    string `json:"display_name"`
	CPU            int    `json:"cpu"`
	MemoryMB       int    `json:"memory_mb"`
	PricePerMinute string `json:"price_per_minute"` // decimal string, IDR/min
	Currency       string `json:"currency"`         // "IDR"
	// Aliases are additional `runs-on` labels that resolve to this same size
	// (e.g. "kumo-ubuntu-latest" → "kumo-ubuntu-24.04"). Empty when none.
	Aliases []string `json:"aliases,omitempty"`
}

// RunnerJobResponse is one CI job that ran (or is queued/running) on Kumo
// runners, returned by GET /api/v1/runner-jobs. This is the LOGICAL view: it
// reports the size requested and the job's status, but never which cloud,
// region, or instance handled it (that's an internal scheduling detail).
//
// Provider is the CI host the job came from ("github" | "gitlab") — a
// user-facing discriminator so a unified job list can branch on origin. It is
// NOT the cloud provider (that stays internal).
//
// The provider-specific identifier fields are a tagged union: GitHub jobs carry
// GithubJobID/RunID/RepoFullName; GitLab jobs carry the GitLab* fields. The
// unused side is zero/omitted. WebURL deep-links to the job on the provider.
//
// State `queued` means the job is waiting for capacity in the admission queue.
// Conclusion mirrors the provider's terminal conclusion (or a Kumo reason like
// no_capacity_timeout) and is empty until the job finishes.
type RunnerJobResponse struct {
	ID        uint           `json:"id"`
	Provider  SourceProvider `json:"provider"`   // "github" | "gitlab"
	SpecLabel string         `json:"spec_label"` // the size requested (runs-on / tag label)

	// GitHub-specific. Always present for github jobs (back-compat: unconditional).
	GithubJobID  int64  `json:"github_job_id"`
	RunID        int64  `json:"run_id"`
	RepoFullName string `json:"repo_full_name"`

	// GitLab-specific. Omitted for github jobs.
	GitLabJobID      *int64 `json:"gitlab_job_id,omitempty"`
	GitLabProjectID  *int64 `json:"gitlab_project_id,omitempty"`
	GitLabPipelineID *int64 `json:"gitlab_pipeline_id,omitempty"`

	// WebURL is a deep-link to the job on the provider (e.g. the GitLab job
	// page). Omitted when unknown.
	WebURL string `json:"web_url,omitempty"`

	State      RunnerJobState `json:"state"`
	Conclusion string         `json:"conclusion,omitempty"`

	QueuedAt   time.Time  `json:"queued_at"`
	StartedAt  *time.Time `json:"started_at,omitempty"`  // when the runner registered
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}
