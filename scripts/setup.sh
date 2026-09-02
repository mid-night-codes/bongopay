#!/usr/bin/env bash
# make setup — report on optional local tooling used by validate/lint/test.
# Nothing here is mandatory: BongoPay's basic validation must work with only
# POSIX shell + python3. See scripts/README.md.
set -euo pipefail

echo "BongoPay local setup check"
echo "---------------------------"

check() {
  local name="$1"
  if command -v "$name" >/dev/null 2>&1; then
    echo "  [ok]      $name"
  else
    echo "  [missing] $name (optional — see below)"
  fi
}

check python3
check node
check npx
check git

cat <<'EOF'

Notes:
  - python3 is used for JSON/YAML/OpenAPI structural checks in scripts/validate-*.sh.
  - node/npx are optional and only used if you opt into richer markdown/YAML linting.
  - No language-specific runtime is required to validate specs/ or contracts/ in Phase 0.

Next: run `make validate` and `make test`.
EOF
