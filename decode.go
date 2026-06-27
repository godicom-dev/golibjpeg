package golibjpeg

import "github.com/godicom-dev/golibjpeg/native"

// DecodeImage decodes JPEG/JPEG-LS/JPEG XT data from stream.
// colourTransform matches pylibjpeg-libjpeg decode(colour_transform=...).
func DecodeImage(stream any, colourTransform ColourTransform) (*Image, error) {
	data, err := ReadStream(stream)
	if err != nil {
		return nil, err
	}
	res, err := native.Decode(data, int(colourTransform))
	if err != nil {
		return nil, err
	}
	return newImage(nativeDecodeResult{
		Output:     res.Output,
		Width:      res.Width,
		Height:     res.Height,
		Components: res.Components,
		Precision:  res.Precision,
	}), nil
}

// Decode decodes with no colour transform (pylibjpeg default colour_transform=0).
func Decode(data []byte) (*Image, error) {
	return DecodeImage(data, ColourTransformNone)
}

// DecodeWithFormat decodes using a legacy Format colour-transform selector.
func DecodeWithFormat(data []byte, format Format) (*Image, error) {
	return DecodeImage(data, format.colourTransform())
}

// GetImageParameters reads JPEG parameters without decoding pixels.
func GetImageParameters(stream any) (*Params, error) {
	data, err := ReadStream(stream)
	if err != nil {
		return nil, err
	}
	p, err := native.GetParameters(data)
	if err != nil {
		return nil, err
	}
	return newParams(nativeImageParams{
		Width:      p.Width,
		Height:     p.Height,
		Components: p.Components,
		Precision:  p.Precision,
	}), nil
}

// GetParameters is an alias for GetImageParameters using a byte buffer.
func GetParameters(data []byte) (*Params, error) {
	return GetImageParameters(data)
}
