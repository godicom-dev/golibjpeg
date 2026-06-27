#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

if ! python3 -c "import ljdata" 2>/dev/null; then
  echo "Installing pylibjpeg-data..."
  pip3 install "git+https://github.com/pydicom/pylibjpeg-data"
fi

python3 - <<'PY'
import pathlib
import shutil

import ljdata

dst = pathlib.Path("testdata")
dst.mkdir(exist_ok=True)
for name in ("10918", "14495"):
    src = ljdata.JPEG_DIRECTORY / name
    target = dst / name
    if target.exists():
        shutil.rmtree(target)
    shutil.copytree(src, target)
    print(f"copied {src} -> {target}")
PY

echo "testdata ready under testdata/10918 and testdata/14495"
