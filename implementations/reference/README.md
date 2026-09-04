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
- `internal/payment/store.go` + `service.go` — an in-memory `Store` and a `Service.Create` that
  enforces [specs/payments/payment-contract.md](../../specs/payments/payment-contract.md)
  "Idempotency": concurrent calls with the same `IdempotencyKey` return the same `Payment`,
  never a second one (covered by a `-race`-clean test).
- `Service.ApplyTransition` — applies canonical state transitions per
  [specs/state-machines/payment-lifecycle.md](../../specs/state-machines/payment-lifecycle.md)
  "Idempotency and Retries": a duplicate delivery of the current status is a no-op, a valid
  transition applies normally, and a genuinely conflicting update (e.g. `FAILED` claimed for an
  already-`SUCCESS` payment) returns a `*TransitionError` without mutating anything.

- `internal/simulator/` — the `SIMULATOR` provider from
  [specs/scenarios/scenario-format.md](../../specs/scenarios/scenario-format.md).
  `Simulator.Initiate` drives `payment.Service` through
  `CREATED → PENDING → SUCCESS|FAILED`. Only the `success` and `failure` scenarios are
  implemented — `TIMEOUT`, `DUPLICATE_CALLBACK`, `OUT_OF_ORDER`, and `INVALID_SIGNATURE` need
  real delay/callback-timing machinery this increment doesn't build yet. A request with an
  unknown scenario or the wrong `Provider.ID` is rejected before any `Payment` is created.

Not yet implemented: the four deferred scenario outcomes above, webhook handling, and the REST
contract — see [ROADMAP.md](../../ROADMAP.md) Phase 1. `Service`'s errors
(`ErrMissingIdempotencyKey`, `ErrPaymentNotFound`, `TransitionError`) and `simulator`'s
(`ErrWrongProvider`, `ErrUnknownScenario`) are provisional and package-local, not the canonical
error taxonomy — see [specs/errors/README.md](../../specs/errors/README.md), still `TODO(ADR)`.

**Known limitation:** two concurrent `Initiate` calls for the same brand-new `IdempotencyKey`
can race — see the doc comment on `Simulator.Initiate` for why, and why sequential replay (the
case the spec actually describes) is unaffected.

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
