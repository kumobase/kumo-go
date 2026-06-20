package client

import (
	"context"
	"fmt"

	"github.com/kumobase/kumo-go/types"
)

// BuildsService backs the git-build surface: creating a git-build app from a
// source connection, and listing/inspecting/triggering/canceling the builds
// of an existing git-build app.
//
// A git-build app's image is produced by the platform from a connected git
// repository on every push to the configured branch (and once on create). The
// image is system-owned — it cannot be set via Apps().Update, which returns
// codes.BuildAppImageImmutable for git-build apps.
type BuildsService struct {
	c *Client
}

// Builds returns the builds service.
func (c *Client) Builds() *BuildsService { return &BuildsService{c: c} }

// ListBuilders returns the platform's selectable builder kinds (auto, railpack,
// dockerfile, static, cnb) and CNB language presets, with the current zero-config
// default flagged. Static metadata — safe to cache briefly. Use it to populate a
// build-config UI instead of hardcoding the list.
func (s *BuildsService) ListBuilders(ctx context.Context) (*types.BuildersResponse, error) {
	var out types.BuildersResponse
	_, _, err := s.c.do(ctx, "GET", "/api/v1/builders", nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateGitBuildApp creates a git-build app bound to the given source
// connection and triggers its initial build (the branch HEAD). The returned
// CreateAppResponse carries the new app's id; the app deploys automatically
// once that first build succeeds. Track build progress via List/Get.
//
// Honors Idempotency-Key. Returns codes.SourceConnectionNotFound when the
// connection doesn't exist or isn't owned by the caller,
// codes.SourceConnectionSuspended / codes.BuildSourceUnavailable when the
// installation is unusable, or codes.BuildProviderError on a provider failure.
func (s *BuildsService) CreateGitBuildApp(ctx context.Context, connID uint, req *types.CreateGitBuildAppRequest, opts ...WriteOption) (*types.CreateAppResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.CreateAppResponse
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/source-connections/%d/apps", connID), req, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// List returns the builds of a git-build app, most recent first, paginated.
func (s *BuildsService) List(ctx context.Context, appID uint, opts ...ListOption) ([]types.BuildResponse, *types.Meta, error) {
	q := resolveListOpts(opts)
	var out []types.BuildResponse
	meta, err := s.c.doList(ctx, "GET",
		withQuery(fmt.Sprintf("/api/v1/apps/%d/builds", appID), q), &out)
	if err != nil {
		return nil, nil, err
	}
	return out, meta, nil
}

// Get fetches a single build. Returns codes.BuildNotFound when the build
// doesn't exist for the app.
//
// The build log URL is no longer part of this response — fetch it on demand
// with GetLogURL when you're about to open the log.
func (s *BuildsService) Get(ctx context.Context, appID, buildID uint) (*types.BuildResponse, error) {
	var out types.BuildResponse
	_, _, err := s.c.do(ctx, "GET",
		fmt.Sprintf("/api/v1/apps/%d/builds/%d", appID, buildID), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLogURL returns a freshly-minted, short-lived presigned URL to the build's
// plain-text log. Call it only when about to open the log (the URL expires
// quickly). Returns codes.BuildLogNotAvailable when the build has no persisted
// log yet (still pending/running, or the upload failed), or codes.BuildNotFound
// when the build doesn't exist for the app.
func (s *BuildsService) GetLogURL(ctx context.Context, appID, buildID uint) (string, error) {
	var out types.BuildLogURLResponse
	_, _, err := s.c.do(ctx, "GET",
		fmt.Sprintf("/api/v1/apps/%d/builds/%d/log-url", appID, buildID), nil, nil, &out)
	if err != nil {
		return "", err
	}
	return out.LogURL, nil
}

// Rebuild triggers a fresh build of the app's configured branch HEAD and
// returns the newly-created build. Returns codes.BuildAlreadyRunning when a
// build is already pending/running, or codes.BuildSourceUnavailable when the
// connection/registry is unusable.
func (s *BuildsService) Rebuild(ctx context.Context, appID uint, opts ...WriteOption) (*types.BuildResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.BuildResponse
	_, _, err = s.c.do(ctx, "POST", fmt.Sprintf("/api/v1/apps/%d/builds", appID), nil, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Cancel stops an in-flight build and returns the updated build (status
// "canceled"). A no-op-safe call: canceling an already-terminal build returns
// codes.BuildNotFound only when the build doesn't exist.
func (s *BuildsService) Cancel(ctx context.Context, appID, buildID uint, opts ...WriteOption) (*types.BuildResponse, error) {
	wopts, err := resolveWriteOpts(opts)
	if err != nil {
		return nil, err
	}
	var out types.BuildResponse
	_, _, err = s.c.do(ctx, "POST",
		fmt.Sprintf("/api/v1/apps/%d/builds/%d/cancel", appID, buildID), nil, &wopts, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
