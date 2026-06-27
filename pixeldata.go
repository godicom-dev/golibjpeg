package golibjpeg

import (
	"fmt"
)

// PixelDataVersion selects decode_pixel_data behaviour.
type PixelDataVersion int

const (
	PixelDataV1 PixelDataVersion = 1
	PixelDataV2 PixelDataVersion = 2
)

// PixelDataOptions configures DecodePixelData (pylibjpeg decode_pixel_data).
type PixelDataOptions struct {
	Version                   PixelDataVersion
	PhotometricInterpretation string
}

// DecodePixelData decodes encapsulated JPEG pixel data for DICOM handlers.
// Version 1 applies a colour transform from PhotometricInterpretation.
// Version 2 returns raw decoded bytes with no colour transform.
func DecodePixelData(src []byte, opts PixelDataOptions) ([]byte, error) {
	if opts.Version == 0 {
		opts.Version = PixelDataV1
	}

	switch opts.Version {
	case PixelDataV1:
		pi := opts.PhotometricInterpretation
		if pi == "" {
			return nil, fmt.Errorf(
				"The (0028,0004) Photometric Interpretation element is missing from the dataset",
			)
		}
		ct, ok := colourTransformForPhotometric(pi)
		if !ok {
			ct = ColourTransformNone
		}
		img, err := DecodeImage(src, ct)
		if err != nil {
			return nil, err
		}
		return img.Pixels, nil
	case PixelDataV2:
		img, err := DecodeImage(src, ColourTransformNone)
		if err != nil {
			return nil, err
		}
		return img.Pixels, nil
	default:
		return nil, fmt.Errorf("golibjpeg: unsupported pixel data version %d", opts.Version)
	}
}
