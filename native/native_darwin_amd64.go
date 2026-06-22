//go:build darwin && amd64

package native

import _ "embed"

//go:embed libs/golibjpeg_darwin_amd64.dylib
var libData []byte
