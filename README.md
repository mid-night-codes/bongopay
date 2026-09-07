# BongoPay

[![CI](https://github.com/mid-night-codes/bongopay/actions/workflows/ci.yml/badge.svg)](https://github.com/mid-night-codes/bongopay/actions/workflows/ci.yml)
[![Commit messages](https://github.com/mid-night-codes/bongopay/actions/workflows/commitlint.yml/badge.svg)](https://github.com/mid-night-codes/bongopay/actions/workflows/commitlint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mid-night-codes/bongopay/implementations/reference)](https://goreportcard.com/report/github.com/mid-night-codes/bongopay/implementations/reference)
[![Go version](https://img.shields.io/github/go-mod/go-version/mid-night-codes/bongopay?filename=implementations%2Freference%2Fgo.mod&label=go)](implementations/reference/go.mod)
[![License](https://img.shields.io/github/license/mid-night-codes/bongopay)](LICENSE)

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

**Early-stage / Phase 0 ("Foundation") and Phase 1 ("Simulator Core") complete, Phase 2
("Developer Tooling") starting.** See [ROADMAP.md](ROADMAP.md). This repository currently
establishes:

- Repository structure, governance, and contribution workflow
- Specification-first architecture (specs → contracts → conformance → implementations)
- Specifications for the payment contract, state machine, provider adapter model, and scenario
  format; the canonical error model and event model are still placeholders (see
  [specs/errors/](specs/errors/README.md), [specs/events/](specs/events/README.md))
- A working Go reference implementation
  ([implementations/reference/](implementations/reference/README.md)) — payment lifecycle,
  the `SIMULATOR` provider (all six scenario outcomes), webhook/callback verification, and a
  REST API you can actually run and `curl`
- CI scaffolding (fast validation, Go build/test, maintainer-approval gating for
  non-contributor PRs, Conventional Commits enforcement) and an AI-agent-friendly contribution
  environment

**BongoPay does not yet:**

- Process real payments or move real money
- Integrate with any real payment provider (`adapters/` is still empty — only the `SIMULATOR`
  is implemented)
- Ship a production-ready reference implementation, SDK, or CLI (the reference implementation
  is in-memory only, with no authentication — see
  [implementations/reference/README.md](implementations/reference/README.md))
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

Validate specs and contracts (no heavy language toolchain required for this part):

```bash
git clone <repository-url>
cd bongopay

make setup      # install local validation tooling
make validate   # validate specs, schemas, and contracts
make test       # run conformance/unit tests, including the Go reference implementation's
```

Run the actual reference implementation (requires Go — see
[implementations/reference/README.md](implementations/reference/README.md)):

```bash
cd implementations/reference
go run ./cmd/server &
curl -X POST localhost:8080/payments \
  -d '{"provider":{"id":"SIMULATOR"},"amount":{"value":5000,"currency":{"code":"TZS"}},"customerReference":{},"idempotencyKey":"demo-1"}'
```

See [docs/development/](docs/development/README.md) for the full local development guide.

## Example API

The shape below is real and runnable today against the `SIMULATOR` provider — see
[implementations/reference/README.md](implementations/reference/README.md) for how to start the
server and `curl` it. Real providers like `MPESA` are still illustrative only: `adapters/` is
empty, so a request naming one is accepted by the contract's shape but has nothing to actually
process it yet. See [specs/payments/](specs/payments/README.md) for the authoritative,
evolving specification this shape derives from.

```json
{
  "provider": { "id": "MPESA" },
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
