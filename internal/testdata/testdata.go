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

// RequireJPEG skips the test when a conformance file is missing.
func RequireJPEG(t *testing.T, parts ...string) []byte {
	t.Helper()
	path := JPEGPath(parts...)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Skipf("reference testdata not installed at %s (run scripts/fetch-testdata.sh)", path)
	}
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
