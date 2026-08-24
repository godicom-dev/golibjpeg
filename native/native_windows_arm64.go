//go:build windows && arm64

package native

import _ "embed"

//go:embed libs/golibjpeg_arm64.dll
var libData []byte
