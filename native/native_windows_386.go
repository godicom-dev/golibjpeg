//go:build windows && 386

package native

import _ "embed"

//go:embed libs/golibjpeg_386.dll
var libData []byte
