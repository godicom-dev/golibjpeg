#ifndef ENCODE_HPP
#define ENCODE_HPP

#include <string>

struct EncodeParameters {
    unsigned int columns;
    unsigned int rows;
    unsigned int samples_per_pixel;
    unsigned int bits_per_sample;
    int frame_type;
    int colour_transform;
    int quality;
    int error_bound;
    int ls_interleaving;
};

extern std::string Encode(
    char *inArray,
    int inLength,
    const struct EncodeParameters *params,
    char **outArray,
    int *outLength
);

#endif
