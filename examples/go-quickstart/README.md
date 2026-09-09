# examples/go-quickstart/

A small Go program using [sdks/go](../../sdks/go/README.md) to demonstrate the flow that's
fully supported end-to-end via the public HTTP API today: initiate a payment against the
`SIMULATOR` (success and failure scenarios), read it back by ID, and show that replaying the
same idempotency key returns the same payment rather than creating a new one.

## What This Doesn't Show, and Why

[examples/README.md](../README.md) suggests "initiate a payment and handle its webhook" as a
model example. This isn't that example, deliberately: `Simulator.Initiate` resolves
`success`/`failure` synchronously in one call today, rather than submitting to `PENDING` and
waiting for a separate callback the way a real async provider would. Building a "webhook"
example would mean either faking a gap that doesn't really exist in the current API, or reaching
into `POST /simulator/callbacks` against a payment that never actually sits in `PENDING` long
enough to matter. Whether `Initiate` should become callback-driven is an open design question
recorded in [ROADMAP.md](../../ROADMAP.md) — a real webhook example belongs here once that's
resolved, not invented ahead of it.

## Running It

Start the reference server first (see
[implementations/reference/README.md](../../implementations/reference/README.md) for
alternatives, including Docker):

```bash
cd implementations/reference
go run ./cmd/server &
```

Then, from this directory:

```bash
cd examples/go-quickstart
go run .
```

Use `-server` if the server isn't on the default `localhost:8080`:

```bash
go run . -server http://localhost:18080
```
