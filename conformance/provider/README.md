# conformance/provider/ — Provider Adapter Conformance Cases

Conformance cases every provider adapter must pass, per capability declared in
[specs/providers/adapter-contract.md](../../specs/providers/adapter-contract.md). This is the
suite referenced by [CONTRIBUTING.md](../../CONTRIBUTING.md#adding-a-provider): an adapter is
accepted once it passes this suite, not once it compiles.

## Status

**Specification-only, Phase 3 ("Provider Ecosystem"), with cases written earlier wherever they
double as simulator conformance** — see [ROADMAP.md](../../ROADMAP.md). No case files exist yet.

## Planned Coverage

Once written, cases here will cover, per declared capability in `ProviderCapabilities`:

- `initiate` — a `PaymentRequest` yields a `Payment` in a valid initial state
  ([specs/state-machines/payment-lifecycle.md](../../specs/state-machines/payment-lifecycle.md)),
  and is idempotent per its `idempotencyKey`.
- `queryStatus` — returns a status translated onto a valid canonical `PaymentStatus`, never a
  raw provider status.
- `refund` / `partialRefund` / `reverse` — only invoked on adapters declaring the capability;
  correctly reject when unsupported.
- `callbacks` — `verifyCallback` rejects an unsigned/invalid callback *before* any canonical
  state transition is attempted; duplicate valid callbacks are a no-op (see
  [specs/state-machines/payment-lifecycle.md](../../specs/state-machines/payment-lifecycle.md)
  "Idempotency and Retries").

See [conformance/README.md](../README.md) for the rules governing everything in this directory.
