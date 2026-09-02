#!/usr/bin/env bash
# make validate (part 2/3) — JSON Schema / JSON syntax validation.
set -euo pipefail
cd "$(dirname "$0")/.."

if ! command -v python3 >/dev/null 2>&1; then
  echo "python3 not found — skipping JSON Schema syntax validation."
  echo "Install python3 to enable this check (see scripts/setup.sh)."
  exit 0
fi

fail=0
count=0

while IFS= read -r -d '' file; do
  count=$((count + 1))
  if ! python3 -c "import json,sys; json.load(open(sys.argv[1]))" "$file" 2>/tmp/bongopay-json-err; then
    echo "  [FAIL] $file"
    sed 's/^/         /' /tmp/bongopay-json-err
    fail=1
  fi
done < <(find contracts/json-schema contracts/examples -type f -name '*.json' -print0 2>/dev/null)

echo "Checked $count JSON file(s) under contracts/json-schema and contracts/examples."

if [ "$fail" -ne 0 ]; then
  echo "JSON Schema validation failed."
  exit 1
fi

echo "JSON Schema validation passed."
