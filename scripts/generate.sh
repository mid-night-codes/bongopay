#!/usr/bin/env bash
# make generate — regenerates generated artifacts from source contracts. See AGENTS.md #12's
# Generated Files table for the source -> generated-location -> validation-command contract
# this script implements.
set -euo pipefail
cd "$(dirname "$0")/.."

if ! command -v go >/dev/null 2>&1; then
  echo "go not found — skipping SDK generation. Install Go to run this (see scripts/setup.sh)."
  exit 0
fi

OAPI_CODEGEN_VERSION="v2.4.1"
SPEC="contracts/openapi/bongopay.yaml"
CONFIG="sdks/go/oapi-codegen-config.yaml"
OUT="sdks/go/generated/client.gen.go"

echo "Generating ${OUT} from ${SPEC}..."

( cd sdks/go && GOFLAGS=-mod=mod go run "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@${OAPI_CODEGEN_VERSION}" \
    -config "$(basename "$CONFIG")" "../../${SPEC}" )

# Prepend this repo's standard generated-file header (AGENTS.md #11) ahead of oapi-codegen's
# own header, so make check-contracts's "DO NOT EDIT MANUALLY" + "Generated from:" check covers
# this file using the same convention as every other generated artifact, present or future.
TMP="$(mktemp)"
trap 'rm -f "$TMP"' EXIT

{
  echo "// DO NOT EDIT MANUALLY."
  echo "// Generated from: ${SPEC}"
  echo "// Generate with: make generate"
  echo "// Validate with: make check-contracts"
  echo "//"
  cat "$OUT"
} > "$TMP"
mv "$TMP" "$OUT"
trap - EXIT

gofmt -w "$OUT"

echo "Generated ${OUT}."
