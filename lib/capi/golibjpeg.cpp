#include "golibjpeg.h"

#include "../interface/decode.hpp"

#include <cstdlib>
#include <cstring>
#include <string>

thread_local std::string g_golibjpeg_last_error;

namespace {

void set_last_error(const std::string& status)
{
    g_golibjpeg_last_error = status;
}

int parse_status_code(const std::string& status, int fallback)
{
    if (status.rfind("0::::", 0) == 0) {
        return GOLIBJPEG_OK;
    }

    const std::size_t sep = status.find("::::");
    if (sep == std::string::npos) {
        return fallback;
    }

    try {
        return std::stoi(status.substr(0, sep));
    } catch (...) {
        return fallback;
    }
}

int bytes_per_pixel(int precision)
{
    return (precision + 7) / 8;
}

int map_interface_code(int code)
{
    switch (code) {
    case GOLIBJPEG_OK:
        return GOLIBJPEG_OK;
    case -8192:
        return GOLIBJPEG_ERR_MEMORY;
    case -8194:
    case -8195:
        return GOLIBJPEG_ERR_PARAM;
    case -8193:
        return GOLIBJPEG_ERR_DECODE;
    default:
        if (code < -1000) {
            return code;
        }
        return GOLIBJPEG_ERR_DECODE;
    }
}

} // namespace

extern "C" {

GOLIBJPEG_EXPORT int golibjpeg_get_parameters(
    const unsigned char* data,
    int data_len,
    int* width,
    int* height,
    int* components,
    int* precision)
{
    if (data == nullptr || data_len <= 0 || width == nullptr || height == nullptr ||
        components == nullptr || precision == nullptr) {
        return GOLIBJPEG_ERR_PARAM;
    }

    JPEGParameters param = {};
    const std::string status = GetJPEGParameters(
        const_cast<char*>(reinterpret_cast<const char*>(data)),
        data_len,
        &param
    );

    const int code = map_interface_code(parse_status_code(status, GOLIBJPEG_ERR_DECODE));
    if (code != GOLIBJPEG_OK) {
        set_last_error(status);
        return code;
    }

    *width = static_cast<int>(param.columns);
    *height = static_cast<int>(param.rows);
    *components = static_cast<int>(param.samples_per_pixel);
    *precision = static_cast<int>(param.bits_per_sample);
    g_golibjpeg_last_error.clear();
    return GOLIBJPEG_OK;
}

GOLIBJPEG_EXPORT int golibjpeg_decode(
    const unsigned char* data,
    int data_len,
    int colour_transform,
    unsigned char** output,
    int* output_len,
    int* width,
    int* height,
    int* components,
    int* precision)
{
    if (data == nullptr || data_len <= 0 || output == nullptr || output_len == nullptr ||
        width == nullptr || height == nullptr || components == nullptr || precision == nullptr) {
        return GOLIBJPEG_ERR_PARAM;
    }

    *output = nullptr;
    *output_len = 0;

    if (colour_transform < 0 || colour_transform > 3) {
        return GOLIBJPEG_ERR_PARAM;
    }

    JPEGParameters param = {};
    std::string status = GetJPEGParameters(
        const_cast<char*>(reinterpret_cast<const char*>(data)),
        data_len,
        &param
    );

    int code = map_interface_code(parse_status_code(status, GOLIBJPEG_ERR_DECODE));
    if (code != GOLIBJPEG_OK) {
        set_last_error(status);
        return code;
    }

    if (param.columns == 0 || param.rows == 0 || param.samples_per_pixel == 0) {
        return GOLIBJPEG_ERR_FORMAT;
    }

    const int bpp = bytes_per_pixel(param.bits_per_sample);
    const int out_len = static_cast<int>(
        param.columns * param.rows * param.samples_per_pixel * bpp
    );
    if (out_len <= 0) {
        return GOLIBJPEG_ERR_FORMAT;
    }

    auto* out_buf = static_cast<unsigned char*>(std::malloc(static_cast<std::size_t>(out_len)));
    if (out_buf == nullptr) {
        return GOLIBJPEG_ERR_MEMORY;
    }

    status = Decode(
        const_cast<char*>(reinterpret_cast<const char*>(data)),
        reinterpret_cast<char*>(out_buf),
        data_len,
        out_len,
        colour_transform
    );

    code = map_interface_code(parse_status_code(status, GOLIBJPEG_ERR_DECODE));
    if (code != GOLIBJPEG_OK) {
        std::free(out_buf);
        set_last_error(status);
        return code;
    }

    g_golibjpeg_last_error.clear();

    *output = out_buf;
    *output_len = out_len;
    *width = static_cast<int>(param.columns);
    *height = static_cast<int>(param.rows);
    *components = static_cast<int>(param.samples_per_pixel);
    *precision = static_cast<int>(param.bits_per_sample);
    g_golibjpeg_last_error.clear();
    return GOLIBJPEG_OK;
}

GOLIBJPEG_EXPORT void golibjpeg_free(unsigned char* p)
{
    std::free(p);
}

GOLIBJPEG_EXPORT const char* golibjpeg_last_error(void)
{
    return g_golibjpeg_last_error.c_str();
}

} // extern "C"
