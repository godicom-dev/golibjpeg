package golibjpeg_test

import (
	"testing"

	"github.com/godicom-dev/golibjpeg"
	"github.com/godicom-dev/golibjpeg/internal/testdata"
)

type refDCMCase struct {
	uid  string
	name string
	rows int
	cols int
	spp  int
	bps  int
}

var refDCMCases = []refDCMCase{
	{"1.2.840.10008.1.2.4.50", "JPEGBaseline_1s_1f_u_08_08.dcm", 100, 100, 1, 8},
	{"1.2.840.10008.1.2.4.50", "SC_rgb_dcmtk_+eb+cy+np.dcm", 100, 100, 3, 8},
	{"1.2.840.10008.1.2.4.50", "color3d_jpeg_baseline.dcm", 480, 640, 3, 8},
	{"1.2.840.10008.1.2.4.50", "SC_rgb_dcmtk_+eb+cr.dcm", 100, 100, 3, 8},
	{"1.2.840.10008.1.2.4.50", "SC_rgb_dcmtk_+eb+cy+n1.dcm", 100, 100, 3, 8},
	{"1.2.840.10008.1.2.4.50", "SC_rgb_dcmtk_+eb+cy+s4.dcm", 100, 100, 3, 8},
	{"1.2.840.10008.1.2.4.51", "RG2_JPLY_fixed.dcm", 2140, 1760, 1, 12},
	{"1.2.840.10008.1.2.4.51", "JPEGExtended_1s_1f_u_16_12.dcm", 1024, 256, 1, 12},
	{"1.2.840.10008.1.2.4.51", "JPEGExtended_3s_1f_u_08_08.dcm", 576, 768, 3, 8},
	{"1.2.840.10008.1.2.4.57", "JPEGLossless_1s_1f_u_16_12.dcm", 1024, 1024, 1, 12},
	{"1.2.840.10008.1.2.4.70", "JPEG-LL.dcm", 1024, 256, 1, 16},
	{"1.2.840.10008.1.2.4.70", "JPEGLosslessP14SV1_1s_1f_u_08_08.dcm", 768, 1024, 1, 8},
	{"1.2.840.10008.1.2.4.70", "JPEGLosslessP14SV1_1s_1f_u_16_16.dcm", 535, 800, 1, 16},
	{"1.2.840.10008.1.2.4.70", "MG1_JPLL.dcm", 4664, 3064, 1, 12},
	{"1.2.840.10008.1.2.4.70", "RG1_JPLL.dcm", 1955, 1841, 1, 15},
	{"1.2.840.10008.1.2.4.70", "RG2_JPLL.dcm", 2140, 1760, 1, 10},
	{"1.2.840.10008.1.2.4.70", "SC_rgb_jpeg_gdcm.dcm", 100, 100, 3, 8},
	{"1.2.840.10008.1.2.4.80", "emri_small_jpeg_ls_lossless.dcm", 64, 64, 1, 16},
	{"1.2.840.10008.1.2.4.80", "MR_small_jpeg_ls_lossless.dcm", 64, 64, 1, 16},
	{"1.2.840.10008.1.2.4.80", "RG1_JLSL.dcm", 1955, 1841, 1, 15},
	{"1.2.840.10008.1.2.4.80", "RG2_JLSL.dcm", 2140, 1760, 1, 10},
	{"1.2.840.10008.1.2.4.81", "CT1_JLSN.dcm", 512, 512, 1, 16},
	{"1.2.840.10008.1.2.4.81", "MG1_JLSN.dcm", 4664, 3064, 1, 12},
	{"1.2.840.10008.1.2.4.81", "RG1_JLSN.dcm", 1955, 1841, 1, 15},
	{"1.2.840.10008.1.2.4.81", "RG2_JLSN.dcm", 2140, 1760, 1, 10},
}

func TestReferenceDCMGetParameters(t *testing.T) {
	for _, tc := range refDCMCases {
		t.Run(tc.uid+"/"+tc.name, func(t *testing.T) {
			frame, meta := testdata.RequireDCMFrame(t, tc.uid, tc.name)
			params, err := golibjpeg.GetImageParameters(frame)
			if err != nil {
				t.Fatalf("GetImageParameters: %v", err)
			}
			if params.Height != tc.rows || params.Width != tc.cols {
				t.Fatalf("dims: got %dx%d want %dx%d", params.Height, params.Width, tc.rows, tc.cols)
			}
			if params.Components != tc.spp {
				t.Fatalf("components: got %d want %d", params.Components, tc.spp)
			}
			if params.Precision != tc.bps {
				t.Fatalf("precision: got %d want %d", params.Precision, tc.bps)
			}
			if meta.Rows != tc.rows || meta.Columns != tc.cols {
				t.Fatalf("metadata mismatch: %+v", meta)
			}
		})
	}
}

func TestReferenceDCMDecode(t *testing.T) {
	for _, tc := range refDCMCases {
		t.Run(tc.uid+"/"+tc.name, func(t *testing.T) {
			frame, _ := testdata.RequireDCMFrame(t, tc.uid, tc.name)
			img, err := golibjpeg.DecodeImage(frame, golibjpeg.ColourTransformNone)
			if err != nil {
				t.Fatalf("DecodeImage: %v", err)
			}
			if img.Height != tc.rows || img.Width != tc.cols {
				t.Fatalf("dims: got %dx%d want %dx%d", img.Height, img.Width, tc.rows, tc.cols)
			}
			if img.Components != tc.spp {
				t.Fatalf("components: got %d want %d", img.Components, tc.spp)
			}
			if img.Precision != tc.bps {
				t.Fatalf("precision: got %d want %d", img.Precision, tc.bps)
			}
		})
	}
}

// Handler spot checks from pylibjpeg-libjpeg test_handler.py (GDCM reference values).
func TestReferenceDCMHandlerPixels(t *testing.T) {
	uid := "1.2.840.10008.1.2.4.50"
	frame, _ := testdata.RequireDCMFrame(t, uid, "JPEGBaseline_1s_1f_u_08_08.dcm")
	img, err := golibjpeg.DecodeImage(frame, golibjpeg.ColourTransformNone)
	if err != nil {
		t.Fatal(err)
	}
	if got := int(img.ByteAt(5, 50, 0)); got != 76 {
		t.Fatalf("arr[5,50]: got %d want 76", got)
	}
	if got := int(img.ByteAt(95, 50, 0)); got != 255 {
		t.Fatalf("arr[95,50]: got %d want 255", got)
	}

	frame, _ = testdata.RequireDCMFrame(t, uid, "SC_rgb_dcmtk_+eb+cy+np.dcm")
	img, err = golibjpeg.DecodeImage(frame, golibjpeg.ColourTransformNone)
	if err != nil {
		t.Fatal(err)
	}
	want := []int{76, 85, 254}
	for c := 0; c < 3; c++ {
		if got := int(img.ByteAt(5, 50, c)); got != want[c] {
			t.Fatalf("arr[5,50,%d]: got %d want %d", c, got, want[c])
		}
	}
}
