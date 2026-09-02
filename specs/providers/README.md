# specs/providers/ — Provider Adapter Contract

This directory defines the language-neutral contract a provider adapter (or the simulator, which
is a peer of adapters) must satisfy, and the mechanism by which provider-specific behavior is
attached without leaking into the canonical domain. Read this before writing or reviewing any
adapter, and before adding a field to `specs/payments/` that "only matters for one provider."

## Status

**Draft / Phase 0.** Directory established to unblock cross-references from
[specs/payments/](../payments/README.md) and [specs/state-machines/](../state-machines/README.md);
the contract documents below are not yet written.

## Contents

- `adapter-contract.md` — **TODO(ADR): not yet written.** Will define the capability contract:
  payment initiation, status query, refund, reversal, callback parsing, and callback
  verification — and the capability-discovery mechanism adapters use to declare which of these
  they support, since none are mandatory for every provider (see
  [ARCHITECTURE.md §5](../../ARCHITECTURE.md#5-provider-boundary)).
- `extensions.md` — **TODO(ADR): not yet written.** Will define the `providerOptions.<provider>`
  namespaced extension mechanism referenced from
  [specs/payments/README.md](../payments/README.md) and
  [ARCHITECTURE.md §11](../../ARCHITECTURE.md#11-extension-points).

## Rules Specific to This Directory

- Adapters translate provider-specific status and error vocabulary into canonical BongoPay state
  ([specs/state-machines/](../state-machines/README.md)) and errors
  ([specs/errors/](../errors/README.md)). Provider vocabulary must never be exposed through the
  canonical interface.
- Capabilities are declared explicitly by each adapter (capability discovery). Orchestration must
  not assume a capability (e.g. refund) is universally supported.
- Provider-specific fields belong only in that provider's namespaced `providerOptions` extension,
  schema-validated, and never required for unrelated providers — never as a first-class canonical
  field in `specs/payments/`.
- The simulator ([specs/scenarios/](../scenarios/README.md)) implements this same contract as a
  peer of real adapters; it is not a special case.
