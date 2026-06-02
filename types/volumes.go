package types

import "time"

// VolumeStatus enumerates the lifecycle states of a Kumo volume. Most
// transitions happen asynchronously via the volume-resize-reconciler worker;
// clients poll GET /api/v1/volumes/:id and watch Status until it leaves the
// transient states (Creating, Resizing, Deleting).
type VolumeStatus string

const (
	VolumeStatusCreating VolumeStatus = "creating"
	VolumeStatusReady    VolumeStatus = "ready"
	// Deprecated: the server no longer emits "detached". An unattached volume is
	// now reported as VolumeStatusReady; attachment is conveyed by AppID/MountPath.
	// Retained so existing references compile and so a client talking to an older
	// server still recognizes the value.
	VolumeStatusDetached VolumeStatus = "detached"
	VolumeStatusResizing VolumeStatus = "resizing"
	VolumeStatusDeleting VolumeStatus = "deleting"
	VolumeStatusFailed   VolumeStatus = "failed"
)

// CreateVolumeRequest is the body for POST /api/v1/volumes. Honors
// Idempotency-Key.
//
// ForceReconfigure: when true and AppID is set, the target app is
// automatically scaled to 1 replica and has autoscaling disabled as part of
// the create. Without it, an app that doesn't already satisfy those
// constraints is rejected (data-loss-safe default).
// Identify the target app (when attaching at create) by exactly one of AppID
// or AppName — the server returns 400 VALIDATION_FAILED if both are set.
type CreateVolumeRequest struct {
	AppID            *uint  `json:"app_id,omitempty"`
	AppName          string `json:"app_name,omitempty"`
	Name             string `json:"name"`         // 1..100 chars
	StorageTier      string `json:"storage_tier"` // slug — see StorageTierResponse.Slug
	SizeGB           int    `json:"size_gb"`      // >= 1
	MountPath        string `json:"mount_path,omitempty"`
	ForceReconfigure bool   `json:"force_reconfigure"`
}

// ResizeVolumeRequest is the body for PATCH /api/v1/volumes/:id/resize.
// Returns 202 Accepted with Retry-After; poll GET /api/v1/volumes/:id until
// Status leaves "resizing".
type ResizeVolumeRequest struct {
	SizeGB int `json:"size_gb"` // >= 1; shrink is rejected by Longhorn (will surface as LastResizeError)
}

// AttachVolumeRequest is the body for POST /api/v1/volumes/:id/attach.
// ForceReconfigure has the same semantics as on CreateVolumeRequest.
//
// Same exactly-one-of contract: set AppID or AppName, never both.
type AttachVolumeRequest struct {
	AppID            uint   `json:"app_id,omitempty"`
	AppName          string `json:"app_name,omitempty"`
	MountPath        string `json:"mount_path,omitempty"`
	ForceReconfigure bool   `json:"force_reconfigure"`
}

// StorageTierResponse describes one purchasable volume tier returned by
// GET /api/v1/volumes/plans (and embedded in VolumeResponse).
//
// PricePerGBHour is a decimal string ("0.0001234"); parse with a decimal
// library if you need arithmetic.
type StorageTierResponse struct {
	ID             uint          `json:"id"`
	Slug           string        `json:"slug"`
	Name           string        `json:"name"`
	PricePerGBHour string        `json:"price_per_gb_hour"`
	MinSizeGB      int           `json:"min_size_gb"`
	MaxSizeGB      int           `json:"max_size_gb"`
	Availability   *Availability `json:"availability,omitempty"`
}

// VolumeResponse is returned by every volume endpoint that handles a single
// volume (GET /api/v1/volumes/:id, POST create, PATCH resize, POST attach,
// POST detach). The server sets ETag from UpdatedAt; echo it back in
// If-Match on PATCH /resize to detect concurrent writes.
//
// LastResizeError / LastResizeAt are populated when a previous resize
// transitioned the volume to "failed" — the user-facing surface for the
// underlying CSI driver's rejection message.
type VolumeResponse struct {
	ID              uint                `json:"id"`
	Name            string              `json:"name"`
	AppID           *uint               `json:"app_id"`
	AppName         *string             `json:"app_name,omitempty"`
	StorageTier     StorageTierResponse `json:"storage_tier"`
	SizeGB          int                 `json:"size_gb"`
	MountPath       string              `json:"mount_path"`
	Status          string              `json:"status"`
	ErrorMessage    *string             `json:"error_message,omitempty"`
	LastResizeError *string             `json:"last_resize_error,omitempty"`
	LastResizeAt    *time.Time          `json:"last_resize_at,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}
