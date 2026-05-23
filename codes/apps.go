package codes

// App-module wire codes returned by /api/v1/apps/* endpoints. Mirrors the
// per-sentinel branches in modules/app/errors.go::handleAppError on the
// server.
const (
	AppNotFound                 = "APP_NOT_FOUND"
	AppOperationNotFound        = "APP_OPERATION_NOT_FOUND"
	AppCustomDomainNotFound     = "APP_CUSTOM_DOMAIN_NOT_FOUND"
	AppRegistryCredentialNotFound = "APP_REGISTRY_CREDENTIAL_NOT_FOUND"

	AppDeploymentInProgress     = "APP_DEPLOYMENT_IN_PROGRESS"
	AppAlreadyStopped           = "APP_ALREADY_STOPPED"
	AppCustomDomainExists       = "APP_CUSTOM_DOMAIN_EXISTS"
	AppDomainAlreadyInUse       = "APP_DOMAIN_ALREADY_IN_USE"
	AppVolumeConflict           = "APP_VOLUME_CONFLICT"
	AppQuotaExceeded            = "APP_QUOTA_EXCEEDED"

	AppInvalidPricingSlug = "APP_INVALID_PRICING_SLUG"
	AppMustBeExposed      = "APP_MUST_BE_EXPOSED"
	AppDomainPlatformZone = "APP_DOMAIN_PLATFORM_ZONE"

	AppInsufficientBalance      = "APP_INSUFFICIENT_BALANCE"
	AppPlatformCapacityExceeded = "APP_PLATFORM_CAPACITY_EXCEEDED"

	AppValidationFailed   = "APP_VALIDATION_FAILED"
	AppUnauthorized       = "APP_UNAUTHORIZED"
	AppInvalidID          = "APP_INVALID_ID"
	AppInvalidRequestBody = "APP_INVALID_REQUEST_BODY"
	AppInternalError      = "APP_INTERNAL_ERROR"
)
