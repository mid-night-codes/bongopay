# Payment Contract (Canonical Concepts)

Status: **Draft**. Version: `0.1.0` (see [specs/compatibility/README.md](../compatibility/README.md)
for what "draft" means for compatibility purposes).

This document defines the canonical concepts BongoPay uses to describe a payment, independent
of any provider, language, or transport. Machine-readable derivations live in
[contracts/json-schema/payments/](../../contracts/json-schema/payments/) and
[contracts/openapi/](../../contracts/openapi/) — this document is authoritative; those are
generated/maintained to match it.

## Design Notes

- Field lists below are **illustrative and non-exhaustive** — this is intentionally a v0.1
  skeleton, not a finished contract. Do not treat omission of a field as a decision; treat it as
  an open TODO unless marked otherwise.
- All monetary values use `Money` (integer minor units + `Currency`) — never floats.
- All identifiers are opaque strings from the canonical domain's point of view; format
  constraints (if any) are provider- or implementation-specific and don't belong here.

## Money

```text
Money
  value:    integer   # amount in the currency's minor unit (e.g., cents). Never a float.
  currency: Currency
```

## Currency

```text
Currency
  code: string   # ISO 4217 alphabetic code, e.g. "TZS", "KES", "USD"
```

TODO(spec): decide whether non-ISO-4217 pseudo-currencies (e.g., loyalty points) are ever
in scope. Default assumption until an RFC says otherwise: **no**.

## CustomerReference

Identifies the customer/payer in a provider-neutral way.

```text
CustomerReference
  msisdn?:    string    # phone number in E.164-like form, where applicable
  externalId?: string   # merchant's own customer identifier
  metadata?:   Metadata
```

TODO(spec): exactly one of `msisdn`/`externalId` required vs. both optional — open question,
depends on which payment methods are prioritized in Phase 1.

## ProviderReference

Identifies how a specific provider refers to a payment, once known.

```text
ProviderReference
  provider:            Provider
  providerTransactionId?: string   # provider's own identifier, once available
  providerRawStatus?:  string      # provider's own status string — for observability/debugging
                                    # only; MUST NOT be used as a canonical status (see
                                    # specs/state-machines/payment-lifecycle.md)
```

## PaymentMethod

```text
PaymentMethod
  type: enum   # e.g. "MOBILE_MONEY", "CARD", "BANK_TRANSFER" — TODO(spec): finalize enum values
```

## Provider

```text
Provider
  id:   string   # e.g. "MPESA", "AIRTEL_MONEY", "SIMULATOR" — see specs/providers/README.md
```

`SIMULATOR` is always a valid `Provider` — the simulator is a first-class peer of real
providers, not a special case (see [ARCHITECTURE.md §6](../../ARCHITECTURE.md#6-simulator-boundary)).

## Metadata

```text
Metadata
  # Free-form, opaque key/value pairs supplied by the caller and echoed back unmodified.
  # MUST NOT be interpreted by BongoPay core or by adapters. Purely a pass-through for
  # caller-side correlation.
  [string]: string
```

## ProviderOptions

```text
ProviderOptions
  [providerNamespace: string]: object
  # e.g. providerOptions.mpesa = { businessCode: "123456" }
  # See specs/providers/extensions.md for the rules governing this field.
```

## PaymentRequest

The input to initiate a payment.

```text
PaymentRequest
  provider:           Provider
  amount:             Money
  customerReference:  CustomerReference
  paymentMethod?:      PaymentMethod
  idempotencyKey:     string   # REQUIRED — see "Idempotency" below
  metadata?:           Metadata
  providerOptions?:    ProviderOptions
```

## Payment

The canonical, persistent representation of a payment, once created.

```text
Payment
  id:                 string          # BongoPay-assigned canonical identifier
  status:             PaymentStatus   # see specs/state-machines/payment-lifecycle.md
  amount:             Money
  provider:           Provider
  providerReference?: ProviderReference
  customerReference:  CustomerReference
  idempotencyKey:     string
  createdAt:          string   # ISO 8601 timestamp
  updatedAt:          string   # ISO 8601 timestamp
  metadata?:           Metadata
  providerOptions?:    ProviderOptions
```

## PaymentResult

Returned synchronously from an initiate/query/refund/reversal operation. Deliberately thin —
most state changes are observed via [events](../events/README.md), not synchronous results,
because provider payments are inherently asynchronous.

```text
PaymentResult
  payment: Payment
  # TODO(spec): do refund/reversal results need a distinct shape from a plain Payment?
  # Leaning yes (they likely reference the original Payment) — open for RFC once refunds
  # are scheduled (see ROADMAP.md Phase 5).
```

## PaymentStatus

See [specs/state-machines/payment-lifecycle.md](../state-machines/payment-lifecycle.md) for the
full canonical state machine. `PaymentStatus` is one of:

```text
CREATED | PENDING | SUCCESS | FAILED | EXPIRED | CANCELLED |
REVERSAL_PENDING | REVERSED | REFUND_PENDING | REFUNDED
```

## Idempotency

Every `PaymentRequest` carries a caller-supplied `idempotencyKey`. Conformance requirement
(tracked in [conformance/api/](../../conformance/api/README.md)):

> Given an idempotency key, multiple identical payment requests must produce the same logical
> `Payment` — not a new one — and must not trigger duplicate provider-side initiation.

TODO(spec): scope of "identical" (same key + same payload vs. same key regardless of payload,
and what happens on same key + different payload — likely a `409`-equivalent canonical error;
see [specs/errors/README.md](../errors/README.md)).

## Open TODOs

- TODO(RFC): finalize `PaymentMethod` enum before any real adapter work begins.
- TODO(ADR): whether `Payment.id` format is opaque-only or has a documented structure
  (e.g., ULID) — affects SDK ergonomics but must not leak provider assumptions.
- TODO(spec): pagination/listing contract for querying multiple payments — deferred to
  Phase 1 REST contract work.
