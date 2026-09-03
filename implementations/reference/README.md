# implementations/reference/ — The Reference Implementation

Demonstrates [specs/](../../specs/README.md) working end-to-end. See
[implementations/README.md](../README.md) for the rules governing this directory.

## Status

**In progress — Phase 1.** Written in **Go**, per
[ADR 0002](../../adr/0002-reference-implementation-language-go.md). So far:

- `internal/payment/types.go` — the canonical domain types from
  [specs/payments/payment-contract.md](../../specs/payments/payment-contract.md) (`Payment`,
  `PaymentRequest`, `Money`, etc.)
- `internal/payment/lifecycle.go` — the canonical state machine from
  [specs/state-machines/payment-lifecycle.md](../../specs/state-machines/payment-lifecycle.md)

Not yet implemented: orchestration, the simulator, webhook handling, and the REST contract —
see [ROADMAP.md](../../ROADMAP.md) Phase 1.

`internal/` is deliberate: this package is not meant to be imported by `adapters/` or `sdks/` —
per [ARCHITECTURE.md §8](../../ARCHITECTURE.md#8-reference-implementation-boundary), the
reference implementation demonstrates the specification, it does not define an API other code
should depend on.

## Requirements

- Go 1.27 or later (see [go.mod](go.mod))

## Build, Test, Format

```bash
cd implementations/reference
go build ./...
go test ./...
gofmt -l .    # should print nothing; run `gofmt -w .` to fix
go vet ./...
```

`go build`/`vet`/`test`/`gofmt` run in CI on every PR via the `go` job in
[.github/workflows/ci.yml](../../.github/workflows/ci.yml) — see
[docs/development/ci.md](../../docs/development/ci.md). `go test` specifically also runs from
the repo root via `make test` (see [scripts/test.sh](../../scripts/test.sh), which dispatches to
`go test ./...` in every directory under `implementations/`, `adapters/`, or `sdks/` containing
a `go.mod`) — `go build`/`vet`/`gofmt`, and non-Go toolchains, are CI-job-specific for now.
