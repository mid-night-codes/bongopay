# Continuous Integration

CI exists to enforce the same checks a contributor is expected to run locally (see
[docs/development/README.md](README.md)) — it is not a separate, stricter standard, and nothing
should pass locally but fail in CI or vice versa.

## What Runs

Defined in [.github/workflows/](../../.github/workflows/), the fast-validation workflow runs on
every push and pull request:

```text
make setup     # confirm tooling assumptions hold in a clean environment
make validate  # specs/ structure, JSON Schema syntax, contract syntax
make lint      # markdown/YAML/JSON linting
make test      # unit/integration tests as they come online
make docs      # internal markdown link resolution
```

Per [AGENTS.md §5](../../AGENTS.md#5-commands-you-should-avoid), none of these may be bypassed
(`--no-verify`, skipping `make validate`) to force a PR through.

Separately, [.github/workflows/commitlint.yml](../../.github/workflows/commitlint.yml) validates
every commit in a pull request against the
[Conventional Commits](../../CONTRIBUTING.md#commit-messages-conventional-commits) format via
[`commitlint`](https://commitlint.js.org/), configured in
[.commitlintrc.json](../../.commitlintrc.json) (which just extends
`@commitlint/config-conventional` — the same `feat`/`fix`/`docs`/`test`/`refactor`/`chore`/
`perf`/`build`/`ci` types [CONTRIBUTING.md](../../CONTRIBUTING.md) already documents, so there is
nothing project-specific to keep in sync between the two). Unlike the `validate` job below, this
one is **not** gated by maintainer approval — see why in the next section.

## Maintainer Approval for Non-Contributor Pull Requests

A pull request's `authorize` job checks the author's GitHub
[`author_association`](https://docs.github.com/en/graphql/reference/enums#commentauthorassociation)
with this repository:

- **`OWNER`, `MEMBER`, `COLLABORATOR`, or `CONTRIBUTOR`** (i.e. anyone with repo access, or
  anyone who has landed a merged PR here before) — CI runs immediately, same as a push to `main`.
- **Anyone else** (`FIRST_TIME_CONTRIBUTOR`, `FIRST_TIMER`, `NONE`, `MANNEQUIN`) — the
  `await-maintainer-approval` job pauses in "Waiting" state, gated by the
  `external-contribution-review` GitHub Environment's required-reviewer protection rule. A
  maintainer must click **Review deployments → Approve** in the Actions tab before the
  `validate` job (and therefore any of the `make` targets above) runs at all.

This exists because a PR's workflow run executes arbitrary code from that PR's branch — running
it unattended for a first-time, unvetted contributor is a real supply-chain exposure, not just a
CI nicety. It does **not** replace [GOVERNANCE.md](../../GOVERNANCE.md)'s review-and-merge
process; it only gates whether CI *executes* before a human has looked at the diff.

The `external-contribution-review` environment's reviewer list is managed in the repository's
Settings → Environments, independent of this workflow file — updating who may approve does not
require a workflow change.

`commitlint.yml` is deliberately excluded from this gate: it only parses commit message text
against a static JSON config (`.commitlintrc.json`, not executable), so there is no code-execution
surface for a malicious PR to exploit the way there is with `scripts/*.sh` running under
`validate`. Every contributor — including a first-time one — gets that feedback immediately,
before a maintainer even looks at the PR.

## Dependency and Supply-Chain Scanning

Automated dependency updates are configured in
[.github/dependabot.yml](../../.github/dependabot.yml), per
[docs/development/dependency-policy.md](dependency-policy.md) and
[SECURITY.md](../../SECURITY.md#dependency-and-supply-chain-hygiene). As of Phase 0 this covers
the `github-actions` ecosystem (the workflow files themselves); an ecosystem entry is added for
each language ecosystem as `implementations/`, `adapters/`, or `sdks/` start using one — never
added speculatively ahead of an actual dependency existing.

## What CI Does Not Do Yet

- **Conformance execution** (`make test-conformance`) — the suite in
  [conformance/](../../conformance/README.md) is specification-only until Phase 1; CI running an
  empty suite would be a false signal, so it is not wired up as a required check yet.
- **Contract generation checks** (`make check-contracts`) — no generation pipeline exists yet
  (see [contracts/README.md](../../contracts/README.md)), so there is nothing to check
  consistency against.
- **Release automation** — releases are manual per
  [GOVERNANCE.md](../../GOVERNANCE.md#release-authority) and
  [VERSIONING.md](../../VERSIONING.md#release-process-current-state) during Phase 0.

Wiring these up as they become meaningful is tracked via [ROADMAP.md](../../ROADMAP.md), not
via ad hoc CI changes made "while in the area."

## Adding a CI Check

A new required CI check is a change to the contributor-facing contract of this repository (a PR
that passed yesterday might fail tomorrow with no code change). Per
[AGENTS.md §8](../../AGENTS.md#8-when-adrs-and-rfcs-are-required), treat adding one as at least
warranting an ADR — it affects every contributor and PR, not just one module.
