package client

import (
	"context"
	"fmt"

	"github.com/kumobase/kumo-go/types"
)

// RunnersService backs the customer surface of the VM-backed CI-runner product:
// discovering available runner sizes and viewing your jobs' status/history.
//
// There is nothing to provision or configure — connect the Kumo GitHub App
// (SourceConnections) and put a `kumo-*` label from ListSpecs in your
// workflow's runs-on. Capacity and cloud placement are managed by Kumo.
type RunnersService struct {
	c *Client
}

// Runners returns the runners service.
func (c *Client) Runners() *RunnersService { return &RunnersService{c: c} }

// ListSpecs returns the available runner sizes (the `kumo-*` labels you can use
// in runs-on). Paginated.
func (s *RunnersService) ListSpecs(ctx context.Context, opts ...ListOption) ([]types.RunnerSpecResponse, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.RunnerSpecResponse
	meta, err := s.c.doList(ctx, "GET", withQuery("/api/v1/runner-specs", q), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, meta, nil
}

// ListJobs returns your CI jobs that ran (or are queued/running) on Kumo
// runners, newest first. Filter with WithExtraQuery("state", …) using a
// types.RunnerJobState value (e.g. "queued", "running", "completed").
func (s *RunnersService) ListJobs(ctx context.Context, opts ...ListOption) ([]types.RunnerJobResponse, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.RunnerJobResponse
	meta, err := s.c.doList(ctx, "GET", withQuery("/api/v1/runner-jobs", q), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, meta, nil
}

// GetJob returns one of your runner jobs by id.
func (s *RunnersService) GetJob(ctx context.Context, id uint) (*types.RunnerJobResponse, error) {
	var out types.RunnerJobResponse
	_, _, err := s.c.do(ctx, "GET", fmt.Sprintf("/api/v1/runner-jobs/%d", id), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Usage returns your aggregated runner usage over an optional window. Pass
// WithExtraQuery("from"/"to", rfc3339) to bound it; defaults to the current
// billing period.
func (s *RunnersService) Usage(ctx context.Context, opts ...ListOption) (*types.RunnerUsageResponse, error) {
	q := resolveListOpts(opts)
	var out types.RunnerUsageResponse
	_, _, err := s.c.do(ctx, "GET", withQuery("/api/v1/runner-usage", q), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
