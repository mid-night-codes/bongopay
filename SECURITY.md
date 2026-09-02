# Security Policy

BongoPay is pre-1.0 and does not yet process real payments or handle real money (see
[ROADMAP.md](ROADMAP.md)). Even so, security discipline starts on day one, because contracts,
adapters, and simulator behavior established now will underpin real integrations later.

## Reporting a Vulnerability

**Do not open a public issue for a security vulnerability.**

Report it privately to the maintainers:

- Preferred: use GitHub's private vulnerability reporting ("Report a vulnerability" under the
  repository's Security tab), if enabled.
- Alternative: email the security contact listed in [MAINTAINERS.md](MAINTAINERS.md).

Please include: a description of the issue, steps to reproduce, affected files/specs/contracts,
and potential impact. We aim to acknowledge reports within a reasonable timeframe and will
coordinate disclosure with you.

## What Must Never Be Committed

Contributors (human or AI) must never commit:

- Real payment provider credentials or API keys
- Production endpoints containing embedded secrets
- Real customer information (names, phone numbers, account identifiers)
- Real transaction data
- Private keys or certificates
- Access tokens, session tokens, or webhook signing secrets

Use obviously-fake placeholders in examples and fixtures (e.g., MSISDNs in the
`255700000000`–`255700000009` test range, keys clearly prefixed `sk_test_fake_...`). If you
discover a real secret already committed, report it per the process above rather than deleting
it silently — deletion alone does not revoke a leaked credential from history.

## Threat Categories BongoPay Must Eventually Consider

These inform the design of specs, contracts, and conformance tests even before an
implementation exists that could be attacked directly:

```text
Webhook spoofing
Replay attacks
Duplicate processing
Credential leakage
Signature bypass
SSRF (e.g., via attacker-controlled callback URLs)
Unsafe callback URLs
Sensitive data in logs
Provider impersonation
Dependency attacks
Supply-chain attacks
```

Specification and contract changes that touch webhook/callback handling, idempotency, or
signature verification should explicitly note how they affect these categories (the pull
request template has a dedicated field for this).

## Supported Versions

Pre-1.0, there is no long-term-support version — security fixes land on the latest `main`.
Once BongoPay reaches 1.0, this section will be updated with a supported-version table per
[VERSIONING.md](VERSIONING.md).

## Dependency and Supply-Chain Hygiene

See [docs/development/dependency-policy.md](docs/development/dependency-policy.md) and
[.github/dependabot.yml](.github/dependabot.yml). Dependencies are kept minimal by policy, and
automated dependency and vulnerability scanning runs in CI (see
[docs/development/ci.md](docs/development/ci.md)).
