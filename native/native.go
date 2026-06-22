package native

import (
	"os"
	"path/filepath"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"
	"golang.org/x/sys/windows"
)

var (
	decodeFn func(data unsafe.Pointer, dataLen int32,
		colourTransform int32,
		output *unsafe.Pointer, outputLen *int32,
		width, height, components, precision *int32) int32
	getParamsFn func(data unsafe.Pointer, dataLen int32,
		width, height, components, precision *int32) int32
	freeFn func(p unsafe.Pointer)
)

func init() {
	tmp := filepath.Join(os.TempDir(), "golibjpeg."+libExt())
	if err := os.WriteFile(tmp, libData, 0755); err != nil {
		panic("golibjpeg: failed to extract native library: " + err.Error())
	}

	handle, err := windows.LoadLibrary(tmp)
	if err != nil {
		panic("golibjpeg: failed to load native library: " + err.Error())
	}

	purego.RegisterLibFunc(&decodeFn, uintptr(handle), "golibjpeg_decode")
	purego.RegisterLibFunc(&getParamsFn, uintptr(handle), "golibjpeg_get_parameters")
	purego.RegisterLibFunc(&freeFn, uintptr(handle), "golibjpeg_free")

	if runtime.GOOS != "windows" {
		os.Remove(tmp)
	}
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
