#!/usr/bin/env bash
# make docs — checks that relative markdown links within the repo resolve to
# real files, so documentation doesn't silently rot as directories move.
set -euo pipefail
cd "$(dirname "$0")/.."

if ! command -v python3 >/dev/null 2>&1; then
  echo "python3 not found — skipping doc link check."
  exit 0
fi

python3 - <<'PYEOF'
import re
import sys
import os

link_re = re.compile(r'\[[^\]]*\]\(([^)]+)\)')
broken = []

for root, dirs, files in os.walk('.'):
    dirs[:] = [d for d in dirs if d not in ('.git', 'node_modules')]
    for f in files:
        if not f.endswith('.md'):
            continue
        path = os.path.join(root, f)
        with open(path, encoding='utf-8') as fh:
            text = fh.read()
        for match in link_re.finditer(text):
            target = match.group(1).split('#')[0].strip()
            if not target or target.startswith(('http://', 'https://', 'mailto:')):
                continue
            resolved = os.path.normpath(os.path.join(root, target))
            if not os.path.exists(resolved):
                broken.append((path, target))

if broken:
    print(f"Found {len(broken)} broken relative link(s):")
    for path, target in broken:
        print(f"  [FAIL] {path} -> {target}")
    sys.exit(1)

print("All relative markdown links resolve.")
PYEOF
