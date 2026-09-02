# contracts/json-schema/ — JSON Schema Derivations

JSON Schema documents derived one-to-one from canonical `specs/` concepts. See
[contracts/README.md](../README.md) for the rules governing everything under `contracts/`.

## Layout

| Directory | Derived from |
|---|---|
| [payments/](payments/) | `Payment`, `PaymentRequest`, `Money`, etc. in [specs/payments/payment-contract.md](../../specs/payments/payment-contract.md) |
| [errors/](errors/) | The canonical error shape in [specs/errors/README.md](../../specs/errors/README.md) |
| [events/](events/) | The event envelope/type shapes in [specs/events/README.md](../../specs/events/README.md) |

## Status

**Empty — Phase 1 ("Simulator Core").** See [ROADMAP.md](../../ROADMAP.md). No schema file
exists in any subdirectory yet: `payments/` is blocked on
[specs/payments/payment-contract.md](../../specs/payments/payment-contract.md)'s open
`TODO(spec)` items, and `errors/`/`events/` are blocked on their respective detail documents
(`specs/errors/error-model.md`, `specs/events/event-envelope.md`) not existing yet.

Each subdirectory's schema files, once written, are validated for syntax by
`make validate` (see [scripts/validate-schemas.sh](../../scripts/validate-schemas.sh)) — that
check currently reports 0 files, which is expected at this phase, not a failure.
