# adapters/ — Provider Adapters

This directory holds provider adapters: code that implements the capability contract in
[specs/providers/adapter-contract.md](../specs/providers/adapter-contract.md) against a specific
real payment provider, translating that provider's vocabulary into canonical BongoPay state,
events, and errors.

```text
specs/  → contracts/  → conformance/  → implementations/ (+ adapters/, sdks/)
```

## Status

**Not started — Phase 3 ("Provider Ecosystem"), with early groundwork acceptable once
[specs/providers/adapter-contract.md](../specs/providers/adapter-contract.md) exists.** See
[ROADMAP.md](../ROADMAP.md). No adapters, sample or otherwise, exist yet. The simulator
(specified in [specs/scenarios/](../specs/scenarios/README.md)) implements the same contract as
a stand-in and is what Phase 1 orchestration work exercises instead of a real adapter.

## Rules for Working in `adapters/`

1. **An adapter implements the contract; it never redefines it.** If a provider needs behavior
   the adapter contract doesn't support, that's a signal to propose a contract change in
   [specs/providers/](../specs/providers/README.md) (with an ADR/RFC as required), not to bolt
   on adapter-specific behavior that other adapters can't express.
2. **Capabilities are declared honestly.** If a provider doesn't support refunds, the adapter
   declares that via capability discovery — it must not fake support or silently no-op.
3. **Provider vocabulary stays inside the adapter.** Provider status codes, error codes, and
   field names are translated to canonical equivalents at the adapter boundary; they must never
   leak through to orchestration, events, or the canonical contract (see
   [ARCHITECTURE.md §5](../ARCHITECTURE.md#5-provider-boundary)).
4. **Must pass the provider conformance suite.** Per
   [CONTRIBUTING.md](../CONTRIBUTING.md#adding-a-provider), an adapter is only accepted once it
   passes [conformance/provider/](../conformance/provider/README.md) — compiling is not
   sufficient.
5. **Never commit real credentials or reach real provider endpoints from anything in this repo**
   (tests, examples, docs) — see [AGENTS.md §10](../AGENTS.md#10-security-restrictions) and
   [SECURITY.md](../SECURITY.md). Adapter tests exercise the simulator or a provider sandbox,
   never production.

## What Does Not Belong Here

- The provider capability contract itself — that's
  [specs/providers/](../specs/providers/README.md).
- The simulator, which is a peer of adapters behind the Provider Interface but is not
  provider-specific — see [specs/scenarios/](../specs/scenarios/README.md).
- Per-language client libraries for application developers — that's
  [sdks/](../sdks/README.md).
