# AGENTS.md — Instructions for AI Coding Agents

This file is the entry point for any AI agent (Claude, GPT, Copilot Workspace, or any other
autonomous or semi-autonomous coding tool) working in this repository. **Read this file before
making any change.** It is a hard requirement, not a suggestion — a change that violates this
file should be treated as invalid regardless of how correct the code looks.

BongoPay treats AI agents as first-class contributors. That means the same discipline expected
of a careful human contributor is expected of you, expressed as explicit rules instead of
tribal knowledge.

## 1. Read Before You Write

In order, before touching any file:

1. This file (`AGENTS.md`).
2. The `README.md` of every directory you intend to modify.
3. Any `ADR` in `adr/` relevant to the area (search by keyword, don't assume none exist).
4. The relevant specification in `specs/` — implementation must conform to spec, never the
   reverse. If behavior you're asked to implement isn't in `specs/`, that's a signal to stop
   (see §7), not to invent it.
5. Existing tests/conformance cases touching the area.

If any of these are missing or ambiguous, say so explicitly in your output rather than guessing.

## 2. Repository Architecture (Summary)

Full detail: [ARCHITECTURE.md](ARCHITECTURE.md). The short version:

```text
specs/  → contracts/  → conformance/  → implementations/ (+ adapters/, sdks/)
```

Specifications are the source of truth. Contracts are machine-readable derivations of specs.
Conformance tests define correctness. Implementations, adapters, and SDKs must satisfy
conformance — they never define architecture on their own.

## 3. Directory Responsibilities

| Directory | Responsibility | Read this before editing |
|---|---|---|
| `specs/` | Language-neutral specifications (source of truth) | [specs/README.md](specs/README.md) |
| `contracts/` | OpenAPI / AsyncAPI / JSON Schema derived from specs | [contracts/README.md](contracts/README.md) |
| `conformance/` | Shared, language-agnostic conformance test definitions | [conformance/README.md](conformance/README.md) |
| `implementations/` | Reference implementation(s) | [implementations/README.md](implementations/README.md) |
| `adapters/` | Provider adapters | [adapters/README.md](adapters/README.md) |
| `sdks/` | Per-language client SDKs | [sdks/README.md](sdks/README.md) |
| `examples/` | Example applications | [examples/README.md](examples/README.md) |
| `docs/` | Human-facing documentation | [docs/README.md](docs/README.md) |
| `adr/` | Architecture Decision Records | [adr/README.md](adr/README.md) |
| `rfcs/` | RFCs for major/breaking change proposals | [rfcs/README.md](rfcs/README.md) |
| `deploy/` | Docker/Compose for local dev only | — |
| `scripts/` | Validation and developer tooling | — |

Each of these READMEs is the next link in the context chain:

```text
AGENTS.md → directory README → architecture docs → specification → ADR → implementation
```

Do not expect (or create) one giant instructions file. Instructions live close to where they
apply.

## 4. Commands You May Run

Prefer the root `Makefile` targets over ad hoc language-specific commands — they are the
stable interface regardless of what implementation languages exist underneath:

```bash
make setup             # install local validation tooling
make validate          # validate specs/contracts/schemas
make lint              # lint markdown, YAML, JSON, schemas
make test              # run available unit/integration tests
make test-conformance  # run conformance suite(s) as they come online
make check-contracts   # verify generated contracts match source specs
make docs              # build/check documentation
make generate          # regenerate generated artifacts from source specs
make clean             # remove local build artifacts
```

These are deterministic and safe to run repeatedly. Prefer them over inventing new tooling.

## 5. Commands You Should Avoid

- Do not run destructive git operations (`git push --force`, `git reset --hard`,
  `git clean -fdx`, branch deletion) without explicit human instruction.
- Do not install new dependencies as a side effect of "trying something" — see
  [Dependency Rules](#9-dependency-rules) below.
- Do not run commands that reach real external payment providers or send real network
  callbacks — BongoPay never touches real money or real provider endpoints (see
  [ROADMAP.md](ROADMAP.md) §"What NOT to Build Yet").
- Do not bypass validation (`--no-verify`, skipping `make validate`) to force a change through.

## 6. Testing, Formatting, and Documentation Expectations

- Every behavior change needs a corresponding test or conformance case. No silent behavior
  changes.
- Every contract change needs a corresponding schema/spec update in the same PR — never update
  a generated artifact by hand (see [§12](#12-generated-files)).
- Every user-visible or contract-visible behavior change needs a documentation update
  (`docs/`, the relevant `specs/` file, or both).
- Match existing formatting/linting conventions in the file you're editing; run `make lint`
  before proposing a change.

## 7. Contract-Change Rules

- Never modify a canonical concept in `specs/payments/`, `specs/state-machines/`,
  `specs/events/`, or `specs/errors/` without first checking whether an ADR or RFC is required
  (see §8 below).
- Never leak provider-specific vocabulary into canonical domain models. Provider specifics
  belong in `providerOptions` (see [specs/providers/extensions.md](specs/providers/extensions.md))
  or inside an adapter.
- Never silently change the meaning of an existing field. Add new fields/versions instead of
  overloading old ones — see [VERSIONING.md](VERSIONING.md).
- If you are asked to implement behavior that isn't in any spec: **stop and document the
  ambiguity** rather than inventing undocumented architecture. Propose the smallest spec change
  needed and flag it for human review.

## 8. When ADRs and RFCs Are Required

### Safe without architectural approval

Documentation fixes, tests, bug fixes that preserve existing contract behavior, examples,
internal refactoring, linting improvements.

### Requires ADR consideration

A new shared abstraction, persistence model changes, new orchestration behavior, changes
affecting multiple modules. Use the template at [adr/0000-template.md](adr/0000-template.md).

### Requires RFC

Breaking a public API/contract, a new provider plugin model, new event architecture, a security
model change, a versioning policy change. Use
[rfcs/0000-template.md](rfcs/0000-template.md) and see [rfcs/README.md](rfcs/README.md) for
process and statuses.

If you are not sure which bucket applies, treat it as the more conservative one and say so in
your output — do not decide silently and proceed.

## 9. Dependency Rules

Before adding any dependency, in any language, answer (and record the answer in the PR):

- Can the standard library solve this instead?
- Is the dependency actively maintained, with a compatible license?
- Does it meaningfully improve maintainability, versus just convenience?
- Could it constrain the language-neutrality of `specs/`/`contracts/`? If it touches those
  directories, the answer should almost always be "don't."

See [docs/development/dependency-policy.md](docs/development/dependency-policy.md).

## 10. Security Restrictions

- Never commit real credentials, tokens, private keys, production endpoints with embedded
  secrets, or real customer/transaction data. Use obviously-fake values in examples (e.g.
  `255700000000`-style test MSISDNs, `sk_test_...`-style placeholder keys).
- Never weaken or remove webhook signature verification, idempotency checks, or input
  validation to "make a test pass."
- Report anything that looks like a real secret already in the repository per
  [SECURITY.md](SECURITY.md) — do not just delete it silently; flag it.

## 11. Identifying Generated Files

Generated files carry a header such as:

```text
DO NOT EDIT MANUALLY.
Generated from: contracts/openapi/bongopay.yaml
Generate with: make generate
Validate with: make check-contracts
```

If you see this header, edit the **source** referenced in the header, then regenerate — never
hand-edit the generated file. See [§12](#12-generated-files) and
[contracts/README.md](contracts/README.md).

## 12. Generated Files

| Source | Generation command | Generated location | Validation command |
|---|---|---|---|
| `contracts/openapi/*.yaml` | `make generate` | (future) `sdks/*/generated/` | `make check-contracts` |
| `contracts/asyncapi/*.yaml` | `make generate` | (future) event schema bindings | `make check-contracts` |
| `contracts/json-schema/**/*.json` | `make generate` | (future) language type stubs | `make check-contracts` |

As of Phase 0, no generation pipeline is wired up yet — this table is the contract for when one
is (see ROADMAP Phase 2). TODO(RFC): SDK generation pipeline.

## 13. How to Make Small, Reviewable Changes

1. Understand the issue/request precisely — restate it if it's ambiguous.
2. Locate the governing specification and any existing tests.
3. Identify the smallest change that satisfies the request.
4. Check whether an ADR/RFC is required (§8). If yes, stop and produce that artifact first,
   or flag that one is needed before code changes.
5. Implement only files relevant to the issue.
6. Run `make validate` / `make lint` / `make test` (whichever apply).
7. Re-read your own diff before finishing — does it touch anything unrelated?
8. Update documentation for anything user-visible or contract-visible that changed.
9. Write a PR description that states your assumptions explicitly (see
   [.github/pull_request_template.md](.github/pull_request_template.md)).

See [.github/CONTRIBUTING_AGENT.md](.github/CONTRIBUTING_AGENT.md) for the detailed workflow
and a self-review checklist to run before opening a PR.

## 14. Guiding Priorities (Tie-Breakers)

When choices conflict, prefer, in order:

1. Architectural portability over implementation convenience.
2. The BongoPay specification over any single framework's convention.
3. A clear contract over a clever abstraction.
4. A small, reviewable change over a large autonomous one.

## 15. If You Are Still Unsure

Stop. Document the ambiguity, the options you considered, and why you didn't pick one, in your
PR description or output. A partially-done, clearly-explained change is more valuable — and
safer for this project — than a fully "finished" one built on an invented assumption.
