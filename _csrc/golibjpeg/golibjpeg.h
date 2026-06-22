#ifndef GOLIBJPEG_H
#define GOLIBJPEG_H

#include <stddef.h>

#ifdef _WIN32
#define GOLIBJPEG_EXPORT __declspec(dllexport)
#else
#define GOLIBJPEG_EXPORT __attribute__((visibility("default")))
#endif

#ifdef __cplusplus
extern "C" {
#endif

#define GOLIBJPEG_OK           0
#define GOLIBJPEG_ERR_MEMORY   -1
#define GOLIBJPEG_ERR_DECODE   -2
#define GOLIBJPEG_ERR_FORMAT   -3
#define GOLIBJPEG_ERR_IO       -4
#define GOLIBJPEG_ERR_PARAM    -5

/// Colour transforms matching JPGFLAG_MATRIX_COLORTRANSFORMATION_*
#define GOLIBJPEG_CT_NONE    0
#define GOLIBJPEG_CT_YCBCR   1
#define GOLIBJPEG_CT_RCT     2
#define GOLIBJPEG_CT_FREEFORM 3

/// Decode JPEG data and return raw pixel data in native precision (8 or 16 bit).
/// Returns GOLIBJPEG_OK on success, negative on error.
GOLIBJPEG_EXPORT int golibjpeg_decode(
    const unsigned char* data, int data_len,
    int colour_transform,
    unsigned char** output, int* output_len,
    int* width, int* height, int* components, int* precision
);

/// Convenience: decode with YCbCr→RGB colour transform (default for JPEG baseline).
GOLIBJPEG_EXPORT int golibjpeg_decode_to_rgb(
    const unsigned char* data, int data_len,
    unsigned char** output, int* output_len,
    int* width, int* height
);

/// Get JPEG image parameters without decoding pixel data.
GOLIBJPEG_EXPORT int golibjpeg_get_parameters(
    const unsigned char* data, int data_len,
    int* width, int* height, int* components, int* precision
);

/// Free output buffer allocated by decode.
GOLIBJPEG_EXPORT void golibjpeg_free(unsigned char* p);

#ifdef __cplusplus
}
#endif

#endif
