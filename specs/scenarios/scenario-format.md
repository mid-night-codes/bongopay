# Simulator Scenario Format

Status: **Draft**. Version: `0.1.0` (see
[specs/compatibility/README.md](../compatibility/README.md) for what "draft" means for
compatibility purposes).

This document defines the declarative scenario format referenced from
[specs/scenarios/README.md](README.md): how a `PaymentRequest` targeting the simulator selects a
simulated outcome, and how each outcome maps onto the canonical state machine in
[specs/state-machines/payment-lifecycle.md](../state-machines/payment-lifecycle.md).

## Design Notes

- Field lists below are **illustrative and non-exhaustive** — this is a v0.1 skeleton, not a
  finished format. Do not treat omission of a field as a decision; treat it as an open TODO
  unless marked otherwise.
- A scenario picks a path through the **existing** canonical state machine — it never adds a
  state or transition. If a scenario can't be expressed as a walk through
  [specs/state-machines/payment-lifecycle.md](../state-machines/payment-lifecycle.md)'s existing
  graph, that's a signal the state machine is missing something, not that the scenario format
  should invent a workaround.
- A scenario is selected **per-request**, not globally configured for the whole simulator — this
  keeps `SIMULATOR` usable like any other `Provider` in a `PaymentRequest`, per
  [ARCHITECTURE.md §6](../../ARCHITECTURE.md#6-simulator-boundary).

## Scenario

```text
Scenario
  name:    string     # identifier, e.g. "instant-success", "timeout-then-expire"
  outcome: Outcome
  delay?:  Duration   # optional simulated processing delay before the outcome applies
```

TODO(spec): whether `name` is free-form or drawn from a registered catalog — open until Phase 1
conformance work ([conformance/scenarios/](../../conformance/scenarios/README.md)) needs to
enumerate a fixed set.

## Outcome

```text
Outcome
  type: enum
    SUCCESS
    FAILURE
    TIMEOUT
    DUPLICATE_CALLBACK
    OUT_OF_ORDER
    INVALID_SIGNATURE
```

These are the six modes already named in
[specs/scenarios/README.md](README.md) and
[ARCHITECTURE.md §6](../../ARCHITECTURE.md#6-simulator-boundary). Each maps onto the canonical
state machine as follows:

| Outcome | Canonical effect |
|---|---|
| `SUCCESS` | `CREATED → PENDING → SUCCESS` |
| `FAILURE` | `CREATED → PENDING → FAILED` |
| `TIMEOUT` | `CREATED → PENDING → EXPIRED` (no further callback ever arrives within the simulated window) |
| `DUPLICATE_CALLBACK` | `CREATED → PENDING → SUCCESS`, then a second, identical `SUCCESS` callback — must be a no-op per [payment-lifecycle.md "Idempotency and Retries"](../state-machines/payment-lifecycle.md#idempotency-and-retries) |
| `OUT_OF_ORDER` | `CREATED → PENDING → SUCCESS`, but the callback claiming `PENDING` is delivered *after* the one claiming `SUCCESS` — exercises the same idempotency handling from the other direction |
| `INVALID_SIGNATURE` | A callback is delivered with an invalid/missing signature and MUST be rejected by `verifyCallback` ([adapter-contract.md](../providers/adapter-contract.md)) before any canonical state transition is attempted |

TODO(spec): `DUPLICATE_CALLBACK` and `OUT_OF_ORDER` as written assume a `SUCCESS`-terminal path
— whether they should also be selectable against a `FAILURE`-terminal path is open.

## Duration

```text
Duration
  # An ISO 8601 duration string, e.g. "PT2S" for two seconds. This is a *simulated* delay —
  # an implementation may collapse it to zero (e.g. in a fast conformance run) as long as the
  # eventual outcome and canonical event ordering are unaffected.
```

## Selecting a Scenario

A `PaymentRequest` targeting `Provider{id: "SIMULATOR"}` selects a scenario through the
namespaced extension mechanism in
[specs/providers/extensions.md](../providers/extensions.md):

```json
{
  "provider": { "id": "SIMULATOR" },
  "providerOptions": {
    "simulator": { "scenario": "instant-success" }
  }
}
```

`providerOptions.simulator` is opaque to core and orchestration, per
[specs/providers/extensions.md](../providers/extensions.md) — only the simulator's own
implementation reads it. `scenario` here is a scenario `name` as defined above; resolving that
name to a concrete `Scenario` (and where such definitions live — inline, a fixture file, a
registry) is left to the implementation for now.

## Rules Specific to This Document

- A scenario is selected per-request, never globally, per Design Notes above.
- Absence of `providerOptions.simulator.scenario` MUST default to a plain `SUCCESS` outcome with
  no delay — TODO(spec): confirm this default rather than requiring every example/test to
  specify one explicitly.
- A scenario must never require a real network call, per
  [ARCHITECTURE.md §6](../../ARCHITECTURE.md#6-simulator-boundary) and
  [ROADMAP.md](../../ROADMAP.md) §"What NOT to Build Yet".
- New outcomes are added additively to this document without changing the canonical payment
  contract, per [ARCHITECTURE.md §11](../../ARCHITECTURE.md#11-extension-points).

## Open TODOs

- TODO(spec): scenario catalog (registered names) vs. fully free-form `name`.
- TODO(spec): `DUPLICATE_CALLBACK`/`OUT_OF_ORDER` against a `FAILURE`-terminal path.
- TODO(spec): the no-scenario-specified default (assumed `SUCCESS`, no delay, above).
- TODO(RFC): whether scenarios need independent versioning beyond this document's own version
  once real conformance cases depend on specific scenario names (see
  [specs/compatibility/README.md](../compatibility/README.md)).
