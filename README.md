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

## Platform support

| OS      | amd64 | arm64 | 386 |
|---------|-------|-------|-----|
| Windows | ✓     |       | ✓   |
| macOS   | ✓     | ✓     |     |
| Linux   | ✓     | ✓     |     |

## Dependencies

- [ebitengine/purego](https://github.com/ebitengine/purego) – FFI without CGO
- [thorfdbg/libjpeg](https://github.com/thorfdbg/libjpeg) – C++ JPEG library (ISO 10918‑1 / 18477)

## References

This Go port takes inspiration from the design and API of **[pylibjpeg](https://github.com/pydicom/pylibjpeg)** — a Python JPEG decoder that also wraps native C libraries via a CPython extension. The concept of a lightweight, cross‑platform JPEG decoder with a minimal public API and format‑aware colour transforms follows pylibjpeg's approach.
