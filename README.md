# kumo-go

Public Go SDK for the [Kumo](https://kumobase.com) platform — the wire-level
contract (DTOs and stable error codes) for everything a customer-facing API
key (`kumo_sk_…`) can call.

## What's in here

```
types/    Request and response DTOs for every user-facing /api/v1/* endpoint.
codes/    Stable wire error codes (UPPER_SNAKE_CASE) returned in the Code
          field of error responses. Clients should branch on Code, not Message.
version/  SDK version constant for User-Agent / compat checks.
```

## What's **not** here

- HTTP client — coming in a follow-up release.
- Admin endpoints (`/api/v1/admin/*`) — internal-only, never exposed.
- Internal model structs (GORM rows, validation tags, etc.).

## Install

```bash
go get github.com/kumobase/kumo-go
```

## Quickstart

Decode an `Idempotency-Key`-friendly response:

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/kumobase/kumo-go/codes"
    "github.com/kumobase/kumo-go/types"
)

func main() {
    body, _ := os.ReadFile("response.json")
    var env types.StructureResponse
    _ = json.Unmarshal(body, &env)

    if env.Code == codes.AppNotFound {
        fmt.Println("app not found")
        return
    }
    // …
}
```

## Compatibility

SDK SemVer follows Go module rules:

- **Patch (v0.1.0 → v0.1.1)** — bug fixes; no field changes.
- **Minor (v0.1.x → v0.2.0)** — additive: new optional fields, new endpoints,
  new codes. Pre-1.0, minor versions MAY include breaking changes; consumers
  should pin tightly until the SDK reaches v1.0.0.
- **Major (v1.x → v2.0.0)** — wire-breaking; rare. New major shipped as
  `github.com/kumobase/kumo-go/v2`.

Server and SDK ship together: server release `vX.Y.Z` is built against
kumo-go `vX.Y.Z`. Older SDK versions are accepted as long as the server's
`MIN_SDK_VERSION` allows.

## License

Apache-2.0. See [LICENSE](./LICENSE).
