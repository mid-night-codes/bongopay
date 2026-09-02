# specs/state-machines/

Defines BongoPay's canonical, provider-neutral lifecycle state machines. Today this is just the
payment lifecycle; future state machines (e.g., a distinct refund or reversal sub-lifecycle, if
they turn out not to fit as `Payment` sub-states) belong here too.

## Contents

- [payment-lifecycle.md](payment-lifecycle.md) — canonical payment states, valid/invalid
  transitions, terminal states, retry/callback/idempotency implications.

## Rules Specific to This Directory

- States and transitions here are canonical — provider-specific statuses must be mapped onto
  these by adapters (see [specs/providers/README.md](../providers/README.md)) and must never
  appear as a `PaymentStatus` value.
- Any change to the state graph (new state, new transition, removing a transition) is a
  breaking or near-breaking change by default — see
  [specs/compatibility/README.md](../compatibility/README.md) and requires at least an ADR,
  likely an RFC if any implementation already exists.
