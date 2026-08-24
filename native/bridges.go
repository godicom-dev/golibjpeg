package native

import "unsafe"

type DecodeResult struct {
	Output     []byte
	Width      int
	Height     int
	Components int
	Precision  int
}

const (
	CTNone     = 0
	CTYCbCr    = 1
	CTRCT      = 2
	CTFreeform = 3

	FrameBaseline   = 0
	FrameSequential = 1
	FrameLossless   = 3
	FrameJPEGLS     = 4

	LSInterleaveNone   = 0
	LSInterleaveLine   = 1
	LSInterleaveSample = 2
)

func DecodeToRGB(data []byte) (*DecodeResult, error) {
	return decode(data, CTYCbCr)
}

func Decode(data []byte, colourTransform int) (*DecodeResult, error) {
	return decode(data, colourTransform)
}

func decode(data []byte, colourTransform int) (*DecodeResult, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyInput
	}

	var outputPtr unsafe.Pointer
	var outputLen int32
	var width, height, components, precision int32

	code := decodeFn(
		unsafe.Pointer(&data[0]),
		int32(len(data)),
		int32(colourTransform),
		&outputPtr,
		&outputLen,
		&width,
		&height,
		&components,
		&precision,
	)
	if code != 0 {
		return nil, errWithCode("Decode()", code)
	}

	if outputPtr == nil || outputLen == 0 {
		return nil, errEmptyOutput
	}

	out := make([]byte, outputLen)
	copy(out, unsafe.Slice((*byte)(outputPtr), outputLen))

	freeFn(outputPtr)

	return &DecodeResult{
		Output:     out,
		Width:      int(width),
		Height:     int(height),
		Components: int(components),
		Precision:  int(precision),
	}, nil
}

type ImageParams struct {
	Width      int
	Height     int
	Components int
	Precision  int
}

func GetParameters(data []byte) (*ImageParams, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyInput
	}

	var width, height, components, precision int32

	code := getParamsFn(
		unsafe.Pointer(&data[0]),
		int32(len(data)),
		&width,
		&height,
		&components,
		&precision,
	)
	if code != 0 {
		return nil, errWithCode("GetJPEGParameters()", code)
	}

	return &ImageParams{
		Width:      int(width),
		Height:     int(height),
		Components: int(components),
		Precision:  int(precision),
	}, nil
}

type EncodeParams struct {
	Columns         int
	Rows            int
	SamplesPerPixel int
	BitsPerSample   int
	FrameType       int
	ColourTransform int
	Quality         int
	ErrorBound      int
	LSInterleaving  int
}

func Encode(src []byte, params EncodeParams) ([]byte, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}
	if len(src) == 0 {
		return nil, errEmptyInput
	}
	if params.Columns <= 0 || params.Rows <= 0 || params.SamplesPerPixel <= 0 ||
		params.BitsPerSample <= 0 {
		return nil, errWithCode("Encode()", -5)
	}

	nativeParams := struct {
		columns         int32
		rows            int32
		samplesPerPixel int32
		bitsPerSample   int32
		frameType       int32
		colourTransform int32
		quality         int32
		errorBound      int32
		lsInterleaving  int32
	}{
		columns:         int32(params.Columns),
		rows:            int32(params.Rows),
		samplesPerPixel: int32(params.SamplesPerPixel),
		bitsPerSample:   int32(params.BitsPerSample),
		frameType:       int32(params.FrameType),
		colourTransform: int32(params.ColourTransform),
		quality:         int32(params.Quality),
		errorBound:      int32(params.ErrorBound),
		lsInterleaving:  int32(params.LSInterleaving),
	}

	var outputPtr unsafe.Pointer
	var outputLen int32

	code := encodeFn(
		unsafe.Pointer(&src[0]),
		int32(len(src)),
		unsafe.Pointer(&nativeParams),
		&outputPtr,
		&outputLen,
	)
	if code != 0 {
		return nil, errWithCode("Encode()", code)
	}
	if outputPtr == nil || outputLen == 0 {
		return nil, errEmptyEncodeOutput
	}

	out := make([]byte, outputLen)
	copy(out, unsafe.Slice((*byte)(outputPtr), outputLen))
	freeFn(outputPtr)
	return out, nil
}
