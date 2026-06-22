//go:build windows && amd64

package native

import _ "embed"

//go:embed libs/golibjpeg_amd64.dll
var libData []byte
