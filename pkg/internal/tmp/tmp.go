package tmp

import (
	"os"
	"path/filepath"
)

// File returns a path to the named file
// if it was stored in the system's temporary directory.
func File(name string) string {
	path := name
	if filepath.Base(name) == name {
		path = filepath.Join(os.TempDir(), name)
	}
	return path
}
