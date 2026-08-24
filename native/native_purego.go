//go:build !wasm

package native

import "github.com/ebitengine/purego"

// bindSymbols resolves the C entry points in a library loaded by loadLibrary.
// It panics on a missing required symbol; loadNative recovers and reports an
// error.
func bindSymbols(handle uintptr) {
	purego.RegisterLibFunc(&decodeFn, handle, "golibjpeg_decode")
	purego.RegisterLibFunc(&encodeFn, handle, "golibjpeg_encode")
	purego.RegisterLibFunc(&getParamsFn, handle, "golibjpeg_get_parameters")
	purego.RegisterLibFunc(&freeFn, handle, "golibjpeg_free")
	registerOptionalLibFunc(handle, "golibjpeg_last_error", &lastErrorFn)
}

// registerOptionalLibFunc binds a symbol that older builds of the library may
// not export. lastErrorDetail nil-checks the result, so a miss just costs the
// error detail.
func registerOptionalLibFunc(handle uintptr, name string, fn any) {
	defer func() {
		_ = recover()
	}()
	purego.RegisterLibFunc(fn, handle, name)
}
