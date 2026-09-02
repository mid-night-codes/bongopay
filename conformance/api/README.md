# conformance/api/ — API Contract Conformance Cases

Conformance cases for the REST/API contract defined in
[contracts/openapi/](../../contracts/openapi/), derived from
[specs/payments/](../../specs/payments/README.md) and [specs/errors/](../../specs/errors/README.md).

## Status

**Specification-only, Phase 0.** No case files exist yet. The first tracked requirement, called
out directly from [specs/payments/payment-contract.md](../../specs/payments/payment-contract.md)
("Idempotency"), will be the first case written here once the REST contract exists (Phase 1):

> Given an idempotency key, multiple identical payment requests must produce the same logical
> `Payment` — not a new one — and must not trigger duplicate provider-side initiation.

See [conformance/README.md](../README.md) for the rules governing everything in this directory
(cases are implementation-neutral, derived from spec not from implementation behavior, and are
the acceptance bar — not "compiles").
