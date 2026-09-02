# contracts/examples/ — Example Payloads

Example payloads validated against [contracts/json-schema/](../json-schema/README.md), used to
sanity-check schemas and to illustrate the contract in generated documentation. See
[contracts/README.md](../README.md) for the rules governing everything under `contracts/`.

## Status

**Empty — Phase 1, following [contracts/json-schema/](../json-schema/README.md).** There is
nothing to validate examples against yet. The illustrative payload in the root
[README.md](../../README.md) "Example API" section is a preview of the eventual shape here, not
a substitute for it — it is explicitly marked not-yet-implemented.

## Rules Specific to This Directory

- Every example here must validate against its corresponding schema in
  [contracts/json-schema/](../json-schema/README.md) — `make validate` is expected to enforce
  this once both exist (see
  [scripts/validate-schemas.sh](../../scripts/validate-schemas.sh)).
- Use only obviously-fake data (test MSISDNs, placeholder keys) per
  [AGENTS.md §10](../../AGENTS.md#10-security-restrictions).
