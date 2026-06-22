//go:build linux && arm64

package native

import _ "embed"

//go:embed libs/golibjpeg_linux_arm64.so
var libData []byte
