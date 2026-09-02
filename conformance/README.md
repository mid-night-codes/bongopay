# conformance/ — Shared, Language-Agnostic Conformance Definitions

This directory defines the tests that any implementation, adapter, or provider integration must
pass to be considered a conforming BongoPay component. It is the mechanism that keeps
"language-neutral" honest (see
[ARCHITECTURE.md §9](../ARCHITECTURE.md#9-conformance-testing)): the same conceptual test — e.g.
"duplicate callbacks must not create duplicate canonical events" — must be expressible against
any implementation, in any language.

```text
specs/  → contracts/  → conformance/  → implementations/ (+ adapters/, sdks/)
```

## Status

**Draft / Phase 0.** As of this phase, conformance is **specification-only**: this directory
describes what must be tested and why, in prose/data form. Executable harnesses that actually
run these cases against a real implementation are a Phase 1+ concern (see
[ROADMAP.md](../ROADMAP.md)).

## Layout

| Directory | Covers |
|---|---|
| [api/](api/) | Conformance cases for the REST/API contract in [contracts/openapi/](../contracts/openapi/) |
| [webhook/](webhook/) | Conformance cases for webhook/callback handling: signature verification, duplicate delivery, out-of-order delivery |
| [provider/](provider/) | Conformance cases every provider adapter must pass, per capability declared in [specs/providers/adapter-contract.md](../specs/providers/adapter-contract.md) |
| [scenarios/](scenarios/) | Conformance cases that exercise the simulator against [specs/scenarios/](../specs/scenarios/README.md) |

No case files exist yet in any of these subdirectories — the layout is established ahead of
content so Phase 1 work has a home without needing a structural PR first.

## Rules for Working in `conformance/`

1. **A conformance case is derived from a spec, not from an implementation's current
   behavior.** If an implementation does something a spec doesn't describe, that's a signal to
   fix the spec or the implementation — never to write a conformance case that just codifies
   whatever the code happens to do today.
2. **Cases must be implementation-neutral.** No language-specific test framework syntax lives
   here; a case describes inputs, actions, and expected canonical outcomes in a form any
   implementation's test suite can adapt.
3. **"Passes conformance" is the bar for acceptance**, not "compiles" or "looks right" — see
   [README.md](../README.md) principle 5 and
   [CONTRIBUTING.md](../CONTRIBUTING.md#adding-a-provider). This applies equally to the
   reference implementation, adapters, and the simulator.
4. **Security-relevant cases are not optional.** Webhook signature verification, idempotency
   handling, and callback validation cases in [webhook/](webhook/) must exist before an adapter
   or implementation claiming that capability is considered conforming — see
   [ARCHITECTURE.md §12](../ARCHITECTURE.md#12-security-boundaries).

## What Does Not Belong Here

- The tests themselves, once written per-language — those live alongside each implementation/
  adapter/SDK, satisfying the case defined here.
- New canonical behavior — a conformance case can only test what a spec already defines.
