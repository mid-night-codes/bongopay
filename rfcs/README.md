# rfcs/ — Requests for Comments

This directory holds proposals for changes significant or risky enough to need broad discussion
and Core Maintainer consensus before implementation — the more weighty counterpart to an
[ADR](../adr/README.md).

## Status

**Draft / Phase 0.** No RFCs have been opened yet. The `TODO(RFC): ...` markers in
[ARCHITECTURE.md](../ARCHITECTURE.md) and [AGENTS.md](../AGENTS.md#12-generated-files) (SDK
generation pipeline, long-term provider plugin/runtime model) are the first candidates once the
project reaches the phase where they need resolving — see [ROADMAP.md](../ROADMAP.md).

## When an RFC Is Required

Per [AGENTS.md §8](../AGENTS.md#8-when-adrs-and-rfcs-are-required): breaking a public
API/contract, a new provider plugin model, new event architecture, a security model change, or a
versioning policy change. If unsure whether something is ADR- or RFC-level, default to RFC and
say so explicitly.

## Process and Statuses

```text
Draft → Discussion → Accepted / Rejected → Implemented
                                 ↘ Withdrawn
```

1. Copy [0000-template.md](0000-template.md) to `NNNN-short-title.md`, using the next sequential
   number (there are none yet, so the first RFC is `0001-*`).
2. Open a PR with status `Draft`. This starts the discussion period.
3. Move to `Discussion` once a maintainer confirms scope; substantive changes to the proposal
   happen here, in the open, not after acceptance.
4. Per [GOVERNANCE.md](../GOVERNANCE.md#rfc-approval), Core Maintainers decide `Accepted` or
   `Rejected` by consensus after the discussion period. Record the outcome and rationale in the
   RFC document itself.
5. Once the change is implemented, update the RFC's status to `Implemented` and link the
   implementing PR(s). An accepted-but-unimplemented RFC is not itself a license to build ahead
   of the phase it belongs to — see [ROADMAP.md](../ROADMAP.md).
6. A proposer may mark their own RFC `Withdrawn` at any point before acceptance.

RFCs are a historical record once merged — a later RFC may supersede an earlier one, but the
earlier document is marked superseded, not deleted or rewritten.

## Naming and Numbering

```text
0001-short-kebab-case-title.md
0002-another-proposal.md
```

Numbers are sequential and never reused, even for a rejected or withdrawn RFC.

## What Does Not Belong Here

- Decisions narrow or reversible enough for an [ADR](../adr/README.md) instead — routing
  everything through RFC slows down decisions that don't need Core Maintainer consensus.
- Implementation-level design detail that doesn't change a contract, API, or policy.
