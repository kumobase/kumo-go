package types

import "time"

// SourceProvider identifies the git host backing a source connection. The
// field is an enum so new providers ship without breaking older clients.
type SourceProvider string

const (
	SourceProviderGitHub SourceProvider = "github"
	SourceProviderGitLab SourceProvider = "gitlab"
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
// secret (it grants nothing without the platform's App private key).
//
// ManageURL is the server-computed deep-link to the provider's "configure
// access" page, where the user adds/removes repositories or switches to "all
// repositories". Changing the grant there needs no reconnect — the platform
// reads the repo list live. Omitted when it can't be built (e.g. a future
// non-GitHub provider).
//
// AppKind identifies which Kumo GitHub App owns this install: "build" (the
// git-build/app product) or "runner" (VM CI-runners). Both products record
// installs in the same surface; clients branch on this to group/filter.
type SourceConnectionResponse struct {
	ID             uint                   `json:"id"`
	Provider       SourceProvider         `json:"provider"`
	InstallationID int64                  `json:"installation_id"`
	AccountLogin   string                 `json:"account_login"`
	AccountType    string                 `json:"account_type"` // "User" | "Organization"
	ManageURL      string                 `json:"manage_url,omitempty"`
	Status         SourceConnectionStatus `json:"status"`
	AppKind        string                 `json:"app_kind"` // "build" | "runner"
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

// GitLabNamespaceKind discriminates what a GitLab runner connection is scoped
// to: a whole group (one webhook covers all child projects) or a single
// project.
type GitLabNamespaceKind string

const (
	GitLabNamespaceGroup   GitLabNamespaceKind = "group"
	GitLabNamespaceProject GitLabNamespaceKind = "project"
)

// GitLabInstanceResponse is a GitLab instance the user may connect against,
// returned by GET /api/v1/source-connections/gitlab/instances. gitlab.com is
// always present (BaseURL "https://gitlab.com"); self-managed instances are
// added by the user with an OAuth application registered on that server. The
// OAuth client secret is never echoed back.
type GitLabInstanceResponse struct {
	ID      uint   `json:"id"`
	BaseURL string `json:"base_url"`
	// Status is "ready" once an OAuth app is configured and reachable, or
	// "unreachable" if a connectivity/credential probe last failed.
	Status string `json:"status"`
}

// GitLabConnectionResponse is the wire shape for a connected GitLab group or
// project, returned by the GitLab source-connection endpoints. It is the GitLab
// analogue of SourceConnectionResponse: GitLab has no App-installation object,
// so a connection is OAuth-backed and scoped to one namespace (group/project).
//
// Namespace is the full path ("acme" for a group, "acme/api" for a project);
// NamespaceID is GitLab's numeric id. WebURL deep-links to the namespace.
// Status mirrors SourceConnectionStatus ("active" | "suspended"); a suspended
// connection means the OAuth grant was revoked/expired and must be reconnected.
type GitLabConnectionResponse struct {
	ID            uint                   `json:"id"`
	Provider      SourceProvider         `json:"provider"` // always "gitlab"
	InstanceID    uint                   `json:"instance_id"`
	BaseURL       string                 `json:"base_url"`
	Kind          GitLabNamespaceKind    `json:"kind"` // "group" | "project"
	NamespaceID   int64                  `json:"namespace_id"`
	NamespacePath string                 `json:"namespace_path"`
	DisplayName   string                 `json:"display_name"`
	WebURL        string                 `json:"web_url,omitempty"`
	Status        SourceConnectionStatus `json:"status"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// GitLabNamespaceResponse is one selectable group or project in the connect
// picker, returned by GET /api/v1/source-connections/gitlab/namespaces. The
// list is fetched live from GitLab using the user's OAuth token, so it reflects
// their current access. ID is GitLab's numeric id; FullPath is what the user
// puts no-where — it's used to create the webhook on connect.
type GitLabNamespaceResponse struct {
	ID       int64               `json:"id"`
	Kind     GitLabNamespaceKind `json:"kind"`
	FullPath string              `json:"full_path"` // "acme" or "acme/api"
	Name     string              `json:"name"`
	WebURL   string              `json:"web_url,omitempty"`
}
