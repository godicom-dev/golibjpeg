package golibjpeg

import (
	"bytes"
	"testing"
)

func makeRGB8(width, height int) []byte {
	buf := make([]byte, width*height*3)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			off := (y*width + x) * 3
			buf[off] = byte(x * 255 / max(width-1, 1))
			buf[off+1] = byte(y * 255 / max(height-1, 1))
			buf[off+2] = byte((x + y) * 255 / max(width+height-2, 1))
		}
	}
	return buf
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func TestEncodeBaselineRoundTrip(t *testing.T) {
	const w, h = 64, 48
	src := makeRGB8(w, h)

	encoded, err := Encode(src, EncodeOptions{
		Columns:         w,
		Rows:            h,
		SamplesPerPixel: 3,
		BitsPerSample:   8,
		FrameType:       FrameBaseline,
		ColourTransform: ColourTransformYCbCr,
		Quality:         90,
	})
	if err != nil {
		t.Fatalf("Encode baseline: %v", err)
	}
	if len(encoded) == 0 {
		t.Fatal("Encode baseline returned empty output")
	}

	params, err := GetParameters(encoded)
	if err != nil {
		t.Fatalf("GetParameters: %v", err)
	}
	if params.Width != w || params.Height != h || params.Components != 3 {
		t.Fatalf("params = %+v, want %dx%d 3 components", params, w, h)
	}

	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatalf("Decode round-trip: %v", err)
	}
	if decoded.Width != w || decoded.Height != h || decoded.Components != 3 {
		t.Fatalf("decoded size = %dx%d x%d", decoded.Width, decoded.Height, decoded.Components)
	}
	if len(decoded.Pixels) != len(src) {
		t.Fatalf("decoded pixel len = %d, want %d", len(decoded.Pixels), len(src))
	}
}

func TestEncodeLosslessRoundTrip(t *testing.T) {
	const w, h = 32, 32
	src := makeRGB8(w, h)

	encoded, err := Encode(src, EncodeOptions{
		Columns:         w,
		Rows:            h,
		SamplesPerPixel: 3,
		BitsPerSample:   8,
		FrameType:       FrameLossless,
	})
	if err != nil {
		t.Fatalf("Encode lossless: %v", err)
	}

	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatalf("Decode lossless round-trip: %v", err)
	}
	if !bytes.Equal(decoded.Pixels, src) {
		t.Fatal("lossless round-trip pixel mismatch")
	}
}

func TestEncodeJPEGLSRoundTrip(t *testing.T) {
	const w, h = 32, 32
	src := makeRGB8(w, h)

	encoded, err := Encode(src, EncodeOptions{
		Columns:         w,
		Rows:            h,
		SamplesPerPixel: 3,
		BitsPerSample:   8,
		FrameType:       FrameJPEGLS,
		ErrorBound:      0,
		LSInterleaving:  LSInterleaveSample,
	})
	if err != nil {
		t.Fatalf("Encode JPEG-LS: %v", err)
	}

	decoded, err := DecodeWithFormat(encoded, FormatJPEGLS)
	if err != nil {
		t.Fatalf("Decode JPEG-LS round-trip: %v", err)
	}
	if !bytes.Equal(decoded.Pixels, src) {
		t.Fatal("JPEG-LS round-trip pixel mismatch")
	}
}

func TestEncodeMonochromeLossless(t *testing.T) {
	const w, h = 16, 16
	src := make([]byte, w*h)
	for i := range src {
		src[i] = byte(i)
	}

	encoded, err := Encode(src, EncodeOptions{
		Columns:         w,
		Rows:            h,
		SamplesPerPixel: 1,
		BitsPerSample:   8,
		FrameType:       FrameLossless,
	})
	if err != nil {
		t.Fatalf("Encode mono lossless: %v", err)
	}

	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatalf("Decode mono lossless: %v", err)
	}
	if !bytes.Equal(decoded.Pixels, src) {
		t.Fatal("mono lossless round-trip pixel mismatch")
	}
}

func TestEncodePixelData(t *testing.T) {
	const w, h = 8, 8
	src := makeRGB8(w, h)

	out, err := EncodePixelData(src, PixelDataDescriptor{
		Columns:         w,
		Rows:            h,
		SamplesPerPixel: 3,
		BitsAllocated:   8,
	}, EncodePixelDataOptions{
		PhotometricInterpretation: PhotometricRGB,
		FrameType:                 FrameBaseline,
		Quality:                   85,
	})
	if err != nil {
		t.Fatalf("EncodePixelData: %v", err)
	}
	if len(out) == 0 {
		t.Fatal("EncodePixelData returned empty output")
	}
}
