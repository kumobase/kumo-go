package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/kumobase/kumo-go/types"
)

// VolumesService backs /api/v1/volumes/*. Create is synchronous (the
// volume row materialises immediately, even though the underlying
// Longhorn PVC takes a moment to bind). Resize is async via the volume's
// Status field — use ResizeAndWait for the common block-until-ready flow.
type VolumesService struct {
	c *Client
}

// Volumes returns the volumes service.
func (c *Client) Volumes() *VolumesService { return &VolumesService{c: c} }

// Get fetches a volume's current state (including transient statuses like
// "creating" / "resizing"). Returns the ETag for use with IfMatch on a
// subsequent Resize.
func (s *VolumesService) Get(ctx context.Context, id uint) (*types.VolumeResponse, string, error) {
	var out types.VolumeResponse
	etag, _, err := s.c.do(ctx, "GET", fmt.Sprintf("/api/v1/volumes/%d", id), nil, nil, &out)
	if err != nil {
		return nil, "", err
	}
	return &out, etag, nil
}

// List returns volumes scoped to the authenticated user. Use
// WithExtraQuery("status", "ready") / ("app_id", "42") / ("attached",
// "true") for endpoint-specific filters. The server rejects mutually
// exclusive combinations (app_id + attached) with 400
// INVALID_FILTER_COMBINATION.
func (s *VolumesService) List(ctx context.Context, opts ...ListOption) ([]types.VolumeResponse, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.VolumeResponse
	meta, err := s.c.doList(ctx, "GET", withQuery("/api/v1/volumes", q), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, meta, nil
}

// Create provisions a new volume. Honors Idempotency-Key. If req.AppID is
// set, the volume is attached on creation — see ForceReconfigure for the
// "auto-scale app to 1 replica and disable autoscaling" opt-in.
func (s *VolumesService) Create(ctx context.Context, req *types.CreateVolumeRequest, opts ...WriteOption) (*types.VolumeResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.VolumeResponse
	_, _, err = s.c.do(ctx, "POST", "/api/v1/volumes", req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Resize asks the underlying CSI driver to expand the volume. Returns
// once the request is accepted (HTTP 202). The volume's Status field
// transitions to "resizing" and back to "ready"/"failed" asynchronously
// via the volume-resize-reconciler worker.
//
// Use ResizeAndWait for the common block-until-terminal flow. Pass
// IfMatch(etag) for optimistic concurrency.
func (s *VolumesService) Resize(ctx context.Context, id uint, req *types.ResizeVolumeRequest, opts ...WriteOption) (*types.VolumeResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.VolumeResponse
	_, _, err = s.c.do(ctx, "PATCH", fmt.Sprintf("/api/v1/volumes/%d/resize", id), req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// ResizeAndWait composes Resize + a status poll until the volume leaves
// the "resizing" state. Returns the final volume on success or the
// surfaced error (e.g. CSI driver rejection) on failure.
func (s *VolumesService) ResizeAndWait(ctx context.Context, id uint, req *types.ResizeVolumeRequest, opts ...PollOption) (*types.VolumeResponse, error) {
	if _, err := s.Resize(ctx, id, req); err != nil {
		return nil, err
	}
	return PollResource(ctx,
		func(ctx context.Context) (*types.VolumeResponse, error) {
			v, _, err := s.Get(ctx, id)
			return v, err
		},
		func(v *types.VolumeResponse) (bool, error) {
			if v == nil {
				return false, nil
			}
			switch types.VolumeStatus(v.Status) {
			case types.VolumeStatusReady, types.VolumeStatusDetached:
				return true, nil
			case types.VolumeStatusFailed:
				msg := "volume resize failed"
				if v.LastResizeError != nil {
					msg = *v.LastResizeError
				}
				return true, errors.New("kumo: " + msg)
			default:
				return false, nil
			}
		},
		opts...,
	)
}

// Attach binds a detached volume to an app. ForceReconfigure on the
// request has the same "auto-reconfigure the target app" semantics as on
// Create.
func (s *VolumesService) Attach(ctx context.Context, id uint, req *types.AttachVolumeRequest, opts ...WriteOption) (*types.VolumeResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.VolumeResponse
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/volumes/%d/attach", id), req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Detach unbinds a volume from its current app. Returns 409 if the volume
// is permanently attached (a server-side policy on certain storage tiers).
func (s *VolumesService) Detach(ctx context.Context, id uint, opts ...WriteOption) (*types.VolumeResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.VolumeResponse
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/volumes/%d/detach", id), nil, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete removes a volume. The actual storage teardown happens
// asynchronously — Get will return the volume with Status="deleting"
// for a brief window before 404'ing.
func (s *VolumesService) Delete(ctx context.Context, id uint, opts ...WriteOption) error {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return err
	}
	_, _, err = s.c.do(ctx, "DELETE", fmt.Sprintf("/api/v1/volumes/%d", id), nil, &wopts, nil)
	return err
}
