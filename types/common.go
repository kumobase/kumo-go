// Package types contains the wire DTOs returned by every user-facing
// /api/v1/* endpoint a Kumo customer API key (kumo_sk_…) can call.
//
// Conventions:
//   - JSON tags here are the wire contract. Server and SDK must agree on
//     every field name; renaming requires a major version bump.
//   - Money / quota amounts are exposed as decimal strings (e.g. "4.99"),
//     not float64 — avoids rounding bugs and removes the shopspring/decimal
//     transitive dep from consumers.
//   - Validation tags from the server (validate:"required,min=3") are
//     intentionally omitted — validation is the server's concern. Consumers
//     that want client-side validation can layer it themselves.
//   - Fields that the server marks json:"-" (DeletedAt, KeyHash, …) are
//     omitted entirely. Admin-only fields (base_price, margin_*) are never
//     included; those endpoints are not part of the public API.
//   - Unknown JSON fields are tolerated by Go's encoder, so the server can
//     add new optional fields without breaking older SDK consumers.
package types

import "encoding/json"

// StructureResponse is the standard envelope every endpoint returns.
//
//   - Code is a stable, machine-readable string (UPPER_SNAKE_CASE).
//     Consumers should branch on Code, not Message. See package
//     github.com/kumobase/kumo-go/codes.
//   - Message is a human-readable summary and may evolve between releases.
//   - Data carries the endpoint-specific payload (a single DTO, an array
//     for list endpoints, or a status object for async POSTs).
//   - Meta is populated on paginated list responses.
type StructureResponse struct {
	Code    string          `json:"code,omitempty"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
	Meta    *Meta           `json:"meta,omitempty"`
}

// Meta holds pagination metadata for list endpoints. Server caps
// PageSize at 100; values above the cap are silently clamped server-side.
type Meta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// EnvironmentVariable is a key/value pair embedded in app and secret
// definitions. Empty Key or Value is rejected server-side.
type EnvironmentVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ValidationError describes one field-level validation failure.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorsResponse is the payload carried in StructureResponse.Data
// when an endpoint rejects a request with code VALIDATION_FAILED.
type ValidationErrorsResponse struct {
	Errors []*ValidationError `json:"errors"`
}

// Availability is a soft hint on public plan / tier listings indicating
// whether the platform currently has capacity for one minimal-footprint
// instance. It is best-effort UX — the create path enforces hard limits and
// will 503 on contention regardless of what was returned here. When platform
// quota cannot be resolved, the server omits this field rather than setting
// Available=false.
type Availability struct {
	Available bool   `json:"available"`
	Reason    string `json:"reason,omitempty"`
}

// AvailabilityReason* are the sentinel Reason values returned with
// Availability. UI uses them to localise the "sold out" badge copy.
const (
	AvailabilityReasonCPUFull     = "platform_cpu_full"
	AvailabilityReasonMemoryFull  = "platform_memory_full"
	AvailabilityReasonStorageFull = "platform_storage_full"
	AvailabilityReasonVPSFull     = "platform_vps_full"
)
