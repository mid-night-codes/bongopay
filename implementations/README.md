# implementations/ — Reference Implementation(s)

This directory holds implementations that demonstrate the BongoPay specification working
end-to-end. Per [ARCHITECTURE.md §8](../ARCHITECTURE.md#8-reference-implementation-boundary):

> The reference implementation demonstrates the BongoPay specification. It does not define the
> specification.

```text
specs/  → contracts/  → conformance/  → implementations/ (+ adapters/, sdks/)
```

## Status

**Starting — Phase 1.** [reference/](reference/) will be written in **Go**, per
[ADR 0002](../adr/0002-reference-implementation-language-go.md).

## Layout

| Directory | Covers |
|---|---|
| [reference/](reference/) | The Go reference implementation (Phase 1) |

## Rules for Working in `implementations/`

1. **Implementation follows spec, never the reverse.** If you need behavior not described in
   `specs/`, stop and flag the gap (see
   [AGENTS.md §7](../AGENTS.md#7-contract-change-rules)) rather than inventing it here.
2. **Must pass conformance.** An implementation is only correct if it passes the suite in
   [conformance/](../conformance/README.md) — compiling or "looking right" is not sufficient.
3. **The reference implementation's language choice must not leak into `specs/` or
   `contracts/`.** Those stay language-neutral regardless of what this directory picks.
4. **Persistence and transport are implementation details** unless a spec says otherwise. Don't
   let a database or message-broker choice creep into canonical domain concepts — see
   [ARCHITECTURE.md §14](../ARCHITECTURE.md#14-what-this-document-does-not-decide),
   `TODO(ADR): Persistence architecture`.
5. **Additional implementations are welcome** (e.g. in a different language) once the reference
   one exists, but each must independently satisfy conformance — one implementation's internals
   never become an implicit second source of truth.

## What Does Not Belong Here

- Provider-specific integration code — that's [adapters/](../adapters/README.md).
- Client libraries for application developers — that's [sdks/](../sdks/README.md).
- Runnable sample applications that use an implementation — that's
  [examples/](../examples/README.md).
