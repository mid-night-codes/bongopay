# Detailed Agent Contribution Workflow

This document is the detailed workflow and self-review checklist referenced from
[AGENTS.md §13](../AGENTS.md#13-how-to-make-small-reviewable-changes). Read `AGENTS.md` first —
this document expands on step 9 of that section ("write a PR description that states your
assumptions explicitly") and does not repeat rules already stated there.

## Detailed Workflow

1. **Restate the request** in your own words in your working notes, including anything
   ambiguous. If it's still ambiguous after restating it, that's the signal described in
   [AGENTS.md §15](../AGENTS.md#15-if-you-are-still-unsure) — stop and say so rather than
   guessing.
2. **Branch from the current tip of `main`** before touching anything — see
   [Branching and Commit Workflow](#branching-and-commit-workflow) below for the exact commands
   and branch naming. Never commit directly to `main`; never build on top of a stale or
   unrelated branch.
3. **Read in the order AGENTS.md §1 specifies** — this file, the relevant directory README(s),
   any relevant ADR, the governing spec, existing tests/conformance cases. Do not skip to
   writing code because the change "looks small."
4. **Locate the smallest change.** If the smallest change that satisfies the request still
   touches multiple directories (e.g. a spec change plus its contract derivation), that's
   expected — "small" means minimal *scope*, not minimal *file count*.
5. **Check the ADR/RFC bar** (AGENTS.md §8) before writing the change, not after — discovering
   partway through an implementation that it needed an RFC wastes the implementation work.
6. **Implement in fine-grained, task-specific commits.** Touch only files relevant to the
   request per commit, and only commits relevant to the request overall. If you notice an
   unrelated issue while working, note it in the PR description rather than fixing it in the
   same branch. Each commit follows Conventional Commits (see
   [CONTRIBUTING.md](../CONTRIBUTING.md#commit-messages-conventional-commits)) and is one
   reviewable logical change, not the whole task at once.
7. **Validate locally**: `make validate`, `make lint`, `make test`, `make docs`, and
   `make test-conformance` if applicable — see
   [docs/development/README.md](../docs/development/README.md).
8. **Self-review your own diff** using the checklist below before writing the PR description.
9. **Push the branch (never `main`) and write the PR description** using
   [.github/pull_request_template.md](pull_request_template.md), stating assumptions explicitly.

## Branching and Commit Workflow

Every task starts from a fresh branch off `main` and ends as a pushed branch with a PR — never
a direct commit to `main`. This applies regardless of which AI agent or tool is doing the work
(see [AGENTS.md](../AGENTS.md)'s tool-agnostic framing).

### 1. Sync with the default branch

```bash
git fetch origin
git checkout main
git pull origin main
```

Always branch from the current tip of `main`, never from a stale local copy or on top of
another feature branch.

### 2. Create a task-specific branch

Name it `<type>/<short-kebab-case-description>` per
[CONTRIBUTING.md](../CONTRIBUTING.md#branch-naming), using the same types as commits:

```bash
git checkout -b docs/provider-adapter-contract
```

One branch per task. Do not pile a second, unrelated task onto an existing branch — cut a new
one from `main` instead, even if it means repeating step 1.

### 3. Make fine-grained, task-specific commits

Each commit is one reviewable logical change — a spec change, its contract derivation, a docs
update — not the whole task squashed into one commit. Follow
[CONTRIBUTING.md](../CONTRIBUTING.md#commit-messages-conventional-commits)'s Conventional
Commits format (enforced by `commitlint.yml`); the [.gitmessage](../.gitmessage) template
documents the exact shape if you enable it locally
(`git config commit.template .gitmessage` — repo-local only, never set globally on someone's
behalf).

**Commit identity:** use whatever author name/email the user has specified for the session. Do
not fall back to a machine's global `git config user.*` without confirming it's the identity
intended for this repo — and never run `git config` to change identity globally; scope an
override to the individual `git commit` invocation instead (e.g.
`git -c user.name="..." -c user.email="..." commit`).

### 4. Push the branch and open a PR

```bash
git push -u origin <branch-name>
gh pr create --fill
```

Never push a feature branch's commits directly to `main`. A PR from a non-contributor pauses in
CI for maintainer approval per [docs/development/ci.md](../docs/development/ci.md) — that is
expected, not a failure to fix.

## Self-Review Checklist

Before opening a PR, confirm each of these — don't just assume they hold:

- [ ] Work happened on a branch cut from the current tip of `main`, not on `main` itself and not
      stacked on a stale or unrelated branch.
- [ ] Every changed file is relevant to the stated request; nothing was touched "while I was in
      there."
- [ ] History is fine-grained: each commit is one logical, reviewable change, not the whole task
      squashed together.
- [ ] If `specs/` or `contracts/` changed, the compatibility impact (patch/minor/breaking, per
      [VERSIONING.md](../VERSIONING.md)) is stated in the PR description.
- [ ] If an ADR/RFC was required, it's attached/linked — not deferred to "a follow-up."
- [ ] No canonical field's existing meaning changed silently (see
      [specs/compatibility/README.md](../specs/compatibility/README.md)).
- [ ] No provider-specific vocabulary leaked into a canonical spec/contract (see
      [AGENTS.md §7](../AGENTS.md#7-contract-change-rules)).
- [ ] No new dependency was added without answering
      [docs/development/dependency-policy.md](../docs/development/dependency-policy.md)'s
      checklist in the PR description.
- [ ] No secrets, real credentials, or real customer/transaction data appear anywhere in the
      diff, including in examples and fixtures (see
      [AGENTS.md §10](../AGENTS.md#10-security-restrictions)).
- [ ] Every behavior change has a corresponding test or conformance case; every user-visible or
      contract-visible change has a documentation update.
- [ ] `make validate && make lint && make test && make docs` all pass locally.
- [ ] Any ambiguity you resolved by assumption is stated explicitly in the PR description, not
      silently decided.

## If You Get Stuck

Per [AGENTS.md §15](../AGENTS.md#15-if-you-are-still-unsure): a partially-done, clearly-explained
change is more valuable than a fully "finished" one built on an invented assumption. Open a
draft PR (or comment on the issue) describing what you tried, what's ambiguous, and the options
you considered — do not keep iterating privately trying to force a confident-looking answer.
