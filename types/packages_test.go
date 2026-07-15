package types

import (
	"testing"
	"time"
)

// pkgTime is a fixed UTC instant. roundTrip compares with reflect.DeepEqual and
// re-marshals byte-for-byte, so a time.Now() value (which carries a monotonic
// reading that does not survive a JSON round-trip) would fail both checks.
var pkgTime = time.Date(2026, 7, 14, 9, 30, 0, 0, time.UTC)

func TestPackagesPricingRoundTrip(t *testing.T) {
	roundTrip(t, "PackagesPlanOption", PackagesPlanOption{
		ID:            72,
		Name:          "Pay-as-you-go",
		Description:   "metered storage, billed monthly",
		Unit:          "GB-month",
		PricePerUnit:  "1060",
		Currency:      "IDR",
		ChargeModel:   "postpaid",
		BillingPeriod: "monthly",
	})
	// Description omitted — exercises the omitempty field.
	roundTrip(t, "PackagesPlanOption+NoDescription", PackagesPlanOption{
		ID:            73,
		Name:          "Free tier",
		Unit:          "GB-month",
		PricePerUnit:  "0",
		Currency:      "IDR",
		ChargeModel:   "postpaid",
		BillingPeriod: "monthly",
	})
	roundTrip(t, "PackagesPricingResponse", PackagesPricingResponse{
		Plans: []PackagesPlanOption{
			{ID: 72, Name: "Pay-as-you-go", Unit: "GB-month", PricePerUnit: "1060", Currency: "IDR", ChargeModel: "postpaid", BillingPeriod: "monthly"},
		},
	})
	// Empty list — the no-plan / all-zero-price case must serialise as {"plans":[]}.
	roundTrip(t, "PackagesPricingResponse+Empty", PackagesPricingResponse{Plans: []PackagesPlanOption{}})
}

func TestPackageResponseRoundTrip(t *testing.T) {
	roundTrip(t, "PackageResponse", PackageResponse{
		ID:             9,
		OrganizationID: 3,
		Format:         string(PackageFormatNPM),
		Name:           "@acme/utils",
		LatestVersion:  "1.2.0",
		VersionCount:   4,
		CreatedAt:      pkgTime,
		UpdatedAt:      pkgTime,
	})
	// LatestVersion omitted — a package whose versions were all unpublished.
	roundTrip(t, "PackageResponse+NoLatest", PackageResponse{
		ID:             10,
		OrganizationID: 3,
		Format:         string(PackageFormatPyPI),
		Name:           "requests",
		VersionCount:   0,
		CreatedAt:      pkgTime,
		UpdatedAt:      pkgTime,
	})
}

func TestPackageVersionResponseRoundTrip(t *testing.T) {
	roundTrip(t, "PackageVersionResponse", PackageVersionResponse{
		Version:     "1.2.0",
		SizeBytes:   2048,
		Shasum:      "aabbcc",
		Integrity:   "sha512-abc==",
		PublishedAt: pkgTime,
	})
	// Deprecated set — exercises the pointer + omitempty pairing.
	roundTrip(t, "PackageVersionResponse+Deprecated", PackageVersionResponse{
		Version:     "1.0.0",
		SizeBytes:   1024,
		Deprecated:  strPtr("use v1.2.0"),
		PublishedAt: pkgTime,
	})
	// Non-npm: no dist metadata, so shasum/integrity drop out.
	roundTrip(t, "PackageVersionResponse+NoDistMeta", PackageVersionResponse{
		Version:     "2.1.0",
		SizeBytes:   4096,
		PublishedAt: pkgTime,
	})
}

func TestPackageDetailResponseRoundTrip(t *testing.T) {
	roundTrip(t, "PackageDetailResponse", PackageDetailResponse{
		Package: PackageResponse{
			ID: 9, OrganizationID: 3, Format: string(PackageFormatNPM),
			Name: "@acme/utils", LatestVersion: "1.2.0", VersionCount: 1,
			CreatedAt: pkgTime, UpdatedAt: pkgTime,
		},
		Versions: []PackageVersionResponse{
			{Version: "1.2.0", SizeBytes: 2048, PublishedAt: pkgTime},
		},
		DistTags: map[string]string{"latest": "1.2.0"},
	})
	// Non-npm formats carry no dist-tags, and a fully-unpublished package has
	// no versions — both must still serialise as [] / {}, never null.
	roundTrip(t, "PackageDetailResponse+Empty", PackageDetailResponse{
		Package: PackageResponse{
			ID: 10, OrganizationID: 3, Format: string(PackageFormatMaven),
			Name: "com.acme:lib", CreatedAt: pkgTime, UpdatedAt: pkgTime,
		},
		Versions: []PackageVersionResponse{},
		DistTags: map[string]string{},
	})
}
