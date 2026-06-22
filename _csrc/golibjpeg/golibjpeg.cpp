#include "golibjpeg.h"

#include "interface/jpeg.hpp"
#include "interface/hooks.hpp"
#include "interface/tagitem.hpp"
#include "interface/parameters.hpp"
#include "tools/traits.hpp"
#include "tools/environment.hpp"

#include <cstring>
#include <cstdlib>
#include <new>
#include <cassert>

#define JPG_TRUE 1

// ---- input stream ----

struct MemStream {
    const unsigned char* data;
    int                  size;
    int                  pos;
};

static LONG input_hook(struct JPG_Hook* hook, struct JPG_TagItem* tags)
{
    auto* in = static_cast<MemStream*>(hook->hk_pData);
    if (!in) return -1;

    LONG action = tags->GetTagData(JPGTAG_FIO_ACTION);

    switch (action) {
    case JPGFLAG_ACTION_READ: {
        UBYTE* buffer = (UBYTE*)tags->GetTagPtr(JPGTAG_FIO_BUFFER);
        LONG   size   = tags->GetTagData(JPGTAG_FIO_SIZE);
        LONG   read   = 0;
        for (LONG i = 0; i < size && in->pos < in->size; i++) {
            *buffer++ = in->data[in->pos++];
            read++;
        }
        tags->SetTagData(JPGTAG_FIO_SIZE, read);
        return read;
    }
    case JPGFLAG_ACTION_SEEK: {
        LONG mode   = tags->GetTagData(JPGTAG_FIO_SEEKMODE);
        LONG offset = tags->GetTagData(JPGTAG_FIO_OFFSET);
        switch (mode) {
        case JPGFLAG_OFFSET_BEGINNING: in->pos = offset; break;
        case JPGFLAG_OFFSET_CURRENT:   in->pos += offset; break;
        case JPGFLAG_OFFSET_END:       in->pos = in->size + offset; break;
        }
        if (in->pos < 0) in->pos = 0;
        if (in->pos > in->size) in->pos = in->size;
        return 0;
    }
    case JPGFLAG_ACTION_QUERY: {
        LONG remaining = in->size - in->pos;
        tags->SetTagData(JPGTAG_FIO_SIZE, remaining);
        return remaining > 0 ? 0 : -1;
    }
    case JPGFLAG_ACTION_WRITE:
        return -1;
    }
    return -1;
}

// ---- stripe output state ----

struct StripeOut {
    unsigned char* output;      // final output buffer (malloc'd, caller-owned)
    int            output_len;  // total bytes in output
    int            pos;         // write position in output

    UBYTE*   stripe;       // temp buffer: width * 8 * depth * bpp
    int      width;
    int      height;
    int      depth;        // components
    UBYTE    pixel_type;   // CTYP_UBYTE or CTYP_UWORD
    int      bpp;          // bytes per pixel per component
    bool     big_endian;
    bool     upsample;
};

static LONG output_hook(struct JPG_Hook* hook, struct JPG_TagItem* tags)
{
    auto* so = static_cast<StripeOut*>(hook->hk_pData);
    if (!so) return 0;

    LONG action = tags->GetTagData(JPGTAG_BIO_ACTION);
    LONG comp   = tags->GetTagData(JPGTAG_BIO_COMPONENT);

    ULONG miny = tags->GetTagData(
        so->upsample ? JPGTAG_BIO_MINY : JPGTAG_BIO_PIXEL_MINY);
    ULONG maxy = tags->GetTagData(
        so->upsample ? JPGTAG_BIO_MAXY : JPGTAG_BIO_PIXEL_MAXY);
    ULONG stripe_w = 1 + tags->GetTagData(
        so->upsample ? JPGTAG_BIO_MAXX : JPGTAG_BIO_PIXEL_MAXX);

    assert(comp < so->depth);
    ULONG stripe_h = (maxy + 1 - miny > 8) ? 8 : (maxy + 1 - miny);

    if (action == JPGFLAG_BIO_REQUEST) {
        // Point the library at our stripe buffer for this component
        UBYTE* mem = so->stripe;
        if (so->pixel_type == CTYP_UBYTE) {
            mem += comp;
            mem -= miny * so->depth * stripe_w;
            tags->SetTagPtr(JPGTAG_BIO_MEMORY, mem);
            tags->SetTagData(JPGTAG_BIO_WIDTH, stripe_w);
            tags->SetTagData(JPGTAG_BIO_HEIGHT, 8 + miny);
            tags->SetTagData(JPGTAG_BIO_BYTESPERROW, so->depth * stripe_w * so->bpp);
            tags->SetTagData(JPGTAG_BIO_BYTESPERPIXEL, so->depth * so->bpp);
            tags->SetTagData(JPGTAG_BIO_PIXELTYPE, so->pixel_type);
        } else if (so->pixel_type == CTYP_UWORD) {
            UWORD* wmem = (UWORD*)so->stripe;
            wmem += comp;
            wmem -= miny * so->depth * stripe_w;
            tags->SetTagPtr(JPGTAG_BIO_MEMORY, wmem);
            tags->SetTagData(JPGTAG_BIO_WIDTH, stripe_w);
            tags->SetTagData(JPGTAG_BIO_HEIGHT, 8 + miny);
            tags->SetTagData(JPGTAG_BIO_BYTESPERROW, so->depth * stripe_w * so->bpp);
            tags->SetTagData(JPGTAG_BIO_BYTESPERPIXEL, so->depth * so->bpp);
            tags->SetTagData(JPGTAG_BIO_PIXELTYPE, so->pixel_type);
        }
    } else if (action == JPGFLAG_BIO_RELEASE) {
        // Last component done → copy stripe to output
        if (comp == so->depth - 1) {
            ULONG count = stripe_w * stripe_h * so->depth;
            ULONG sz    = so->bpp;
            UBYTE* mem  = so->stripe;
            for (ULONG i = 0; i < count; i++) {
                if (so->pos + (int)sz > so->output_len) break;
                for (ULONG j = 0; j < sz; j++) {
                    so->output[so->pos++] = *mem++;
                }
            }
        }
    }

    return 0;
}

