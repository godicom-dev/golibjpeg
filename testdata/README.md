# testdata

Conformance JPEG files used by [pylibjpeg-libjpeg](https://github.com/pydicom/pylibjpeg-libjpeg) tests.

## Layout

```
testdata/
├── 10918/p1/A1.JPG    # ISO 10918 compliance streams
└── 14495/JLS/...      # JPEG-LS compliance streams
```

## Fetch (local or CI)

```bash
bash scripts/fetch-testdata.sh
```

This installs [pylibjpeg-data](https://github.com/pydicom/pylibjpeg-data) and copies `10918/` and `14495/` here.

Conformance image binaries are **not** committed to git; CI fetches them before `go test`.
