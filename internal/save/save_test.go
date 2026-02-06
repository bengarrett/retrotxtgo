package save_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bengarrett/retrotxtgo/internal/save"
	"github.com/nalgeon/be"
)

func TestSave(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "testfile.txt")
	content := []byte("Hello, World!")

	// Test successful save
	n, path, err := save.Save(testFile, content...)
	be.Err(t, err, nil)
	be.Equal(t, n, len(content)-1) // Returns index of last byte written
	be.Equal(t, path, testFile)

	// Verify file was created and contains correct content
	data, err := os.ReadFile(path)
	be.Err(t, err, nil)
	be.Equal(t, string(data), string(content))

	// Test with empty content
	n, path, err = save.Save(filepath.Join(tempDir, "empty.txt"))
	be.Err(t, err, nil)
	be.Equal(t, n, 0) // No bytes written, returns 0

	// Test with nested directory creation
	nestedFile := filepath.Join(tempDir, "subdir", "nested", "file.txt")
	n, path, err = save.Save(nestedFile, []byte("nested content")...)
	be.Err(t, err, nil)
	be.Equal(t, n, len("nested content")-1)

	// Verify nested directory was created
	_, err = os.Stat(filepath.Dir(path))
	be.Err(t, err, nil)
}

func TestSaveErrors(t *testing.T) {
	t.Parallel()

	// Test with invalid path (should still work due to directory creation)
	tempDir := t.TempDir()
	invalidPath := filepath.Join(tempDir, "nonexistent", "file.txt")
	n, _, err := save.Save(invalidPath, []byte("test")...)
	be.Err(t, err, nil)
	be.Equal(t, n, 3) // Returns index of last byte written

	// Test with empty filename
	n, _, err = save.Save("", []byte("test")...)
	be.True(t, err != nil) // Should fail with empty filename
	be.Equal(t, n, 0)
}
