# specs/scenarios/ — Simulator Scenario Format

This directory defines the declarative, versioned format used to describe simulator behavior:
how a simulated provider should respond to a payment, including the failure modes real providers
rarely let you trigger on demand.

## Status

**Draft / Phase 0.** Directory established to unblock cross-references from
[ARCHITECTURE.md §6](../../ARCHITECTURE.md#6-simulator-boundary); the scenario format itself is
not yet written.

## Contents

- `scenario-format.md` — **TODO(ADR): not yet written.** Will define the declarative scenario
  schema: how a scenario selects outcomes (success, failure, timeout, duplicate callback,
  out-of-order event, invalid signature) and any timing/sequencing controls.

## Rules Specific to This Directory

- Scenarios describe **simulator** behavior only. They never define or alter canonical domain
  concepts in [specs/payments/](../payments/README.md) or
  [specs/state-machines/](../state-machines/README.md) — a scenario picks a path through the
  existing canonical state machine, it does not add states to it.
- A scenario must be fully expressible and executable without a network call to a real provider
  (see [ARCHITECTURE.md §6](../../ARCHITECTURE.md#6-simulator-boundary) and
  [ROADMAP.md](../../ROADMAP.md) §"What NOT to Build Yet").
- The simulator implements the same provider contract as any real adapter
  ([specs/providers/adapter-contract.md](../providers/adapter-contract.md)); scenarios configure
  that implementation, they do not bypass the contract.
- New simulated failure modes are added additively to this format without changing the canonical
  payment contract (see [ARCHITECTURE.md §11](../../ARCHITECTURE.md#11-extension-points)).
