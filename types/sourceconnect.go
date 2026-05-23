package types

import "time"

// SourceProvider identifies the git host backing a source connection. Only
// GitHub is supported today; the field is an enum so new providers can ship
// without breaking older clients.
type SourceProvider string

const (
	SourceProviderGitHub SourceProvider = "github"
)

// SourceConnectionStatus reflects whether a connection can currently be used
// to fetch source. "suspended" means the provider-side install was suspended
// (e.g. by the account owner) and must be reactivated before use.
type SourceConnectionStatus string

const (
	SourceConnectionStatusActive    SourceConnectionStatus = "active"
	SourceConnectionStatusSuspended SourceConnectionStatus = "suspended"
)

// SourceConnectionResponse is the wire shape for a connected git account,
// returned by GET /api/v1/source-connections.
//
// InstallationID is the provider-side installation identifier. It is NOT a
// secret (it grants nothing without the platform's App private key) and is
// exposed so the frontend can deep-link to the provider's "configure access"
// page for adding/removing repositories.
type SourceConnectionResponse struct {
	ID             uint                   `json:"id"`
	Provider       SourceProvider         `json:"provider"`
	InstallationID int64                  `json:"installation_id"`
	AccountLogin   string                 `json:"account_login"`
	AccountType    string                 `json:"account_type"` // "User" | "Organization"
	Status         SourceConnectionStatus `json:"status"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// SourceRepoResponse is one repository a connection has been granted access
// to, returned by GET /api/v1/source-connections/:id/repos and used to
// populate the repo picker. The list is fetched live from the provider, so it
// always reflects the account owner's current grant. ID is the provider's
// numeric repository id.
type SourceRepoResponse struct {
	ID            int64  `json:"id"`
	FullName      string `json:"full_name"` // "owner/repo"
	Private       bool   `json:"private"`
	DefaultBranch string `json:"default_branch"`
}
