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

## Phase 1 — Simulator Core (current)

```text
Payment lifecycle (implemented against the canonical state machine)
Scenario specification (executable, not just documented)
Webhook simulation
Deterministic test scenarios
REST contract (first working implementation)
```

## Phase 2 — Developer Tooling

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