// ---- JPEG construction via Read (no tags needed) ----

static JPEG* create_jpeg(const MemStream& in)
{
    JPG_Hook hook(input_hook, const_cast<MemStream*>(&in));
    JPG_TagItem tags[] = {
        JPG_PointerTag(JPGTAG_HOOK_IOHOOK, &hook),
        JPG_PointerTag(JPGTAG_HOOK_IOSTREAM, const_cast<MemStream*>(&in)),
        JPG_EndTag
    };
    JPEG* jpeg = JPEG::Construct(tags);
    if (!jpeg) return nullptr;

    if (!jpeg->Read(tags)) {
        JPEG::Destruct(jpeg);
        return nullptr;
    }
    return jpeg;
}

// ---- public API ----

int golibjpeg_decode_to_rgb(
    const unsigned char* data, int data_len,
    unsigned char** output, int* output_len,
    int* width, int* height)
{
    int components = 0, precision = 0;
    return golibjpeg_decode(data, data_len, GOLIBJPEG_CT_YCBCR,
                            output, output_len,
                            width, height, &components, &precision);
}

int golibjpeg_get_parameters(
    const unsigned char* data, int data_len,
    int* width, int* height, int* components, int* precision)
{
    if (!data || data_len <= 0 || !width || !height || !components || !precision)
        return GOLIBJPEG_ERR_PARAM;

    MemStream in;
    in.data = data;
    in.size = data_len;
    in.pos  = 0;

    JPEG* jpeg = create_jpeg(in);
    if (!jpeg) return GOLIBJPEG_ERR_DECODE;

    JPG_TagItem itags[] = {
        JPG_ValueTag(JPGTAG_IMAGE_WIDTH, 0),
        JPG_ValueTag(JPGTAG_IMAGE_HEIGHT, 0),
        JPG_ValueTag(JPGTAG_IMAGE_DEPTH, 0),
        JPG_ValueTag(JPGTAG_IMAGE_PRECISION, 0),
        JPG_EndTag
    };

    int ok = jpeg->GetInformation(itags);
    if (!ok) {
        JPEG::Destruct(jpeg);
        return GOLIBJPEG_ERR_DECODE;
    }

    *width      = (int)itags->GetTagData(JPGTAG_IMAGE_WIDTH);
    *height     = (int)itags->GetTagData(JPGTAG_IMAGE_HEIGHT);
    *components = (int)itags->GetTagData(JPGTAG_IMAGE_DEPTH);
    *precision  = (int)itags->GetTagData(JPGTAG_IMAGE_PRECISION);

    JPEG::Destruct(jpeg);
    return GOLIBJPEG_OK;
}

