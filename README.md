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
- Stripe‑based decoding processes 8 lines at a time, reducing memory pressure.
- Output pixels are in native precision (8‑bit or 16‑bit), planar‑interleaved.

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
| macOS   | ✓     | ✓     |
| Linux   | ✓     | ✓     |

## Dependencies

- [ebitengine/purego](https://github.com/ebitengine/purego) – FFI without CGO
- [thorfdbg/libjpeg](https://github.com/thorfdbg/libjpeg) – C++ JPEG library (ISO 10918‑1 / 18477)

## Development

```bash
# Build native library for current platform
make build-native

# Run tests
make test
```

Prebuilt native libraries (`.dll`/`.so`/`.dylib`) are stored in `native/libs/` and committed to git. On push to `main` affecting `_csrc/`, a CI workflow rebuilds all platform libraries and commits the result automatically.

### Release workflow

1. Make changes, push to `main`.
2. CI rebuilds and commits native libraries to `native/libs/`.
3. Create and push a tag: `git tag v1.0.1 && git push origin v1.0.1`.
4. CI runs tests on all platforms, then attaches the built libraries to a GitHub Release.

## References

This Go port takes inspiration from the design and API of **[pylibjpeg](https://github.com/pydicom/pylibjpeg)** — a Python JPEG decoder that also wraps native C libraries via a CPython extension. The concept of a lightweight, cross‑platform JPEG decoder with a minimal public API and format‑aware colour transforms follows pylibjpeg's approach.
