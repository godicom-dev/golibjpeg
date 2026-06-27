package golibjpeg

// Format selects a colour transform for DecodeWithFormat.
// Deprecated: prefer ColourTransform with DecodeImage.
type Format int

const (
	FormatAuto   Format = -1
	FormatJPEG   Format = 1
	FormatJPEGLS Format = 2
	FormatJPEGXT Format = 3
)

func (f Format) colourTransform() ColourTransform {
	switch f {
	case FormatJPEG:
		return ColourTransformYCbCr
	case FormatJPEGLS:
		return ColourTransformRCT
	case FormatJPEGXT:
		return ColourTransformFreeform
	default:
		return ColourTransformNone
	}
}
