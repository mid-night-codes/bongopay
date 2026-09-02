#!/usr/bin/env bash
# make validate (part 1/3) — structural checks on specs/.
# Every spec subdirectory must have a README.md explaining its responsibility
# (see AGENTS.md and specs/README.md).
set -euo pipefail
cd "$(dirname "$0")/.."

fail=0

echo "Validating specs/ structure..."
for dir in specs/*/; do
  name="$(basename "$dir")"
  if [ ! -f "${dir}README.md" ]; then
    echo "  [FAIL] specs/${name}/ is missing a README.md"
    fail=1
  else
    echo "  [ok]   specs/${name}/README.md"
  fi
done

if [ "$fail" -ne 0 ]; then
  echo "specs/ validation failed."
  exit 1
fi

echo "specs/ validation passed."