int golibjpeg_decode(
    const unsigned char* data, int data_len,
    int colour_transform,
    unsigned char** output, int* output_len,
    int* width, int* height, int* components, int* precision)
{
    if (!data || data_len <= 0 || !output || !output_len ||
        !width || !height || !components || !precision)
        return GOLIBJPEG_ERR_PARAM;

    *output    = nullptr;
    *output_len = 0;

    if (colour_transform < 0 || colour_transform > 3)
        return GOLIBJPEG_ERR_PARAM;

    MemStream in;
    in.data = data;
    in.size = data_len;
    in.pos  = 0;

    JPEG* jpeg = create_jpeg(in);
    if (!jpeg) return GOLIBJPEG_ERR_DECODE;

    // Get image info
    UBYTE subx[4] = {1,1,1,1}, suby[4] = {1,1,1,1};
    JPG_TagItem itags[] = {
        JPG_ValueTag(JPGTAG_IMAGE_WIDTH, 0),
        JPG_ValueTag(JPGTAG_IMAGE_HEIGHT, 0),
        JPG_ValueTag(JPGTAG_IMAGE_DEPTH, 0),
        JPG_ValueTag(JPGTAG_IMAGE_PRECISION, 0),
        JPG_ValueTag(JPGTAG_IMAGE_IS_FLOAT, false),
        JPG_PointerTag(JPGTAG_IMAGE_SUBX, subx),
        JPG_PointerTag(JPGTAG_IMAGE_SUBY, suby),
        JPG_ValueTag(JPGTAG_IMAGE_SUBLENGTH, 4),
        JPG_EndTag
    };

    if (!jpeg->GetInformation(itags)) {
        JPEG::Destruct(jpeg);
        return GOLIBJPEG_ERR_DECODE;
    }

    int w      = (int)itags->GetTagData(JPGTAG_IMAGE_WIDTH);
    int h      = (int)itags->GetTagData(JPGTAG_IMAGE_HEIGHT);
    int depth  = (int)itags->GetTagData(JPGTAG_IMAGE_DEPTH);
    int prec   = (int)itags->GetTagData(JPGTAG_IMAGE_PRECISION);

    if (w <= 0 || h <= 0 || depth <= 0) {
        JPEG::Destruct(jpeg);
        return GOLIBJPEG_ERR_FORMAT;
    }

    UBYTE pixel_type = CTYP_UBYTE;
    int   bpp        = 1;
    if (prec > 8) {
        pixel_type = CTYP_UWORD;
        bpp        = 2;
    }

    // Allocate stripe buffer
    int   stripe_w = w;
    ULONG stripe_sz = (ULONG)stripe_w * 8 * depth * bpp;
    auto* stripe = (UBYTE*)malloc(stripe_sz);
    if (!stripe) {
        JPEG::Destruct(jpeg);
        return GOLIBJPEG_ERR_MEMORY;
    }
    memset(stripe, 0, stripe_sz);

    // Allocate final output
    ULONG out_sz = (ULONG)w * h * depth * bpp;
    auto* output_buf = (unsigned char*)malloc(out_sz);
    if (!output_buf) {
        free(stripe);
        JPEG::Destruct(jpeg);
        return GOLIBJPEG_ERR_MEMORY;
    }

    StripeOut so;
    so.output     = output_buf;
    so.output_len = (int)out_sz;
    so.pos        = 0;
    so.stripe     = stripe;
    so.width      = w;
    so.height     = h;
    so.depth      = depth;
    so.pixel_type = pixel_type;
    so.bpp        = bpp;
    so.big_endian = false;
    so.upsample   = true;

    JPG_Hook outhook(output_hook, &so);

    // Stripe-based reconstruction (8 lines at a time)
    ULONG y = 0;
    int ok = 1;
    while (y < (ULONG)h && ok) {
        ULONG lastline = y + 8;
        if (lastline > (ULONG)h) lastline = (ULONG)h;

        JPG_TagItem dtags[] = {
            JPG_PointerTag(JPGTAG_BIH_HOOK, &outhook),
            JPG_ValueTag(JPGTAG_DECODER_MINY, (LONG)y),
            JPG_ValueTag(JPGTAG_DECODER_MAXY, (LONG)(lastline - 1)),
            JPG_ValueTag(JPGTAG_DECODER_UPSAMPLE, true),
            JPG_ValueTag(JPGTAG_MATRIX_LTRAFO, colour_transform),
            JPG_EndTag
        };

        ok = jpeg->DisplayRectangle(dtags);
        y  = lastline;
    }

    free(stripe);

    if (!ok) {
        const char* err_str = nullptr;
        jpeg->LastError(err_str);
        free(output_buf);
        JPEG::Destruct(jpeg);
        return GOLIBJPEG_ERR_DECODE;
    }

    JPEG::Destruct(jpeg);

    *output      = output_buf;
    *output_len  = (int)out_sz;
    *width       = w;
    *height      = h;
    *components  = depth;
    *precision   = prec;

    return GOLIBJPEG_OK;
}

void golibjpeg_free(unsigned char* p)
{
    free(p);
}
