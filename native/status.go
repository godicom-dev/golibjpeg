package native

import "unsafe"

func lastErrorDetail() string {
	if lastErrorFn == nil {
		return ""
	}
	ptr := lastErrorFn()
	if ptr == 0 {
		return ""
	}
	return string(cString(ptr))
}

func cString(ptr uintptr) []byte {
	if ptr == 0 {
		return nil
	}
	var length int
	for {
		if *(*byte)(unsafe.Pointer(ptr + uintptr(length))) == 0 {
			break
		}
		length++
	}
	if length == 0 {
		return nil
	}
	buf := make([]byte, length)
	copy(buf, unsafe.Slice((*byte)(unsafe.Pointer(ptr)), length))
	return buf
}
