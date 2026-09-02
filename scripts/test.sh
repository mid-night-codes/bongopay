#!/usr/bin/env bash
# make test — runs unit/integration tests as implementations come online.
# There is no implementation yet (Phase 0/1 — see ROADMAP.md), so this is a
# clearly-labeled no-op rather than a fabricated pass.
set -euo pipefail
cd "$(dirname "$0")/.."

if [ -z "$(find implementations sdks adapters -type f \( -name '*_test.*' -o -name '*.test.*' -o -path '*/tests/*' \) 2>/dev/null)" ]; then
  echo "No implementation tests found yet (expected pre-Phase 1 — see ROADMAP.md)."
  echo "make test currently only guarantees: specs/contracts are structurally valid."
  echo "Run 'make validate' for that check."
  exit 0
fi

echo "TODO: dispatch to per-implementation test runners once implementations/ exists."
exit 1
