package golibjpeg

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

func createTestJPEG(t *testing.T) []byte {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))

	red := color.RGBA{255, 0, 0, 255}
	green := color.RGBA{0, 255, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}

	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			switch {
			case x < 32 && y < 32:
				img.SetRGBA(x, y, red)
			case x >= 32 && y < 32:
				img.SetRGBA(x, y, green)
			case x < 32 && y >= 32:
				img.SetRGBA(x, y, blue)
			default:
				img.SetRGBA(x, y, color.RGBA{255, 255, 0, 255})
			}
		}
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("failed to create test JPEG: %v", err)
	}
	return buf.Bytes()
}

func TestDecode(t *testing.T) {
	data := createTestJPEG(t)
	if len(data) == 0 {
		t.Fatal("test JPEG data is empty")
	}

	img, err := Decode(data)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if img.Width != 64 {
		t.Errorf("expected width 64, got %d", img.Width)
	}
	if img.Height != 64 {
		t.Errorf("expected height 64, got %d", img.Height)
	}
	if img.Components < 3 {
		t.Errorf("expected at least 3 components, got %d", img.Components)
	}
	if img.Precision != 8 {
		t.Errorf("expected precision 8, got %d", img.Precision)
	}
	expectedLen := 64 * 64 * 3
	if len(img.Pixels) != expectedLen {
		t.Errorf("expected pixel data length %d, got %d", expectedLen, len(img.Pixels))
	}

	t.Logf("Decoded JPEG: %dx%d, %d components, precision %d, %d bytes",
		img.Width, img.Height, img.Components, img.Precision, len(img.Pixels))

	pixelIndex := 0
	r := img.Pixels[pixelIndex*3]
	g := img.Pixels[pixelIndex*3+1]
	b := img.Pixels[pixelIndex*3+2]
	t.Logf("Pixel (0,0) = R:%d G:%d B:%d", r, g, b)
}

func TestDecodeWithFormat(t *testing.T) {
	data := createTestJPEG(t)
	img, err := DecodeWithFormat(data, FormatJPEG)
	if err != nil {
		t.Fatalf("DecodeWithFormat failed: %v", err)
	}
	if img.Width != 64 || img.Height != 64 {
		t.Errorf("unexpected dimensions: %dx%d", img.Width, img.Height)
	}
}

func TestDecodeFromFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.jpg")
	data := createTestJPEG(t)
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		t.Fatal(err)
	}

	fileData, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	img, err := Decode(fileData)
	if err != nil {
		t.Fatalf("Decode from file failed: %v", err)
	}

	if img.Width != 64 || img.Height != 64 {
		t.Errorf("unexpected dimensions: %dx%d", img.Width, img.Height)
	}
}

func TestGetParameters(t *testing.T) {
	data := createTestJPEG(t)
	params, err := GetParameters(data)
	if err != nil {
		t.Fatalf("GetParameters failed: %v", err)
	}
	if params.Width != 64 {
		t.Errorf("expected width 64, got %d", params.Width)
	}
	if params.Height != 64 {
		t.Errorf("expected height 64, got %d", params.Height)
	}
	if params.Components < 1 {
		t.Errorf("expected at least 1 component, got %d", params.Components)
	}
	if params.Precision < 8 {
		t.Errorf("expected precision >= 8, got %d", params.Precision)
	}
	t.Logf("JPEG params: %dx%d, %d components, precision %d",
		params.Width, params.Height, params.Components, params.Precision)
}

func TestGetParametersInvalid(t *testing.T) {
	_, err := GetParameters([]byte{0xFF, 0xD8, 0xFF})
	if err == nil {
		t.Fatal("expected error for invalid JPEG")

	}
	t.Logf("Expected error: %v", err)
}
