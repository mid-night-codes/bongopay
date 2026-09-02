# Local Development Guide

This is the full local development guide referenced from the root [README.md](../../README.md)
Quick Start and the [Makefile](../../Makefile). It expands on, but never contradicts,
[AGENTS.md §4](../../AGENTS.md#4-commands-you-may-run) — if the two disagree, `AGENTS.md` wins.

## Prerequisites

Per [scripts/setup.sh](../../scripts/setup.sh):

- `python3` — required. Used for JSON/YAML/OpenAPI structural checks in `scripts/validate-*.sh`
  and the doc-link checker in `scripts/docs.sh`.
- `node` / `npx` — optional. Only used for richer Markdown linting (`markdownlint-cli2`) via
  `make lint`; everything else degrades gracefully without it.
- `git` — required.

No language-specific runtime (JVM, Go, Rust, etc.) is required to work in `specs/`, `contracts/`,
or `docs/` during Phase 0 — see [ROADMAP.md](../../ROADMAP.md). Once
`implementations/reference/` picks a language (`TODO(ADR)`, see
[ARCHITECTURE.md §14](../../ARCHITECTURE.md#14-what-this-document-does-not-decide)), that
language's toolchain becomes a prerequisite for working in that directory specifically, not
repository-wide.

## Getting Started

```bash
git clone <repository-url>
cd bongopay

make setup      # checks for the prerequisites above and reports what's available
make validate   # validates specs/ structure, JSON Schema syntax, contract syntax
make lint       # markdown/YAML/JSON linting
make test       # runs unit/integration tests as they come online
make docs       # checks that internal markdown links resolve
```

All of the above are safe to run repeatedly and never mutate source files — see
[scripts/README.md](../../scripts/README.md) for what each target does under the hood.

## Everyday Workflow

1. Branch from the current tip of `main` — never commit directly to `main`. See
   [CONTRIBUTING.md](../../CONTRIBUTING.md#branch-naming) for naming, or
   [.claude/skills/feature-branch/SKILL.md](../../.claude/skills/feature-branch/SKILL.md) for
   the exact commands if you're an AI agent.
2. Read [AGENTS.md](../../AGENTS.md) §1 ("Read Before You Write") for the order to read things
   in before changing anything — this applies to human contributors too, not just AI agents.
3. Make the smallest change that satisfies the issue/request (see
   [CONTRIBUTING.md](../../CONTRIBUTING.md#keep-pull-requests-small)), as a series of
   fine-grained, task-specific commits rather than one large one.
4. Before opening a PR, run:

   ```bash
   make validate && make lint && make test && make docs
   ```

5. If your change touches `specs/` or `contracts/`, check whether it needs an
   [ADR](../../adr/README.md) or [RFC](../../rfcs/README.md) per
   [AGENTS.md §8](../../AGENTS.md#8-when-adrs-and-rfcs-are-required) *before* writing the change,
   not after.

## Editor / Tooling Notes

- Markdown is the primary "source code" of this repository during Phase 0 — treat broken
  internal links (`make docs`) as build failures, not nitpicks.
- No IDE or editor configuration is currently mandated. If one is added later (e.g.
  `.editorconfig`), it belongs here as a documented convention, not silently assumed.

## Dependency Policy

See [dependency-policy.md](dependency-policy.md) before adding any dependency, in any language.

## Continuous Integration

See [ci.md](ci.md) for what CI runs and how it maps to the `make` targets above.
