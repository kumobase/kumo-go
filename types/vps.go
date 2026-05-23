package types

// VPSStatus enumerates the lifecycle states of a rented VPS instance.
// The "deleted" state is filtered server-side from list/get responses, so
// clients should not need to handle it in practice — included here for
// completeness of the wire enum.
type VPSStatus string

const (
	VPSStatusProvisioning VPSStatus = "provisioning"
	VPSStatusRunning      VPSStatus = "running"
	VPSStatusStopped      VPSStatus = "stopped"
	VPSStatusExpired      VPSStatus = "expired"
	VPSStatusDeleted      VPSStatus = "deleted"
)

// VPSActionStatus tracks an in-flight async action on a VPS instance.
// Empty string means no action is running; clients poll
// GET /api/v1/vps/servers/:id and watch for the field to flip back to "".
type VPSActionStatus string

const (
	VPSActionStatusNone         VPSActionStatus = ""
	VPSActionStatusRebooting    VPSActionStatus = "rebooting"
	VPSActionStatusReinstalling VPSActionStatus = "reinstalling"
	VPSActionStatusPoweringOn   VPSActionStatus = "powering_on"
	VPSActionStatusPoweringOff  VPSActionStatus = "powering_off"
)

// RentServerRequest is the body for POST /api/v1/vps/servers. Honors
// Idempotency-Key — retrying the same key + body replays the cached response
// rather than provisioning a second server.
type RentServerRequest struct {
	Provider string `json:"provider"`
	Region   string `json:"region"`
	Plan     string `json:"plan"`
	Name     string `json:"name"` // 1..100 chars
}

// UpdateServerNameRequest is the body for PATCH /api/v1/vps/servers/:id/name.
type UpdateServerNameRequest struct {
	Name string `json:"name"` // 1..100 chars
}

// VPSRegionResponse is one entry returned by GET /api/v1/vps/regions.
type VPSRegionResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PublicVPSPlanResponse is the customer-facing plan DTO. Internal pricing
// detail (base_price, margin_usd, exchange_rate) is intentionally absent;
// only the final SellingPrice is exposed.
//
// SellingPrice is a decimal string (e.g. "4.99") — parse with
// strconv.ParseFloat or a decimal library if you need arithmetic.
type PublicVPSPlanResponse struct {
	ProviderName         string        `json:"provider_name"`
	PlanID               uint          `json:"plan_id"`
	ExternalPlanID       string        `json:"external_plan_id"`
	Name                 string        `json:"name"`
	CPU                  int           `json:"cpu"`
	Memory               int           `json:"memory"`
	Disk                 int           `json:"disk"`
	Egress               int           `json:"egress,omitempty"`
	MaxOutboundBandwidth *int          `json:"max_outbound_bandwidth,omitempty"`
	SellingPrice         string        `json:"selling_price"`
	Availability         *Availability `json:"availability,omitempty"`
}

// VPSServerResponse is the detail returned by GET /api/v1/vps/servers/:id
// and items of GET /api/v1/vps/servers.
//
// ExpiresAt / CreatedAt / ActionStatusUpdatedAt are RFC3339 strings (the
// server serialises time.Time directly; using string here lets clients
// inspect raw values without needing time.Parse for fields they may not use).
type VPSServerResponse struct {
	ID              uint   `json:"id"`
	ExternalID      string `json:"external_id,omitempty"`
	DisplayName     string `json:"display_name,omitempty"`
	DisplayProvider string `json:"display_provider"`
	RegionID        string `json:"region_id"`
	OS              string `json:"os,omitempty"`
	Status          string `json:"status"`
	IPAddress       string `json:"ip_address,omitempty"`
	SSHPort         int    `json:"ssh_port"`
	AutoRenew       bool   `json:"auto_renew"`
	ExpiresAt       string `json:"expires_at,omitempty"`
	CreatedAt       string `json:"created_at"`

	SSHSetupCompleted bool `json:"ssh_setup_completed"`

	// Async action tracking. ActionStatus is "" when idle, otherwise one of
	// the VPSActionStatus* constants. Poll GET /api/v1/vps/servers/:id and
	// watch for ActionStatus to flip back to "" — ActionError will be set
	// if the action terminally failed.
	ActionStatus          string `json:"action_status"`
	ActionStatusUpdatedAt string `json:"action_status_updated_at,omitempty"`
	ActionError           string `json:"action_error,omitempty"`

	Plan *PublicVPSPlanResponse `json:"plan,omitempty"`
}
