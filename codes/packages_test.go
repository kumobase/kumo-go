package codes

import "testing"

// Wire codes are a public contract. These assert the exact string values
// so an accidental rename is caught here before release.
func TestPackagesCodeValues(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"PackageNotFound", PackageNotFound, "PACKAGE_NOT_FOUND"},
		{"PackageVersionNotFound", PackageVersionNotFound, "PACKAGE_VERSION_NOT_FOUND"},
		{"PackageVersionExists", PackageVersionExists, "PACKAGE_VERSION_EXISTS"},
		{"PackageInvalidName", PackageInvalidName, "PACKAGE_INVALID_NAME"},
		{"PackageInvalidVersion", PackageInvalidVersion, "PACKAGE_INVALID_VERSION"},
		{"PackageIntegrityMismatch", PackageIntegrityMismatch, "PACKAGE_INTEGRITY_MISMATCH"},
		{"PackageMalformedPublish", PackageMalformedPublish, "PACKAGE_MALFORMED_PUBLISH"},
		{"PackageTarballTooLarge", PackageTarballTooLarge, "PACKAGE_TARBALL_TOO_LARGE"},
		{"PackageTagNotFound", PackageTagNotFound, "PACKAGE_TAG_NOT_FOUND"},
		{"PackageTagProtected", PackageTagProtected, "PACKAGE_TAG_PROTECTED"},
		{"PackageOrgSuspended", PackageOrgSuspended, "PACKAGE_ORG_SUSPENDED"},
		{"PackageAuthRequired", PackageAuthRequired, "PACKAGE_AUTH_REQUIRED"},
		{"PackageForbidden", PackageForbidden, "PACKAGE_FORBIDDEN"},
		{"PackageUnpublishForbidden", PackageUnpublishForbidden, "PACKAGE_UNPUBLISH_FORBIDDEN"},
		{"PackageInternalError", PackageInternalError, "PACKAGE_INTERNAL_ERROR"},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}
