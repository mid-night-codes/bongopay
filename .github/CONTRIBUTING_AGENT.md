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
2. **Find or create an open issue for this task** — before branching, before any commit. See
   [Branching and Commit Workflow](#branching-and-commit-workflow) step 0 below. A commit or PR
   with no issue behind it is incomplete, not just under-documented.
3. **Branch from the current tip of `main`** before touching anything — see
   [Branching and Commit Workflow](#branching-and-commit-workflow) below for the exact commands
   and branch naming. Never commit directly to `main`; never build on top of a stale or
   unrelated branch.
4. **Read in the order AGENTS.md §1 specifies** — this file, the relevant directory README(s),
   any relevant ADR, the governing spec, existing tests/conformance cases. Do not skip to
   writing code because the change "looks small."
5. **Locate the smallest change.** If the smallest change that satisfies the request still
   touches multiple directories (e.g. a spec change plus its contract derivation), that's
   expected — "small" means minimal *scope*, not minimal *file count*.
6. **Check the ADR/RFC bar** (AGENTS.md §8) before writing the change, not after — discovering
   partway through an implementation that it needed an RFC wastes the implementation work.
7. **Implement in fine-grained, task-specific commits.** Touch only files relevant to the
   request per commit, and only commits relevant to the request overall. If you notice an
   unrelated issue while working, note it in the PR description rather than fixing it in the
   same branch. Each commit follows Conventional Commits (see
   [CONTRIBUTING.md](../CONTRIBUTING.md#commit-messages-conventional-commits)) and is one
   reviewable logical change, not the whole task at once.
8. **Validate locally**: `make validate`, `make lint`, `make test`, `make docs`, and
   `make test-conformance` if applicable — see
   [docs/development/README.md](../docs/development/README.md).
9. **Self-review your own diff** using the checklist below before writing the PR description.
10. **Push the branch (never `main`) and write the PR description** using
    [.github/pull_request_template.md](pull_request_template.md), stating assumptions
    explicitly, and linking the issue from step 2 (`Closes #N` or `Refs #N`).

## Branching and Commit Workflow

Every task starts from a fresh branch off `main` and ends as a pushed branch with a PR — never
a direct commit to `main`. This applies regardless of which AI agent or tool is doing the work
(see [AGENTS.md](../AGENTS.md)'s tool-agnostic framing).

### 0. Find or create an open issue

No commit happens without an open issue behind it. Before doing anything else:

```bash
gh issue list -R <owner>/<repo> --search "<keywords>" --state open
```

- **A matching open issue exists** — use it; note its number for step 4 and the PR.
- **Nothing matches** — create one first, using the closest-fitting template in
  [.github/ISSUE_TEMPLATE/](../.github/ISSUE_TEMPLATE/) (`gh issue create` also works directly
  with `--title`/`--body` when a task doesn't cleanly fit a template):

  ```bash
  gh issue create -R <owner>/<repo> --title "..." --body "..."
  ```

This applies even to work you were just asked to do directly in conversation — the issue is
what makes the *why* discoverable independent of the commit history or a chat transcript, and
gives step 1 above ("restate the request") something durable to restate *into*, not just an
ephemeral instruction. Skipping this because the task "is obviously small" is exactly the case
this step exists to catch — see [AGENTS.md §15](../AGENTS.md#15-if-you-are-still-unsure)'s
general principle of not silently deciding something doesn't need the process.

### 1. Sync with the default branch

```bash
git fetch origin
git checkout main
git pull origin main
```

Branch from the current tip of `main`, never from a stale local copy — **unless** the task
genuinely depends on unmerged work in another open PR, in which case branch from that PR's
branch instead. See [Stacked PRs](#stacked-prs) below for when that applies and how to keep it
manageable.

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

Reference the issue from step 0 in at least one commit's footer (`Refs #N`, or `Closes #N` on
the commit that actually resolves it) — the [.gitmessage](../.gitmessage) footer section shows
the exact form.

### 4. Push the branch and open a PR

```bash
git push -u origin <branch-name>
gh pr create --fill
```

Then link the issue from step 0 in the PR body if it isn't already covered by a commit footer
GitHub picked up (`Closes #N` auto-closes it on merge; use `Refs #N` if the PR doesn't fully
resolve it). Never push a feature branch's commits directly to `main`. A PR from a
non-contributor pauses in CI for maintainer approval per
[docs/development/ci.md](../docs/development/ci.md) — that is expected, not a failure to fix.

## Stacked PRs

Sometimes a task genuinely needs code or docs that only exist on another open, unmerged PR's
branch — not just "it would be convenient to skip a merge wait." Use a stacked PR for that case
instead of either blocking on the other PR merging first, or duplicating its unmerged content.

**When it applies:** the new task imports a type, calls a function, or edits a section that only
exists on `some-open-pr-branch`, and rewriting around that dependency would be artificial (e.g.
extending a "Branching and Commit Workflow" section a still-open PR is introducing, the way this
change stacks on [PR #9](https://github.com/mid-night-codes/bongopay/pull/9)). It does not apply
just because two tasks happen to be in flight at the same time with no real dependency between
them — those still each branch from `main` per step 1.

**How to do it:**

```bash
git fetch origin
git checkout -b <new-branch> origin/<base-pr-branch>
git branch --unset-upstream   # avoid accidentally pushing over the base PR's branch, see below
```

`git checkout -b <new> <remote-ref>` sets up the new local branch to track that remote ref by
default — meaning a bare `git push` would push to the *base PR's* remote branch, not create
yours. Always follow with `git branch --unset-upstream`, and always push with an explicit
branch name:

```bash
git push -u origin <new-branch>
gh pr create --fill --base <base-pr-branch> --head <new-branch>
```

The `--base` is what makes it a stacked PR on GitHub — the diff shown is only *your* commits on
top of the base branch, not the base branch's changes too.

**Keeping it manageable:**

- Prefer a stack depth of one. If task C depends on unmerged B which depends on unmerged A,
  that's a sign to slow down and get A merged before starting C, not to build a three-deep
  stack — conflicts compound with depth.
- Merge bottom-up: the base PR merges first, then retarget the stacked PR's base to `main`
  (`gh pr edit <n> --base main`) and rebase (`git rebase main`) before it merges too.
- If the base PR changes after you've branched from it (a fixup commit, a force-push from
  review feedback), rebase your stacked branch onto the new tip rather than leaving it stale.
- State the stacking explicitly in the stacked PR's description (which PR/issue it depends on
  and why) — a reviewer seeing an unfamiliar base branch with no explanation will assume it's a
  mistake.

## Self-Review Checklist

Before opening a PR, confirm each of these — don't just assume they hold:

- [ ] An open issue exists for this work and is referenced (`Refs #N`/`Closes #N`) in a commit
      or the PR body — created before the first commit, not added after the fact.
- [ ] Work happened on a branch cut from the current tip of `main`, not on `main` itself — or,
      if genuinely [stacked](#stacked-prs), the PR's `--base` matches the real dependency and
      that dependency is stated in the PR description.
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
