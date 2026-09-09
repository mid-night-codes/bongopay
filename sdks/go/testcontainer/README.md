# sdks/go/testcontainer/

A [testcontainers-go](https://github.com/testcontainers/testcontainers-go) helper that boots
[implementations/reference/](../../../implementations/reference/README.md) as a real Docker
container — built from [deploy/docker/reference.Dockerfile](../../../deploy/docker/reference.Dockerfile)
— for use in Go test suites that want to exercise BongoPay against a real running server instead
of hand-rolling `exec.Command`/port-polling lifecycle management, or a mock.

## Why a Separate Module

This is the first external dependency in the project's history — everything else, including
[sdks/go](../README.md)'s own client, is stdlib-only by design (see
[ADR 0002](../../../adr/0002-reference-implementation-language-go.md)). Someone importing
`sdks/go` to make API calls has no reason to also pull in Docker-orchestration machinery, so
this package has its own `go.mod`: only importing `sdks/go/testcontainer` specifically takes on
`testcontainers-go` and its (fairly large) transitive dependency tree.
[docs/development/dependency-policy.md](../../../docs/development/dependency-policy.md)'s
checklist, answered:

- **Can stdlib solve this instead?** Not without reimplementing container build/start/
  readiness/cleanup — exactly the complexity a dedicated, widely-used library exists to get
  right.
- **Actively maintained, compatible license?** Yes — `testcontainers-go`, MIT, backed by the
  Testcontainers organization.
- **Meaningfully improves maintainability vs. convenience?** Yes — container lifecycle and
  readiness-waiting is genuine complexity, not boilerplate this project would casually
  reinvent well.
- **Constrains `specs/`/`contracts/` neutrality?** No — isolated to this one opt-in nested
  module, nowhere near either.

## Usage

```go
import (
    "context"

    bongopay "github.com/mid-night-codes/bongopay/sdks/go"
    "github.com/mid-night-codes/bongopay/sdks/go/testcontainer"
)

func TestSomething(t *testing.T) {
    ctx := context.Background()
    container, err := testcontainer.RunReference(ctx, "/path/to/bongopay/repo")
    if err != nil {
        t.Fatal(err)
    }
    defer container.Terminate(ctx)

    client := bongopay.New(container.BaseURL)
    // ... use client as normal
}
```

`repoRoot` is the BongoPay repository root (the Dockerfile's build context) — resolve it
relative to your test file, or skip the test if it isn't found there, the way
[reference_test.go](reference_test.go) does.

## A Real Race, Found by Actually Running This

The first version waited on `cmd/server`'s own startup log line
(`wait.ForLog("bongopay reference server listening")`). That log line is printed *just before*
`http.ListenAndServe` is called, and even once the process is genuinely listening, Docker's own
port-mapping/proxy can lag slightly behind — this flaked with a connection `EOF` on the very
first real run against an actual container, not in theory. Switched to
`wait.ForListeningPort("8080/tcp")`, which checks that the host-mapped port is actually
reachable — the real thing a caller needs, since they connect the same way — and re-ran it
several times (`go test -count=1`) to confirm it wasn't a coincidence before trusting it.

## Testing

`reference_test.go` builds the real image and starts a real container, then drives it via
`sdks/go`'s client. Skips (not fails) if Docker isn't reachable or the repo root can't be
resolved, since those are environment facts, not bugs.

```bash
go test ./...
```
