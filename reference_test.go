package golibjpeg

import (
	"os"
	"path/filepath"
	"testing"
)

// Reference cases from ref/pylibjpeg-libjpeg/libjpeg/tests/test_parameters.py
func TestGetParametersReferenceJPEG(t *testing.T) {
	path := filepath.Join("testdata", "10918", "p1", "A1.JPG")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skipf("reference testdata not installed: %v", err)
	}
	if err != nil {
		t.Fatal(err)
	}

	params, err := GetParameters(data)
	if err != nil {
		t.Fatalf("GetParameters: %v", err)
	}

	// A1.JPG: 257 rows, 255 columns, 4 components, 8-bit precision
	if params.Height != 257 {
		t.Errorf("rows: got %d, want 257", params.Height)
	}
	if params.Width != 255 {
		t.Errorf("columns: got %d, want 255", params.Width)
	}
	if params.Components != 4 {
		t.Errorf("components: got %d, want 4", params.Components)
	}
	if params.Precision != 8 {
		t.Errorf("precision: got %d, want 8", params.Precision)
	}
}

// Reference cases from ref/pylibjpeg-libjpeg/libjpeg/tests/test_decode.py (baseline frame pixel)
func TestDecodeReferenceJPEGPixels(t *testing.T) {
	path := filepath.Join("testdata", "10918", "p1", "A1.JPG")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skipf("reference testdata not installed: %v", err)
	}
	if err != nil {
		t.Fatal(err)
	}

	img, err := DecodeWithFormat(data, FormatJPEG)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	if img.Height != 257 || img.Width != 255 {
		t.Fatalf("dims: got %dx%d, want 257x255", img.Height, img.Width)
	}

	// Corner pixels from pylibjpeg-libjpeg test_decode.py colourspace tables
	want := [4]byte{138, 76, 239, 216}
	if len(img.Pixels) < 4 {
		t.Fatal("pixel buffer too short")
	}
	for i := range want {
		if img.Pixels[i] != want[i] {
			t.Errorf("pixel[%d]: got %d, want %d", i, img.Pixels[i], want[i])
		}
	}
}
