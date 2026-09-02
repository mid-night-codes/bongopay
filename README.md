# BongoPay

**A language-agnostic payment orchestration and simulation platform.**

BongoPay defines a common, provider-neutral contract for initiating payments, tracking their
lifecycle, handling callbacks and webhooks, and simulating provider behavior for testing —
so applications can integrate payments once and swap provider implementations without
rewriting business logic.

> Integrate payments once, simulate multiple providers, and switch provider implementations
> without changing core application business logic.

---

## Why BongoPay?

Payment integration work is repeated, provider-specific, and hard to test. Every mobile money
or card provider has its own status codes, callback shapes, and quirks, and most projects end
up hard-coding against one provider's model. BongoPay separates **what a payment is** (a stable,
canonical contract) from **how a provider implements it** (an adapter), and provides a
**simulator** so the canonical contract can be exercised — success, failure, timeouts, duplicate
callbacks, chaos scenarios — without touching a real provider or moving real money.

## Current Project Status

**Early-stage / Phase 0 — Foundation complete, Phase 1 ("Simulator Core") starting.** See
[ROADMAP.md](ROADMAP.md). This repository currently establishes:

- Repository structure, governance, and contribution workflow
- Specification-first architecture (specs → contracts → conformance → implementations)
- Placeholder specifications for the payment contract, state machine, provider adapter
  model, scenario format, error model, and event model
- CI scaffolding and an AI-agent-friendly contribution environment

**BongoPay does not yet:**

- Process real payments or move real money
- Integrate with any real payment provider
- Ship a production-ready reference implementation, SDK, or CLI
- Guarantee stability of any contract (everything is pre-1.0 and may change; changes are
  tracked via [ADRs](adr/) and [RFCs](rfcs/))

See [ROADMAP.md](ROADMAP.md) for what comes next and [Non-Goals](docs/architecture/non-goals.md)
for what is explicitly out of scope right now.

## Core Principles

1. **Contract-first, not language-first.** Specifications and contracts are the source of
   truth. Implementations conform to them — they don't define them. See [ARCHITECTURE.md](ARCHITECTURE.md).
2. **Language, framework, and provider neutrality.** The core is not tied to any programming
   language, database, message broker, or Mobile Network Operator.
3. **Canonical domain, provider-specific adapters.** Provider quirks stay in adapters.
   Provider-specific statuses never leak into the canonical payment state machine.
4. **Simulation before integration.** The simulator and scenario system let contributors and
   AI agents exercise realistic payment behavior — including failure modes — without a real
   provider.
5. **Conformance over trust.** An implementation or adapter is only considered correct if it
   passes the shared conformance suite, not because it compiles or "looks right."
6. **Small, reviewable changes.** Humans and AI agents are expected to make minimal, well-scoped
   changes backed by tests and documentation, with ADRs/RFCs for anything architectural.

## Architecture Overview

```text
Specifications (specs/)
      ↓
Contracts (contracts/: OpenAPI, AsyncAPI, JSON Schema)
      ↓
Conformance Tests (conformance/)
      ↓
Implementations (implementations/reference, adapters/, sdks/)
```

See [ARCHITECTURE.md](ARCHITECTURE.md) for the full architecture, boundaries, and diagrams.

## Quick Start

> The reference implementation does not exist yet (Phase 1). Today, "running BongoPay" means
> validating specifications and contracts.

```bash
git clone <repository-url>
cd bongopay

make setup      # install local validation tooling (no heavy language toolchains required)
make validate   # validate specs, schemas, and contracts
make test       # run conformance and unit tests (as they come online)
```

See [docs/development/](docs/development/README.md) for the full local development guide.

## Example API (Illustrative — Not Yet Implemented)

The shape below illustrates the target payment contract. It is **not** a working API yet —
see [specs/payments/](specs/payments/README.md) for the authoritative, evolving specification.

```json
{
  "provider": "MPESA",
  "amount": {
    "value": 50000,
    "currency": "TZS"
  },
  "customerReference": {
    "msisdn": "255700000005"
  },
  "providerOptions": {
    "mpesa": {
      "businessCode": "123456"
    }
  }
}
```

## Repository Structure

```text
bongopay/
├── specs/            # Language-neutral specifications (source of truth)
├── contracts/        # OpenAPI / AsyncAPI / JSON Schema derived from specs
├── conformance/       # Shared conformance test definitions
├── implementations/  # Reference implementation(s) that demonstrate the spec
├── adapters/         # Provider adapters (future — mpesa, airtel-money, etc.)
├── sdks/             # Client SDKs for multiple languages (future)
├── examples/         # Example applications and integrations
├── deploy/           # Docker / Compose for local development
├── docs/             # Architecture, concepts, and contributor documentation
├── adr/              # Architecture Decision Records
├── rfcs/             # RFCs for major/breaking changes
├── scripts/          # Validation and developer tooling scripts
└── tests/, tools/    # Shared test fixtures and developer tools
```

Every directory listed above has its own `README.md` explaining its responsibility — start
there before making changes in that area.

## Roadmap

See [ROADMAP.md](ROADMAP.md) for the staged plan: Foundation → Simulator Core → Developer
Tooling → Provider Ecosystem → Reliability Testing → Extended Payment Tooling.

## Contributing

Contributions from humans **and AI coding agents** are welcome. Start with
[CONTRIBUTING.md](CONTRIBUTING.md). If you are an AI agent, read [AGENTS.md](AGENTS.md) first —
it is a hard requirement, not a suggestion.

## Security

Please report vulnerabilities privately as described in [SECURITY.md](SECURITY.md). Never open
a public issue for a security report.

## License

BongoPay is licensed under the [Apache License 2.0](LICENSE). See that file, and the note in
[GOVERNANCE.md](GOVERNANCE.md#licensing-rationale), for why.
