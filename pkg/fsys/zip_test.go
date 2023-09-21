package fsys_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"github.com/bengarrett/retrotxtgo/pkg/internal/tmp"
)

func ExampleUniqueName() {
	name := "retrotxtgo_uniquetest.txt"

	// Create a temporary 1 byte file in the temporary directory
	tmpFile, err := fsys.SaveTemp(tmp.File(name), []byte("x")...)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile)

	// Use UniqueName to find a new unique filename
	// so not to conflict with the previously saved file
	u, err := fsys.UniqueName(tmpFile)
	if err != nil {
		log.Print(err)
		return
	}

	// In Linux the new name will be retrotxtgo_uniquetest_1.txt
	// In Windows the name be retrotxtgo_uniquetest (1).txt
	newName := filepath.Base(u)

	// As the new unique names vary based on the host operating system
	// Compare the name lengths to confirm the creation of a new filename
	unique := bool(len(newName) > len(name))
	fmt.Fprint(os.Stdout, unique)
	// Output: true
}
