package golibjpeg

import "github.com/godicom-dev/golibjpeg/native"

type Image struct {
	Pixels     []byte
	Width      int
	Height     int
	Components int
	Precision  int
}

func Decode(data []byte) (*Image, error) {
	return DecodeWithFormat(data, FormatAuto)
}

func DecodeWithFormat(data []byte, format Format) (*Image, error) {
	ct := int(format)
	if ct < 0 {
		ct = native.CTYCbCr
	}
	res, err := native.Decode(data, ct)
	if err != nil {
		return nil, err
	}
	return &Image{
		Pixels:     res.Output,
		Width:      res.Width,
		Height:     res.Height,
		Components: res.Components,
		Precision:  res.Precision,
	}, nil
}

type Params struct {
	Width      int
	Height     int
	Components int
	Precision  int
}

func GetParameters(data []byte) (*Params, error) {
	p, err := native.GetParameters(data)
	if err != nil {
		return nil, err
	}
	return &Params{
		Width:      p.Width,
		Height:     p.Height,
		Components: p.Components,
		Precision:  p.Precision,
	}, nil
}
