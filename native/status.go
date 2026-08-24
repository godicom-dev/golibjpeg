package native

// lastErrorDetail reports the message the last failing C call left behind, or ""
// when the library predates the symbol and lastErrorFn was never bound.
//
// purego converts the C char* to a Go string itself — including the NULL case —
// so nothing here has to touch unsafe. Doing that conversion by hand is what
// made `go vet ./...` flag this file: turning the returned uintptr into an
// unsafe.Pointer is the pattern vet rejects, and vet cannot tell that this
// particular address belongs to the loaded library rather than the Go heap.
func lastErrorDetail() string {
	if lastErrorFn == nil {
		return ""
	}
	return lastErrorFn()
}
