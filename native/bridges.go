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
	CTNone    = 0
	CTYCbCr   = 1
	CTRCT     = 2
	CTFreeform = 3
)

func DecodeToRGB(data []byte) (*DecodeResult, error) {
	return decode(data, CTYCbCr)
}

func Decode(data []byte, colourTransform int) (*DecodeResult, error) {
	return decode(data, colourTransform)
}

func decode(data []byte, colourTransform int) (*DecodeResult, error) {
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
