# Provider Adapter Contract

Status: **Draft**. Version: `0.1.0` (see
[specs/compatibility/README.md](../compatibility/README.md)).

This document defines the language-neutral capability contract a provider adapter — or the
simulator, which implements the same contract as a peer (see
[ARCHITECTURE.md §6](../../ARCHITECTURE.md#6-simulator-boundary)) — must satisfy. It is the
contract [conformance/provider/](../../conformance/provider/README.md) tests against.

## Design Notes

- Every capability below is **optional at the adapter level** — a provider that can't support
  refunds must say so via capability discovery, not silently no-op or error unpredictably.
- Every capability's inputs and outputs are canonical types from
  [specs/payments/payment-contract.md](../payments/payment-contract.md); an adapter never
  requires a caller to pass provider-specific shapes to invoke a capability.

## Capabilities

```text
initiate(PaymentRequest) -> PaymentResult
  Submits a new payment to the provider. Must be idempotent with respect to
  PaymentRequest.idempotencyKey (see "Idempotency" in payment-contract.md).

queryStatus(providerTransactionId | Payment.id) -> PaymentResult
  Actively polls the provider for a payment's current status, for providers/situations where
  callbacks alone aren't sufficient or timely.

refund(Payment.id, amount?: Money) -> PaymentResult
  Requests a refund of a SUCCESS payment. Partial-refund support is provider-specific and must
  be declared via capability discovery (see below).

reverse(Payment.id) -> PaymentResult
  Requests a reversal of a SUCCESS payment. Distinct from refund() where a provider
  distinguishes the two (see specs/state-machines/payment-lifecycle.md).

parseCallback(rawRequest) -> ProviderCallbackEvent
  Parses a provider's raw webhook/callback payload into a provider-scoped intermediate event,
  prior to translation into a canonical event.

verifyCallback(rawRequest) -> boolean
  Verifies a callback's authenticity (signature, shared secret, source IP allow-list, etc.)
  before parseCallback's output is trusted. MUST be called, and MUST reject on failure, before
  any canonical state transition is applied — see ARCHITECTURE.md §12.
```

TODO(ADR): exact method signatures above are illustrative pseudocode, not a finalized interface
— finalizing them is blocked on `TODO(ADR): Reference implementation language` in
[ARCHITECTURE.md §14](../../ARCHITECTURE.md#14-what-this-document-does-not-decide).

## Capability Discovery

An adapter declares which of the capabilities above it supports, rather than orchestration
assuming universal support:

```text
ProviderCapabilities
  initiate:       boolean   # effectively always true — a provider that can't initiate isn't a provider
  queryStatus:    boolean
  refund:         boolean
  partialRefund:  boolean   # only meaningful if refund is true
  reverse:        boolean
  callbacks:      boolean   # false means the adapter is poll-only and relies on queryStatus
```

Orchestration must check `ProviderCapabilities` before invoking a capability and surface a
canonical error (see [specs/errors/README.md](../errors/README.md)) rather than calling an
unsupported capability and letting the adapter fail unpredictably.

## Translation Requirements

- **Status translation.** An adapter MUST map every provider status it can observe onto a valid
  `PaymentStatus` transition per
  [specs/state-machines/payment-lifecycle.md](../state-machines/payment-lifecycle.md). A
  provider status with no sensible canonical mapping is a bug in the adapter, not a reason to
  invent a new canonical status.
- **Error translation.** An adapter MUST map provider-specific error codes onto the canonical
  taxonomy in [specs/errors/README.md](../errors/README.md). Raw provider error codes may be
  preserved as diagnostic detail but must never stand in for the canonical categorization.
- **No leakage.** Provider-specific field names, status strings, or error codes must never be
  exposed through the canonical `Payment`/`PaymentResult`/event/error shapes — only through
  `providerReference.providerRawStatus` (observability-only, see
  [specs/payments/payment-contract.md](../payments/payment-contract.md)) or a namespaced
  `providerOptions` extension (see [specs/providers/extensions.md](extensions.md)).

## Conformance

An adapter is conforming only once it passes every applicable case in
[conformance/provider/](../../conformance/provider/README.md) for the capabilities it declares
— declaring a capability it doesn't actually satisfy is worse than not declaring it at all.

## Open TODOs

- TODO(ADR): finalize method signatures once a reference implementation language is chosen.
- TODO(spec): whether `queryStatus` is required for every adapter or may be omitted when
  `callbacks: true` is sufficient — currently unresolved.
