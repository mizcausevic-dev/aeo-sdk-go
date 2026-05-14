# aeo-sdk-go

Go SDK for the [AEO Protocol v0.1](https://github.com/mizcausevic-dev/aeo-protocol-spec) — parse, build, validate, and fetch AEO declaration documents. **Zero non-stdlib dependencies.**

[![Go Reference](https://pkg.go.dev/badge/github.com/mizcausevic-dev/aeo-sdk-go.svg)](https://pkg.go.dev/github.com/mizcausevic-dev/aeo-sdk-go)

## Install

```bash
go get github.com/mizcausevic-dev/aeo-sdk-go
```

## Quickstart

```go
package main

import (
    "context"
    "fmt"
    "log"

    aeo "github.com/mizcausevic-dev/aeo-sdk-go"
)

func main() {
    // Fetch and parse from a live well-known URL
    doc, err := aeo.FetchWellKnown("https://mizcausevic-dev.github.io")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(doc.Entity.Name)                              // "Miz Causevic"
    fmt.Println(doc.ClaimIDs())                               // [current-role location ...]
    fmt.Println(doc.FindClaim("years-experience").Value)      // 30

    // Parse from bytes
    doc2, _ := aeo.ParseDocumentString(`{"aeo_version":"0.1", ...}`)

    // Custom HTTP client with context
    ctx := context.Background()
    client := aeo.DefaultClient()
    _, _ = client.FetchWellKnown(ctx, "https://example.com")

    _ = doc2
}
```

## What it does

- **Parse** — `ParseDocument(bytes)`, `ParseDocumentString(s)`, `LoadDocument(path)` — all use `json.DisallowUnknownFields` for strict conformance
- **Build** — `Document`, `Entity`, `Authority`, `Claim`, `Verification`, `CitationPreferences`, `AnswerConstraints`, `Audit` are all exported struct types with proper JSON tags
- **Serialize** — `doc.Marshal()` returns pretty-printed JSON
- **Fetch** — `aeo.FetchWellKnown(origin)` does HTTP discovery against `/.well-known/aeo.json` with `Accept: application/aeo+json, application/json`. Custom `Client` lets you bring your own `*http.Client` and pass a `context.Context`.
- **Query** — `doc.ClaimIDs()` and `doc.FindClaim(id)` helpers

## Conformance

Supports the AEO Protocol at **conformance Level 1 (Declare)**. Signature verification (L2) and audit-endpoint posting (L3) deferred to v0.2.

## Dependencies

None. Uses only the standard library: `encoding/json`, `net/http`, `context`, `io`.

## Development

```bash
go vet ./...
go test -race -v ./...
go build ./...
```

## Specification

Full spec at [github.com/mizcausevic-dev/aeo-protocol-spec](https://github.com/mizcausevic-dev/aeo-protocol-spec).

## License

MIT-licensed. Free for commercial and non-commercial use with attribution. The AEO Protocol specification this SDK implements is also MIT (see [aeo-protocol-spec](https://github.com/mizcausevic-dev/aeo-protocol-spec)).

## Kinetic Gain Protocol Suite

| Spec | Implementation |
|---|---|
| [AEO Protocol](https://github.com/mizcausevic-dev/aeo-protocol-spec) | [aeo-sdk-python](https://github.com/mizcausevic-dev/aeo-sdk-python) · [aeo-sdk-typescript](https://github.com/mizcausevic-dev/aeo-sdk-typescript) · [aeo-sdk-rust](https://github.com/mizcausevic-dev/aeo-sdk-rust) · **aeo-sdk-go** (this) · [aeo-cli](https://github.com/mizcausevic-dev/aeo-cli) · [aeo-crawler](https://github.com/mizcausevic-dev/aeo-crawler) |
| [Prompt Provenance](https://github.com/mizcausevic-dev/prompt-provenance-spec) | — |
| [Agent Cards](https://github.com/mizcausevic-dev/agent-cards-spec) | — |
| [AI Evidence Format](https://github.com/mizcausevic-dev/ai-evidence-format-spec) | — |
| [MCP Tool Cards](https://github.com/mizcausevic-dev/mcp-tool-card-spec) | — |

---

**Connect:** [LinkedIn](https://www.linkedin.com/in/mirzacausevic/) · [Kinetic Gain](https://kineticgain.com) · [Medium](https://medium.com/@mizcausevic/) · [Skills](https://mizcausevic.com/skills/)
