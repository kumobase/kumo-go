package client

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/kumobase/kumo-go/types"
)

// RDSService backs /api/v1/rds/* — Kumo's managed relational database service
// (PostgreSQL at launch). Reads (flavors, instances) are synchronous.
// Mutations (create, scale, resize, delete) are asynchronous: each returns an
// operation_id, and the instance moves through transient Status values until
// the operation settles. The *AndWait helpers compose a mutation with
// PollResource until Status reaches a terminal state.
type RDSService struct {
	c *Client
}

// RDS returns the RDS service.
func (c *Client) RDS() *RDSService { return &RDSService{c: c} }

// ── Reads ──────────────────────────────────────────────────────────

// ListPlans returns the public database plan (instance class) catalogue
// (sanitised — no internal cost/margin). Filter by engine with
// WithExtraQuery("engine", "postgresql").
func (s *RDSService) ListPlans(ctx context.Context, opts ...ListOption) ([]types.PublicRDSPlanResponse, error) {
	q := resolveListOpts(opts)
	var out []types.PublicRDSPlanResponse
	_, _, err := s.c.do(ctx, "GET", withQuery("/api/v1/rds/plans", q), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// List returns the user's database instances, paginated. Filter via
// WithExtraQuery: "status", "engine".
func (s *RDSService) List(ctx context.Context, opts ...ListOption) ([]types.RDSInstanceResponse, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.RDSInstanceResponse
	meta, err := s.c.doList(ctx, "GET", withQuery("/api/v1/rds", q), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, meta, nil
}

// Get fetches a single database instance by id.
func (s *RDSService) Get(ctx context.Context, id uint) (*types.RDSInstanceResponse, error) {
	var out types.RDSInstanceResponse
	_, _, err := s.c.do(ctx, "GET", fmt.Sprintf("/api/v1/rds/%d", id), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetByName fetches a single database instance by its name. The server
// resolves a non-numeric path segment as a name. Returns 404 if no instance in
// the caller's scope matches.
func (s *RDSService) GetByName(ctx context.Context, name string) (*types.RDSInstanceResponse, error) {
	var out types.RDSInstanceResponse
	_, _, err := s.c.do(ctx, "GET", "/api/v1/rds/"+url.PathEscape(name), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetConnection returns connection details (incl. the live password from the
// credentials secret). Returns 409 RDS_INSTANCE_NOT_READY if the instance is
// not yet ready.
func (s *RDSService) GetConnection(ctx context.Context, id uint) (*types.RDSConnectionResponse, error) {
	var out types.RDSConnectionResponse
	_, _, err := s.c.do(ctx, "GET", fmt.Sprintf("/api/v1/rds/%d/connection", id), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetOperation polls the status of an async lifecycle operation.
func (s *RDSService) GetOperation(ctx context.Context, id uint, operationID string) (*types.RDSOperationResponse, error) {
	var out types.RDSOperationResponse
	_, _, err := s.c.do(ctx, "GET",
		fmt.Sprintf("/api/v1/rds/%d/operations/%s", id, url.PathEscape(operationID)), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ── Lifecycle (async writes) ───────────────────────────────────────

// Create provisions a new database instance. Honors Idempotency-Key — the
// server caches the response so duplicate creates are impossible (billing is
// involved). Returns 202 with an operation_id; poll Get until Status="ready".
func (s *RDSService) Create(ctx context.Context, req *types.CreateRDSInstanceRequest, opts ...WriteOption) (*types.RDSMutationResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.RDSMutationResponse
	_, _, err = s.c.do(ctx, "POST", "/api/v1/rds", req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Scale vertically changes the instance's compute plan (instance class). Async
// (202 + operation_id). Pass WithIfMatch(etag) to guard against concurrent
// writes.
func (s *RDSService) Scale(ctx context.Context, id uint, plan string, opts ...WriteOption) (*types.RDSMutationResponse, error) {
	return s.patch(ctx, id, &types.UpdateRDSInstanceRequest{Plan: plan}, opts...)
}

// Start resumes a suspended database (e.g. one stopped after a failed charge).
// The server balance-checks before starting; returns 409 RDS_INSTANCE_NOT_SUSPENDED
// if the instance isn't suspended, or 409 RDS_INSUFFICIENT_BALANCE if the wallet
// can't cover the next hour. Async (202 + operation_id).
func (s *RDSService) Start(ctx context.Context, id uint, opts ...WriteOption) (*types.RDSMutationResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.RDSMutationResponse
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/rds/%d/start", id), nil, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Resize grows the data disk to sizeGB. Shrink is rejected. Async (202 +
// operation_id).
func (s *RDSService) Resize(ctx context.Context, id uint, sizeGB int, opts ...WriteOption) (*types.RDSMutationResponse, error) {
	return s.patch(ctx, id, &types.UpdateRDSInstanceRequest{StorageGB: &sizeGB}, opts...)
}

// Delete tears down the instance and its storage. Async (202 + operation_id).
func (s *RDSService) Delete(ctx context.Context, id uint, opts ...WriteOption) (*types.RDSMutationResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.RDSMutationResponse
	_, _, err = s.c.do(ctx, "DELETE", fmt.Sprintf("/api/v1/rds/%d", id), nil, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *RDSService) patch(ctx context.Context, id uint, req *types.UpdateRDSInstanceRequest, opts ...WriteOption) (*types.RDSMutationResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.RDSMutationResponse
	_, _, err = s.c.do(ctx, "PATCH", fmt.Sprintf("/api/v1/rds/%d", id), req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ── Polling helpers ────────────────────────────────────────────────

// WaitForReady blocks until the instance reaches a terminal state. Returns the
// instance when Status="ready"; returns an error if Status settles on "failed"
// or "suspended".
func (s *RDSService) WaitForReady(ctx context.Context, id uint, opts ...PollOption) (*types.RDSInstanceResponse, error) {
	return PollResource(ctx,
		func(ctx context.Context) (*types.RDSInstanceResponse, error) {
			return s.Get(ctx, id)
		},
		func(v *types.RDSInstanceResponse) (bool, error) {
			if v == nil {
				return false, nil
			}
			switch types.RDSStatus(v.Status) {
			case types.RDSStatusReady:
				return true, nil
			case types.RDSStatusFailed:
				return true, errors.New("kumo: rds instance provisioning failed: " + v.StatusMessage)
			case types.RDSStatusSuspended:
				return true, errors.New("kumo: rds instance is suspended")
			default:
				return false, nil
			}
		},
		opts...,
	)
}

// CreateAndWait composes Create + WaitForReady.
func (s *RDSService) CreateAndWait(ctx context.Context, req *types.CreateRDSInstanceRequest, opts ...PollOption) (*types.RDSInstanceResponse, error) {
	created, err := s.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return s.WaitForReady(ctx, created.ID, opts...)
}
