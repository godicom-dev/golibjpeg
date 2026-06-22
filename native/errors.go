package native

import "fmt"

var errEmptyInput = fmt.Errorf("golibjpeg: empty input data")

var errEmptyOutput = fmt.Errorf("golibjpeg: decode returned empty output")

func errWithCode(code int32) error {
	switch code {
	case -1:
		return fmt.Errorf("golibjpeg: memory allocation failed")
	case -2:
		return fmt.Errorf("golibjpeg: decode error")
	case -3:
		return fmt.Errorf("golibjpeg: unsupported format")
	case -4:
		return fmt.Errorf("golibjpeg: I/O error")
	case -5:
		return fmt.Errorf("golibjpeg: invalid parameter")
	default:
		if code < -1000 {
			// libjpeg internal error codes (e.g. -1024..-2046)
			return fmt.Errorf("golibjpeg: libjpeg error (code %d)", code)
		}
		return fmt.Errorf("golibjpeg: unknown error (code %d)", code)
	}
}
