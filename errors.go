package golibjpeg

import "github.com/godicom-dev/golibjpeg/native"

// ErrUnsupportedPlatform reports that this GOOS/GOARCH has no prebuilt native
// library, so JPEG and JPEG-LS data can be neither decoded nor encoded here.
// Every function in this package returns an error wrapping it rather than
// panicking, which keeps the module importable everywhere Go builds. Test for it
// with errors.Is. The README lists the platforms that do have a library.
var ErrUnsupportedPlatform = native.ErrUnsupportedPlatform
