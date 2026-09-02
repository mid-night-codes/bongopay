# specs/

**This directory is the architectural source of truth for BongoPay.** If a concept, field, or
behavior isn't described here, it is not yet part of BongoPay — no matter what any
implementation does.

```text
specs/  → contracts/  → conformance/  → implementations/ (+ adapters/, sdks/)
```

See [ARCHITECTURE.md](../ARCHITECTURE.md) for the full picture. This README covers how to work
safely inside `specs/` specifically.

## Layout

| Directory | Covers |
|---|---|
| [payments/](payments/README.md) | Canonical payment domain concepts: `Payment`, `PaymentRequest`, `Money`, etc. |
| [providers/](providers/README.md) | Provider adapter contract, capability discovery, extension mechanism |
| [events/](events/README.md) | Canonical event envelope and event types |
| [errors/](errors/README.md) | Canonical, transport-agnostic error model |
| [scenarios/](scenarios/README.md) | Declarative simulator scenario format |
| [state-machines/](state-machines/README.md) | Canonical payment lifecycle (states and transitions) |
| [compatibility/](compatibility/README.md) | Cross-cutting backward-compatibility rules for all of the above |

## Rules for Working in `specs/`

1. **Prose first, schema second.** A spec here is a Markdown document explaining intent,
   rationale, and constraints in human language. Machine-readable derivations belong in
   [contracts/](../contracts/README.md), generated *from* what's decided here — never the
   reverse.
2. **No language- or framework-specific vocabulary.** No Java classes, no Spring annotations,
   no ORM entities, no Kafka topics. If you find yourself describing a class hierarchy, you're
   probably writing an implementation doc — it belongs elsewhere (see
   [implementations/README.md](../implementations/README.md)).
3. **No provider-specific vocabulary in canonical specs.** Provider quirks belong in
   [providers/](providers/README.md) as adapter-mapping guidance, or under a provider's
   namespaced `providerOptions` extension (see
   [providers/extensions.md](providers/extensions.md)) — never as a first-class canonical field.
4. **Changes here are contract changes.** Per [AGENTS.md §8](../AGENTS.md#8-when-adrs-and-rfcs-are-required),
   a new shared concept usually warrants an ADR; a breaking change to an existing concept
   requires an RFC. Small clarifications/typo fixes do not.
5. **Leave TODOs visible.** Where a decision hasn't been made yet, write
   `TODO(ADR): ...` or `TODO(RFC): ...` directly in the spec rather than silently picking an
   answer. This is what lets AI agents and new contributors tell "decided" from "open" at a
   glance.
6. **Versioned independently.** Each spec area has its own versioning story — see
   [VERSIONING.md](../VERSIONING.md) and [compatibility/README.md](compatibility/README.md).

## What Does Not Belong Here

- Executable code of any kind (that's `implementations/`, `adapters/`, `sdks/`).
- Generated OpenAPI/AsyncAPI/JSON Schema files (that's `contracts/`, generated *from* these
  specs — see [contracts/README.md](../contracts/README.md)).
- Test assertions/harnesses (that's `conformance/`, though conformance *scenarios* are
  described conceptually here where relevant).
