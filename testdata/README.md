# testdata

Conformance JPEG files used by [pylibjpeg-libjpeg](https://github.com/pydicom/pylibjpeg-libjpeg) tests.

## Layout (from pylibjpeg-libjpeg)

```
testdata/
├── 10918/p1/A1.JPG    # ISO 10918 compliance stream
└── 14495/             # JPEG-LS compliance streams
```

## Obtaining files

Test images ship with the `pylibjpeg-libjpeg` Python package (`libjpeg.data.JPEG_DIRECTORY`).
To populate this directory locally:

```bash
pip install pylibjpeg-libjpeg
python -c "import libjpeg.data as d; print(d.JPEG_DIRECTORY)"
# copy 10918/ and 14495/ from that path into testdata/
```

Go tests in `reference_test.go` skip when these files are absent.
