package golibjpeg

import (
	"encoding/binary"
	"math"
)

// Image holds decoded pixel data in native precision (8- or 16-bit), planar-interleaved.
type Image struct {
	Pixels     []byte
	Width      int
	Height     int
	Components int
	Precision  int
}

// Params holds JPEG image parameters without decoding pixels.
type Params struct {
	Width      int
	Height     int
	Components int
	Precision  int
}

func (p *Params) Rows() int       { return p.Height }
func (p *Params) Columns() int    { return p.Width }
func (p *Params) NrComponents() int { return p.Components }

func (img *Image) BytesPerSample() int {
	return int(math.Ceil(float64(img.Precision) / 8))
}

func (img *Image) stride() int {
	return img.Width * img.Components * img.BytesPerSample()
}

func (img *Image) offset(y, x, c int) int {
	bps := img.BytesPerSample()
	return (y*img.Width+x)*img.Components*bps + c*bps
}

// ByteAt returns one byte from the interleaved pixel buffer.
func (img *Image) ByteAt(y, x, c int) byte {
	return img.Pixels[img.offset(y, x, c)]
}

// Uint16At returns a little-endian sample when precision > 8.
func (img *Image) Uint16At(y, x, c int) uint16 {
	off := img.offset(y, x, c)
	return binary.LittleEndian.Uint16(img.Pixels[off:])
}

// CornerSamples returns top-left and bottom-right samples as integers.
// For multi-byte precision each component is truncated to the low byte for 8-bit cases.
func (img *Image) CornerSamples() (topLeft, bottomRight []int) {
	topLeft = img.sampleAt(0, 0)
	bottomRight = img.sampleAt(img.Height-1, img.Width-1)
	return topLeft, bottomRight
}

func (img *Image) sampleAt(y, x int) []int {
	samples := make([]int, img.Components)
	for c := 0; c < img.Components; c++ {
		if img.BytesPerSample() == 1 {
			samples[c] = int(img.ByteAt(y, x, c))
			continue
		}
		samples[c] = int(img.Uint16At(y, x, c))
	}
	return samples
}

func newImage(res nativeDecodeResult) *Image {
	return &Image{
		Pixels:     res.Output,
		Width:      res.Width,
		Height:     res.Height,
		Components: res.Components,
		Precision:  res.Precision,
	}
}

func newParams(p nativeImageParams) *Params {
	return &Params{
		Width:      p.Width,
		Height:     p.Height,
		Components: p.Components,
		Precision:  p.Precision,
	}
}

// nativeDecodeResult avoids importing native in this file's doc - type alias in decode.go
type nativeDecodeResult struct {
	Output     []byte
	Width      int
	Height     int
	Components int
	Precision  int
}

type nativeImageParams struct {
	Width      int
	Height     int
	Components int
	Precision  int
}
