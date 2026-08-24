//go:build !(linux && (amd64 || arm64)) && !(darwin && (amd64 || arm64)) && !(windows && (amd64 || arm64))

package native

// No prebuilt library is built for this platform, so there is nothing to embed.
// loadNative turns the empty libData into ErrUnsupportedPlatform, which keeps
// the module compiling everywhere Go and purego do — decoding and encoding then
// fail with that error instead of the whole build failing.
//
// Adding a platform means adding its native_GOOS_GOARCH.go *and* excluding that
// platform from the constraint above. Skipping the second half declares libData
// twice, which does not compile, so the compiler enforces the pairing.
var libData []byte
