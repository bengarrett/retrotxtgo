package tmp_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bengarrett/retrotxtgo/internal/tmp"
	"github.com/nalgeon/be"
)

func TestFile(t *testing.T) {
	t.Parallel()

	// Test with absolute path - should return the same path
	tempDir := t.TempDir()
	abspath := filepath.Join(tempDir, "testfile.txt")
	result := tmp.File(abspath)
	be.Equal(t, result, abspath)

	// Test with relative path that has directory components
	relpath := "subdir/testfile.txt"
	result = tmp.File(relpath)
	be.Equal(t, result, relpath)

	// Test with just a filename - should return path in temp directory
	filename := "testfile.txt"
	result = tmp.File(filename)
	expected := filepath.Join(os.TempDir(), filename)
	be.Equal(t, result, expected)

	// Test with empty string
	result = tmp.File("")
	be.Equal(t, result, "")

	// Test with path that looks like filename but has path separators
	mixedPath := "test/file.txt"
	result = tmp.File(mixedPath)
	be.Equal(t, result, mixedPath)
}
