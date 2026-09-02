# conformance/webhook/ — Webhook and Callback Conformance Cases

Conformance cases for webhook/callback handling: signature verification, duplicate delivery, and
out-of-order delivery, against
[specs/providers/adapter-contract.md](../../specs/providers/adapter-contract.md)'s
`verifyCallback`/`parseCallback` capabilities and
[specs/events/README.md](../../specs/events/README.md).

## Status

**Specification-only, Phase 1 ("Simulator Core") for the security-relevant cases, extended in
Phase 4 ("Reliability Testing") for chaos/replay coverage.** See
[ROADMAP.md](../../ROADMAP.md). No case files exist yet.

## Planned Coverage

These are treated as **not optional** per
[ARCHITECTURE.md §12](../../ARCHITECTURE.md#12-security-boundaries) — an adapter or
implementation declaring the `callbacks` capability is not conforming without passing them:

- **Signature verification** — an invalid or missing signature is rejected by
  `verifyCallback` *before* any canonical state transition is attempted; the rejection is
  surfaced as a canonical error (see [specs/errors/README.md](../../specs/errors/README.md)),
  not a silent drop.
- **Duplicate delivery** — a second, valid delivery of a callback already processed is a no-op
  with respect to canonical state and does not re-emit a duplicate event (see
  [specs/state-machines/payment-lifecycle.md](../../specs/state-machines/payment-lifecycle.md)
  "Idempotency and Retries").
- **Out-of-order delivery** — a callback implying an earlier lifecycle state arriving after a
  later one is already recorded does not regress canonical state.
- **Replay** (Phase 4) — a callback replayed well after original delivery is handled the same as
  ordinary duplicate delivery, not treated as new.

See [conformance/README.md](../README.md) for the rules governing everything in this directory.
