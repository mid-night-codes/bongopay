# examples/ — Example Applications

This directory holds small, runnable applications that demonstrate integrating BongoPay against
the simulator (and, later, real provider adapters) — the thing a new contributor or evaluator
runs to see the contract working, rather than reading specs in the abstract.

## Status

**Not started — Phase 2 ("Developer Tooling"),** and dependent on
[implementations/reference/](../implementations/README.md) and
[sdks/](../sdks/README.md) existing first. See [ROADMAP.md](../ROADMAP.md). No examples exist
yet.

## Rules for Working in `examples/`

1. **Examples run against the simulator by default**, never a real provider or real money — see
   [ARCHITECTURE.md §6](../ARCHITECTURE.md#6-simulator-boundary) and
   [ROADMAP.md](../ROADMAP.md) §"What NOT to Build Yet".
2. **An example demonstrates the canonical contract, not a workaround for it.** If making an
   example work requires bypassing something in `specs/`/`contracts/`, that's a signal the
   contract has a gap — fix that, don't route around it here.
3. **Keep examples minimal.** An example should demonstrate one flow clearly (e.g. "initiate a
   payment and handle its webhook") rather than building a full sample product.
4. **Use only obviously-fake data** — test MSISDNs, placeholder API keys, fake customer
   references — per [AGENTS.md §10](../AGENTS.md#10-security-restrictions).
5. **Keep examples runnable.** An example that no longer runs against the current contract is a
   bug — update it in the same PR that changes the contract it demonstrates, or remove it.

## What Does Not Belong Here

- Conformance assertions — that's [conformance/](../conformance/README.md); examples may reuse
  conformance scenarios for realism, but examples themselves are not test suites.
- Anything reaching a real payment provider or moving real money.
