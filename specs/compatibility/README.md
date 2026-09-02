# specs/compatibility/ — Cross-Cutting Backward-Compatibility Rules

This directory holds the backward-compatibility rules that apply across all of `specs/` —
payments, providers, events, errors, scenarios, and state machines — rather than duplicating the
same "what counts as breaking" reasoning inside each of those directories.

## Status

**Draft / Phase 0.** Directory established to unblock cross-references from
[specs/state-machines/README.md](../state-machines/README.md) and
[specs/events/README.md](../events/README.md); the consolidated rules document is not yet
written. [VERSIONING.md](../../VERSIONING.md) at the repository root currently covers
per-artifact-type versioning (releases, REST APIs, event schemas, provider adapter contracts,
scenario specs, SDKs) and is authoritative until this directory's document exists.

## Contents

- `compatibility-rules.md` — **TODO(ADR): not yet written.** Will consolidate the shared
  backward-compatibility rules referenced from each spec area (e.g. "fields are additive only,
  never repurposed," "removing a state-machine transition is breaking by default") so they are
  defined once instead of restated per directory.

## Rules Specific to This Directory

- This directory clarifies and cross-references [VERSIONING.md](../../VERSIONING.md); it does
  not replace it. Where the two disagree, treat that as a documentation bug to flag, not a
  license to pick either.
- A rule here must apply to more than one spec area. A rule specific to a single area (e.g. "the
  payment state graph is breaking to change") belongs in that area's own README instead.
- Changes to this directory are themselves compatibility-policy changes — per
  [AGENTS.md §8](../../AGENTS.md#8-when-adrs-and-rfcs-are-required), treat them as requiring at
  least ADR consideration, likely an RFC once any implementation exists.
