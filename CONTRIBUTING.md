# Contributing to BongoPay

Thank you for considering a contribution to BongoPay. This document covers the human
contribution workflow. If you are an AI coding agent, read [AGENTS.md](AGENTS.md) first — it is
mandatory, not optional — then return here for the workflow details below.

## Before You Start

- Read [README.md](README.md) and [ARCHITECTURE.md](ARCHITECTURE.md) to understand what
  BongoPay is and isn't.
- Check [ROADMAP.md](ROADMAP.md) to understand what phase the project is in.
- Skim open issues and [rfcs/](rfcs/) to avoid duplicating in-flight work.

## Contribution Flow

```text
Find/Create Issue
      ↓
Discuss Scope
      ↓
Fork / Branch
      ↓
Implement Small Change
      ↓
Run Validation
      ↓
Open Pull Request
      ↓
Automated Checks
      ↓
Maintainer Review
      ↓
Merge
```

1. **Find or create an issue.** Use the templates in
   [.github/ISSUE_TEMPLATE/](.github/ISSUE_TEMPLATE/). For anything touching a canonical
   contract (specs/ or contracts/), expect a discussion before implementation starts.
2. **Discuss scope** in the issue, especially if the change might need an ADR or RFC (see
   [AGENTS.md §8](AGENTS.md#8-when-adrs-and-rfcs-are-required), which applies to humans too).
3. **Fork the repository** (or branch directly, for maintainers) and create a branch — always
   cut from the current tip of `main`, never from a stale copy or on top of another feature
   branch. AI agents: see
   [.claude/skills/feature-branch/SKILL.md](.claude/skills/feature-branch/SKILL.md) for the
   exact commands.
4. **Implement the smallest reasonable change.** See "Keep Pull Requests Small" below.
5. **Run validation locally**: `make validate`, `make lint`, `make test` (and
   `make test-conformance` if you touched a spec, contract, or adapter).
6. **Open a pull request** using the template — it will be pre-filled from
   [.github/pull_request_template.md](.github/pull_request_template.md).
7. **Automated checks** run in CI (see [docs/development/ci.md](docs/development/ci.md)).
8. **Maintainer review** — see [GOVERNANCE.md](GOVERNANCE.md) for review expectations.
9. **Merge**, typically squash-merged with a Conventional Commit message.

## Branch Naming

```text
feat/payment-events
fix/webhook-signature
docs/provider-contract
refactor/scenario-parser
chore/ci
```

Use `<type>/<short-kebab-case-description>`. Types mirror
[Conventional Commits](https://www.conventionalcommits.org/) types (see below).

## Commit Messages: Conventional Commits

```text
feat: add payment expiration event
fix: prevent duplicate callback delivery
docs: document provider capabilities
test: add idempotency conformance case
refactor: simplify scenario matching
chore: bump CI action versions
```

Common types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`, `perf`, `build`, `ci`. Add a
`!` (e.g. `feat!:`) or a `BREAKING CHANGE:` footer for anything breaking a published contract —
this should be rare and should already be backed by an RFC.

This is enforced in CI, not just documented here — see
[docs/development/ci.md](docs/development/ci.md) for the `commitlint` check that runs on every
pull request.

Each commit should be one fine-grained, task-specific, reviewable change — not an entire feature
squashed into a single commit. Optionally scaffold messages in this format with a local commit
template: `git config commit.template .gitmessage` (repo-local only; see
[.gitmessage](.gitmessage)).

## Keep Pull Requests Small

Prefer several small PRs over one large one. Explicitly avoid:

- Huge refactors mixed with feature work.
- Unrelated formatting changes bundled into a functional change.
- Silent contract changes (a schema or spec change without discussion, an ADR, or an RFC as
  required).
- Unnecessary new dependencies — see [docs/development/dependency-policy.md](docs/development/dependency-policy.md).
- Premature abstraction — three similar lines beat a speculative helper that only has one caller.

## Adding a Provider

See [docs/providers/README.md](docs/providers/README.md) and
[specs/providers/README.md](specs/providers/README.md). In short: a provider adapter is only
accepted once it implements the [adapter contract](specs/providers/adapter-contract.md),
declares its capabilities honestly, and **passes the conformance suite** in
[conformance/provider/](conformance/provider/README.md) — compiling is not sufficient.

## Adding an SDK

See [sdks/README.md](sdks/README.md). SDKs should derive from `contracts/` wherever
practical rather than being written from scratch against ad hoc assumptions.

## Proposing Architectural Changes

- Smaller, reversible architectural decisions → an [ADR](adr/README.md).
- Larger or breaking changes (new provider model, breaking contract change, new orchestration
  capability, persistence/eventing architecture, security model, plugin system) → an
  [RFC](rfcs/README.md).

## Releases and Compatibility

See [VERSIONING.md](VERSIONING.md) for how releases, API versions, event schema versions,
provider contract versions, scenario spec versions, and SDK versions are tracked
independently, and what backward-compatibility guarantees apply pre- and post-1.0.

## Definition of Done

A contribution isn't complete until, where relevant:

```text
Implementation complete
Tests passing
Contract validated
Conformance tests passing
Documentation updated
No unintended breaking change
Security reviewed where relevant
No secrets
No unnecessary dependencies
PR is reviewable (small, focused, well-described)
```

## Code of Conduct

Participation in BongoPay is governed by the [Code of Conduct](CODE_OF_CONDUCT.md).
