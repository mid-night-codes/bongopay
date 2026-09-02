#!/usr/bin/env bash
# make lint — markdown/YAML/JSON linting. Uses npx-based linters if available,
# otherwise falls back to the syntax checks in validate-*.sh so `make lint`
# never hard-fails just because a contributor hasn't installed Node tooling.
set -euo pipefail
cd "$(dirname "$0")/.."

if command -v npx >/dev/null 2>&1; then
  echo "Running markdownlint-cli2 (via npx)..."
  npx --yes markdownlint-cli2 "**/*.md" "#node_modules" || {
    echo "markdownlint reported issues (see above)."
    exit 1
  }
else
  echo "npx not found — skipping markdownlint. Falling back to structural/syntax checks."
fi

bash scripts/validate-schemas.sh
bash scripts/validate-contracts.sh

echo "Lint checks complete."
