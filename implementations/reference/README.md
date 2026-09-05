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
  `Simulator.Initiate` drives `payment.Service` through `CREATED → PENDING → SUCCESS|FAILED`
  for the `success`/`failure` scenarios, via `Service.CreateAndAdvance`, which performs
  create-or-lookup plus every transition under a single lock acquisition so concurrent
  `Initiate` calls for the same brand-new `IdempotencyKey` can't interleave mid-sequence. A
  request with an unknown scenario or the wrong `Provider.ID` is rejected before any `Payment`
  is created.
- `internal/simulator/callback.go` — `specs/providers/adapter-contract.md`'s `parseCallback`
  and `verifyCallback` capabilities: `Callback`, `ParseCallback`, and a `CallbackVerifier` doing
  HMAC-SHA256 over the raw callback body (a simulator-specific signing scheme for exercising
  verification *behavior* — not any real provider's actual scheme, which is adapter-specific).
  `Simulator.HandleCallback` verifies before ever calling `Service.ApplyTransition`, per
  [ARCHITECTURE.md §12](../../ARCHITECTURE.md#12-security-boundaries). Reusing
  `ApplyTransition`'s idempotency semantics through this same path is what implements
  `scenario-format.md`'s `DUPLICATE_CALLBACK` (a repeat is a no-op) and `OUT_OF_ORDER` (a stale
  conflicting claim is rejected, not silently applied) outcomes, plus `INVALID_SIGNATURE` (a
  wrong signature or a tampered body is rejected before any mutation). `Simulator.Initiate`
  itself does not yet route through `HandleCallback` — it's still synchronous end-to-end for
  `success`/`failure`; whether/how to make it callback-driven is an open design question for a
  later increment.

- `internal/httpapi/` + `cmd/server/` — an HTTP server implementing
  [contracts/openapi/bongopay.yaml](../../contracts/openapi/bongopay.yaml)'s `POST /payments`
  and `GET /payments/{id}` against stdlib `net/http` pattern routing (Go 1.22+ — no router
  dependency). `cmd/server` wires an `InMemoryStore` → `Service` → `Simulator` → `httpapi.Server`
  and listens on `:8080` by default. Every `PaymentRequest` reaches the simulator, since there
  are no adapters yet — `Simulator.Initiate`'s own `Provider.ID` check is what rejects anything
  else. Package-local errors map onto HTTP status codes per the contract's documented responses
  (missing idempotency key / wrong provider / unknown scenario → 400, not found → 404, a
  transition conflict → 409, anything else → 500, with the message withheld to avoid leaking
  internals). Run it and try it:

  ```bash
  cd implementations/reference
  go run ./cmd/server &
  curl -X POST localhost:8080/payments \
    -d '{"provider":{"id":"SIMULATOR"},"amount":{"value":5000,"currency":{"code":"TZS"}},"customerReference":{},"idempotencyKey":"demo-1"}'
  ```

Not yet implemented: `TIMEOUT` (needs real delay machinery) and wiring the `DUPLICATE_CALLBACK`/
`OUT_OF_ORDER`/`INVALID_SIGNATURE` behaviors into `Initiate`'s scenario selection (they're only
reachable via `HandleCallback` directly today, and there's no HTTP route for callback delivery
either — see [contracts/openapi/README.md](../../contracts/openapi/README.md) on why that's
deliberately not part of the canonical contract). See [ROADMAP.md](../../ROADMAP.md) Phase 1.
`Service`'s errors (`ErrMissingIdempotencyKey`, `ErrPaymentNotFound`, `TransitionError`) and
`simulator`'s (`ErrWrongProvider`, `ErrUnknownScenario`, `ErrInvalidCallbackSignature`) are
provisional and package-local, not the canonical error taxonomy — see
[specs/errors/README.md](../../specs/errors/README.md), still `TODO(ADR)`.

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
