# conformance/scenarios/ — Simulator Conformance Cases

Conformance cases that exercise the simulator against
[specs/scenarios/](../../specs/scenarios/README.md) — verifying that a declared scenario
(success, failure, timeout, duplicate callback, out-of-order event, invalid signature) actually
drives the simulator through the corresponding canonical state transitions in
[specs/state-machines/payment-lifecycle.md](../../specs/state-machines/payment-lifecycle.md).

## Status

**Specification-only, Phase 1 ("Simulator Core").** See [ROADMAP.md](../../ROADMAP.md). No case
files exist yet — they depend on
[specs/scenarios/scenario-format.md](../../specs/scenarios/README.md) being written first.

## Planned Coverage

- Each declared scenario outcome (success, failure, timeout, duplicate callback, out-of-order
  event, invalid signature — see
  [ARCHITECTURE.md §6](../../ARCHITECTURE.md#6-simulator-boundary)) drives the simulator to the
  matching canonical `PaymentStatus`, never an undeclared one.
- The simulator satisfies
  [specs/providers/adapter-contract.md](../../specs/providers/adapter-contract.md) as a peer of
  real adapters — cases here overlap with, and may be reused by,
  [conformance/provider/](../provider/README.md) for that reason.
- A scenario never causes an invalid state-machine transition, even under a "chaos" failure mode
  — the simulator is bound by the same lifecycle rules as any real adapter.

See [conformance/README.md](../README.md) for the rules governing everything in this directory.
