package golibjpeg

import (
	"fmt"
	"io"
	"os"
)

// ReadStream reads JPEG data from bytes, a file path, or an io.Reader.
// This mirrors pylibjpeg-libjpeg stream handling for decode() and get_parameters().
func ReadStream(stream any) ([]byte, error) {
	switch v := stream.(type) {
	case nil:
		return nil, fmt.Errorf(
			"invalid type 'nil' - must be the path to a JPEG file, a buffer containing the JPEG data or an open JPEG file-like",
		)
	case string:
		return os.ReadFile(v)
	case []byte:
		return v, nil
	case io.Reader:
		return io.ReadAll(v)
	default:
		return nil, fmt.Errorf(
			"invalid type '%T' - must be the path to a JPEG file, a buffer containing the JPEG data or an open JPEG file-like",
			stream,
		)
	}
}
