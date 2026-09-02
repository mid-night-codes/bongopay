# tests/ — Shared Test Fixtures

This directory holds test fixtures and data shared across multiple implementations, adapters,
or SDKs — the kind of thing that would otherwise be duplicated per-language. It is distinct from
[conformance/](../conformance/README.md), which defines *what* must be tested; this directory
holds reusable *inputs* for those tests once they're executable (Phase 1+, see
[ROADMAP.md](../ROADMAP.md)).

## Status

**Empty — Phase 0.** No shared fixtures exist yet; there is nothing to share across
implementations until at least one implementation exists.

## What Belongs Here

- Fixture payloads (e.g. sample `PaymentRequest`/`Payment` JSON) reused by more than one
  language's test suite, validated against [contracts/json-schema/](../contracts/json-schema/).
- Golden files for conformance cases in [conformance/](../conformance/README.md), once those
  cases are executable rather than specification-only.

## What Does Not Belong Here

- Language-specific unit tests — those live alongside the code they test in
  [implementations/](../implementations/README.md), [adapters/](../adapters/README.md), or
  [sdks/](../sdks/README.md).
- Conformance case *definitions* — those live in [conformance/](../conformance/README.md); this
  directory only holds data those cases might consume.
- Real customer or transaction data, ever — see
  [AGENTS.md §10](../AGENTS.md#10-security-restrictions). Fixtures use obviously-fake values
  only.
