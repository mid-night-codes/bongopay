#!/usr/bin/env bash
# make test-conformance — runs the conformance suite as it comes online.
# See conformance/README.md: today conformance cases are specification
# documents, not yet executable — this script validates that structure.
set -euo pipefail
cd "$(dirname "$0")/.."

fail=0
for dir in conformance/*/; do
  name="$(basename "$dir")"
  [ "$name" = "README.md" ] && continue
  if [ ! -f "${dir}README.md" ]; then
    echo "  [FAIL] conformance/${name}/ is missing a README.md"
    fail=1
  else
    echo "  [ok]   conformance/${name}/README.md"
  fi
done

if [ "$fail" -ne 0 ]; then
  exit 1
fi

echo "No executable conformance harness exists yet (see ROADMAP.md Phase 1/3)."
echo "Conformance case definitions validated structurally."
