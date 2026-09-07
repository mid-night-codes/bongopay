# Roadmap

BongoPay is being built in stages, deliberately deferring real integrations until the
contract-first foundation is solid. This roadmap will evolve — significant changes to phase
scope should be reflected here alongside the RFC/ADR that drove them.

## Phase 0 — Foundation (complete)

```text
Repository structure
Governance
Specifications (initial placeholders)
Contribution system (human + AI agent)
AI-agent support (AGENTS.md, directory READMEs, PR/issue templates)
CI (fast validation)
Contracts (initial placeholders)
```

**Exit criteria:** a new contributor (human or AI agent) can clone the repo, run
`make setup && make validate && make test`, and understand what to do next without asking a
maintainer a clarifying question that documentation should have answered.

## Phase 1 — Simulator Core (complete)

```text
Payment lifecycle (implemented against the canonical state machine)
Scenario specification (executable, not just documented)
Webhook simulation
Deterministic test scenarios
REST contract (first working implementation)
```

**Exit criteria:** every item above has a working Go implementation in
[implementations/reference/](implementations/reference/README.md), covered by tests
(including `-race` concurrency tests where relevant), matching a written spec — see
[ADR 0002](adr/0002-reference-implementation-language-go.md) for the language choice and
[specs/scenarios/scenario-format.md](specs/scenarios/scenario-format.md) for the six scenario
outcomes, all of which now have real behavior: `success`/`failure`/`timeout` through
`Simulator.Initiate`, and `duplicate_callback`/`out_of_order`/`invalid_signature` through
`Simulator.HandleCallback` (exposed over HTTP as `POST /simulator/callbacks`).

**Left open, not blocking exit:** whether `Simulator.Initiate` should become callback-driven for
`success`/`failure` too — i.e. always submit to `PENDING` and require a separate callback to
resolve the outcome, matching how a real async provider behaves — instead of today's synchronous
one-call convenience. This is a real behavior change for existing callers (including the REST
contract's demo UX), not an additive increment, so it's recorded here as an open question for a
future ADR rather than decided as a side effect of closing this phase.

## Phase 2 — Developer Tooling (current)

```text
Docker image
CLI
Testcontainers support
SDK generation
Example applications
```

## Phase 3 — Provider Ecosystem

```text
Provider contract (stabilized)
Provider conformance suite (executable)
Sample adapters
Community adapters
```

## Phase 4 — Reliability Testing

```text
Chaos simulation
Callback replay
Latency simulation
Duplicate event simulation
Failure injection
```

## Phase 5 — Extended Payment Tooling

```text
Refunds
Reversals
Reconciliation fixtures
Observability
Advanced test automation
```

## What NOT to Build Yet

The following are explicit non-goals for the current and near-term phases. They may become
in-scope far later, but building them now would compromise the contract-first foundation this
roadmap depends on:

- A real payment switch
- Real money movement
- Settlement infrastructure
- PCI card processing
- Large web dashboards
- Kubernetes operators
- Complex microservice architecture
- Multiple real provider integrations
- Production persistence architecture
- Distributed transaction infrastructure
- A service mesh
- An elaborate plugin runtime

See also [docs/architecture/non-goals.md](docs/architecture/non-goals.md) for the rationale
behind these exclusions.

## How This Roadmap Is Maintained

Phase scope changes, additions, or reprioritizations should be proposed via issue or RFC (for
anything touching architecture) and merged as a normal PR to this file, reviewed per
[GOVERNANCE.md](GOVERNANCE.md).
