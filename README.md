# golibjpeg

Go JPEG decoder — baseline JPEG, JPEG-LS, JPEG XT. No CGO dependency.

## Overview

`golibjpeg` is a Go library for decoding JPEG images with native precision (8‑bit and 16‑bit). It bundles a platform-specific shared library extracted at runtime via FFI (`ebitengine/purego`), avoiding the need for CGO.

Supported formats:
- **JPEG** (ISO 10918‑1, baseline)
- **JPEG‑LS** (ISO 14495)
- **JPEG XT** (ISO 18477, HDR)

## API

```go
func Decode(data []byte) (*Image, error)
func DecodeWithFormat(data []byte, format Format) (*Image, error)
func GetParameters(data []byte) (*Params, error)
```

`Format` constants: `FormatAuto`, `FormatJPEG`, `FormatJPEGLS`, `FormatJPEGXT`.

## How it works

- Go wraps a C++ shared library via `purego` (no CGO).
- The native library is embedded per platform with `//go:embed` and extracted to a temp directory on first use.
- C++ decode logic follows [pylibjpeg-libjpeg](https://github.com/pydicom/pylibjpeg-libjpeg) (`lib/interface/` + `thorfdbg/libjpeg`).
- Stripe‑based decoding processes 8 lines at a time, reducing memory pressure.
- Output pixels are in native precision (8‑bit or 16‑bit), planar‑interleaved.

## Project layout

```
golibjpeg.go          # public API
native/               # purego loader + embedded prebuilt libs
lib/
  libjpeg/            # submodule → thorfdbg/libjpeg
  interface/          # decode + streamhook (from pylibjpeg-libjpeg)
  capi/               # C ABI for purego
ref/pylibjpeg-libjpeg # read-only reference submodule
testdata/             # optional conformance JPEGs (see testdata/README.md)
```

## Installation

```bash
go get github.com/godicom-dev/golibjpeg
```

## Usage

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/godicom-dev/golibjpeg"
)

func main() {
	data, err := os.ReadFile("image.jpg")
	if err != nil {
		log.Fatal(err)
	}

	// Decode with auto-detection of format
	img, err := golibjpeg.Decode(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%dx%d, %d components, precision %d\n",
		img.Width, img.Height, img.Components, img.Precision)

	// img.Pixels is RGB bytes (or grayscale if source is grayscale)
	// Process pixels as needed...
	_ = img.Pixels
}
```

With explicit format:

```go
import "github.com/godicom-dev/golibjpeg"

// Force JPEG-LS decoding
img, err := golibjpeg.DecodeWithFormat(data, golibjpeg.FormatJPEGLS)
```

Read image parameters without decoding pixels:

```go
params, err := golibjpeg.GetParameters(data)
if err != nil {
	log.Fatal(err)
}
fmt.Printf("%dx%d, %d components, precision %d\n",
	params.Width, params.Height, params.Components, params.Precision)
```

## Platform support

| OS      | amd64 | arm64 |
|---------|-------|-------|
| Windows | ✓     |       |
| macOS   |       | ✓     |
| Linux   | ✓     | ✓     |

## Dependencies

- [ebitengine/purego](https://github.com/ebitengine/purego) – FFI without CGO
- [thorfdbg/libjpeg](https://github.com/thorfdbg/libjpeg) – C++ JPEG library (ISO 10918‑1 / 18477)

## Development

Native libraries in `native/libs/` are **not built locally by default**. They are produced by GitHub Actions (`build-libs.yml`) when `lib/**` changes on `main`, then auto-committed to the repository.

```bash
git clone --recurse-submodules https://github.com/godicom-dev/golibjpeg.git
cd golibjpeg
go test ./...
```

To rebuild native libraries on CI without changing `lib/`:

```bash
gh workflow run build-libs.yml
```

Optional local native build (requires CMake):

```bash
make build-native
```

### CI workflows

| Workflow | Trigger | Action |
|----------|---------|--------|
| `build-libs.yml` | push `lib/**` to `main`, or manual | Build 5 platform libs → commit to `native/libs/` |
| `build.yml` | push (except `lib/**`), PR, tags | `go test` using committed `native/libs/` |

Reference tests (`reference_test.go`) align with `ref/pylibjpeg-libjpeg/libjpeg/tests/`; install conformance JPEGs as described in `testdata/README.md`.

### Release workflow

1. Merge changes to `main` (if `lib/**` changed, wait for `build-libs` to commit updated `native/libs/`).
2. `build` workflow runs `go test` on all platforms using embedded libs.
3. Create and push a tag: `git tag v1.0.1 && git push origin v1.0.1`.
4. CI attaches the committed libraries from `native/libs/` to a GitHub Release.

## References

This Go port follows **[pylibjpeg-libjpeg](https://github.com/pydicom/pylibjpeg-libjpeg)** for native decode behaviour and tests, and **[pylibjpeg](https://github.com/pydicom/pylibjpeg)** for the overall plugin-style integration model used by pydicom.
