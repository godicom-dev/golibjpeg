package golibjpeg

type Format int

const (
	FormatAuto    Format = -1
	FormatJPEG    Format = 1  // CT_YCbCr
	FormatJPEGLS  Format = 2  // CT_RCT
	FormatJPEGXT  Format = 3  // CT_Freeform
)
