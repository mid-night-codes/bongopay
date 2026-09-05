# contracts/openapi/ — REST API Contract

The OpenAPI document(s) derived from [specs/payments/](../../specs/payments/README.md) and
[specs/errors/](../../specs/errors/README.md). See [contracts/README.md](../README.md) for the
rules governing everything under `contracts/` (never author a change here before the spec,
generated files are never hand-edited once a pipeline exists, no provider-specific vocabulary).

## Status

**Draft v0.1 — Phase 1 ("Simulator Core"), "REST contract (first working implementation)".**
See [ROADMAP.md](../../ROADMAP.md) and [bongopay.yaml](bongopay.yaml): `POST /payments`
(initiate) and `GET /payments/{id}` (query) only, matching
[internal/payment/types.go](../../implementations/reference/internal/payment/types.go)
exactly.

This corrects an earlier version of this document, which said no OpenAPI document should exist
until [specs/payments/payment-contract.md](../../specs/payments/payment-contract.md)'s open
`TODO(spec)`/`TODO(RFC)` items (e.g. finalizing `PaymentMethod`) resolved. That turned out to be
inconsistent with how this project actually operates: the Go reference implementation was
already built directly against that same Draft spec, propagating its open TODOs forward as code
comments rather than waiting for them. `bongopay.yaml` does the same — its `PaymentMethodType`
enum and `CustomerReference` schema carry the identical `TODO(spec)` notes forward rather than
guessing an answer.

No HTTP server implements this contract yet — that's a separate, larger follow-up. No
simulator-specific callback-delivery endpoint is in this document either: each real provider has
its own webhook payload shape (see
[specs/providers/adapter-contract.md](../../specs/providers/adapter-contract.md)'s
`parseCallback`), so there's no single canonical wire shape to standardize here — a callback
endpoint, if one is added for local testing, belongs to
[implementations/reference/](../../implementations/reference/README.md) specifically, not this
canonical contract.
