// Package client is the HTTP client for the Kumo API. It composes with
// github.com/kumobase/kumo-go/types (wire DTOs) and
// github.com/kumobase/kumo-go/codes (stable wire error codes) — those
// packages are the public contract; this package is the transport layer
// over them.
//
// Quickstart:
//
//	c, err := client.New("https://api.kumo.run",
//	    client.WithAPIKey("kumo_sk_…"))
//	if err != nil { … }
//
//	app, err := c.Apps().Create(ctx, &types.CreateAppRequest{
//	    BaseCreateApp: types.BaseCreateApp{Name: "demo", …},
//	})
//	if errors.Is(err, client.ErrIdempotencyKeyConflict) { … }
//	if client.IsCode(err, codes.AppQuotaExceeded) { … }
package client

import (
	"errors"
	"fmt"

	"github.com/kumobase/kumo-go/codes"
)

// APIError is the typed wrapper returned for every non-2xx HTTP response.
// It carries the server's stable wire code (UPPER_SNAKE_CASE; see package
// github.com/kumobase/kumo-go/codes), the human-readable message, the HTTP
// status, and the raw body for debugging.
//
// Callers branch on Code (stable across SDK releases). Message can evolve.
type APIError struct {
	StatusCode int    // HTTP status code (4xx or 5xx)
	Code       string // server-emitted Code constant; "" if the response had no body or the body was malformed
	Message    string // human-readable Message from the server, or a synthesised fallback
	Body       []byte // raw response body — useful when Code is empty or the failure is unexpected
}

// Error returns a compact human-readable string. Callers should not parse it.
func (e *APIError) Error() string {
	if e == nil {
		return "<nil APIError>"
	}
	if e.Code != "" {
		return fmt.Sprintf("kumo: %s (%s, http %d)", e.Message, e.Code, e.StatusCode)
	}
	return fmt.Sprintf("kumo: http %d: %s", e.StatusCode, e.Message)
}

// Unwrap exposes the sentinel for errors.Is matching. Returns nil if the
// error doesn't map to one of the well-known cross-cutting sentinels.
func (e *APIError) Unwrap() error {
	if e == nil {
		return nil
	}
	switch e.Code {
	case codes.IdempotencyKeyConflict:
		return ErrIdempotencyKeyConflict
	case codes.IdempotencyInProgress:
		return ErrIdempotencyInProgress
	case codes.ETagMismatch:
		return ErrETagMismatch
	case codes.ValidationFailed:
		return ErrValidationFailed
	case codes.InvalidFilterCombination:
		return ErrInvalidFilterCombination
	}
	return nil
}

// Cross-cutting sentinel errors. Match with errors.Is for code-aware
// branching without string comparison:
//
//	if errors.Is(err, client.ErrETagMismatch) {
//	    // refetch, get new ETag, retry the PATCH
//	}
//
// Module-specific codes (codes.AppNotFound, codes.VolumeResizing, etc.)
// don't have dedicated sentinels — use the IsCode predicate or branch on
// (*APIError).Code directly.
var (
	ErrIdempotencyKeyConflict   = errors.New("idempotency key conflict: same key, different body")
	ErrIdempotencyInProgress    = errors.New("idempotency in progress: same key still running")
	ErrETagMismatch             = errors.New("etag mismatch: resource was modified concurrently")
	ErrValidationFailed         = errors.New("validation failed")
	ErrInvalidFilterCombination = errors.New("invalid filter combination")
)

// IsCode reports whether err is an *APIError carrying the given code.
// Returns false on nil or non-APIError errors. Use it for module-specific
// codes that don't have a dedicated sentinel:
//
//	if client.IsCode(err, codes.AppDeploymentInProgress) { … }
func IsCode(err error, code string) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	return apiErr.Code == code
}

// IsNotFound is a convenience: true for any APIError whose code is a
// well-known "not found" code (AppNotFound, SecretNotFound, …) OR whose
// StatusCode is 404. The latter catches 404s from proxies that strip the
// JSON body.
func IsNotFound(err error) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	if apiErr.StatusCode == 404 {
		return true
	}
	switch apiErr.Code {
	case codes.AppNotFound, codes.AppOperationNotFound, codes.AppCustomDomainNotFound, codes.AppRegistryCredentialNotFound,
		codes.SecretNotFound,
		codes.VolumeNotFound, codes.StorageTierNotFound, codes.StorageClassNotFound,
		codes.InstanceNotFound, codes.PlanNotFound, codes.ProviderNotFound,
		codes.OrgNotFound,
		codes.RegistryRepositoryNotFound, codes.RegistryManifestNotFound, codes.RegistryBlobNotFound,
		codes.APIKeyNotFound:
		return true
	}
	return false
}

// IsConflict is a convenience: true for any APIError whose code is a
// well-known conflict code OR whose StatusCode is 409.
func IsConflict(err error) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	if apiErr.StatusCode == 409 {
		return true
	}
	switch apiErr.Code {
	case codes.IdempotencyKeyConflict, codes.IdempotencyInProgress,
		codes.AppDeploymentInProgress, codes.AppAlreadyStopped,
		codes.AppCustomDomainExists, codes.AppDomainAlreadyInUse, codes.AppVolumeConflict,
		codes.AppQuotaExceeded,
		codes.SecretInUse,
		codes.VolumeAttached, codes.VolumePermanentlyAttached, codes.VolumeResizing,
		codes.ActionInProgress, codes.AutoRenewAlreadyCancelled,
		codes.OrgSlugTaken, codes.OrgMaxOrganizationsReached, codes.OrgHasRepos, codes.OrgCannotDeleteDefault,
		codes.RegistryRepositoryAlreadyExists, codes.RegistryTagImmutable, codes.RegistryMaxRepositoriesReached:
		return true
	}
	return false
}
