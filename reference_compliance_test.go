package golibjpeg_test

import (
	"reflect"
	"testing"

	"github.com/godicom-dev/golibjpeg"
	"github.com/godicom-dev/golibjpeg/internal/testdata"
)

// refJPG mirrors REF_JPG from ref/pylibjpeg-libjpeg/libjpeg/tests/test_decode.py
type refJPGCase struct {
	dir      string // p1, p2, p4, p14, JLS, JNL
	filename string
	rows     int
	cols     int
	comp     int
	prec     int
	topLeft  []int
	bottomRt []int
}

var refJPGCases = []refJPGCase{
	// 10918 p1
	{"p1", "A1.JPG", 257, 255, 4, 8, []int{138, 76, 239, 216}, []int{155, 217, 191, 115}},
	// 10918 p2
	{"p2", "A1.JPG", 257, 255, 4, 8, []int{138, 76, 239, 216}, []int{155, 217, 191, 115}},
	{"p2", "C1.JPG", 257, 255, 4, 8, []int{138, 76, 239, 216}, []int{155, 217, 191, 115}},
	{"p2", "C2.JPG", 257, 255, 4, 8, []int{138, 76, 239, 216}, []int{155, 217, 191, 115}},
	// 10918 p4
	{"p4", "A1.JPG", 257, 255, 4, 8, []int{138, 76, 239, 216}, []int{155, 217, 191, 115}},
	{"p4", "C1.JPG", 257, 255, 4, 8, []int{138, 76, 239, 216}, []int{155, 217, 191, 115}},
	{"p4", "C2.JPG", 257, 255, 4, 8, []int{138, 76, 239, 216}, []int{155, 217, 191, 115}},
	{"p4", "E1.JPG", 257, 255, 4, 12, []int{2119, 1183, 3907, 3487}, []int{2502, 3402, 3041, 1872}},
	{"p4", "E2.JPG", 257, 255, 4, 12, []int{2119, 1183, 3907, 3487}, []int{2502, 3402, 3041, 1872}},
	// 10918 p14
	{"p14", "O1.JPG", 257, 255, 4, 8, []int{132, 76, 245, 218}, []int{156, 211, 191, 116}},
	{"p14", "O2.JPG", 257, 255, 4, 16, []int{33792, 19456, 62720, 55808}, []int{39936, 53888, 48768, 29696}},
	// 14495 JLS
	{"JLS", "T8C0E0.JLS", 256, 256, 3, 8, []int{161, 122, 108}, []int{101, 99, 95}},
	{"JLS", "T8C1E0.JLS", 256, 256, 3, 8, []int{161, 122, 108}, []int{101, 99, 95}},
	{"JLS", "T8C2E0.JLS", 256, 256, 3, 8, []int{161, 122, 108}, []int{101, 99, 95}},
	{"JLS", "T8NDE0.JLS", 128, 128, 1, 8, []int{108}, []int{231}},
	{"JLS", "T8SSE0.JLS", 256, 256, 3, 8, []int{161, 122, 108}, []int{101, 171, 231}},
	{"JLS", "T16E0.JLS", 256, 256, 1, 12, []int{1963}, []int{1596}},
	// 14495 JNL
	{"JNL", "T8C0E3.JLS", 256, 256, 3, 8, []int{161, 119, 105}, []int{98, 96, 93}},
	{"JNL", "T8C1E3.JLS", 256, 256, 3, 8, []int{161, 119, 105}, []int{101, 100, 97}},
	{"JNL", "T8C2E3.JLS", 256, 256, 3, 8, []int{161, 119, 105}, []int{98, 96, 94}},
	{"JNL", "T8NDE3.JLS", 128, 128, 1, 8, []int{105}, []int{229}},
	{"JNL", "T8SSE3.JLS", 256, 256, 3, 8, []int{161, 119, 105}, []int{102, 169, 234}},
	{"JNL", "T16E3.JLS", 256, 256, 1, 12, []int{1960}, []int{1593}},
}

func jpegDataPath(dir, filename string) []string {
	if dir == "JLS" || dir == "JNL" {
		return []string{"14495", dir, filename}
	}
	return []string{"10918", dir, filename}
}

