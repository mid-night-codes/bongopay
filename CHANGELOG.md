# Changelog

All notable changes to BongoPay are documented in this file. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and versioning follows the rules in
[VERSIONING.md](VERSIONING.md).

## [Unreleased]

### Added

- Initial repository structure: governance, contribution workflow, and AI-agent contribution
  environment (`AGENTS.md`, `.github/CONTRIBUTING_AGENT.md`).
- Placeholder specifications for payments, provider adapters, scenarios, errors, events, and
  the canonical payment state machine (`specs/`).
- Placeholder contracts (`contracts/openapi/`, `contracts/asyncapi/`, `contracts/json-schema/`).
- Conformance philosophy and directory scaffolding (`conformance/`).
- ADR and RFC processes (`adr/`, `rfcs/`), including `adr/0001-record-architecture-decisions.md`.
- CI scaffolding for fast PR validation: `.github/workflows/ci.yml` running
  `make setup && make validate && make lint && make test && make docs`, plus
  `.github/dependabot.yml` for automated dependency updates.
- `.github/ISSUE_TEMPLATE/` (bug report, feature request, provider request, SDK request,
  documentation), `.github/pull_request_template.md`, and `.github/CODEOWNERS`.
- Root `Makefile` developer interface (`make setup`, `make validate`, `make lint`, `make test`,
  `make test-conformance`, `make check-contracts`, `make docs`, `make generate`, `make clean`).
- READMEs for every directory referenced from `AGENTS.md` §3 (`contracts/`, `conformance/`,
  `implementations/`, `adapters/`, `sdks/`, `examples/`, `docs/`, `adr/`, `rfcs/`) and for every
  `specs/` subdirectory (`payments/`, `providers/`, `events/`, `errors/`, `scenarios/`,
  `state-machines/`, `compatibility/`), closing out Phase 0's directory-structure requirement.
- First substantive spec documents: `specs/state-machines/payment-lifecycle.md`,
  `specs/providers/adapter-contract.md`, `specs/providers/extensions.md`.
- `docs/architecture/non-goals.md`, `docs/development/` (README, `ci.md`,
  `dependency-policy.md`), and `docs/providers/README.md`.
- `.markdownlint-cli2.jsonc`, disabling `MD013` (line-length), `MD060` (table pipe style), and
  `MD041` (first-line-heading) to match this repository's existing long-form prose and template
  conventions rather than the linter's code-comment-oriented defaults — see the file's inline
  comments for the rationale per rule.
- READMEs for the remaining directories in the root `README.md` "Repository Structure" list that
  had none (`tests/`, `tools/`, `deploy/`, `deploy/docker/`, `deploy/compose/`) and for every
  `contracts/` subdirectory (`openapi/`, `asyncapi/`, `json-schema/`, `examples/`), so that
  claim is now accurate for every directory the README lists.
- Founding Core Maintainer entry in `MAINTAINERS.md` and `.github/CODEOWNERS`, replacing the
  `@TBD-core-maintainer` placeholder.
- Maintainer-approval gate in `.github/workflows/ci.yml`: PRs from non-contributors (anyone
  without prior repo access or a merged PR) pause on the `external-contribution-review`
  environment's required reviewer before `make validate`/`lint`/`test`/`docs` execute any code
  from the PR branch.
- `.github/workflows/commitlint.yml` and `.commitlintrc.json`, enforcing the Conventional
  Commits format documented in `CONTRIBUTING.md` on every pull request's commits — deliberately
  not behind the approval gate above, since it only parses commit text against a static JSON
  config.
- `.gitmessage` commit template and a "Branching and Commit Workflow" section in
  `.github/CONTRIBUTING_AGENT.md`, formalizing (and, for AI agents, giving exact commands for)
  the required workflow: branch from the current tip of `main` for every task, commit in
  fine-grained task-specific steps, then push the branch and open a PR — never commit straight
  to `main`. Deliberately placed in `.github/`, not a tool-specific `.claude/skills/` directory,
  to stay consistent with `AGENTS.md`'s tool-agnostic framing (Claude, GPT, Copilot Workspace,
  or any other agent).

### Fixed

- The `ci.yml` trust check now keys on fork status and merged-PR history instead of GitHub's
  `author_association`, which returned `NONE` for an actual repository member with private org
  membership visibility — this was caught dogfooding the new branch/PR workflow on
  [PR #4](https://github.com/mid-night-codes/bongopay/pull/4).
- Added the READMEs for `contracts/json-schema/{payments,errors,events}/`,
  `implementations/reference/`, and `docs/{concepts,contributing,maintainers}/` — these
  directories existed only on local disk and were never actually committed (git does not track
  empty directories), so a fresh clone was silently missing them and every link pointing at
  them, undetected by `make docs` run locally against the pre-existing local tree. Also caught
  via PR #4's CI run, not local validation — a reminder that "passes locally" and "passes on a
  fresh checkout" aren't always the same claim.

This closes Phase 0 ("Foundation") per [ROADMAP.md](ROADMAP.md): every `make` target
(`setup`, `validate`, `lint`, `test`, `test-conformance`, `check-contracts`, `docs`, `generate`,
`clean`) passes, and every directory referenced from the root `README.md` and `AGENTS.md` has a
README.

### Added (Phase 1)

- [ADR 0002](adr/0002-reference-implementation-language-go.md): the reference implementation
  (`implementations/reference/`) will be written in Go, unblocking implementation work per
  [implementations/README.md](implementations/README.md).
- A new required step in `.github/CONTRIBUTING_AGENT.md`'s branching/commit workflow: find or
  create an open issue *before* the first commit, and reference it (`Refs #N`/`Closes #N`) in a
  commit footer or the PR body. `.gitmessage`'s footer section reflects this as required, not
  an optional example — and its stale reference to the removed `.claude/skills/` path (missed
  when that directory was deleted, since `make docs`'s link checker only scans `*.md` files) is
  fixed to point at `.github/CONTRIBUTING_AGENT.md`. Opened retroactively as
  [#7](https://github.com/mid-night-codes/bongopay/issues/7), with
  [#8](https://github.com/mid-night-codes/bongopay/issues/8) opened to bring the then-already-open
  PR #6 into compliance.
- `implementations/reference/internal/payment`: the canonical domain types from
  `specs/payments/payment-contract.md` and the state machine from
  `specs/state-machines/payment-lifecycle.md`, with exhaustive tests over every `(from, to)`
  status pair.
- A `go` job in `.github/workflows/ci.yml` running `go build`/`vet`/`test` and a `gofmt` check
  against `implementations/reference/`, gated by the same maintainer-approval check as
  `validate`.

### Changed

- `scripts/test.sh` (`make test`) now actually dispatches to `go test ./...` for every
  `go.mod` found under `implementations/`, `adapters/`, or `sdks/`, instead of exiting 1 with a
  `TODO` the moment any `*_test.*` file existed anywhere in those directories — the placeholder
  it replaces was written before there was real code to hit that branch.

Nothing has been released yet. This entry will move under a version heading (e.g. `[0.1.0]`)
at the first tagged release, per [VERSIONING.md](VERSIONING.md).
