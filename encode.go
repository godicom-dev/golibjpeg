package golibjpeg

import (
	"fmt"

	"github.com/godicom-dev/golibjpeg/native"
)

// FrameType selects the JPEG encoding process (libjpeg JPGFLAG_* frame types).
type FrameType int

const (
	FrameBaseline   FrameType = native.FrameBaseline
	FrameSequential FrameType = native.FrameSequential
	FrameLossless   FrameType = native.FrameLossless
	FrameJPEGLS     FrameType = native.FrameJPEGLS
)

// LSInterleaving configures JPEG-LS component interleaving.
type LSInterleaving int

const (
	LSInterleaveNone   LSInterleaving = native.LSInterleaveNone
	LSInterleaveLine   LSInterleaving = native.LSInterleaveLine
	LSInterleaveSample LSInterleaving = native.LSInterleaveSample
)

// EncodeOptions configures Encode / EncodePixelData.
type EncodeOptions struct {
	Columns         int
	Rows            int
	SamplesPerPixel int
	BitsPerSample   int
	FrameType       FrameType
	ColourTransform ColourTransform
	// Quality applies to lossy JPEG (1–100). Ignored for lossless / JPEG-LS.
	Quality int
	// ErrorBound is the JPEG-LS NEAR parameter (0 = lossless).
	ErrorBound int
	// LSInterleaving applies when FrameType is FrameJPEGLS.
	LSInterleaving LSInterleaving
}

// Encode compresses interleaved little-endian pixel samples to JPEG / JPEG-LS.
func Encode(src []byte, opts EncodeOptions) ([]byte, error) {
	if opts.FrameType == 0 {
		opts.FrameType = FrameBaseline
	}
	if opts.Quality <= 0 {
		opts.Quality = 75
	}
	if opts.SamplesPerPixel == 1 {
		opts.ColourTransform = ColourTransformNone
	}
	return native.Encode(src, native.EncodeParams{
		Columns:         opts.Columns,
		Rows:            opts.Rows,
		SamplesPerPixel: opts.SamplesPerPixel,
		BitsPerSample:   opts.BitsPerSample,
		FrameType:       int(opts.FrameType),
		ColourTransform: int(opts.ColourTransform),
		Quality:         opts.Quality,
		ErrorBound:      opts.ErrorBound,
		LSInterleaving:  int(opts.LSInterleaving),
	})
}

// EncodePixelDataOptions configures EncodePixelData for DICOM handlers.
type EncodePixelDataOptions struct {
	PhotometricInterpretation string
	FrameType                 FrameType
	Quality                   int
	ErrorBound                int
	LSInterleaving            LSInterleaving
}

// PixelDataDescriptor holds image geometry for EncodePixelData.
type PixelDataDescriptor struct {
	Columns         int
	Rows            int
	SamplesPerPixel int
	BitsAllocated   int
	BitsStored      int
}

// EncodePixelData encodes one DICOM frame for JPEG / JPEG-LS transfer syntaxes.
func EncodePixelData(src []byte, desc PixelDataDescriptor, opts EncodePixelDataOptions) ([]byte, error) {
	bits := desc.BitsStored
	if bits == 0 {
		bits = desc.BitsAllocated
	}
	if opts.FrameType == 0 {
		opts.FrameType = FrameBaseline
	}

	ct := ColourTransformNone
	switch opts.PhotometricInterpretation {
	case PhotometricRGB:
		if opts.FrameType == FrameBaseline || opts.FrameType == FrameSequential {
			ct = ColourTransformYCbCr
		}
	}

	encOpts := EncodeOptions{
		Columns:         desc.Columns,
		Rows:            desc.Rows,
		SamplesPerPixel: desc.SamplesPerPixel,
		BitsPerSample:   bits,
		FrameType:       opts.FrameType,
		ColourTransform: ct,
		Quality:         opts.Quality,
		ErrorBound:      opts.ErrorBound,
		LSInterleaving:  opts.LSInterleaving,
	}
	if encOpts.Quality <= 0 {
		encOpts.Quality = 75
	}

	out, err := Encode(src, encOpts)
	if err != nil {
		return nil, fmt.Errorf("golibjpeg: encode pixel data: %w", err)
	}
	return out, nil
}
