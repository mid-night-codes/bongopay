# scripts/

Developer and CI tooling invoked by the root [Makefile](../Makefile). These scripts are the
implementation detail behind the stable `make <target>` interface described in
[AGENTS.md](../AGENTS.md) and [docs/development/README.md](../docs/development/README.md).

## Rules for scripts in this directory

- Must run with only POSIX shell + widely available tools (`python3`, `node` if present) —
  no heavy toolchain installs required for basic validation (see
  [ARCHITECTURE.md](../ARCHITECTURE.md) on language neutrality).
- Must degrade gracefully and explain itself when an optional tool isn't installed, rather than
  failing with an opaque error.
- Must be safe to run repeatedly and must not mutate source files as a side effect (except
  `generate.sh`, whose entire job is to regenerate declared generated artifacts).
- Should exit non-zero on real failure so CI (see [.github/workflows/](../.github/workflows/))
  can gate on them.

## Current scripts

| Script | Called by | Purpose |
|---|---|---|
| `setup.sh` | `make setup` | Checks for optional local tooling and reports what's available |
| `lint.sh` | `make lint` | Markdown/YAML/JSON linting |
| `validate-specs.sh` | `make validate` | Structural checks on `specs/` (required READMEs, headings) |
| `validate-schemas.sh` | `make validate` | JSON Schema syntax validation under `contracts/json-schema/` |
| `validate-contracts.sh` | `make validate` | OpenAPI/AsyncAPI syntax validation under `contracts/` |
| `check-contracts.sh` | `make check-contracts` | Verifies generated artifacts declare their source (see [contracts/README.md](../contracts/README.md)) |
| `test.sh` | `make test` | Runs unit/integration tests as they come online |
| `test-conformance.sh` | `make test-conformance` | Runs the conformance suite as it comes online |
| `docs.sh` | `make docs` | Checks internal markdown links resolve |
| `generate.sh` | `make generate` | Regenerates generated artifacts from source contracts |
| `clean.sh` | `make clean` | Removes local validation/build artifacts |

As of Phase 0, most of these are intentionally thin — they validate structure and syntax, not
business behavior, because there is no implementation yet to exercise. Do not add
implementation-specific logic here; language-specific tooling belongs behind these scripts, not
mixed into the Makefile itself.
