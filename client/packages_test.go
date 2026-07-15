package client_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/kumobase/kumo-go/client"
	"github.com/kumobase/kumo-go/codes"
	"github.com/kumobase/kumo-go/types"
)

func TestPackages_ListPlans(t *testing.T) {
	var seenPath string
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		writeStruct(w, 200, "", "ok", &types.PackagesPricingResponse{
			Plans: []types.PackagesPlanOption{{ID: 72, Name: "Pay-as-you-go", Unit: "GB-month", PricePerUnit: "1060"}},
		})
	})
	plans, err := c.Packages().ListPlans(context.Background())
	if err != nil {
		t.Fatalf("ListPlans: %v", err)
	}
	// The server wraps the catalogue in {"plans":[…]}; the method must flatten it.
	if len(plans) != 1 || plans[0].PricePerUnit != "1060" {
		t.Errorf("plans = %+v, want 1 plan priced 1060", plans)
	}
	if want := "/api/v1/packages/plans"; seenPath != want {
		t.Errorf("path = %q, want %q", seenPath, want)
	}
}

func TestOrgPackages_List(t *testing.T) {
	var seenPath, seenQuery string
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath, seenQuery = r.URL.Path, r.URL.RawQuery
		// writeStruct cannot emit meta, so hand-write the envelope.
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"message":"ok","data":[{"id":1,"name":"lodash","format":"npm"}],`+
			`"meta":{"page":2,"page_size":50,"total_items":60,"total_pages":2}}`)
	})
	items, meta, err := c.Packages().Org("acme").List(context.Background(),
		client.WithPage(2), client.WithPageSize(50), client.WithSort("name", "asc"))
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(items) != 1 || items[0].Name != "lodash" {
		t.Errorf("items = %+v", items)
	}
	if meta == nil || meta.TotalItems != 60 || meta.Page != 2 {
		t.Errorf("meta = %+v, want page 2 / 60 items", meta)
	}
	if want := "/api/v1/packages/organizations/acme/packages/"; seenPath != want {
		t.Errorf("path = %q, want %q", seenPath, want)
	}
	for _, want := range []string{"page=2", "page_size=50", "sort=name", "sort_order=asc"} {
		if !strings.Contains(seenQuery, want) {
			t.Errorf("query %q missing %q", seenQuery, want)
		}
	}
}

// TestOrgPackages_Get_ScopedNameIsEncoded is the highest-value test in this
// file. The server routes the name through a greedy wildcard and URL-decodes
// it, so a scoped name MUST arrive %2F-encoded or it would split into two path
// segments.
//
// It asserts on r.RequestURI, NOT r.URL.Path: net/http decodes %2F into
// URL.Path, so asserting on Path would pass even if the escaping were dropped
// entirely — proving nothing.
func TestOrgPackages_Get_ScopedNameIsEncoded(t *testing.T) {
	var seenURI string
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seenURI = r.RequestURI
		writeStruct(w, 200, "", "ok", &types.PackageDetailResponse{
			Package:  types.PackageResponse{ID: 1, Name: "@acme/utils", Format: "npm"},
			DistTags: map[string]string{"latest": "1.2.0"},
		})
	})
	detail, _, err := c.Packages().Org("acme").Get(context.Background(), types.PackageFormatNPM, "@acme/utils")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	want := "/api/v1/packages/organizations/acme/packages/npm/@acme%2Futils"
	if seenURI != want {
		t.Errorf("RequestURI = %q, want %q", seenURI, want)
	}
	if detail.Package.Name != "@acme/utils" {
		t.Errorf("name = %q", detail.Package.Name)
	}
}

// Format is part of the package's identity, so it must land in the path for
// every non-npm ecosystem too.
func TestOrgPackages_Get_NonNPMFormats(t *testing.T) {
	for _, tc := range []struct {
		format types.PackageFormat
		name   string
		want   string
	}{
		{types.PackageFormatPyPI, "requests", "/api/v1/packages/organizations/acme/packages/pypi/requests"},
		{types.PackageFormatMaven, "com.acme:lib", "/api/v1/packages/organizations/acme/packages/maven/com.acme:lib"},
		{types.PackageFormatNuGet, "Newtonsoft.Json", "/api/v1/packages/organizations/acme/packages/nuget/Newtonsoft.Json"},
		{types.PackageFormatRubyGems, "rails", "/api/v1/packages/organizations/acme/packages/rubygems/rails"},
	} {
		t.Run(string(tc.format), func(t *testing.T) {
			var seenURI string
			c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				seenURI = r.RequestURI
				writeStruct(w, 200, "", "ok", &types.PackageDetailResponse{})
			})
			if _, _, err := c.Packages().Org("acme").Get(context.Background(), tc.format, tc.name); err != nil {
				t.Fatalf("Get: %v", err)
			}
			if seenURI != tc.want {
				t.Errorf("RequestURI = %q, want %q", seenURI, tc.want)
			}
		})
	}
}

func TestOrgPackages_Get_ReturnsETag(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `W/"17f1e2a3b4c"`)
		writeStruct(w, 200, "", "ok", &types.PackageDetailResponse{})
	})
	_, etag, err := c.Packages().Org("acme").Get(context.Background(), types.PackageFormatNPM, "lodash")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if want := `W/"17f1e2a3b4c"`; etag != want {
		t.Errorf("etag = %q, want %q", etag, want)
	}
}

// The version segment must NOT be escaped: the server unescapes only the name,
// so an escaped "+" would arrive as %2B and never match the stored version.
func TestOrgPackages_GetVersion_VersionNotEscaped(t *testing.T) {
	for _, tc := range []struct{ version, want string }{
		{"1.0.0", "/api/v1/packages/organizations/acme/packages/npm/lodash/versions/1.0.0"},
		{"1.0.0+build", "/api/v1/packages/organizations/acme/packages/npm/lodash/versions/1.0.0+build"},
	} {
		t.Run(tc.version, func(t *testing.T) {
			var seenURI string
			c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				seenURI = r.RequestURI
				writeStruct(w, 200, "", "ok", &types.PackageVersionResponse{Version: tc.version})
			})
			got, err := c.Packages().Org("acme").GetVersion(context.Background(), types.PackageFormatNPM, "lodash", tc.version)
			if err != nil {
				t.Fatalf("GetVersion: %v", err)
			}
			if seenURI != tc.want {
				t.Errorf("RequestURI = %q, want %q", seenURI, tc.want)
			}
			if got.Version != tc.version {
				t.Errorf("version = %q, want %q", got.Version, tc.version)
			}
		})
	}
}

func TestOrgPackages_GetVersion_NotFound(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeStruct(w, 404, codes.PackageVersionNotFound, "package version not found", nil)
	})
	_, err := c.Packages().Org("acme").GetVersion(context.Background(), types.PackageFormatNPM, "lodash", "9.9.9")
	if err == nil {
		t.Fatal("GetVersion: want error, got nil")
	}
	if !client.IsNotFound(err) {
		t.Errorf("IsNotFound = false, want true (err=%v)", err)
	}
	if !client.IsCode(err, codes.PackageVersionNotFound) {
		t.Errorf("IsCode(PackageVersionNotFound) = false (err=%v)", err)
	}
}

func TestOrgPackages_Delete(t *testing.T) {
	for _, tc := range []struct {
		name    string
		call    func(*client.Client) error
		wantURI string
	}{
		{
			name: "package",
			call: func(c *client.Client) error {
				return c.Packages().Org("acme").Delete(context.Background(), types.PackageFormatNPM, "@acme/utils")
			},
			wantURI: "/api/v1/packages/organizations/acme/packages/npm/@acme%2Futils",
		},
		{
			name: "version",
			call: func(c *client.Client) error {
				return c.Packages().Org("acme").DeleteVersion(context.Background(), types.PackageFormatPyPI, "requests", "2.0.0")
			},
			wantURI: "/api/v1/packages/organizations/acme/packages/pypi/requests/versions/2.0.0",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var seenMethod, seenURI, seenIdem string
			c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				seenMethod, seenURI = r.Method, r.RequestURI
				seenIdem = r.Header.Get("Idempotency-Key")
				writeStruct(w, 200, "", "successfully scheduled package deletion", nil)
			})
			if err := tc.call(c); err != nil {
				t.Fatalf("delete: %v", err)
			}
			if seenMethod != "DELETE" {
				t.Errorf("method = %q, want DELETE", seenMethod)
			}
			if seenURI != tc.wantURI {
				t.Errorf("RequestURI = %q, want %q", seenURI, tc.wantURI)
			}
			// resolveWriteOpts auto-generates a key even when the caller passes
			// none. The packages module ignores it today, but the header must
			// still be well-formed.
			if seenIdem == "" {
				t.Error("Idempotency-Key header is empty, want auto-generated value")
			}
		})
	}
}

// The detail response must decode from the server's snake_case wire shape.
// Guards against a regression to the PascalCase form that shipped before the
// json tags were added server-side.
func TestOrgPackages_Get_DecodesSnakeCase(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"message":"ok","data":{`+
			`"package":{"id":1,"name":"lodash","format":"npm","version_count":2},`+
			`"versions":[{"version":"1.0.0","size_bytes":10}],`+
			`"dist_tags":{"latest":"1.0.0"}}}`)
	})
	detail, _, err := c.Packages().Org("acme").Get(context.Background(), types.PackageFormatNPM, "lodash")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if detail.Package.Name != "lodash" || detail.Package.VersionCount != 2 {
		t.Errorf("package = %+v", detail.Package)
	}
	if len(detail.Versions) != 1 || detail.Versions[0].SizeBytes != 10 {
		t.Errorf("versions = %+v", detail.Versions)
	}
	if detail.DistTags["latest"] != "1.0.0" {
		t.Errorf("dist_tags = %+v", detail.DistTags)
	}
}