func TestReferenceGetParametersJPG(t *testing.T) {
	for _, tc := range refJPGCases {
		t.Run(tc.dir+"/"+tc.filename, func(t *testing.T) {
			data := testdata.RequireJPEG(t, jpegDataPath(tc.dir, tc.filename)...)
			params, err := golibjpeg.GetImageParameters(data)
			if err != nil {
				t.Fatalf("GetImageParameters: %v", err)
			}
			if params.Height != tc.rows {
				t.Errorf("rows: got %d want %d", params.Height, tc.rows)
			}
			if params.Width != tc.cols {
				t.Errorf("columns: got %d want %d", params.Width, tc.cols)
			}
			if params.Components != tc.comp {
				t.Errorf("components: got %d want %d", params.Components, tc.comp)
			}
			if params.Precision != tc.prec {
				t.Errorf("precision: got %d want %d", params.Precision, tc.prec)
			}
		})
	}
}

func TestReferenceDecodeJPG(t *testing.T) {
	for _, tc := range refJPGCases {
		t.Run(tc.dir+"/"+tc.filename, func(t *testing.T) {
			data := testdata.RequireJPEG(t, jpegDataPath(tc.dir, tc.filename)...)
			img, err := golibjpeg.DecodeImage(data, golibjpeg.ColourTransformNone)
			if err != nil {
				t.Fatalf("DecodeImage: %v", err)
			}
			if img.Height != tc.rows || img.Width != tc.cols {
				t.Fatalf("dims: got %dx%d want %dx%d", img.Height, img.Width, tc.rows, tc.cols)
			}
			if img.Components != tc.comp {
				t.Fatalf("components: got %d want %d", img.Components, tc.comp)
			}
			if img.Precision != tc.prec {
				t.Fatalf("precision: got %d want %d", img.Precision, tc.prec)
			}

			topLeft, bottomRight := img.CornerSamples()
			if !reflect.DeepEqual(topLeft, tc.topLeft) {
				t.Errorf("top-left: got %v want %v", topLeft, tc.topLeft)
			}
			if !reflect.DeepEqual(bottomRight, tc.bottomRt) {
				t.Errorf("bottom-right: got %v want %v", bottomRight, tc.bottomRt)
			}

			raw, err := golibjpeg.DecodePixelData(data, golibjpeg.PixelDataOptions{Version: golibjpeg.PixelDataV2})
			if err != nil {
				t.Fatalf("DecodePixelData v2: %v", err)
			}
			if !reflect.DeepEqual(raw, img.Pixels) {
				t.Errorf("DecodePixelData v2 length mismatch: got %d want %d", len(raw), len(img.Pixels))
			}
		})
	}
}

func TestReferenceGetParametersA1Path(t *testing.T) {
	path := testdata.JPEGPath("10918", "p1", "A1.JPG")
	params, err := golibjpeg.GetImageParameters(path)
	if err != nil {
		t.Fatalf("GetImageParameters(path): %v", err)
	}
	if params.Height != 257 || params.Width != 255 || params.Components != 4 || params.Precision != 8 {
		t.Fatalf("unexpected params: %+v", params)
	}
}

func TestReferenceDecodeA1Reader(t *testing.T) {
	data := testdata.RequireJPEG(t, "10918", "p1", "A1.JPG")
	img, err := golibjpeg.DecodeImage(data, golibjpeg.ColourTransformNone)
	if err != nil {
		t.Fatal(err)
	}
	if img.Height != 257 || img.Width != 255 || img.Components != 4 {
		t.Fatalf("unexpected image: %+v", img)
	}
}

func TestDecodePixelDataMissingPhotometric(t *testing.T) {
	data := testdata.RequireJPEG(t, "10918", "p1", "A1.JPG")
	_, err := golibjpeg.DecodePixelData(data, golibjpeg.PixelDataOptions{Version: golibjpeg.PixelDataV1})
	if err == nil {
		t.Fatal("expected error for missing photometric interpretation")
	}
}

func TestDecodeInvalidStreamType(t *testing.T) {
	_, err := golibjpeg.DecodeImage(nil, golibjpeg.ColourTransformNone)
	if err == nil {
		t.Fatal("expected type error")
	}
}

func TestGetParametersInvalidJPEG(t *testing.T) {
	_, err := golibjpeg.GetImageParameters([]byte{0xFF, 0xD8, 0xFF})
	if err == nil {
		t.Fatal("expected error for truncated JPEG")
	}
}
