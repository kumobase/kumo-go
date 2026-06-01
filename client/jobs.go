package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/kumobase/kumo-go/types"
)

// JobsService backs /api/v1/jobs/*. Create/Update/Delete/RunNow are
// asynchronous (return 202 + operation_id); the rest are synchronous.
type JobsService struct {
	c *Client
}

// Jobs returns the jobs service.
func (c *Client) Jobs() *JobsService { return &JobsService{c: c} }

// Get fetches a single job by numeric id. Returns the response and the ETag
// for use with IfMatch on a subsequent Update.
func (s *JobsService) Get(ctx context.Context, id uint) (*types.JobResponse, string, error) {
	var out types.JobResponse
	etag, _, err := s.c.do(ctx, "GET", fmt.Sprintf("/api/v1/jobs/%d", id), nil, nil, &out)
	if err != nil {
		return nil, "", err
	}
	return &out, etag, nil
}

// GetByName fetches a job by name (the server resolves a non-numeric path
// segment as a name). Returns 409 AMBIGUOUS_NAME if more than one job
// matches in the caller's scope — fall back to Get(ctx, id) to disambiguate.
func (s *JobsService) GetByName(ctx context.Context, name string) (*types.JobResponse, string, error) {
	var out types.JobResponse
	etag, _, err := s.c.do(ctx, "GET", "/api/v1/jobs/"+url.PathEscape(name), nil, nil, &out)
	if err != nil {
		return nil, "", err
	}
	return &out, etag, nil
}

// List returns paginated jobs (list-item shape, heavy fields omitted).
// Endpoint-specific filters available via WithExtraQuery: "kind"
// (app_attached|standalone), "app_id" (numeric), "suspended" (true|false).
// Passing mutually-exclusive filters returns 400 INVALID_FILTER_COMBINATION.
func (s *JobsService) List(ctx context.Context, opts ...ListOption) ([]types.JobListItem, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.JobListItem
	meta, err := s.c.doList(ctx, "GET", withQuery("/api/v1/jobs", q), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, meta, nil
}

// Create persists a new job. Honors Idempotency-Key via WriteOption. Returns
// 202 + ResponseJobAsync; poll the operation via OperationID until the
// underlying CronJob is applied.
func (s *JobsService) Create(ctx context.Context, req *types.CreateJobRequest, opts ...WriteOption) (*types.ResponseJobAsync, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.ResponseJobAsync
	_, _, err = s.c.do(ctx, "POST", "/api/v1/jobs", req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Update mutates an existing job. Kind is immutable. Pass IfMatch(etag)
// for optimistic concurrency. Returns 202 + ResponseJobAsync.
func (s *JobsService) Update(ctx context.Context, id uint, req *types.UpdateJobRequest, opts ...WriteOption) (*types.ResponseJobAsync, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.ResponseJobAsync
	_, _, err = s.c.do(ctx, "PATCH", fmt.Sprintf("/api/v1/jobs/%d", id), req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete removes a job. Async: returns 202 + ResponseJobAsync; the
// CronJob and any running Jobs are torn down by the deployment consumer.
// In-flight executions are closed as "failed" and metered to that point.
func (s *JobsService) Delete(ctx context.Context, id uint, opts ...WriteOption) (*types.ResponseJobAsync, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.ResponseJobAsync
	_, _, err = s.c.do(ctx, "DELETE", fmt.Sprintf("/api/v1/jobs/%d", id), nil, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// RunNow triggers an ad-hoc execution of the job, independent of its
// schedule. Honors Idempotency-Key. Concurrent with any running scheduled
// execution (manual is explicit user intent).
func (s *JobsService) RunNow(ctx context.Context, id uint, opts ...WriteOption) (*types.RunJobResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.RunJobResponse
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/jobs/%d/run", id), nil, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Suspend sets the job's suspended flag to true (CronJob.spec.suspend=true).
// In-flight pods continue to completion and bill normally. Idempotent at
// the application level: calling on an already-suspended job returns 409
// JOB_ALREADY_SUSPENDED, which most clients can treat as a no-op.
func (s *JobsService) Suspend(ctx context.Context, id uint, opts ...WriteOption) (*types.JobResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.JobResponse
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/jobs/%d/suspend", id), nil, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Resume clears the suspended flag.
func (s *JobsService) Resume(ctx context.Context, id uint, opts ...WriteOption) (*types.JobResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.JobResponse
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/jobs/%d/resume", id), nil, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ListExecutions returns paginated execution history for a job. Filters via
// WithExtraQuery: "status" (one of pending|running|succeeded|failed|timeout),
// "from"/"to" (RFC3339 timestamps).
func (s *JobsService) ListExecutions(ctx context.Context, jobID uint, opts ...ListOption) ([]types.JobExecution, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.JobExecution
	meta, err := s.c.doList(ctx, "GET", withQuery(fmt.Sprintf("/api/v1/jobs/%d/executions", jobID), q), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, meta, nil
}

// GetExecution fetches a single execution by id.
func (s *JobsService) GetExecution(ctx context.Context, jobID, executionID uint) (*types.JobExecution, error) {
	var out types.JobExecution
	_, _, err := s.c.do(ctx, "GET",
		fmt.Sprintf("/api/v1/jobs/%d/executions/%d", jobID, executionID), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
