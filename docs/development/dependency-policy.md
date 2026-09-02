# Dependency Policy

Referenced from [AGENTS.md §9](../../AGENTS.md#9-dependency-rules),
[CONTRIBUTING.md](../../CONTRIBUTING.md#keep-pull-requests-small), and
[SECURITY.md](../../SECURITY.md#dependency-and-supply-chain-hygiene). This document is the
detail behind those pointers.

## Default Answer: Don't

BongoPay's core value is language- and framework-neutrality (see
[ARCHITECTURE.md](../../ARCHITECTURE.md)). A dependency added casually in one place tends to
either (a) quietly become load-bearing everywhere, or (b) constrain a choice
(`specs/`/`contracts/` language-neutrality, reference-implementation language) that hasn't been
made yet. The default answer to "should I add this dependency" is **don't**, until the checklist
below says otherwise.

## Before Adding Any Dependency

Answer all of these, in the PR description, before adding a dependency in any language:

1. **Can the standard library solve this instead?** If yes, use it — even if a library would be
   marginally more convenient.
2. **Is it actively maintained, with a compatible license?** Check last-release recency and
   license compatibility with [Apache License 2.0](../../LICENSE) (copyleft licenses like GPL
   are generally incompatible with BongoPay's distribution model).
3. **Does it meaningfully improve maintainability, versus just convenience?** "It saves me 10
   lines" is not sufficient justification on its own.
4. **Could it constrain the language-neutrality of `specs/`/`contracts/`?** If the dependency
   would touch either directory, the answer should almost always be "don't" — those directories
   are prose/schema, not code, and should need nothing beyond `python3` (see
   [scripts/README.md](../../scripts/README.md)).

## Where This Applies Differently

- **`specs/` and `contracts/`** — no runtime dependencies, ever; these are language-neutral by
  design (see [specs/README.md](../../specs/README.md) rule 2).
- **`scripts/`** — POSIX shell plus `python3`/`node` only, per
  [scripts/README.md](../../scripts/README.md); no dependency installs required for basic
  validation.
- **`implementations/`, `adapters/`, `sdks/`** — normal per-language dependency hygiene applies
  once a language is chosen (`TODO(ADR)`, see
  [ARCHITECTURE.md §14](../../ARCHITECTURE.md#14-what-this-document-does-not-decide)), subject
  to the four questions above.
- **`examples/`** — kept minimal; an example pulling in a large framework to demonstrate a small
  integration point is usually a sign the example is doing too much.

## Recording the Decision

- A dependency added to satisfy the checklist above is recorded in the PR, not in a separate
  document — there is no repository-wide dependency manifest to update by hand during Phase 0.
- A dependency significant enough to affect multiple modules or constrain future architecture
  (e.g. a shared serialization library used by both an adapter and an SDK) warrants an
  [ADR](../../adr/README.md) per [AGENTS.md §8](../../AGENTS.md#8-when-adrs-and-rfcs-are-required),
  not just a PR description note.

## Automated Scanning

See [.github/dependabot.yml](../../.github/dependabot.yml) and
[docs/development/ci.md](ci.md#dependency-and-supply-chain-scanning) for how dependency updates
and known-vulnerability scanning are automated once an ecosystem exists to scan.
