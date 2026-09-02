# specs/events/ — Canonical Event Model

This directory defines the canonical event envelope and the catalog of event types emitted
during a payment's lifecycle. Events are how orchestration communicates canonical state changes
outward (to consumer applications, webhooks, and internal subscribers) without leaking provider
or transport vocabulary.

## Status

**Draft / Phase 0.** Directory established to unblock cross-references from
[specs/payments/](../payments/README.md) and
[ARCHITECTURE.md §11](../../ARCHITECTURE.md#11-extension-points); the documents below are not
yet written.

## Contents

- `event-envelope.md` — **TODO(ADR): not yet written.** Will define the canonical event
  envelope shape (identity, timestamp, event type, canonical payload, versioning fields) shared
  by every event type.
- `event-types.md` — **TODO(ADR): not yet written.** Will define the catalog of canonical event
  types emitted across the payment lifecycle in
  [specs/state-machines/payment-lifecycle.md](../state-machines/payment-lifecycle.md).

## Rules Specific to This Directory

- New event types are additive only. An existing event type must never be repurposed to mean
  something different — add a new type or version instead (see
  [specs/compatibility/README.md](../compatibility/README.md) and
  [VERSIONING.md](../../VERSIONING.md)).
- Event payloads carry canonical fields only — no provider-specific vocabulary. Provider detail
  travels in `providerOptions`, per [specs/providers/extensions.md](../providers/extensions.md).
- Events reflect transitions already defined in
  [specs/state-machines/](../state-machines/README.md); this directory does not define new
  states on its own.
- Duplicate delivery must be assumed by any consumer of these events — the canonical model does
  not guarantee exactly-once delivery, only that duplicates are safely identifiable (see
  [specs/errors/README.md](../errors/README.md) for idempotency-related failure modes).
