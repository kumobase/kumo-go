package codes

// Container-registry wire codes returned by
// /api/v1/registry/organizations/:slug/repositories/* endpoints.
const (
	RegistryRepositoryNotFound      = "REGISTRY_REPOSITORY_NOT_FOUND"
	RegistryRepositoryAlreadyExists = "REGISTRY_REPOSITORY_ALREADY_EXISTS"
	RegistryInvalidRepositoryName   = "REGISTRY_INVALID_REPOSITORY_NAME"
	RegistryInvalidTagMutability    = "REGISTRY_INVALID_TAG_MUTABILITY"
	RegistryInvalidSoftDeleteDays   = "REGISTRY_INVALID_SOFT_DELETE_DAYS"

	RegistryTagImmutable     = "REGISTRY_TAG_IMMUTABLE"
	RegistryManifestNotFound = "REGISTRY_MANIFEST_NOT_FOUND"
	RegistryBlobNotFound     = "REGISTRY_BLOB_NOT_FOUND"

	RegistrySuspended         = "REGISTRY_SUSPENDED"
	RegistryMaxRepositoriesReached = "REGISTRY_MAX_REPOSITORIES_REACHED"

	RegistryUnauthorized  = "REGISTRY_UNAUTHORIZED"
	RegistryForbidden     = "REGISTRY_FORBIDDEN"
	RegistryInternalError = "REGISTRY_INTERNAL_ERROR"
)
