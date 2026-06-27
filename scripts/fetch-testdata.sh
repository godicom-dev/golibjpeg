#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

if ! python3 -c "import ljdata" 2>/dev/null; then
  echo "Installing pylibjpeg-data..."
  pip3 install "git+https://github.com/pydicom/pylibjpeg-data"
fi
if ! python3 -c "import pydicom" 2>/dev/null; then
  echo "Installing pydicom..."
  pip3 install pydicom
fi

python3 - <<'PY'
import json
import pathlib
import shutil

import ljdata
from pydicom.encaps import generate_frames

dst = pathlib.Path("testdata")
dst.mkdir(exist_ok=True)
for name in ("10918", "14495"):
    src = ljdata.JPEG_DIRECTORY / name
    target = dst / name
    if target.exists():
        shutil.rmtree(target)
    shutil.copytree(src, target)
    print(f"copied {src} -> {target}")

# REF_DCM from pylibjpeg-libjpeg test_parameters.py / test_decode.py
REF = {
    "1.2.840.10008.1.2.4.50": [
        "JPEGBaseline_1s_1f_u_08_08.dcm",
        "SC_rgb_dcmtk_+eb+cy+np.dcm",
        "color3d_jpeg_baseline.dcm",
        "SC_rgb_dcmtk_+eb+cr.dcm",
        "SC_rgb_dcmtk_+eb+cy+n1.dcm",
        "SC_rgb_dcmtk_+eb+cy+s4.dcm",
    ],
    "1.2.840.10008.1.2.4.51": [
        "RG2_JPLY_fixed.dcm",
        "JPEGExtended_1s_1f_u_16_12.dcm",
        "JPEGExtended_3s_1f_u_08_08.dcm",
    ],
    "1.2.840.10008.1.2.4.57": [
        "JPEGLossless_1s_1f_u_16_12.dcm",
    ],
    "1.2.840.10008.1.2.4.70": [
        "JPEG-LL.dcm",
        "JPEGLosslessP14SV1_1s_1f_u_08_08.dcm",
        "JPEGLosslessP14SV1_1s_1f_u_16_16.dcm",
        "MG1_JPLL.dcm",
        "RG1_JPLL.dcm",
        "RG2_JPLL.dcm",
        "SC_rgb_jpeg_gdcm.dcm",
    ],
    "1.2.840.10008.1.2.4.80": [
        "emri_small_jpeg_ls_lossless.dcm",
        "MR_small_jpeg_ls_lossless.dcm",
        "RG1_JLSL.dcm",
        "RG2_JLSL.dcm",
    ],
    "1.2.840.10008.1.2.4.81": [
        "CT1_JLSN.dcm",
        "MG1_JLSN.dcm",
        "RG1_JLSN.dcm",
        "RG2_JLSN.dcm",
    ],
}

dcm_root = dst / "dcm"
for uid, names in REF.items():
    index = ljdata.get_indexed_datasets(uid)
    out = dcm_root / uid
    out.mkdir(parents=True, exist_ok=True)
    for fname in names:
        ds = index[fname]["ds"]
        frame = next(generate_frames(ds.PixelData, number_of_frames=1))
        frame_path = out / (fname + ".frame")
        meta_path = out / (fname + ".json")
        frame_path.write_bytes(frame)
        meta = {
            "rows": int(ds.Rows),
            "columns": int(ds.Columns),
            "samples_per_pixel": int(ds.SamplesPerPixel),
            "bits_allocated": int(ds.BitsAllocated),
            "bits_stored": int(getattr(ds, "BitsStored", ds.BitsAllocated)),
            "pixel_representation": int(getattr(ds, "PixelRepresentation", 0)),
            "photometric_interpretation": str(ds.PhotometricInterpretation),
        }
        meta_path.write_text(json.dumps(meta, indent=2))
        print(f"exported {uid}/{fname}")
PY

echo "testdata ready under testdata/{10918,14495,dcm}"
