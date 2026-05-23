package codes

// Volume-module wire codes returned by /api/v1/volumes/* endpoints.
const (
	VolumeNotFound         = "VOLUME_NOT_FOUND"
	VolumeNotAttached      = "VOLUME_NOT_ATTACHED"
	VolumeAttached         = "VOLUME_ATTACHED"
	VolumePermanentlyAttached = "VOLUME_PERMANENTLY_ATTACHED"
	VolumeCreating         = "VOLUME_CREATING"
	VolumeDeleting         = "VOLUME_DELETING"
	VolumeResizing         = "VOLUME_RESIZING"
	VolumeFailed           = "VOLUME_FAILED"
	VolumeNotReady         = "VOLUME_NOT_READY"
	VolumeNeverAttached    = "VOLUME_NEVER_ATTACHED"
	VolumeOutOfSync        = "VOLUME_OUT_OF_SYNC"

	StorageTierNotFound       = "STORAGE_TIER_NOT_FOUND"
	StorageClassNotFound      = "STORAGE_CLASS_NOT_FOUND"
	StorageClassNotExpandable = "STORAGE_CLASS_NOT_EXPANDABLE"

	SizeBelowMinimum    = "SIZE_BELOW_MINIMUM"
	SizeAboveMaximum    = "SIZE_ABOVE_MAXIMUM"
	CannotShrinkVolume  = "CANNOT_SHRINK_VOLUME"
	InvalidVolumeSize   = "INVALID_VOLUME_SIZE"

	TargetAppAlreadyHasVolume = "TARGET_APP_ALREADY_HAS_VOLUME"
	AppMustHaveOneReplica     = "APP_MUST_HAVE_ONE_REPLICA"
	AutoscalingWithVolume     = "AUTOSCALING_WITH_VOLUME"
	// VolumeDeploymentInProgress is emitted by volume routes when the target
	// app's deployment is mid-flight. Distinct constant from
	// codes/apps.go::AppDeploymentInProgress (same wire value
	// "DEPLOYMENT_IN_PROGRESS" — apps uses "APP_DEPLOYMENT_IN_PROGRESS"; this
	// is the legacy short form that the volume handler still emits).
	VolumeDeploymentInProgress = "DEPLOYMENT_IN_PROGRESS"
	VolumeAppDeploymentFailed  = "APP_DEPLOYMENT_FAILED"

	VolumeInvalidVolumeID     = "INVALID_VOLUME_ID"
	VolumeInvalidAppID        = "INVALID_APP_ID"
	VolumeInvalidRequestBody  = "INVALID_REQUEST_BODY"
	VolumeInvalidStatusFilter = "INVALID_STATUS_FILTER"
	VolumeInvalidAttachedFilter = "INVALID_ATTACHED_FILTER"
)
