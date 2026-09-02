---
name: feature-branch
description: Start any BongoPay task from a fresh branch off main, make fine-grained task-specific Conventional Commits, and push/open a PR — the repo's required contribution flow instead of committing straight to main.
---

# Feature Branch Workflow

Use this before making any commit in this repository. BongoPay's CI
(`.github/workflows/ci.yml`, `.github/workflows/commitlint.yml`) and
[GOVERNANCE.md](../../../GOVERNANCE.md) both assume changes arrive via pull request — commits
pushed straight to `main` skip the maintainer-approval gate and any review entirely.

## 1. Sync with the default branch

```bash
git fetch origin
git checkout main
git pull origin main
```

Always branch from the current tip of `main`, never from a stale local copy or on top of
another feature branch.

## 2. Create a task-specific branch

Name it `<type>/<short-kebab-case-description>` per
[CONTRIBUTING.md](../../../CONTRIBUTING.md#branch-naming), using the same types as commits:

```bash
git checkout -b docs/provider-adapter-contract
```

One branch per task. Do not pile a second, unrelated task onto an existing branch — cut a new
one from `main` instead, even if it means repeating step 1.

## 3. Make fine-grained, task-specific commits

Each commit is one reviewable logical change — a spec change, its contract derivation, a docs
update — not the whole task squashed into one commit. Follow
[CONTRIBUTING.md](../../../CONTRIBUTING.md#commit-messages-conventional-commits)'s Conventional
Commits format (enforced by `commitlint.yml`); the [.gitmessage](../../../.gitmessage) template
documents the exact shape if you enable it locally.

Before writing any file, follow [AGENTS.md §1](../../../AGENTS.md#1-read-before-you-write) (what
to read first) and scope the change per
[AGENTS.md §13](../../../AGENTS.md#13-how-to-make-small-reviewable-changes).

**Commit identity:** use whatever author name/email the user has specified for the session. Do
not fall back to a machine's global `git config user.*` without confirming it's the identity
intended for this repo — and never run `git config` to change identity globally; scope an
override to the individual `git commit` invocation instead (e.g.
`git -c user.name="..." -c user.email="..." commit`).

## 4. Push the branch and open a PR

```bash
git push -u origin <branch-name>
gh pr create --fill
```

Never push a feature branch's commits directly to `main`. A PR from a non-contributor pauses in
CI for maintainer approval per [docs/development/ci.md](../../../docs/development/ci.md) — that
is expected, not a failure to fix.

## 5. Before opening the PR

Run the same checks CI runs:

```bash
make validate && make lint && make test && make docs
```

Then self-review against
[.github/CONTRIBUTING_AGENT.md](../../../.github/CONTRIBUTING_AGENT.md)'s checklist before
handing the PR to a maintainer.
