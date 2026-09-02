# specs/errors/ — Canonical Error Model

This directory defines BongoPay's canonical, transport-agnostic error taxonomy: the errors that
can arise during payment operations, independent of whether they surface over REST, an event, or
an SDK exception.

## Status

**Draft / Phase 0.** Directory established to unblock cross-references from
[specs/payments/](../payments/README.md) and
[ARCHITECTURE.md §12](../../ARCHITECTURE.md#12-security-boundaries); the error taxonomy itself is
not yet written.

## Contents

- `error-model.md` — **TODO(ADR): not yet written.** Will define the canonical error taxonomy,
  including error categories for validation failures, provider/adapter failures, and
  security-relevant failures (webhook signature verification, callback validation, idempotency
  violations).

## Rules Specific to This Directory

- Errors here are transport-agnostic: no HTTP status codes, gRPC codes, or provider error codes
  baked into the canonical taxonomy. A transport binding (e.g. REST in
  [contracts/openapi/](../../contracts/openapi/)) maps canonical errors onto its own wire format,
  never the reverse.
- Security-relevant failures — invalid webhook signatures, callback URL validation failures,
  idempotency-key conflicts — are modeled as canonical errors here, not left as
  implementation-specific exceptions (see
  [ARCHITECTURE.md §12](../../ARCHITECTURE.md#12-security-boundaries) and
  [SECURITY.md](../../SECURITY.md)).
- Adapters translate provider-specific error codes into a canonical error from this taxonomy;
  raw provider error codes/messages must never be exposed as the canonical error itself
  (they may be preserved as diagnostic detail, not as the categorization).
- Adding a new error category is additive; narrowing or changing the meaning of an existing one
  is a breaking change — see [specs/compatibility/README.md](../compatibility/README.md).
