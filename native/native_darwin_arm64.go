//go:build darwin && arm64

package native

import _ "embed"

//go:embed libs/golibjpeg_darwin_arm64.dylib
var libData []byte
