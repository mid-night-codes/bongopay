# 0002. Reference Implementation Language: Go

- **Status:** Accepted
- **Date:** 2026-09-02
- **Deciders:** Mid Night Coder (Core Maintainer)
- **Related:** [ARCHITECTURE.md §14](../ARCHITECTURE.md#14-what-this-document-does-not-decide),
  [ARCHITECTURE.md §8](../ARCHITECTURE.md#8-reference-implementation-boundary),
  [implementations/README.md](../implementations/README.md),
  [implementations/reference/README.md](../implementations/reference/README.md)

## Context

[ROADMAP.md](../ROADMAP.md) Phase 1 ("Simulator Core") requires implementing the payment
lifecycle against the canonical state machine, an executable scenario engine, webhook
simulation, and a first working REST contract. All of that lands under
`implementations/reference/`, which — per
[implementations/README.md](../implementations/README.md) and
[ARCHITECTURE.md §14](../ARCHITECTURE.md#14-what-this-document-does-not-decide) — has been
explicitly blocked pending this decision, since writing code before it exists would either
presuppose the answer or need to be redone once it's decided.

The choice must not constrain the language-neutrality of `specs/` or `contracts/` (see
[ARCHITECTURE.md §8](../ARCHITECTURE.md#8-reference-implementation-boundary): "The reference
implementation demonstrates the BongoPay specification. It does not define the specification.").
It does, however, shape `adapters/`, `sdks/`, `deploy/`, and the initial conformance harness,
since those will most naturally start in whatever language the reference implementation uses.

Candidates considered: Go, TypeScript (Node.js), Python (FastAPI/Pydantic), Rust, and
Java/Kotlin. Evaluation criteria, in order of weight for this decision:

1. Fit with [docs/development/dependency-policy.md](../docs/development/dependency-policy.md)'s
   minimal-dependency default.
2. Fit for simulating the async, race-prone behaviors Phase 1 needs — timeouts, duplicate
   callbacks, out-of-order events (see
   [ARCHITECTURE.md §6](../ARCHITECTURE.md#6-simulator-boundary)).
3. Long-term maintainability: how well the language holds up under years of multi-contributor
   change, not just initial development velocity.
4. Contract tooling maturity (JSON Schema / OpenAPI generation and validation).
5. Accessibility to the broad human *and* AI-agent contributor base this project explicitly
   courts (see [AGENTS.md](../AGENTS.md), [README.md](../README.md) "Contributing").

## Decision

The reference implementation (`implementations/reference/`) will be written in **Go**.

## Alternatives Considered

- **TypeScript (Node.js)** — the strongest contract tooling (ajv, zod, openapi-typescript) and
  the most natural first-SDK story, with Node already a soft dependency in this repo's own
  tooling. Not chosen because its type system is erased at runtime — schema enforcement depends
  on a validation library layered on top rather than the language itself — and npm's dependency
  graphs run against this project's minimal-dependency default more than Go's sufficient
  standard library does. Close second; revisit if contract/SDK-generation needs outweigh this.
- **Python (FastAPI/Pydantic)** — fastest to prototype, and FastAPI generating an OpenAPI
  document directly from Pydantic models is a genuine advantage for the REST contract work. Not
  chosen because real type enforcement requires sustained `mypy --strict` discipline rather than
  being structural, and large Python codebases tend to degrade faster under refactoring without
  that discipline maintained indefinitely — a long-term-maintainability risk specifically.
- **Rust** — the strongest correctness guarantees of any candidate (exhaustive enum matching
  would make an unhandled `PaymentStatus` transition a compile error, not a runtime bug). Not
  chosen because its learning curve is the steepest of the group, which cuts against the broad
  human-and-AI-agent contribution this project explicitly wants, and its async ecosystem is more
  complex than Go's concurrency model for the same webhook-simulation use case. Worth
  reconsidering later for a narrow, performance-critical component, not the whole reference
  implementation.
- **Java/Kotlin** — the closest match to how real-world payment infrastructure is actually
  built, with mature OpenAPI codegen and strong large-team refactoring tooling. Not chosen
  because its build tooling (Maven/Gradle) and runtime footprint are heavier than this project's
  stated "no heavy toolchain required" ethos calls for.

## Consequences

- `implementations/reference/README.md`'s "Blocked" status is lifted; Go code may now be written
  there.
- The standard library covers HTTP, JSON, and concurrency primitives sufficient for most of
  Phase 1 without a third-party dependency — new dependencies still go through
  [docs/development/dependency-policy.md](../docs/development/dependency-policy.md)'s checklist,
  not a blanket exemption.
- `adapters/` and the initial `conformance/` harness will most naturally start in Go too, though
  neither is required to stay in Go long-term — conformance cases are language-neutral by
  design (see [conformance/README.md](../conformance/README.md)).
- `sdks/`'s first entry is not implied by this decision — a Go reference implementation does not
  make Go the first client SDK language; that remains a separate, later decision.
- CI (`.github/workflows/ci.yml`) will need a Go toolchain setup step once
  `implementations/reference/` has code to build/test — tracked as follow-up work, not part of
  this ADR.
- `specs/` and `contracts/` remain language-neutral; nothing here changes their content or
  review rules (see [AGENTS.md §7](../AGENTS.md#7-contract-change-rules)).

## Scope Check

- [x] Confirmed this is an ADR-level decision, not an RFC-level one — it's a foundational choice
      but not a breaking change to a published contract, provider plugin model, event
      architecture, security model, or versioning policy (see
      [AGENTS.md §8](../AGENTS.md#8-when-adrs-and-rfcs-are-required)).
- [x] Confirmed this does not change the meaning of any existing canonical field, event, or
      error — it constrains an implementation, not a spec.
