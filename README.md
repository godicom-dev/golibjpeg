# golibjpeg

Go JPEG codec — baseline JPEG, JPEG-LS, JPEG XT decode; JPEG / JPEG-LS encode. No CGO dependency.

## Overview

`golibjpeg` is a Go library for decoding and encoding JPEG images with native precision (8‑bit and 16‑bit). It bundles a platform-specific shared library extracted at runtime via FFI (`ebitengine/purego`), avoiding the need for CGO.

Supported formats:
- **JPEG** (ISO 10918‑1, baseline / lossless)
- **JPEG‑LS** (ISO 14495, lossless / near‑lossless)
- **JPEG XT** (ISO 18477, HDR — decode only)

## API

Aligned with [pylibjpeg-libjpeg](https://github.com/pydicom/pylibjpeg-libjpeg) `libjpeg.utils`:

```go
// Decode JPEG/JPEG-LS/JPEG XT (colour_transform matches Python default 0)
func DecodeImage(stream any, colourTransform ColourTransform) (*Image, error)

// Encode interleaved little-endian pixels to JPEG / JPEG-LS
func Encode(src []byte, opts EncodeOptions) ([]byte, error)

// DICOM encapsulated pixel data
func DecodePixelData(src []byte, opts PixelDataOptions) ([]byte, error)
func EncodePixelData(src []byte, desc PixelDataDescriptor, opts EncodePixelDataOptions) ([]byte, error)

// Read parameters without decoding
func GetImageParameters(stream any) (*Params, error)

// Shorthands
func Decode(data []byte) (*Image, error)
func GetParameters(data []byte) (*Params, error)
```

`stream` may be `[]byte`, file path (`string`), or `io.Reader`.

`ColourTransform` constants: `ColourTransformNone` (0), `ColourTransformYCbCr` (1), `ColourTransformRCT` (2), `ColourTransformFreeform` (3).

No CGO: native code is loaded via `purego` + `//go:embed` prebuilt libraries.

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

Encode to JPEG baseline:

```go
out, err := golibjpeg.Encode(pixels, golibjpeg.EncodeOptions{
	Columns: 512, Rows: 512, SamplesPerPixel: 3, BitsPerSample: 8,
	FrameType: golibjpeg.FrameBaseline, ColourTransform: golibjpeg.ColourTransformYCbCr,
	Quality: 90,
})
```

JPEG-LS lossless:

```go
out, err := golibjpeg.Encode(frame, golibjpeg.EncodeOptions{
	Columns: 512, Rows: 512, SamplesPerPixel: 1, BitsPerSample: 16,
	FrameType: golibjpeg.FrameJPEGLS, LSInterleaving: golibjpeg.LSInterleaveSample,
})
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
| Windows | ✓     | ✓     |
| macOS   | ✓     | ✓     |
| Linux   | ✓     | ✓     |

Anywhere else this module still **builds** — it just cannot decode or encode.
Every function returns an error wrapping `ErrUnsupportedPlatform` instead, so a
program that imports `golibjpeg` (or `godicom`, which does) keeps compiling and
running on a platform with no prebuilt library, and only JPEG and JPEG-LS fail:

```go
img, err := golibjpeg.Decode(data)
if errors.Is(err, golibjpeg.ErrUnsupportedPlatform) {
	// no library for this GOOS/GOARCH; err names which one
}
```

The `cross-build` CI job compiles the module for a spread of platforms outside
the table — `js/wasm` and `wasip1/wasm` among them — so this stays true. Loading
is lazy and never panics: a read-only or `noexec` `TMPDIR` also surfaces as an
error from the first call.

### What each platform costs your binary

The six libraries together are about 10 MB, but a binary only ever carries the
one it can load — the `//go:embed` directives are behind per-platform build
tags, so the other five are not compiled in:

```
$ GOOS=windows GOARCH=arm64 go build -o app .
$ GOOS=linux   GOARCH=386   go build -o app .   # off the matrix: nothing embedded
```

| target | embedded | added to the binary |
|--------|----------|---------------------|
| `linux/amd64` | `golibjpeg_linux_amd64.so` | ~2.2 MB |
| `linux/arm64` | `golibjpeg_linux_arm64.so` | ~2.1 MB |
| `darwin/amd64` | `golibjpeg_darwin_amd64.dylib` | ~1.5 MB |
| `darwin/arm64` | `golibjpeg_darwin_arm64.dylib` | ~1.4 MB |
| `windows/amd64` | `golibjpeg_amd64.dll` | ~1.5 MB |
| `windows/arm64` | `golibjpeg_arm64.dll` | ~1.4 MB |
| anything else | — | nothing |

The `checks` CI job asserts that set per platform with `go list -f
'{{.EmbedFiles}}'`, because a `libs/*` glob or a forgotten build tag would put
all six into every binary and nothing else would notice.

`go get` does download all six, since they live in one module — that cost is
paid once in the module cache, not per build and not per user binary.

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

`build.yml` runs two jobs on their own — **checks** (gofmt, `go vet`, and the
one-embedded-library-per-platform assertion) and **cross-build** (compile for 7
platforms with no prebuilt library) — plus this chain in order:

1. **build-native** — build shared library on 6 platforms, upload artifacts  
2. **commit-native** — on push to `main`, write artifacts into `native/libs/` and commit  
3. **test** — download artifact per platform, then `go test`  
4. **release** — on `v*` tags, attach libraries to GitHub Release  

Committed files in `native/libs/` let `go get` work without a local CMake install.

Reference tests (`reference_compliance_test.go`) mirror `ref/pylibjpeg-libjpeg/libjpeg/tests/test_parameters.py` and `test_decode.py` (`REF_JPG` table, 23 images). Fetch testdata before running:

```bash
bash scripts/fetch-testdata.sh
go test ./...
```

### Release workflow

1. Merge changes to `main` and wait for `build` workflow (build → commit `native/libs/` → test).
2. `build` workflow runs `go test` on all platforms using embedded libs.
3. Create and push a tag: `git tag v1.0.1 && git push origin v1.0.1`.
4. CI attaches the committed libraries from `native/libs/` to a GitHub Release.

## References

This Go port follows **[pylibjpeg-libjpeg](https://github.com/pydicom/pylibjpeg-libjpeg)** for native decode behaviour and tests, and **[pylibjpeg](https://github.com/pydicom/pylibjpeg)** for the overall plugin-style integration model used by pydicom.
