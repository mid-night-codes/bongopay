# Governance

BongoPay uses a lightweight, three-tier governance model appropriate for an early-stage
project. This will evolve — changes to this document should themselves go through the normal
PR process with core maintainer review.

## Roles

### Contributors

Anyone who opens an issue, PR, RFC, or ADR. No special access required. Contributors' PRs are
subject to review as described in [CONTRIBUTING.md](CONTRIBUTING.md).

### Maintainers

Contributors trusted with review and merge rights over specific areas (reflected in
[.github/CODEOWNERS](.github/CODEOWNERS)). Maintainers:

- Review and merge PRs in their owned areas.
- Triage issues.
- May approve ADRs for their area.
- Cannot unilaterally approve RFCs or changes to `specs/`, `contracts/`, `adr/`, or `.github/`
  governance/security files — those require Core Maintainer review per CODEOWNERS.

### Core Maintainers

Maintainers with repository-wide responsibility. Core Maintainers:

- Have final say on RFC acceptance/rejection.
- Approve changes to canonical contracts (`specs/`, `contracts/`), governance files, and
  security policy.
- Hold release authority (cutting and publishing releases per [VERSIONING.md](VERSIONING.md)).
- Are responsible for coordinating security response per [SECURITY.md](SECURITY.md).
- Nominate new Maintainers and Core Maintainers (see below).

Current maintainers are listed in [MAINTAINERS.md](MAINTAINERS.md).

## Review Expectations

- Every PR needs at least one Maintainer approval; PRs touching CODEOWNERS-protected paths
  need at least one Core Maintainer approval.
- Reviewers should focus on: contract stability, conformance, security impact, and whether the
  change matches [ARCHITECTURE.md](ARCHITECTURE.md) and any governing spec — not just style.
- Reviewers should request an ADR or RFC (per [AGENTS.md §8](AGENTS.md#8-when-adrs-and-rfcs-are-required))
  if a PR introduces architectural change without one.

## Release Authority

Core Maintainers hold release authority. Releases follow [VERSIONING.md](VERSIONING.md) and are
recorded in [CHANGELOG.md](CHANGELOG.md).

## Security Response Responsibility

Core Maintainers triage and coordinate response to reports made per [SECURITY.md](SECURITY.md),
including embargoed disclosure timing when appropriate.

## RFC Approval

RFCs (see [rfcs/README.md](rfcs/README.md)) are approved or rejected by Core Maintainer
consensus after a discussion period. Status transitions (`Draft → Discussion → Accepted /
Rejected → Implemented`, or `Withdrawn`) are recorded directly in the RFC document.

## Maintainer Nomination

Any Core Maintainer may nominate a Contributor as Maintainer (or a Maintainer as Core
Maintainer) based on sustained, high-quality contribution and review activity. Nomination
requires no objection from existing Core Maintainers within a reasonable discussion period.
There is no fixed term; roles are re-evaluated as the project grows.

## Licensing Rationale

BongoPay is licensed under the [Apache License 2.0](LICENSE). Apache 2.0 was chosen over MIT
because it includes an explicit patent grant and patent-retaliation clause, which matters for
payment/financial infrastructure where patent risk is non-trivial, while remaining permissive
enough for broad commercial adoption — important for a project that expects both community and
commercial provider adapters and SDKs.

## Amending This Document

Changes to governance require a PR reviewed and approved by Core Maintainers, following the
same [ADR](adr/README.md)/[RFC](rfcs/README.md) guidance as other architectural changes —
governance changes generally warrant an RFC given their project-wide impact.
