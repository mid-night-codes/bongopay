# Non-Goals

This document is the rationale behind [ROADMAP.md](../../ROADMAP.md) §"What NOT to Build Yet".
It exists so "why doesn't BongoPay just do X" has a durable answer instead of being re-litigated
in every issue and PR that proposes X.

## Why Non-Goals Matter Here Specifically

BongoPay's core bet (see [ARCHITECTURE.md](../../ARCHITECTURE.md)) is that a solid,
contract-first foundation only stays solid if scope is deliberately narrow while that foundation
is being laid. Each item below would pull design attention toward a concrete deployment/
integration reality before the canonical contract it should conform to even exists — inverting
the "specs first" principle this project is built on.

## Current Non-Goals and Why

- **A real payment switch / real money movement.** Until the canonical payment lifecycle
  ([specs/state-machines/](../../specs/state-machines/README.md)) and error model
  ([specs/errors/](../../specs/errors/README.md)) are stable, moving real money against them
  would mean either freezing the contract prematurely or risking real financial harm from a
  contract that's still supposed to change. See
  [SECURITY.md](../../SECURITY.md#threat-categories-bongopay-must-eventually-consider) for the
  threat model this eventually needs to satisfy first.
- **Settlement infrastructure / PCI card processing.** Both carry compliance and security
  obligations that only make sense to take on once there's something real to protect. Building
  them speculatively risks a compliance surface with no product behind it.
- **Large web dashboards.** A dashboard is a consumer of the canonical contract, not part of it.
  Building one now would create implicit pressure to shape the contract around a specific UI
  instead of around the domain (see [ARCHITECTURE.md §3](../../ARCHITECTURE.md#3-canonical-domain)).
- **Kubernetes operators / a service mesh / complex microservice architecture.** These are
  deployment-topology decisions. [ARCHITECTURE.md §13](../../ARCHITECTURE.md#13-logical-architecture-diagram)
  is explicit that BongoPay's architecture is logical, not a required topology — committing to
  infrastructure this heavy now would constrain implementations that should be free to colocate
  responsibilities early on.
- **Multiple real provider integrations.** One real integration built too early tends to leak
  that provider's vocabulary into the canonical model by accident. The simulator
  ([ARCHITECTURE.md §6](../../ARCHITECTURE.md#6-simulator-boundary)) exists specifically so the
  contract can be exercised realistically — including failure modes — before that risk is taken
  on.
- **Production persistence architecture / distributed transaction infrastructure.** Both are
  reference-implementation concerns, and the reference implementation's language and persistence
  model are still open (`TODO(ADR)` in
  [ARCHITECTURE.md §14](../../ARCHITECTURE.md#14-what-this-document-does-not-decide)). Building
  production-grade persistence against an undecided language would be wasted, possibly
  misleading, work.
- **An elaborate plugin runtime for provider adapters.** The long-term plugin/runtime model is
  explicitly deferred to RFC (`TODO(RFC)` in
  [ARCHITECTURE.md §14](../../ARCHITECTURE.md#14-what-this-document-does-not-decide)) once there
  are enough real adapters to generalize from. Designing a plugin runtime for zero existing
  adapters risks over-engineering for guessed requirements.

## How This List Changes

An item leaves this list only via the phase it belongs to in [ROADMAP.md](../../ROADMAP.md)
being reached, and — for anything architectural — the corresponding ADR or RFC being accepted
per [AGENTS.md §8](../../AGENTS.md#8-when-adrs-and-rfcs-are-required). Removing an item from this
document without both is itself a governance violation, not just a documentation edit.
