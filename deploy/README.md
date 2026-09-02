# deploy/ — Local Development Deployment Only

This directory holds Docker/Compose configuration for running BongoPay components locally
during development — it is explicitly **not** a production deployment story. Per
[AGENTS.md §3](../AGENTS.md#3-directory-responsibilities), this directory has no dedicated
context chain beyond this README: it's operational configuration, not architecture.

## Status

**Empty — Phase 2 ("Developer Tooling"), following
[implementations/reference/](../implementations/README.md).** See [ROADMAP.md](../ROADMAP.md).
There is nothing to containerize yet.

## Layout

| Directory | Covers |
|---|---|
| [docker/](docker/) | Dockerfile(s) for local components (simulator, reference implementation) once they exist |
| [compose/](compose/) | Docker Compose files to run those components together locally |

## Rules for Working in `deploy/`

1. **Local development only.** Per
   [ROADMAP.md](../ROADMAP.md) §"What NOT to Build Yet", this is explicitly not a Kubernetes
   operator, service mesh, or production deployment topology — see
   [docs/architecture/non-goals.md](../docs/architecture/non-goals.md).
2. **No real credentials, ever**, including in example `.env` files or Compose environment
   blocks — see [AGENTS.md §10](../AGENTS.md#10-security-restrictions). Use obviously-fake
   placeholders.
3. **Reflects `implementations/` and `adapters/`, never leads them.** Containerizing something
   that doesn't exist yet, or inventing a deployment shape ahead of the reference
   implementation's own choices, isn't useful scaffolding.

## What Does Not Belong Here

- Any production infrastructure-as-code — out of scope for the foreseeable future (see
  [docs/architecture/non-goals.md](../docs/architecture/non-goals.md)).
