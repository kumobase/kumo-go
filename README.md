# kumo-go

Public Go SDK for the [Kumo](https://kumobase.com) platform — the wire-level
contract (DTOs, stable error codes, and HTTP client) for Kumo's customer-facing
API.

## What's in here

```
types/    Request and response DTOs for every user-facing /api/v1/* endpoint.
codes/    Stable wire error codes (UPPER_SNAKE_CASE) returned in the Code
          field of error responses. Clients should branch on Code, not Message.
client/   Authenticated API client and restricted unauthenticated catalogue client.
version/  SDK version constant for User-Agent / compat checks.
```

## What's **not** here

- Admin endpoints (`/api/v1/admin/*`) — internal-only, never exposed.
- Internal model structs (GORM rows, validation tags, etc.).

## Install

```bash
go get github.com/kumobase/kumo-go
```

## Quickstart

Create an authenticated client with exactly one API key or JWT:

```go
import (
	"github.com/kumobase/kumo-go/client"
)

c, err := client.New("https://api.kumo.run",
	client.WithAPIKey("kumo_sk_…"))
```

For signed-out pricing and discovery, create a client that cannot access any
protected endpoint and never sends an `Authorization` header:

```go
c, err := client.NewPublic("https://api.kumo.run")
plans, err := c.Apps().ListPlans(ctx)
```

`NewPublic` permits only exact GET requests to these paths (query parameters
are allowed):

- `/api/v1/apps/plans`
- `/api/v1/vps/regions`
- `/api/v1/vps/providers`
- `/api/v1/vps/plans`
- `/api/v1/volumes/plans`
- `/api/v1/registry/plans`
- `/api/v1/packages/plans`
- `/api/v1/runners/plans`

Authentication options are rejected by `NewPublic`. Other options, including
custom HTTP clients, retries, loggers, user agents, and `WithBaseURL`, remain
available. Disallowed calls fail locally with `client.ErrPublicClientRestricted`.

## Compatibility

SDK SemVer follows Go module rules:

- **Patch (v0.1.0 → v0.1.1)** — bug fixes; no field changes.
- **Minor (v0.1.x → v0.2.0)** — additive: new optional fields, new endpoints,
  new codes. Pre-1.0, minor versions MAY include breaking changes; consumers
  should pin tightly until the SDK reaches v1.0.0.
- **Major (v1.x → v2.0.0)** — wire-breaking; rare. New major shipped as
  `github.com/kumobase/kumo-go/v2`.

Kumo pins an explicit kumo-go release. Older SDK versions remain compatible
unless a server endpoint explicitly documents a newer minimum.

## License

Apache-2.0. See [LICENSE](./LICENSE).
