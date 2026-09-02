# 0001. Record Architecture Decisions

- **Status:** Accepted
- **Date:** 2026-09-02
- **Deciders:** BongoPay founding maintainer(s)
- **Related:** [AGENTS.md §8](../AGENTS.md#8-when-adrs-and-rfcs-are-required), [rfcs/README.md](../rfcs/README.md)

## Context

BongoPay is contract-first (see [ARCHITECTURE.md](../ARCHITECTURE.md)) and expects contributions
from both humans and AI coding agents (see [AGENTS.md](../AGENTS.md)). Both audiences need a way
to tell "this was decided, and here's why" apart from "this is how the code happens to work
today." Without a durable record, architectural rationale lives only in PR discussions, chat
history, or a maintainer's memory — none of which survive well, and none of which an AI agent
starting a fresh session can read before proposing a change that contradicts a decision it never
saw.

## Decision

BongoPay records architecturally significant decisions as Architecture Decision Records (ADRs)
in this directory, using [0000-template.md](0000-template.md), numbered sequentially and never
renumbered or deleted once merged. A decision is ADR-level (as opposed to RFC-level, or not
requiring either) per the bar set in
[AGENTS.md §8](../AGENTS.md#8-when-adrs-and-rfcs-are-required): a new shared abstraction,
persistence model changes, new orchestration behavior, or changes affecting multiple modules.

This ADR is itself `0001` — the first record establishes the mechanism the rest of the project
then uses.

## Alternatives Considered

- **No formal record; rely on PR descriptions and issue discussion.** Rejected — PR discussions
  are not indexed or discoverable the way a numbered, directory-listed ADR is, and don't survive
  being superseded cleanly (a PR thread can't easily say "this is now superseded by #234" in a
  way a future reader will find).
- **A single running `DECISIONS.md` file instead of one file per decision.** Rejected — a single
  growing file is harder to review incrementally (every decision's PR diff touches the same
  file) and harder to mark individual decisions as superseded without disturbing unrelated
  entries.
- **Only use RFCs, skip the lighter ADR tier entirely.** Rejected — per
  [AGENTS.md §8](../AGENTS.md#8-when-adrs-and-rfcs-are-required), routing every reversible,
  narrow decision through full RFC discussion and Core Maintainer consensus
  ([GOVERNANCE.md](../GOVERNANCE.md#rfc-approval)) would slow down decisions that don't need
  that weight, and would create pressure to skip documentation entirely rather than go through
  a heavyweight process for a small decision.

## Consequences

- Every ADR- or RFC-level decision now has one canonical, linkable location.
- Contributors (human or AI) are expected to check `adr/` and `rfcs/` for existing decisions
  before proposing architecture that might contradict one — see
  [AGENTS.md §1](../AGENTS.md#1-read-before-you-write).
- This adds process overhead for genuinely small decisions; the bar in
  [AGENTS.md §8](../AGENTS.md#8-when-adrs-and-rfcs-are-required) exists specifically to keep
  that overhead proportional, and "not sure" defaults to the more conservative tier rather than
  skipping the record.

## Scope Check

- [x] Confirmed this is an ADR-level decision, not an RFC-level one — it establishes a
      contribution-process mechanism, not a contract, API, or security-model change.
- [x] Confirmed this does not change the meaning of any existing canonical field, event, or
      error.
