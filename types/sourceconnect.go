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
//
// AvatarURL and AccountID describe the GitHub account the App is installed on.
// RepoSelection is the install's grant scope: "all" (every repo, current and
// future) or "selected" (an explicit subset). RepoCount is the size of that
// subset and is only meaningful for "selected" — it is omitted/zero for "all",
// where the count is unbounded. These reflect the install as last synced via
// the installation webhook; all are omitted when unknown.
//
// GitLab is a non-nil sub-object only when Provider == "gitlab"; it carries the
// GitLab-specific shape (instance, namespace, group/project) that does not fit
// the GitHub-centric top-level fields. For "gitlab" rows the GitHub-only fields
// (InstallationID, AccountType, ManageURL, repo summary) are not populated —
// read GitLab instead.
type SourceConnectionResponse struct {
	ID             uint                      `json:"id"`
	Provider       SourceProvider            `json:"provider"`
	InstallationID int64                     `json:"installation_id"`
	AccountLogin   string                    `json:"account_login"`
	AccountType    string                    `json:"account_type"` // "User" | "Organization"
	AccountID      int64                     `json:"account_id,omitempty"`
	AvatarURL      string                    `json:"avatar_url,omitempty"`
	ManageURL      string                    `json:"manage_url,omitempty"`
	RepoSelection  string                    `json:"repo_selection,omitempty"` // "all" | "selected"
	RepoCount      int                       `json:"repo_count,omitempty"`
	Status         SourceConnectionStatus    `json:"status"`
	AppKind        string                    `json:"app_kind"` // "build" | "runner"
	GitLab         *GitLabConnectionResponse `json:"gitlab,omitempty"`
	CreatedAt      time.Time                 `json:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at"`
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
//
// DisplayName is the namespace's human name ("My Team"), distinct from
// NamespacePath; it falls back to the path for connections created before the
// name was captured. InstanceKind tags where the GitLab server lives ("saas"
// for gitlab.com, "self_managed" for any other host) and Host is that host
// (e.g. "gitlab.com", "gitlab.acme.com") — both derived from BaseURL so clients
// can render a "gitlab.com vs self-hosted" badge without parsing URLs. AvatarURL
// and Visibility ("private" | "internal" | "public") mirror the namespace as
// last seen on connect; both are omitted when unknown.
type GitLabConnectionResponse struct {
	ID            uint                   `json:"id"`
	Provider      SourceProvider         `json:"provider"` // always "gitlab"
	InstanceID    uint                   `json:"instance_id"`
	BaseURL       string                 `json:"base_url"`
	InstanceKind  string                 `json:"instance_kind,omitempty"` // "saas" | "self_managed"
	Host          string                 `json:"host,omitempty"`
	Kind          GitLabNamespaceKind    `json:"kind"` // "group" | "project"
	NamespaceID   int64                  `json:"namespace_id"`
	NamespacePath string                 `json:"namespace_path"`
	DisplayName   string                 `json:"display_name"`
	AvatarURL     string                 `json:"avatar_url,omitempty"`
	Visibility    string                 `json:"visibility,omitempty"` // "private" | "internal" | "public"
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
