#include "encode.hpp"

#include "decode.hpp"

#include "../libjpeg/cmd/bitmaphook.hpp"
#include "../libjpeg/interface/hooks.hpp"
#include "../libjpeg/interface/jpeg.hpp"
#include "../libjpeg/interface/parameters.hpp"
#include "../libjpeg/interface/tagitem.hpp"
#include "../libjpeg/interface/types.hpp"
#include "../libjpeg/tools/traits.hpp"

#include "std/stdio.hpp"
#include "std/stdlib.hpp"
#include "std/string.hpp"

#include <cstring>
#include <sstream>

namespace {

struct OutputBuffer {
    unsigned char *data;
    int length;
    int capacity;
};

struct EncodeStreamData {
    const unsigned char *pixels;
    int width;
    int height;
    int depth;
    int bytes_per_pixel;
    unsigned char *stripe;
    int stripe_bytes;
};

bool output_buffer_init(OutputBuffer *buf, int initial_capacity)
{
    buf->data = static_cast<unsigned char *>(std::malloc(static_cast<std::size_t>(initial_capacity)));
    if (buf->data == nullptr) {
        buf->length = 0;
        buf->capacity = 0;
        return false;
    }
    buf->length = 0;
    buf->capacity = initial_capacity;
    return true;
}

void output_buffer_free(OutputBuffer *buf)
{
    std::free(buf->data);
    buf->data = nullptr;
    buf->length = 0;
    buf->capacity = 0;
}

bool output_buffer_grow(OutputBuffer *buf, int needed)
{
    if (needed <= buf->capacity) {
        return true;
    }
    int new_capacity = buf->capacity > 0 ? buf->capacity : 4096;
    while (new_capacity < needed) {
        new_capacity *= 2;
    }
    unsigned char *next = static_cast<unsigned char *>(
        std::realloc(buf->data, static_cast<std::size_t>(new_capacity))
    );
    if (next == nullptr) {
        return false;
    }
    buf->data = next;
    buf->capacity = new_capacity;
    return true;
}

JPG_LONG MemIOHook(struct JPG_Hook *hook, struct JPG_TagItem *tags)
{
    OutputBuffer *out = static_cast<OutputBuffer *>(hook->hk_pData);

    switch (tags->GetTagData(JPGTAG_FIO_ACTION)) {
    case JPGFLAG_ACTION_READ:
    {
        UBYTE *buffer = static_cast<UBYTE *>(tags->GetTagPtr(JPGTAG_FIO_BUFFER));
        ULONG size = static_cast<ULONG>(tags->GetTagData(JPGTAG_FIO_SIZE));
        ULONG available = static_cast<ULONG>(out->length);
        ULONG offset = 0;
        if (available == 0) {
            return 0;
        }
        if (size > available) {
            size = available;
        }
        std::memcpy(buffer, out->data + offset, size);
        return static_cast<JPG_LONG>(size);
    }
    case JPGFLAG_ACTION_WRITE:
    {
        UBYTE *buffer = static_cast<UBYTE *>(tags->GetTagPtr(JPGTAG_FIO_BUFFER));
        ULONG size = static_cast<ULONG>(tags->GetTagData(JPGTAG_FIO_SIZE));
        int needed = out->length + static_cast<int>(size);
        if (!output_buffer_grow(out, needed)) {
            return -1;
        }
        std::memcpy(out->data + out->length, buffer, size);
        out->length += static_cast<int>(size);
        return static_cast<JPG_LONG>(size);
    }
    case JPGFLAG_ACTION_SEEK:
    {
        LONG mode = tags->GetTagData(JPGTAG_FIO_SEEKMODE);
        LONG offset = tags->GetTagData(JPGTAG_FIO_OFFSET);
        switch (mode) {
        case JPGFLAG_OFFSET_CURRENT:
            out->length += offset;
            if (out->length < 0) {
                out->length = 0;
                return -1;
            }
            if (out->length > out->capacity && !output_buffer_grow(out, out->length)) {
                return -1;
            }
            return 0;
        case JPGFLAG_OFFSET_BEGINNING:
            out->length = offset;
            if (out->length < 0) {
                out->length = 0;
                return -1;
            }
            if (out->length > out->capacity && !output_buffer_grow(out, out->length)) {
                return -1;
            }
            return 0;
        case JPGFLAG_OFFSET_END:
            out->length = out->capacity + offset;
            if (out->length < 0) {
                out->length = 0;
                return -1;
            }
            if (out->length > out->capacity && !output_buffer_grow(out, out->length)) {
                return -1;
            }
            return 0;
        }
        return -1;
    }
    case JPGFLAG_ACTION_QUERY:
        return 0;
    }
    return -1;
}

int bytes_per_sample(int precision)
{
    return (precision + 7) / 8;
}

int expected_input_length(const struct EncodeParameters *params)
{
    const int bpp = bytes_per_sample(static_cast<int>(params->bits_per_sample));
    return static_cast<int>(params->columns * params->rows * params->samples_per_pixel * bpp);
}

JPG_LONG EncodeBitmapHook(struct JPG_Hook *hook, struct JPG_TagItem *tags)
{
    struct EncodeStreamData *enc = static_cast<struct EncodeStreamData *>(hook->hk_pData);
    struct BitmapMemory bmm;
    std::memset(&bmm, 0, sizeof(bmm));

    bmm.bmm_pMemPtr = enc->stripe;
    bmm.bmm_ulWidth = static_cast<ULONG>(enc->width);
    bmm.bmm_ulHeight = static_cast<ULONG>(enc->height);
    bmm.bmm_usDepth = static_cast<UWORD>(enc->depth);
    bmm.bmm_ucPixelType = (enc->bytes_per_pixel > 1) ? CTYP_UWORD : CTYP_UBYTE;
    bmm.bmm_bUpsampling = true;

    UWORD comp = tags->GetTagData(JPGTAG_BIO_COMPONENT);
    ULONG miny = tags->GetTagData(JPGTAG_BIO_MINY);
    ULONG maxy = tags->GetTagData(JPGTAG_BIO_MAXY);
    ULONG width = 1 + tags->GetTagData(JPGTAG_BIO_MAXX);

    if (tags->GetTagData(JPGTAG_BIO_ACTION) == JPGFLAG_BIO_REQUEST && comp == 0) {
        ULONG height = maxy + 1 - miny;
        if (height > 8) {
            height = 8;
        }
        const int row_bytes = enc->width * enc->depth * enc->bytes_per_pixel;
        const unsigned char *src = enc->pixels + static_cast<int>(miny) * row_bytes;
        std::memcpy(enc->stripe, src, static_cast<std::size_t>(height) * row_bytes);
    }

    struct JPG_Hook inner(BitmapHook, &bmm);
    return BitmapHook(&inner, tags);
}

} // namespace

