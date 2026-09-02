#!/usr/bin/env bash
# make generate — regenerates generated artifacts from source contracts.
# No generation pipeline exists yet (see AGENTS.md #12, ROADMAP.md Phase 2).
# This script exists so `make generate` is always a valid, safe no-op today
# and gets real logic added in one place later.
set -euo pipefail
cd "$(dirname "$0")/.."

echo "No generation pipeline is configured yet."
echo "When one exists, this script will regenerate artifacts declared in AGENTS.md #12"
echo "(currently: none — contracts/ has no downstream generated code yet)."
exit 0
