package codes

// Secret-module wire codes returned by /api/v1/secrets/* endpoints.
const (
	SecretNotFound        = "SECRET_NOT_FOUND"
	SecretInUse           = "SECRET_IN_USE"
	SecretTypeImmutable   = "SECRET_TYPE_IMMUTABLE"
	SecretUnsupportedType = "SECRET_UNSUPPORTED_TYPE"
	SecretEnvVarsEmpty    = "SECRET_ENV_VARS_EMPTY"

	SecretValidationFailed   = "SECRET_VALIDATION_FAILED"
	SecretInvalidRequestBody = "SECRET_INVALID_REQUEST_BODY"
	SecretInvalidType        = "SECRET_INVALID_TYPE"
	SecretInvalidID          = "SECRET_INVALID_ID"
	SecretUnauthorized       = "SECRET_UNAUTHORIZED"
	SecretInternalError      = "SECRET_INTERNAL_ERROR"
)