std::string Encode(
    char *inArray,
    int inLength,
    const struct EncodeParameters *params,
    char **outArray,
    int *outLength)
{
    if (inArray == nullptr || params == nullptr || outArray == nullptr || outLength == nullptr) {
        return "-8194::::Invalid encode parameters";
    }

    *outArray = nullptr;
    *outLength = 0;

    if (params->columns == 0 || params->rows == 0 || params->samples_per_pixel == 0 ||
        params->bits_per_sample == 0) {
        return "-8194::::Invalid image dimensions";
    }

    const int expected = expected_input_length(params);
    if (inLength != expected) {
        return "-8195::::Invalid input array size";
    }

    if (params->colour_transform < 0 || params->colour_transform > 3) {
        return "-8194::::Invalid colourTransform value";
    }

    int frame_type = params->frame_type;
    int colour_transform = params->colour_transform;
    if (params->samples_per_pixel == 1) {
        colour_transform = JPGFLAG_MATRIX_COLORTRANSFORMATION_NONE;
    }

    if (frame_type == JPGFLAG_LOSSLESS || frame_type == JPGFLAG_JPEG_LS) {
        colour_transform = JPGFLAG_MATRIX_COLORTRANSFORMATION_NONE;
    }

    int quality = params->quality;
    if (quality <= 0) {
        quality = 75;
    }
    if (quality > 100) {
        quality = 100;
    }

    if (frame_type == JPGFLAG_LOSSLESS || frame_type == JPGFLAG_JPEG_LS) {
        quality = 100;
    }

    int ls_interleaving = params->ls_interleaving;
    if (frame_type == JPGFLAG_JPEG_LS && ls_interleaving < 0) {
        ls_interleaving = JPGFLAG_SCAN_LS_INTERLEAVING_SAMPLE;
    }

    int residual_frame_type = JPGFLAG_RESIDUAL;
    if (frame_type == JPGFLAG_LOSSLESS || frame_type == JPGFLAG_JPEG_LS) {
        residual_frame_type = JPGFLAG_RESIDUAL;
    } else if (quality >= 100) {
        residual_frame_type = JPGFLAG_RESIDUAL;
    } else {
        residual_frame_type = JPGFLAG_SEQUENTIAL;
    }

    UBYTE subx[4] = {1, 1, 1, 1};
    UBYTE suby[4] = {1, 1, 1, 1};

    const int bpp = bytes_per_sample(static_cast<int>(params->bits_per_sample));
    const int stripe_bytes = static_cast<int>(params->columns * 8 * params->samples_per_pixel * bpp);
    auto *stripe = static_cast<unsigned char *>(std::malloc(static_cast<std::size_t>(stripe_bytes)));
    if (stripe == nullptr) {
        return "-8192::::Unable to allocate memory to buffer the image";
    }

    EncodeStreamData stream = {
        reinterpret_cast<const unsigned char *>(inArray),
        static_cast<int>(params->columns),
        static_cast<int>(params->rows),
        static_cast<int>(params->samples_per_pixel),
        bpp,
        stripe,
        stripe_bytes,
    };

    struct JPG_Hook bmhook(EncodeBitmapHook, &stream);

    struct JPG_TagItem tags[] = {
        JPG_PointerTag(JPGTAG_BIH_HOOK, &bmhook),
        JPG_ValueTag(JPGTAG_ENCODER_LOOP_ON_INCOMPLETE, true),
        JPG_ValueTag(JPGTAG_IMAGE_WIDTH, params->columns),
        JPG_ValueTag(JPGTAG_IMAGE_HEIGHT, params->rows),
        JPG_ValueTag(JPGTAG_IMAGE_DEPTH, params->samples_per_pixel),
        JPG_ValueTag(JPGTAG_IMAGE_PRECISION, params->bits_per_sample),
        JPG_ValueTag(JPGTAG_IMAGE_FRAMETYPE, frame_type),
        JPG_ValueTag(JPGTAG_RESIDUAL_FRAMETYPE, residual_frame_type),
        JPG_ValueTag(JPGTAG_IMAGE_QUALITY, quality),
        JPG_ValueTag(JPGTAG_QUANTIZATION_MATRIX, JPGFLAG_QUANTIZATION_ANNEX_K),
        JPG_ValueTag(JPGTAG_RESIDUALQUANT_MATRIX, JPGFLAG_QUANTIZATION_ANNEX_K),
        JPG_ValueTag(JPGTAG_IMAGE_ERRORBOUND, params->error_bound),
        JPG_ValueTag(JPGTAG_MATRIX_LTRAFO, colour_transform),
        JPG_PointerTag(JPGTAG_IMAGE_SUBX, subx),
        JPG_PointerTag(JPGTAG_IMAGE_SUBY, suby),
        JPG_ValueTag(
            (frame_type == JPGFLAG_JPEG_LS) ? JPGTAG_SCAN_LS_INTERLEAVING : JPGTAG_TAG_IGNORE,
            ls_interleaving
        ),
        JPG_ValueTag(JPGTAG_IMAGE_IS_FLOAT, false),
        JPG_ValueTag(JPGTAG_IMAGE_OUTPUT_CONVERSION, false),
        JPG_EndTag
    };

    OutputBuffer output;
    if (!output_buffer_init(&output, expected > 0 ? expected : 4096)) {
        std::free(stripe);
        return "-8192::::Unable to allocate memory for encoded output";
    }

    class JPEG *jpeg = JPEG::Construct(NULL);
    if (jpeg == nullptr) {
        output_buffer_free(&output);
        std::free(stripe);
        return "-8193::::Failed to construct the JPEG object";
    }

    int ok = jpeg->ProvideImage(tags);
    if (ok) {
        struct JPG_Hook iohook(MemIOHook, &output);
        struct JPG_TagItem iotags[] = {
            JPG_PointerTag(JPGTAG_HOOK_IOHOOK, &iohook),
            JPG_PointerTag(JPGTAG_HOOK_IOSTREAM, &output),
            JPG_EndTag
        };
        ok = jpeg->Write(iotags);
    }

    std::free(stripe);

    if (!ok || output.length <= 0) {
        const char *error = nullptr;
        int code = jpeg->LastError(error);
        output_buffer_free(&output);
        JPEG::Destruct(jpeg);
        std::ostringstream status;
        status << code << "::::" << (error != nullptr ? error : "encoding failed");
        return status.str();
    }

    JPEG::Destruct(jpeg);

    *outArray = reinterpret_cast<char *>(output.data);
    *outLength = output.length;
    return "0::::";
}
