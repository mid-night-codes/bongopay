# Provider Options Extension Mechanism

Status: **Draft**. Version: `0.1.0` (see
[specs/compatibility/README.md](../compatibility/README.md)).

This document defines the `providerOptions` extension mechanism referenced from
[specs/payments/payment-contract.md](../payments/payment-contract.md) — the only sanctioned way
for provider-specific data to travel through the canonical `PaymentRequest`/`Payment` shapes
without becoming a first-class canonical field.

## Why This Exists

Per [specs/payments/README.md](../payments/README.md): "A canonical field must make sense for
*any* provider. If a field only makes sense for one provider, it belongs in that provider's
`providerOptions` namespace instead." This document is what makes that rule concrete and
enforceable rather than a matter of taste.

## Shape

```text
ProviderOptions
  [providerNamespace: string]: object
```

- `providerNamespace` MUST match the `Provider.id` the options apply to (see
  [specs/payments/payment-contract.md](../payments/payment-contract.md) `Provider`), lowercased
  — e.g. `providerOptions.mpesa`, `providerOptions.airtel_money`.
- The object under each namespace is entirely provider-defined. BongoPay core and orchestration
  MUST NOT read, validate, or branch on its contents — only the owning adapter may.
- `SIMULATOR` may also have a `providerOptions.simulator` namespace, primarily to select a
  [scenario](../scenarios/README.md) for a given request.

## Example

```json
{
  "provider": "MPESA",
  "providerOptions": {
    "mpesa": {
      "businessCode": "123456"
    }
  }
}
```

## Rules

1. **Never required for unrelated providers.** A `PaymentRequest` targeting `MPESA` must never
   require anything under `providerOptions.airtel_money`, and vice versa.
2. **Never changes canonical meaning.** A `providerOptions` entry may add provider-specific
   behavior an adapter chooses to honor; it must never redefine what a canonical field (e.g.
   `amount`, `status`) means for that request.
3. **Schema-validated per namespace.** Once [contracts/json-schema/](../../contracts/json-schema/)
   exists for a given provider, its `providerOptions.<namespace>` shape should be validated
   against that provider's own schema — not against a single global "provider options" schema,
   since namespaces are independent.
4. **Additive within a namespace.** Changes to one provider's `providerOptions` namespace follow
   that provider's own adapter versioning (see [VERSIONING.md](../../VERSIONING.md)) and never
   require a canonical contract version bump on their own.
5. **Not a workaround for a missing canonical field.** If multiple providers would plausibly
   want the same option, that's a signal it belongs in the canonical contract instead (with an
   ADR/RFC as required) — not that every provider should reinvent it under its own namespace.

## Open TODOs

- TODO(spec): whether `providerOptions` namespaces need a registry (to catch typos like
  `mpesa` vs `m_pesa`) or whether adapter-level validation is sufficient.
