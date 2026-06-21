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

// ListEngineVersions returns the public engine-version catalogue — the PG
// versions offered for new instances. The Version field of each entry is what
// you pass as CreateRDSInstanceRequest.EngineVersion. Filter by engine with
// WithExtraQuery("engine", "postgresql").
func (s *RDSService) ListEngineVersions(ctx context.Context, opts ...ListOption) ([]types.PublicRDSEngineVersionResponse, error) {
	q := resolveListOpts(opts)
	var out []types.PublicRDSEngineVersionResponse
	_, _, err := s.c.do(ctx, "GET", withQuery("/api/v1/rds/engine-versions", q), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ListEngineVersionParameters returns the editable parameter catalogue
// (allowlist) for an engine version — what a parameter template may set.
func (s *RDSService) ListEngineVersionParameters(ctx context.Context, engineVersionID uint, opts ...ListOption) ([]types.RDSPgParameterResponse, error) {
	q := resolveListOpts(opts)
	var out []types.RDSPgParameterResponse
	_, _, err := s.c.do(ctx, "GET",
		withQuery(fmt.Sprintf("/api/v1/rds/engine-versions/%d/parameters", engineVersionID), q), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ListParameterTemplates returns the caller's parameter templates plus the
// read-only system templates.
func (s *RDSService) ListParameterTemplates(ctx context.Context, opts ...ListOption) ([]types.PublicRDSParameterTemplateResponse, error) {
	q := resolveListOpts(opts)
	var out []types.PublicRDSParameterTemplateResponse
	_, _, err := s.c.do(ctx, "GET", withQuery("/api/v1/rds/parameter-templates", q), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetParameterTemplate fetches one template by id.
func (s *RDSService) GetParameterTemplate(ctx context.Context, id uint) (*types.PublicRDSParameterTemplateResponse, error) {
	var out types.PublicRDSParameterTemplateResponse
	_, _, err := s.c.do(ctx, "GET", fmt.Sprintf("/api/v1/rds/parameter-templates/%d", id), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateParameterTemplate creates a custom parameter template.
func (s *RDSService) CreateParameterTemplate(ctx context.Context, req *types.CreateRDSParameterTemplateRequest, opts ...WriteOption) (*types.PublicRDSParameterTemplateResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.PublicRDSParameterTemplateResponse
	_, _, err = s.c.do(ctx, "POST", "/api/v1/rds/parameter-templates", req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateParameterTemplate edits a custom template. System/default templates are
// read-only (409 RDS_PARAMETER_TEMPLATE_READ_ONLY). Pass WithIfMatch(etag) to
// guard against concurrent edits.
func (s *RDSService) UpdateParameterTemplate(ctx context.Context, id uint, req *types.UpdateRDSParameterTemplateRequest, opts ...WriteOption) (*types.PublicRDSParameterTemplateResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.PublicRDSParameterTemplateResponse
	_, _, err = s.c.do(ctx, "PATCH", fmt.Sprintf("/api/v1/rds/parameter-templates/%d", id), req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteParameterTemplate removes a custom template. Fails with 409
// RDS_PARAMETER_TEMPLATE_IN_USE if any instance still references it.
func (s *RDSService) DeleteParameterTemplate(ctx context.Context, id uint, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "DELETE", fmt.Sprintf("/api/v1/rds/parameter-templates/%d", id), nil, &wopts, nil)
	return err
}

// SetParameterTemplate attaches a parameter template to a running instance and
// live-reconfigures it. Async (202 + operation_id).
func (s *RDSService) SetParameterTemplate(ctx context.Context, id uint, templateSlug string, opts ...WriteOption) (*types.RDSMutationResponse, error) {
	return s.patch(ctx, id, &types.UpdateRDSInstanceRequest{ParameterTemplate: templateSlug}, opts...)
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

// ScaleReplicas changes the topology mode and/or read-replica count on a running
// instance (KubeBlocks HorizontalScaling + sync-mode reconfigure). Pass an empty
// mode to leave it unchanged, or nil readReplicas to leave the count unchanged
// (at least one must change). Async (202 + operation_id).
func (s *RDSService) ScaleReplicas(ctx context.Context, id uint, mode string, readReplicas *int, opts ...WriteOption) (*types.RDSMutationResponse, error) {
	return s.patch(ctx, id, &types.UpdateRDSInstanceRequest{Mode: mode, ReadReplicas: readReplicas}, opts...)
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

// Switchover triggers a planned, operator-initiated HA role swap: the
// synchronous standby is promoted to primary and the old primary demoted while
// both are healthy — a graceful counterpart to automatic failover, for node
// maintenance, AZ drains, or rebalancing. Only valid on HA instances (409
// RDS_SWITCHOVER_NOT_HA on standalone), gated behind the platform
// RDS_SWITCHOVER_ENABLED flag (409 RDS_SWITCHOVER_DISABLED when off), and
// rejected with 409 RDS_SWITCHOVER_NOT_READY when no healthy sync standby
// exists. Async read replicas are never promotion candidates. Async (202 +
// operation_id). Pass WithIfMatch(etag) to guard against concurrent writes.
func (s *RDSService) Switchover(ctx context.Context, id uint, opts ...WriteOption) (*types.RDSMutationResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.RDSMutationResponse
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/rds/%d/switchover", id), nil, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateTLS changes the connection TLS enforcement of a running instance between
// "optional" and "required" (server-side pg_hba reload, no restart). Transitions
// to/from "disabled" are rejected (409 RDS_TLS_MODE_CHANGE_UNSUPPORTED); asking
// for "required" while the platform enforcement flag is off is rejected (400
// RDS_TLS_ENFORCEMENT_DISABLED). Async (202 + operation_id). Pass WithIfMatch(etag)
// to guard against concurrent writes.
func (s *RDSService) UpdateTLS(ctx context.Context, id uint, mode string, opts ...WriteOption) (*types.RDSMutationResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.RDSMutationResponse
	_, _, err = s.c.do(ctx, "PATCH", fmt.Sprintf("/api/v1/rds/%d/tls", id), &types.UpdateRDSTLSRequest{TLSMode: mode}, &wopts, &out)
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

// ── Backups (to object storage) ────────────────────────────────────

// CreateBackup starts an on-demand backup of the instance to object storage.
// Honors Idempotency-Key (retrying the same key returns the same backup rather
// than starting a second). Gated behind the platform RDS_BACKUP_ENABLED flag
// (409 RDS_BACKUP_DISABLED when off); rejected with 409 RDS_BACKUP_IN_PROGRESS
// when a backup is already running, or 409 RDS_INSTANCE_NOT_READY unless the
// instance is ready. Async (202 + operation_id); poll GetBackup until
// Status="completed".
func (s *RDSService) CreateBackup(ctx context.Context, id uint, req *types.CreateRDSBackupRequest, opts ...WriteOption) (*types.RDSMutationResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	if req == nil {
		req = &types.CreateRDSBackupRequest{}
	}
	var out types.RDSMutationResponse
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/rds/%d/backups", id), req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ListBackups returns the instance's backups, paginated (newest first). Filter
// by status with WithExtraQuery("status", "completed").
func (s *RDSService) ListBackups(ctx context.Context, id uint, opts ...ListOption) ([]types.RDSBackupResponse, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.RDSBackupResponse
	meta, err := s.c.doList(ctx, "GET", withQuery(fmt.Sprintf("/api/v1/rds/%d/backups", id), q), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, meta, nil
}

// ListAllBackups returns ALL the caller's backups across every database
// (including backups whose source database was deleted) — the global backups
// view. Paginated, newest first. Filter by status with
// WithExtraQuery("status", "completed").
func (s *RDSService) ListAllBackups(ctx context.Context, opts ...ListOption) ([]types.RDSBackupResponse, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.RDSBackupResponse
	meta, err := s.c.doList(ctx, "GET", withQuery("/api/v1/rds/backups", q), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, meta, nil
}

// RestoreBackup provisions a NEW database from any backup the caller owns, keyed
// on the backup id (not a source database) — so a RETAINED backup whose source
// database was deleted is still restorable. Point-in-time recovery is not
// available here (use Restore on a live source for PITR). Async (202 +
// operation_id on the new instance).
func (s *RDSService) RestoreBackup(ctx context.Context, req *types.RestoreRDSBackupRequest, opts ...WriteOption) (*types.RDSMutationResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.RDSMutationResponse
	_, _, err = s.c.do(ctx, "POST", "/api/v1/rds/restore", req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetBackup fetches a single backup by id.
func (s *RDSService) GetBackup(ctx context.Context, id, backupID uint) (*types.RDSBackupResponse, error) {
	var out types.RDSBackupResponse
	_, _, err := s.c.do(ctx, "GET", fmt.Sprintf("/api/v1/rds/%d/backups/%d", id, backupID), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteBackup removes a backup from object storage and stops billing for it.
// Async (202 + operation_id). Pass WithIfMatch(etag) to guard against
// concurrent writes.
func (s *RDSService) DeleteBackup(ctx context.Context, id, backupID uint, opts ...WriteOption) (*types.RDSMutationResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.RDSMutationResponse
	_, _, err = s.c.do(ctx, "DELETE", fmt.Sprintf("/api/v1/rds/%d/backups/%d", id, backupID), nil, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// SetBackupConfig enables/disables automatic scheduled backups for an instance
// and sets the schedule, retention, and backup tier. Gated behind the platform
// RDS_BACKUP_ENABLED flag (409 RDS_BACKUP_DISABLED when off). Returns the
// resulting configuration.
func (s *RDSService) SetBackupConfig(ctx context.Context, id uint, req *types.UpdateRDSBackupConfigRequest, opts ...WriteOption) (*types.RDSBackupConfigResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.RDSBackupConfigResponse
	_, _, err = s.c.do(ctx, "PUT", fmt.Sprintf("/api/v1/rds/%d/backup-config", id), req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Restore provisions a NEW database instance from a completed backup (the
// source instance is untouched and billed independently). The backup must be
// completed (409 RDS_BACKUP_NOT_READY otherwise) and the target storage must be
// >= the source's (400 RDS_RESTORE_STORAGE_TOO_SMALL). Async (202 +
// operation_id on the new instance); poll the returned instance id until ready.
func (s *RDSService) Restore(ctx context.Context, id uint, req *types.RestoreRDSBackupRequest, opts ...WriteOption) (*types.RDSMutationResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.RDSMutationResponse
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/rds/%d/restore", id), req, &wopts, &out)
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
