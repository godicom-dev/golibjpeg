package native

import (
	"os"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
)

var (
	decodeFn func(data unsafe.Pointer, dataLen int32,
		colourTransform int32,
		output *unsafe.Pointer, outputLen *int32,
		width, height, components, precision *int32) int32
	getParamsFn func(data unsafe.Pointer, dataLen int32,
		width, height, components, precision *int32) int32
	freeFn func(p unsafe.Pointer)
	lastErrorFn func() uintptr
)

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

func init() {
	f, err := os.CreateTemp("", "golibjpeg-*."+libExt())
	if err != nil {
		panic("golibjpeg: failed to create temp file: " + err.Error())
	}
	path := f.Name()
	_ = f.Close()
	handle, err := extractAndLoad(path)
	if err != nil {
		panic("golibjpeg: failed to load native library: " + err.Error())
	}
	purego.RegisterLibFunc(&decodeFn, uintptr(handle), "golibjpeg_decode")
	purego.RegisterLibFunc(&getParamsFn, uintptr(handle), "golibjpeg_get_parameters")
	purego.RegisterLibFunc(&freeFn, uintptr(handle), "golibjpeg_free")
	registerOptionalLibFunc(uintptr(handle), "golibjpeg_last_error", &lastErrorFn)
}

func registerOptionalLibFunc(handle uintptr, name string, fn any) {
	defer func() {
		recover()
	}()
	purego.RegisterLibFunc(fn, handle, name)
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
