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
2. **Read in the order AGENTS.md §1 specifies** — this file, the relevant directory README(s),
   any relevant ADR, the governing spec, existing tests/conformance cases. Do not skip to
   writing code because the change "looks small."
3. **Locate the smallest change.** If the smallest change that satisfies the request still
   touches multiple directories (e.g. a spec change plus its contract derivation), that's
   expected — "small" means minimal *scope*, not minimal *file count*.
4. **Check the ADR/RFC bar** (AGENTS.md §8) before writing the change, not after — discovering
   partway through an implementation that it needed an RFC wastes the implementation work.
5. **Implement.** Touch only files relevant to the request. If you notice an unrelated issue
   while working, note it in the PR description rather than fixing it in the same PR.
6. **Validate locally**: `make validate`, `make lint`, `make test`, `make docs`, and
   `make test-conformance` if applicable — see
   [docs/development/README.md](../docs/development/README.md).
7. **Self-review your own diff** using the checklist below before writing the PR description.
8. **Write the PR description** using
   [.github/pull_request_template.md](pull_request_template.md), stating assumptions explicitly.

## Self-Review Checklist

Before opening a PR, confirm each of these — don't just assume they hold:

- [ ] Every changed file is relevant to the stated request; nothing was touched "while I was in
      there."
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
