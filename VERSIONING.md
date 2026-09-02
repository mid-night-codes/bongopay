# Versioning Strategy

BongoPay versions several kinds of artifacts **independently**, because they change at
different rates and a change to one does not necessarily imply a change to another. All use
[Semantic Versioning (SemVer)](https://semver.org/) as the underlying scheme unless noted.

| Artifact | Versioned as | Where tracked |
|---|---|---|
| Project releases (this repo, as a whole) | `MAJOR.MINOR.PATCH` git tags | [CHANGELOG.md](CHANGELOG.md) |
| REST API contract | `MAJOR.MINOR.PATCH`, surfaced in the OpenAPI `info.version` | [contracts/openapi/](contracts/openapi/) |
| Event schemas | `MAJOR.MINOR.PATCH` per event type, surfaced in the event envelope `version` field | [specs/events/](specs/events/README.md) |
| Provider adapter contract | `MAJOR.MINOR.PATCH` | [specs/providers/adapter-contract.md](specs/providers/adapter-contract.md) |
| Scenario specification | `MAJOR.MINOR.PATCH` | [specs/scenarios/README.md](specs/scenarios/README.md) |
| SDKs (per language) | `MAJOR.MINOR.PATCH`, independent per SDK/package ecosystem | [sdks/README.md](sdks/README.md) |

## What Counts as Patch / Minor / Breaking

These definitions apply uniformly across the artifacts above unless a specific spec overrides
them for a documented reason.

### Patch (`x.y.PATCH`)

- Editorial/documentation fixes with no behavior change.
- Bug fixes that bring behavior in line with the existing spec (i.e., the old behavior was
  itself a bug).
- Internal refactoring with no observable contract change.

### Minor (`x.MINOR.z`)

- Additive, backward-compatible changes: new optional fields, new event types, new error
  codes, new scenario behaviors, new optional capabilities.
- New provider adapters or SDKs (these don't affect existing consumers).

### Breaking / Major (`MAJOR.y.z`)

- Removing or renaming a field, event type, error code, or capability.
- Changing the meaning or type of an existing field.
- Changing default behavior in a way existing consumers would observe.
- Any change that requires consumers to update their integration to keep working.

Breaking changes to a stable (post-1.0) contract require an [RFC](rfcs/README.md) — see
[AGENTS.md §8](AGENTS.md#8-when-adrs-and-rfcs-are-required) — and must document a migration
path.

## Pre-1.0 Expectations

Everything in this repository is currently pre-1.0 (`0.y.z`). Per SemVer convention, this means:

- Breaking changes may still occur in minor (`0.MINOR.0`) releases.
- Contracts under `specs/` and `contracts/` are not yet stable, but are **not** casually
  modified either — see [specs/compatibility/README.md](specs/compatibility/README.md) for the
  compatibility discipline that applies even pre-1.0.
- Reaching `1.0.0` for a given artifact is itself an architectural decision and should be
  proposed via RFC.

## Backward Compatibility Expectations

- Contracts must not be casually modified. Any change to `specs/` or `contracts/` should state
  its compatibility impact explicitly (the PR template requires this).
- Additive changes are preferred over breaking ones even pre-1.0, to build the habit and
  tooling needed for post-1.0 stability.
- See [specs/compatibility/README.md](specs/compatibility/README.md) for consumer-facing
  compatibility guidance and [ADR 0001](adr/0001-record-architecture-decisions.md) for how
  such decisions get recorded.

## Release Process (Current State)

Pre-1.0, releases are lightweight: a tagged commit plus a [CHANGELOG.md](CHANGELOG.md) entry,
authorized by a Core Maintainer per [GOVERNANCE.md](GOVERNANCE.md). A more formal release
process (automated changelog generation, per-artifact release trains) is a Phase 2+ concern —
see [ROADMAP.md](ROADMAP.md).
