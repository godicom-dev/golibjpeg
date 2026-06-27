package native

import (
	"fmt"
	"strconv"
	"strings"
)

var errEmptyInput = fmt.Errorf("golibjpeg: empty input data")

var errEmptyOutput = fmt.Errorf("golibjpeg: decode returned empty output")

func errWithCode(op string, code int32) error {
	detail := lastErrorDetail()
	if detail != "" {
		if msg := detailMessage(detail); msg != "" {
			return &StatusError{Op: op, Code: int(code), Detail: msg}
		}
	}
	return &StatusError{Op: op, Code: int(code)}
}

func detailMessage(status string) string {
	parts := strings.SplitN(status, "::::", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

func statusCode(status string) (int, bool) {
	parts := strings.SplitN(status, "::::", 2)
	if len(parts) == 0 {
		return 0, false
	}
	code, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, false
	}
	return code, true
}

// StatusError mirrors pylibjpeg-libjpeg RuntimeError messages.
type StatusError struct {
	Op     string
	Code   int
	Detail string
}

func (e *StatusError) Error() string {
	if e.Code == 0 {
		return fmt.Sprintf("golibjpeg: %s failed", e.Op)
	}

	if known, ok := libjpegErrorMessages[e.Code]; ok {
		if e.Detail != "" {
			return fmt.Sprintf(
				"libjpeg error code '%d' returned from %s(): %s - %s",
				e.Code, e.Op, known, e.Detail,
			)
		}
		return fmt.Sprintf(
			"libjpeg error code '%d' returned from %s(): %s",
			e.Code, e.Op, known,
		)
	}

	if e.Detail != "" {
		return fmt.Sprintf(
			"Unknown error code '%d' returned from %s(): %s",
			e.Code, e.Op, e.Detail,
		)
	}
	return fmt.Sprintf("Unknown error code '%d' returned from %s()", e.Code, e.Op)
}

var libjpegErrorMessages = map[int]string{
	-1:    "memory allocation failed",
	-2:    "decode error",
	-3:    "unsupported format",
	-4:    "I/O error",
	-5:    "invalid parameter",
	-1024: "A parameter for a function was out of range",
	-1025: "Stream run out of data",
	-1026: "A code block run out of data",
	-1027: "Tried to perform an unputc or or an unget on an empty stream",
	-1028: "Some parameter run out of range",
	-1029: "The requested operation does not apply",
	-1030: "Tried to create an already existing object",
	-1031: "Tried to access a non-existing object",
	-1032: "A non-optional parameter was left out",
	-1033: "Forgot to delay a 0xFF",
	-1034: "Internal error: the requested operation is not available",
	-1035: "Internal error: an item computed on a former pass does not coincide with the same item on a later pass",
	-1036: "The stream passed in is no valid jpeg stream",
	-1037: "A unique marker turned up more than once. The input stream is most likely corrupt",
	-1038: "A misplaced marker segment was found",
	-1040: "The specified parameters are valid, but are not supported by the selected profile. Either use a higher profile, or use simpler parameters (encoder only). ",
	-1041: "Internal error: the worker thread that was currently active had to terminate unexpectedly",
	-1042: "The encoder tried to emit a symbol for which no Huffman code was defined. This happens if the standard Huffman table is used for an alphabet for which it was not defined. The reaction to this exception should be to create a custom huffman table instead",
	-2046: "Failed to construct the JPEG object",
	-8194: "Invalid colourTransform value",
}
