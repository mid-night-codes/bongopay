# contracts/ — Machine-Readable Derivations of specs/

This directory holds the machine-readable artifacts derived from [specs/](../specs/README.md):
OpenAPI, AsyncAPI, and JSON Schema. Nothing here is authored from scratch — every contract here
should be traceable to a specific spec document that motivates it.

```text
specs/  → contracts/  → conformance/  → implementations/ (+ adapters/, sdks/)
```

## Status

**Draft / Phase 0.** Directory structure and JSON Schema subdirectories exist as placeholders;
no OpenAPI/AsyncAPI documents or schema files have been written yet, since the canonical specs
they'd derive from ([specs/payments/](../specs/payments/README.md) and friends) are themselves
still Draft.

## Layout

| Directory | Covers | Derived from |
|---|---|---|
| [openapi/](openapi/) | REST API contract | [specs/payments/](../specs/payments/README.md), [specs/errors/](../specs/errors/README.md) |
| [asyncapi/](asyncapi/) | Event/webhook contract | [specs/events/](../specs/events/README.md) |
| [json-schema/payments/](json-schema/payments/) | `Payment`, `PaymentRequest`, `Money`, etc. as JSON Schema | [specs/payments/](../specs/payments/README.md) |
| [json-schema/events/](json-schema/events/) | Event envelope/type schemas | [specs/events/](../specs/events/README.md) |
| [json-schema/errors/](json-schema/errors/) | Canonical error shape as JSON Schema | [specs/errors/](../specs/errors/README.md) |
| [examples/](examples/) | Example payloads validated against the schemas above | — |

## Rules for Working in `contracts/`

1. **Never author a contract change here first.** A field, event, or error that doesn't exist
   in `specs/` yet doesn't belong in a contract either — update the spec first (with an ADR/RFC
   if required, per [AGENTS.md §8](../AGENTS.md#8-when-adrs-and-rfcs-are-required)), then derive
   the contract change from it in the same PR.
2. **Generated files are never hand-edited.** Per
   [AGENTS.md §11–12](../AGENTS.md#11-identifying-generated-files), once `make generate` is
   wired up (see [ROADMAP.md](../ROADMAP.md) Phase 2), anything under here carrying the
   generated-file header described there must be regenerated from its source, not patched
   directly. As of Phase 0 no generation pipeline exists, so any file added here today is
   hand-authored —
   treat that as temporary scaffolding, not a precedent.
3. **No provider-specific vocabulary.** Same rule as `specs/`: provider quirks live in an
   adapter or a namespaced `providerOptions` extension
   ([specs/providers/extensions.md](../specs/providers/extensions.md)), never as a first-class
   field in a shared contract.
4. **Validate before proposing a change.** `make validate` and `make check-contracts` (once
   generation exists) must pass — see [scripts/README.md](../scripts/README.md).
5. **Independent versioning.** Each contract type versions independently — see
   [VERSIONING.md](../VERSIONING.md).

## What Does Not Belong Here

- Prose explaining *why* a concept exists — that's `specs/`.
- Executable test assertions — that's `conformance/`.
- Language-specific generated client code — that's `sdks/` (future, generated *from* here).
