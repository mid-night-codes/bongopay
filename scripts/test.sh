#!/usr/bin/env bash
# make test — runs unit/integration tests as implementations come online.
set -euo pipefail
cd "$(dirname "$0")/.."

fail=0
ran=0

while IFS= read -r -d '' gomod; do
  dir="$(dirname "$gomod")"
  echo "Running 'go test ./...' in ${dir}..."
  if command -v go >/dev/null 2>&1; then
    if ! (cd "$dir" && go test ./...); then
      fail=1
    fi
    ran=1
  else
    echo "  go not found — skipping. Install Go to run tests here (see ${dir}/README.md)."
  fi
done < <(find implementations adapters sdks -type f -name 'go.mod' -print0 2>/dev/null)

if [ "$ran" -eq 0 ]; then
  echo "No implementation tests found yet (expected pre-Phase 1 — see ROADMAP.md)."
  echo "make test currently only guarantees: specs/contracts are structurally valid."
  echo "Run 'make validate' for that check."
  exit 0
fi

if [ "$fail" -ne 0 ]; then
  echo "make test failed."
  exit 1
fi

echo "make test passed."
