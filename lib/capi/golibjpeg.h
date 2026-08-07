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
#define GOLIBJPEG_ERR_ENCODE   -6

#define GOLIBJPEG_CT_NONE     0
#define GOLIBJPEG_FRAME_BASELINE   0
#define GOLIBJPEG_FRAME_SEQUENTIAL 1
#define GOLIBJPEG_FRAME_LOSSLESS   3
#define GOLIBJPEG_FRAME_JPEG_LS    4

typedef struct golibjpeg_encode_params {
    int columns;
    int rows;
    int samples_per_pixel;
    int bits_per_sample;
    int frame_type;
    int colour_transform;
    int quality;
    int error_bound;
    int ls_interleaving;
} golibjpeg_encode_params;
#define GOLIBJPEG_CT_YCBCR    1
#define GOLIBJPEG_CT_RCT      2
#define GOLIBJPEG_CT_FREEFORM 3

GOLIBJPEG_EXPORT int golibjpeg_decode(
    const unsigned char* data,
    int data_len,
    int colour_transform,
    unsigned char** output,
    int* output_len,
    int* width,
    int* height,
    int* components,
    int* precision
);

GOLIBJPEG_EXPORT int golibjpeg_get_parameters(
    const unsigned char* data,
    int data_len,
    int* width,
    int* height,
    int* components,
    int* precision
);

GOLIBJPEG_EXPORT int golibjpeg_encode(
    const unsigned char* src,
    int src_len,
    const golibjpeg_encode_params* params,
    unsigned char** output,
    int* output_len
);

GOLIBJPEG_EXPORT void golibjpeg_free(unsigned char* p);

GOLIBJPEG_EXPORT const char* golibjpeg_last_error(void);

#ifdef __cplusplus
}
#endif

#endif
