# specs/payments/ — Canonical Payment Contract

This directory defines the canonical, provider-neutral payment domain. It is the most
frequently referenced spec in BongoPay — read it before touching anything in `contracts/`,
`implementations/`, or `adapters/` that deals with payments.

## Status

**Draft / Phase 0.** Enough structure exists to establish the contract-first approach and unblock
`contracts/openapi/` and `contracts/json-schema/` placeholder work. Fields are intentionally not
finalized — see the TODOs inline in [payment-contract.md](payment-contract.md).

## Contents

- [payment-contract.md](payment-contract.md) — the canonical concepts: `Payment`,
  `PaymentRequest`, `PaymentResult`, `PaymentStatus`, `Money`, `Currency`,
  `CustomerReference`, `ProviderReference`, `PaymentMethod`, `Provider`, `Metadata`,
  `ProviderOptions`.

Related specs, kept separate because they're independently versioned and reused outside
payments specifically:

- Lifecycle/state transitions → [specs/state-machines/payment-lifecycle.md](../state-machines/payment-lifecycle.md)
- Provider capability and extension model → [specs/providers/README.md](../providers/README.md)
- Errors raised during payment operations → [specs/errors/README.md](../errors/README.md)
- Events emitted during a payment's lifecycle → [specs/events/README.md](../events/README.md)

## Rules Specific to This Directory

- A canonical field must make sense for *any* provider. If a field only makes sense for one
  provider, it belongs in that provider's `providerOptions` namespace instead (see
  [specs/providers/extensions.md](../providers/extensions.md)), not here.
- `Money` amounts are integers in the currency's minor unit (e.g., cents), never floats — this
  avoids an entire category of rounding bugs and is non-negotiable without an RFC.
- Do not add a field "just in case." Every field here should be traceable to a concrete use case
  described in this directory or linked from an issue/RFC.
