#!/usr/bin/env bash
# make clean — removes local validation/build artifacts. Never touches
# source files under specs/, contracts/, docs/, adr/, or rfcs/.
set -euo pipefail
cd "$(dirname "$0")/.."

rm -rf node_modules .cache /tmp/bongopay-json-err
echo "Cleaned local build/validation artifacts."
