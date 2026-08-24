package native

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"unsafe"
)

var (
	decodeFn func(data unsafe.Pointer, dataLen int32,
		colourTransform int32,
		output *unsafe.Pointer, outputLen *int32,
		width, height, components, precision *int32) int32
	encodeFn func(src unsafe.Pointer, srcLen int32, params unsafe.Pointer,
		output *unsafe.Pointer, outputLen *int32) int32
	getParamsFn func(data unsafe.Pointer, dataLen int32,
		width, height, components, precision *int32) int32
	freeFn      func(p unsafe.Pointer)
	lastErrorFn func() uintptr
)

// ensureLoaded extracts the embedded library and binds its entry points, once
// per process, and reports the same result to every later caller.
//
// The work is deferred to the first call instead of running in init() so that
// failure is an error the caller can handle rather than a panic during package
// initialisation. Two failures are worth handling: a platform with no prebuilt
// library, and a host whose temporary directory is read-only or mounted noexec.
// Importing this package must stay safe in both cases, because anything that
// imports golibjpeg — godicom's pixels package, for one — would otherwise fail
// to start even when it only ever touches uncompressed pixel data.
var ensureLoaded = sync.OnceValue(loadNative)

func loadNative() (err error) {
	defer func() {
		// bindSymbols panics on a missing required symbol, which would otherwise
		// surface from whichever call happened to be the first.
		if r := recover(); r != nil {
			err = fmt.Errorf("golibjpeg: binding the native library failed: %v", r)
		}
	}()

	if len(libData) == 0 {
		return fmt.Errorf("%w (%s/%s)", ErrUnsupportedPlatform, runtime.GOOS, runtime.GOARCH)
	}

	f, err := os.CreateTemp("", "golibjpeg-*."+libExt())
	if err != nil {
		return fmt.Errorf("golibjpeg: creating a temp file for the native library failed: %w", err)
	}
	path := f.Name()
	_ = f.Close()

	handle, err := extractAndLoad(path)
	if err != nil {
		return fmt.Errorf("golibjpeg: loading the native library failed: %w", err)
	}

	bindSymbols(handle)
	return nil
}

func extractAndLoad(path string) (uintptr, error) {
	if err := os.WriteFile(path, libData, 0o755); err != nil {
		return 0, err
	}
	handle, err := loadLibrary(path)
	if err != nil {
		_ = os.Remove(path)
		return 0, err
	}
	if runtime.GOOS != "windows" {
		_ = os.Remove(path)
	}
	return handle, nil
}

func libExt() string {
	switch runtime.GOOS {
	case "windows":
		return "dll"
	case "darwin":
		return "dylib"
	default:
		return "so"
	}
}
