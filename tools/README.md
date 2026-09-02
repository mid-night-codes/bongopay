# tools/ — Developer Tools

This directory holds standalone developer tools that are more substantial than a
[scripts/](../scripts/README.md) validation script but aren't themselves part of the canonical
contract, an implementation, or an SDK — for example, a future scaffolding generator for new
provider adapters, or a local scenario-authoring helper.

## Status

**Empty — Phase 0.** No tools exist yet. `scripts/` covers all current validation/developer
needs (see [scripts/README.md](../scripts/README.md)); this directory is reserved for tooling
that outgrows a single shell script.

## Rules for Working in `tools/`

1. **A tool here supports contributors; it never becomes a hidden second source of truth.** If a
   tool starts encoding decisions that belong in `specs/` (e.g. "valid provider IDs are..."),
   that logic should read from `specs/`/`contracts/`, not duplicate it.
2. **Prefer extending `scripts/` first.** Only graduate something to `tools/` once it's too
   substantial for a `scripts/*.sh` entry — see
   [scripts/README.md](../scripts/README.md) for what belongs there instead.
3. **Same dependency discipline as everywhere else** — see
   [docs/development/dependency-policy.md](../docs/development/dependency-policy.md) before
   adding a tool that requires a new runtime or package.

## What Does Not Belong Here

- Validation invoked by `make validate`/`make lint`/`make docs` — that's
  [scripts/](../scripts/README.md).
- Anything that is itself part of the canonical contract, an implementation, an adapter, or an
  SDK — those have their own directories.
