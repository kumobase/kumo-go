package codes

// Kumo Packages wire codes. Returned both by the language-package registry
// protocol endpoints (npm/maven/pypi/nuget/rubygems; served at /<fmt>/:org/*)
// and by the dashboard/management API
// (/api/v1/packages/organizations/:slug/packages/*).
//
// Note: the npm protocol endpoints return npm-native JSON bodies and do NOT
// embed these codes in the response body — npm dictates the wire shape. The
// codes are still assigned + logged server-side, and are used verbatim by the
// management API. They remain a public contract: never change a string value.
const (
	PackageNotFound        = "PACKAGE_NOT_FOUND"
	PackageVersionNotFound = "PACKAGE_VERSION_NOT_FOUND"
	PackageVersionExists   = "PACKAGE_VERSION_EXISTS"
	PackageInvalidName     = "PACKAGE_INVALID_NAME"
	PackageInvalidVersion  = "PACKAGE_INVALID_VERSION"

	PackageIntegrityMismatch = "PACKAGE_INTEGRITY_MISMATCH"
	PackageMalformedPublish  = "PACKAGE_MALFORMED_PUBLISH"
	PackageTarballTooLarge   = "PACKAGE_TARBALL_TOO_LARGE"

	PackageTagNotFound  = "PACKAGE_TAG_NOT_FOUND"
	PackageTagProtected = "PACKAGE_TAG_PROTECTED"

	PackageOrgSuspended       = "PACKAGE_ORG_SUSPENDED"
	PackageAuthRequired       = "PACKAGE_AUTH_REQUIRED"
	PackageForbidden          = "PACKAGE_FORBIDDEN"
	PackageUnpublishForbidden = "PACKAGE_UNPUBLISH_FORBIDDEN"

	PackageInternalError = "PACKAGE_INTERNAL_ERROR"
)
