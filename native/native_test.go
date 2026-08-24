package native

import (
	"errors"
	"runtime"
	"strings"
	"sync"
	"testing"
)

// TestEncodeLoadsTheEmbeddedLibrary covers the lazy load itself: nothing runs in
// init() any more, so the first call has to be the one that extracts and binds
// the library.
func TestEncodeLoadsTheEmbeddedLibrary(t *testing.T) {
	const w, h = 8, 8
	encoded, err := Encode(make([]byte, w*h*3), EncodeParams{
		Columns:         w,
		Rows:            h,
		SamplesPerPixel: 3,
		BitsPerSample:   8,
		FrameType:       FrameBaseline,
		ColourTransform: CTYCbCr,
		Quality:         90,
	})
	if err != nil {
		t.Fatalf("Encode on %s/%s: %v", runtime.GOOS, runtime.GOARCH, err)
	}

	params, err := GetParameters(encoded)
	if err != nil {
		t.Fatalf("GetParameters: %v", err)
	}
	if params.Width != w || params.Height != h {
		t.Errorf("params = %+v, want %dx%d", params, w, h)
	}
}

// simulateUnsupported empties libData and resets the loader, which is what a
// GOOS/GOARCH with no prebuilt library looks like from loadNative's side. Tests
// cannot run on such a platform — CI only cross-builds for those — so this is
// the only way to exercise the behaviour.
func simulateUnsupported(t *testing.T) {
	t.Helper()
	saved := libData
	libData = nil
	ensureLoaded = sync.OnceValue(loadNative)
	t.Cleanup(func() {
		libData = saved
		ensureLoaded = sync.OnceValue(loadNative)
	})
}

// Reaching the end of this test at all is half the point: before the lazy load
// an unloadable library was a panic, and a panic here would fail the test.
func TestUnsupportedPlatformIsAnErrorFromEveryEntryPoint(t *testing.T) {
	simulateUnsupported(t)

	jpeg := []byte{0xFF, 0xD8, 0xFF, 0xE0}

	if _, err := Decode(jpeg, CTYCbCr); !errors.Is(err, ErrUnsupportedPlatform) {
		t.Errorf("Decode: got %v, want ErrUnsupportedPlatform", err)
	}
	if _, err := DecodeToRGB(jpeg); !errors.Is(err, ErrUnsupportedPlatform) {
		t.Errorf("DecodeToRGB: got %v, want ErrUnsupportedPlatform", err)
	}
	if _, err := GetParameters(jpeg); !errors.Is(err, ErrUnsupportedPlatform) {
		t.Errorf("GetParameters: got %v, want ErrUnsupportedPlatform", err)
	}
	params := EncodeParams{Columns: 4, Rows: 4, SamplesPerPixel: 1, BitsPerSample: 8}
	if _, err := Encode(make([]byte, 16), params); !errors.Is(err, ErrUnsupportedPlatform) {
		t.Errorf("Encode: got %v, want ErrUnsupportedPlatform", err)
	}
}

// The platform belongs in the message; "unsupported platform" alone leaves the
// reader guessing which one the build was for.
func TestUnsupportedPlatformErrorNamesThePlatform(t *testing.T) {
	simulateUnsupported(t)

	_, err := GetParameters([]byte{0xFF, 0xD8})
	if err == nil {
		t.Fatal("GetParameters succeeded with no embedded library")
	}
	for _, want := range []string{runtime.GOOS, runtime.GOARCH} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not mention %q", err, want)
		}
	}
}

// A failed load must stay failed. Retrying per call would risk binding a
// half-initialised library, and would repeat the extraction on every call.
func TestLoadResultIsCached(t *testing.T) {
	simulateUnsupported(t)

	first, second := ensureLoaded(), ensureLoaded()
	if first == nil {
		t.Fatal("ensureLoaded() succeeded with no embedded library")
	}
	if first != second {
		t.Errorf("ensureLoaded() returned a new error value on the second call: %v then %v", first, second)
	}
}

// The detail half of a decode failure comes from a separate C symbol bound as
// func() string, and nothing else in the suite reads it — so a rename, a
// signature change, or a library built without the symbol would all go
// unnoticed while the exported behaviour quietly lost half its message.
//
// The input is a well-formed JFIF header truncated to an immediate EOI: past the
// magic check, so the failure happens inside libjpeg with something to say about
// it, rather than being rejected up front.
func TestDecodeErrorCarriesTheDetailFromTheLibrary(t *testing.T) {
	truncated := []byte{
		0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 'J', 'F', 'I', 'F', 0x00,
		0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0xFF, 0xD9,
	}

	_, err := Decode(truncated, CTNone)
	if err == nil {
		t.Fatal("Decode accepted a truncated JPEG")
	}

	var status *StatusError
	if !errors.As(err, &status) {
		t.Fatalf("Decode returned %T (%v), want *StatusError", err, err)
	}
	if status.Detail == "" {
		t.Fatalf("StatusError has no Detail, so golibjpeg_last_error did not bind: %v", err)
	}

	// pylibjpeg-libjpeg raises
	//   libjpeg error code '-1038' returned from Decode(): <known> - <detail>
	// so Op supplies the "Decode()" and the format string must not add a second
	// pair of parentheses.
	msg := err.Error()
	for _, want := range []string{
		"returned from Decode(): ",
		" - " + status.Detail,
	} {
		if !strings.Contains(msg, want) {
			t.Errorf("error %q does not contain %q", msg, want)
		}
	}
	if strings.Contains(msg, "()()") {
		t.Errorf("error %q doubles the parentheses after the operation name", msg)
	}
}
