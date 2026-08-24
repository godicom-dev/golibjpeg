//go:build wasm

package native

// wasm has no dynamic loader, and purego does not provide RegisterLibFunc there,
// so neither of these can be implemented. Nothing calls them: libData is empty
// on wasm, so loadNative returns ErrUnsupportedPlatform before it reaches either.
// They exist so that the module still compiles for js/wasm and wasip1/wasm —
// everything in godicom that does not decode JPEG or JPEG-LS works there.

func loadLibrary(string) (uintptr, error) {
	return 0, ErrUnsupportedPlatform
}

func bindSymbols(uintptr) {}
