package codes

// Container-registry wire codes returned by
// /api/v1/registry/organizations/:slug/repositories/* endpoints.
const (
	RegistryRepositoryNotFound      = "REGISTRY_REPOSITORY_NOT_FOUND"
	RegistryRepositoryAlreadyExists = "REGISTRY_REPOSITORY_ALREADY_EXISTS"
	RegistryInvalidRepositoryName   = "REGISTRY_INVALID_REPOSITORY_NAME"
	RegistryInvalidTagMutability    = "REGISTRY_INVALID_TAG_MUTABILITY"
	// Deprecated: soft_delete_days is no longer user-configurable (the purge
	// window is a fixed internal policy). This code is retained for wire
	// stability — it is never returned by the server anymore.
	RegistryInvalidSoftDeleteDays = "REGISTRY_INVALID_SOFT_DELETE_DAYS"

	RegistryTagImmutable     = "REGISTRY_TAG_IMMUTABLE"
	RegistryManifestNotFound = "REGISTRY_MANIFEST_NOT_FOUND"
	RegistryBlobNotFound     = "REGISTRY_BLOB_NOT_FOUND"

	RegistrySuspended         = "REGISTRY_SUSPENDED"
	RegistryMaxRepositoriesReached = "REGISTRY_MAX_REPOSITORIES_REACHED"

	// RegistryRepositorySystemOwned — the repository is owned by the platform
	// (e.g. auto-provisioned for a git-build app) and cannot be deleted via the
	// user API. It is removed automatically when its owning resource is deleted.
	RegistryRepositorySystemOwned = "REGISTRY_REPOSITORY_SYSTEM_OWNED"

	RegistryUnauthorized  = "REGISTRY_UNAUTHORIZED"
	RegistryForbidden     = "REGISTRY_FORBIDDEN"
	RegistryInternalError = "REGISTRY_INTERNAL_ERROR"
)
