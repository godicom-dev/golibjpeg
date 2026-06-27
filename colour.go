package golibjpeg

// ColourTransform matches libjpeg JPGFLAG_MATRIX_COLORTRANSFORMATION_* values
// used by pylibjpeg-libjpeg decode().
type ColourTransform int

const (
	ColourTransformNone     ColourTransform = 0
	ColourTransformYCbCr    ColourTransform = 1
	ColourTransformRCT      ColourTransform = 2
	ColourTransformFreeform ColourTransform = 3
)

// PhotometricInterpretation values supported by DecodePixelData (v1).
const (
	PhotometricMonochrome1  = "MONOCHROME1"
	PhotometricMonochrome2  = "MONOCHROME2"
	PhotometricRGB          = "RGB"
	PhotometricYBRFull      = "YBR_FULL"
	PhotometricYBRFull422   = "YBR_FULL_422"
)

var photometricColourTransform = map[string]ColourTransform{
	PhotometricMonochrome1: ColourTransformNone,
	PhotometricMonochrome2: ColourTransformNone,
	PhotometricRGB:         ColourTransformYCbCr,
	PhotometricYBRFull:     ColourTransformNone,
	PhotometricYBRFull422:  ColourTransformNone,
}

func colourTransformForPhotometric(pi string) (ColourTransform, bool) {
	ct, ok := photometricColourTransform[pi]
	return ct, ok
}
