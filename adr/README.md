# adr/ — Architecture Decision Records

This directory records architectural decisions that are narrower or more reversible than an
[RFC](../rfcs/README.md)-level change, but still significant enough to need a documented
rationale rather than living only in a PR discussion or a maintainer's memory.

## Status

**Active.** Two ADRs recorded so far —
[0001](0001-record-architecture-decisions.md) (this process itself) and
[0002](0002-reference-implementation-language-go.md) (reference implementation language). The
remaining `TODO(ADR): ...` markers across [ARCHITECTURE.md](../ARCHITECTURE.md),
[specs/](../specs/README.md), and this repository's other docs are the backlog of decisions
this directory expects to keep accumulating.

## When an ADR Is Required

Per [AGENTS.md §8](../AGENTS.md#8-when-adrs-and-rfcs-are-required): a new shared abstraction,
persistence model changes, new orchestration behavior, or a change affecting multiple modules.
If you're unsure whether something needs an ADR or the more conservative RFC, treat it as
requiring the RFC and say so explicitly rather than deciding silently.

## Process

1. Copy [0000-template.md](0000-template.md) to `NNNN-short-title.md`, using the next sequential
   number (check existing files here to find it — the next one after
   [0002](0002-reference-implementation-language-go.md) is `0003-*`).
2. Fill in context, decision, and consequences. Leave open sub-decisions as
   `TODO(ADR): ...` / `TODO(RFC): ...` rather than guessing.
3. Open a PR. Per [GOVERNANCE.md](../GOVERNANCE.md), an ADR may be approved by a Maintainer for
   their owned area; ADRs touching `specs/` or `contracts/` need Core Maintainer review per
   CODEOWNERS.
4. Once merged, an ADR's status is `Accepted`. A later ADR may supersede an earlier one — mark
   the superseded ADR's status as `Superseded by NNNN` rather than deleting it; ADRs are a
   historical record, not living documents to edit in place.

## Naming and Numbering

```text
0001-short-kebab-case-title.md
0002-another-decision.md
```

Numbers are sequential and never reused, even if an ADR is later superseded or rejected.

## What Does Not Belong Here

- Breaking contract changes, new provider plugin models, new event architecture, security model
  changes, or versioning policy changes — those require an [RFC](../rfcs/README.md) instead.
- Implementation-level decisions with no architectural weight (e.g. internal variable naming) —
  those don't need either an ADR or an RFC.
