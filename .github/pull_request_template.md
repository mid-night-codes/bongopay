<!--
Thanks for contributing to BongoPay. Fill in every section — "N/A" is a fine answer where a
section genuinely doesn't apply, but don't delete a section instead of answering it.
Human and AI-agent contributors: see CONTRIBUTING.md and AGENTS.md before opening this PR.
-->

## Summary

What does this PR do, and why? Link the issue/RFC/ADR it addresses, if any.

## Type of Change

- [ ] `feat` — new capability
- [ ] `fix` — bug fix, no contract change
- [ ] `docs` — documentation only
- [ ] `test` / `conformance` — test or conformance case addition
- [ ] `refactor` — internal change, no observable behavior change
- [ ] `chore` / `ci` / `build` — tooling
- [ ] Breaking change (requires a linked, accepted RFC — see AGENTS.md §8)

## Contract Impact

- [ ] This PR does **not** touch `specs/` or `contracts/`.
- [ ] This PR touches `specs/` or `contracts/`. Compatibility impact:
      <!-- Patch / Minor / Breaking, per VERSIONING.md — state which and why. -->
- [ ] An ADR is attached/linked (required for a new shared abstraction, persistence change, or
      cross-module behavior — see AGENTS.md §8).
- [ ] An RFC is attached/linked (required for a breaking contract change, new provider plugin
      model, new event architecture, security model change, or versioning policy change).

## Security Impact

- [ ] This PR does not touch webhook/callback handling, idempotency, or signature verification.
- [ ] This PR touches one of the above. Explain how it affects the threat categories in
      [SECURITY.md](../SECURITY.md#threat-categories-bongopay-must-eventually-consider):
      <!-- e.g. "adds signature verification for provider X callbacks" -->

## Assumptions

State any assumption you made where the spec, an existing test, or a directory README was
ambiguous or silent, per [AGENTS.md §1](../AGENTS.md#1-read-before-you-write). If you made no
assumptions, say so explicitly.

## Testing

- [ ] `make validate` passes
- [ ] `make lint` passes
- [ ] `make test` passes
- [ ] `make docs` passes
- [ ] `make test-conformance` passes (if this PR touches a spec, contract, or adapter)

## Dependencies

- [ ] No new dependencies added.
- [ ] New dependency added. Answered the checklist in
      [docs/development/dependency-policy.md](../docs/development/dependency-policy.md):
      <!-- stdlib alternative considered? maintained + compatible license? worth it vs.
           convenience? does it constrain specs/contracts neutrality? -->

## Checklist

- [ ] Documentation updated for any user-visible or contract-visible change.
- [ ] No secrets, real credentials, or real customer/transaction data committed.
- [ ] PR is small and focused (see [CONTRIBUTING.md](../CONTRIBUTING.md#keep-pull-requests-small)).
