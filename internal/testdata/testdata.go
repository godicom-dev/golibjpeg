package testdata

import (
	"os"
	"path/filepath"
	"testing"
)

// Root returns testdata root, or empty if conformance JPEGs are not installed.
func Root() string {
	return filepath.Join("testdata")
}

// JPEGPath joins paths under testdata (e.g. "10918", "p1", "A1.JPG").
func JPEGPath(parts ...string) string {
	all := append([]string{Root()}, parts...)
	return filepath.Join(all...)
}

// RequireJPEGPath returns the path of a conformance file, skipping the test when
// it is missing. Tests that exercise the file-path overloads need the path, not
// the bytes RequireJPEG hands back, and reaching for JPEGPath instead turns a
// missing tree into a failure rather than a skip.
func RequireJPEGPath(t *testing.T, parts ...string) string {
	t.Helper()
	path := JPEGPath(parts...)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			t.Skipf("reference testdata not installed at %s (run scripts/fetch-testdata.sh)", path)
		}
		t.Fatal(err)
	}
	return path
}

// RequireJPEG skips the test when a conformance file is missing.
func RequireJPEG(t *testing.T, parts ...string) []byte {
	t.Helper()
	data, err := os.ReadFile(RequireJPEGPath(t, parts...))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// Available reports whether conformance JPEG tree is present.
func Available() bool {
	_, err := os.Stat(JPEGPath("10918", "p1", "A1.JPG"))
	return err == nil
}
