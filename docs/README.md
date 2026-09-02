# docs/ — Human-Facing Documentation

This directory holds documentation aimed at people (and agents) trying to understand, run, or
extend BongoPay — as distinct from [specs/](../specs/README.md), which is the normative,
contract-defining source of truth. If something here disagrees with `specs/`, `specs/` wins;
treat the disagreement as a bug in `docs/` to fix.

## Layout

| Directory | Covers |
|---|---|
| [architecture/](architecture/) | Deeper architectural explanation and rationale beyond [ARCHITECTURE.md](../ARCHITECTURE.md), including [non-goals.md](architecture/non-goals.md) |
| [concepts/](concepts/) | Conceptual/tutorial explanations of BongoPay ideas (payment lifecycle, simulation, provider adapters) for readers new to the project |
| [development/](development/) | Local development, CI, and dependency-policy guides — see [development/README.md](development/README.md) |
| [contributing/](contributing/) | Extended contribution guidance that doesn't fit in the root [CONTRIBUTING.md](../CONTRIBUTING.md) |
| [providers/](providers/) | Guidance for building and submitting a provider adapter — see [providers/README.md](providers/README.md) |
| [maintainers/](maintainers/) | Maintainer-facing process docs (release process, triage, security response coordination) |

## Status

**Draft / Phase 0.** Directory structure exists; most subdirectories have no content yet beyond
what is directly linked from [README.md](../README.md), [AGENTS.md](../AGENTS.md), and
[CONTRIBUTING.md](../CONTRIBUTING.md) (see those files' own TODOs).

## Rules for Working in `docs/`

1. **Explain, don't redefine.** `docs/` clarifies and teaches; it never introduces a canonical
   concept, field, or behavior that isn't already in `specs/`. If you find yourself defining new
   behavior while writing a doc, stop and route it to `specs/` first.
2. **Keep it current or mark it clearly stale.** A doc describing behavior that no longer
   matches `specs/` or the current Phase (see [ROADMAP.md](../ROADMAP.md)) is worse than no doc
   — fix or flag it rather than leaving it silently wrong.
3. **Prefer linking to the source of truth over restating it.** E.g. link to
   [specs/payments/README.md](../specs/payments/README.md) rather than re-describing the
   `Payment` model here; docs rot faster than specs when duplicated.
4. **Every user-visible or contract-visible behavior change needs a doc update** somewhere in
   here or in the relevant spec — see [AGENTS.md §6](../AGENTS.md#6-testing-formatting-and-documentation-expectations).

## What Does Not Belong Here

- Normative contract definitions — that's [specs/](../specs/README.md).
- Machine-readable schemas — that's [contracts/](../contracts/README.md).
