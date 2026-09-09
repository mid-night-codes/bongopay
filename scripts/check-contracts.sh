#!/usr/bin/env bash
# make check-contracts — verifies any file claiming to be generated declares
# its source, per the Generated Files policy in AGENTS.md and contracts/README.md.
set -euo pipefail
cd "$(dirname "$0")/.."

fail=0
while IFS= read -r -d '' file; do
  if grep -q "DO NOT EDIT MANUALLY" "$file" && ! grep -q "Generated from:" "$file"; then
    echo "  [FAIL] $file: has a 'DO NOT EDIT MANUALLY' header but no 'Generated from:' source reference"
    fail=1
  fi
done < <(find . -type f \( -name '*.md' -o -name '*.json' -o -name '*.yaml' -o -name '*.yml' -o -name '*.go' \) \
  -not -path './node_modules/*' -not -path './.git/*' -print0)

if [ "$fail" -ne 0 ]; then
  echo "check-contracts failed."
  exit 1
fi

echo "No generated-file policy violations found."
echo "(sdks/go/generated/client.gen.go, from contracts/openapi/bongopay.yaml via 'make generate',"
echo "is the first real entry in this policy — see AGENTS.md #12.)"
