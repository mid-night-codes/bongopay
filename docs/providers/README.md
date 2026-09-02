# Building a Provider Adapter

This is the practical, human-facing companion to
[specs/providers/README.md](../../specs/providers/README.md) and
[specs/providers/adapter-contract.md](../../specs/providers/adapter-contract.md) — read both of
those first; this document doesn't restate their rules, only walks through applying them.

## Status

**Guidance-only, Phase 3 ("Provider Ecosystem").** See [ROADMAP.md](../../ROADMAP.md). Early
groundwork (reading the specs, understanding capability discovery) is useful ahead of that
phase; writing an actual adapter before
[specs/providers/adapter-contract.md](../../specs/providers/adapter-contract.md) leaves Draft is
premature.

## Steps to Add a Provider Adapter

1. **Read the contract, not just this guide.**
   [specs/providers/adapter-contract.md](../../specs/providers/adapter-contract.md) defines the
   capabilities (`initiate`, `queryStatus`, `refund`, `reverse`, `parseCallback`,
   `verifyCallback`) and the `ProviderCapabilities` discovery shape.
2. **Map the provider's states honestly.** Every state your provider can report must map onto a
   valid `PaymentStatus` in
   [specs/state-machines/payment-lifecycle.md](../../specs/state-machines/payment-lifecycle.md).
   If it doesn't map cleanly, that's a signal to raise the gap (see
   [AGENTS.md §7](../../AGENTS.md#7-contract-change-rules)), not to invent a new canonical
   status.
3. **Map the provider's errors honestly**, onto
   [specs/errors/README.md](../../specs/errors/README.md)'s taxonomy.
4. **Declare capabilities truthfully.** If the provider can't do partial refunds, declare
   `partialRefund: false` — don't approximate it by silently refunding the full amount.
5. **Keep provider-specific data in `providerOptions`.** See
   [specs/providers/extensions.md](../../specs/providers/extensions.md) for the namespacing
   rules; nothing provider-specific belongs on the canonical `Payment` shape.
6. **Verify callbacks before trusting them.** `verifyCallback` must reject before any canonical
   state transition is attempted — see
   [ARCHITECTURE.md §12](../../ARCHITECTURE.md#12-security-boundaries).
7. **Write it against the simulator's test doubles first**, then against a provider sandbox —
   never against production, and never with real credentials committed to this repository (see
   [AGENTS.md §10](../../AGENTS.md#10-security-restrictions)).
8. **Pass conformance.** Your adapter is accepted once it passes
   [conformance/provider/](../../conformance/provider/README.md) for every capability it
   declares — per [CONTRIBUTING.md](../../CONTRIBUTING.md#adding-a-provider), compiling is not
   sufficient.

## Where Your Code Lives

The adapter itself lives under [adapters/](../../adapters/README.md), one directory per
provider. It does not live in `specs/`, `contracts/`, or `conformance/` — those define and test
the contract; the adapter satisfies it.

## Common Mistakes

- Leaking a provider status string through `PaymentStatus` instead of
  `providerReference.providerRawStatus`.
- Treating `providerOptions` as a place to smuggle a field that should really be canonical (if
  two unrelated providers would both want it, raise that as a spec change instead).
- Skipping `verifyCallback` "just for local testing" and forgetting to re-enable it — see
  [AGENTS.md §10](../../AGENTS.md#10-security-restrictions): never weaken signature verification
  to make a test pass.
