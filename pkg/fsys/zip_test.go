package fsys_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"github.com/bengarrett/retrotxtgo/pkg/internal/tmp"
)

func ExampleZip() {
	// Create a temporary directory
	tmpZip := tmp.File("retrotxtgo_zip_directory_test")
	err := os.MkdirAll(tmpZip, 0o755)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpZip)

	// Create a temporary 1 byte file in the temporary directory
	tmpFile, err := fsys.SaveTemp(path.Join(tmpZip, "temp.zip"), []byte("x")...)
	if err != nil {
		log.Print(err)
		return
	}
	defer os.Remove(tmpFile)

	// Initialize the Zip archive file
	name := tmp.File("exampleZip.zip")
	zip := fsys.Zip{
		Name:      name,
		Root:      tmpZip,
		Comment:   "",
		Overwrite: true,
		Writer:    io.Discard,
	}

	// Create the Zip archive file
	if err := zip.Create(); err != nil {
		log.Print(err)
		return
	}

	// Check the Zip archive exists
	s, err := os.Stat(name)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Fprintf(os.Stdout, "%s, %d", s.Name(), s.Size())
	// Output: exampleZip.zip, 149
}

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
