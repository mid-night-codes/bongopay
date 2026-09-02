# sdks/ — Per-Language Client SDKs

This directory holds thin, idiomatic client libraries that give application developers access to
the canonical BongoPay contract in their language of choice, generated from
[contracts/](../contracts/README.md) wherever practical.

```text
specs/  → contracts/  → conformance/  → implementations/ (+ adapters/, sdks/)
```

## Status

**Not started — Phase 2 ("Developer Tooling").** See [ROADMAP.md](../ROADMAP.md),
`TODO(RFC): SDK generation pipeline and target languages for v1` in
[ARCHITECTURE.md §14](../ARCHITECTURE.md#14-what-this-document-does-not-decide) and
[AGENTS.md §12](../AGENTS.md#12-generated-files). No SDK, generated or hand-written, exists yet.

## Rules for Working in `sdks/`

1. **Prefer generated over hand-written.** Once a generation pipeline exists (Phase 2), an SDK
   under here derived from `contracts/` should be treated as generated output — see
   [AGENTS.md §11](../AGENTS.md#11-identifying-generated-files) — and edited at its source, not
   by hand.
2. **No business logic.** SDKs must not encode business logic or provider-specific behavior
   beyond what the canonical contract and provider capability discovery already expose (see
   [ARCHITECTURE.md §7](../ARCHITECTURE.md#7-sdk-boundary)).
3. **One directory per language/ecosystem**, versioned independently per
   [VERSIONING.md](../VERSIONING.md) (e.g. `sdks/typescript/`, `sdks/python/`), once any SDK
   work starts.
4. **An SDK is not a place to work around a contract gap.** If the canonical contract is missing
   something an SDK needs, fix the contract (with an ADR/RFC as required) rather than adding an
   SDK-only affordance that isn't backed by `contracts/`.

## What Does Not Belong Here

- Server-side implementation code — that's [implementations/](../implementations/README.md).
- Provider-specific integration code — that's [adapters/](../adapters/README.md).
- Full sample applications built with an SDK — that's [examples/](../examples/README.md).
