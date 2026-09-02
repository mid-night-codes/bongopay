#!/usr/bin/env bash
# make validate (part 3/3) — OpenAPI / AsyncAPI YAML syntax validation.
set -euo pipefail
cd "$(dirname "$0")/.."

if ! command -v python3 >/dev/null 2>&1; then
  echo "python3 not found — skipping OpenAPI/AsyncAPI syntax validation."
  exit 0
fi

python3 - <<'PYEOF'
import sys
import os

try:
    import yaml  # type: ignore
except ImportError:
    print("PyYAML not installed — falling back to a minimal presence check only.")
    print("Install with: pip install pyyaml   (optional; see scripts/setup.sh)")
    yaml = None

fail = False
count = 0
for base in ("contracts/openapi", "contracts/asyncapi"):
    if not os.path.isdir(base):
        continue
    for root, _dirs, files in os.walk(base):
        for f in files:
            if f.endswith((".yaml", ".yml")):
                path = os.path.join(root, f)
                count += 1
                if yaml is None:
                    continue
                try:
                    with open(path) as fh:
                        yaml.safe_load(fh)
                    print(f"  [ok]   {path}")
                except Exception as e:  # noqa: BLE001
                    print(f"  [FAIL] {path}: {e}")
                    fail = True

print(f"Checked {count} contract YAML file(s).")
sys.exit(1 if fail else 0)
PYEOF
