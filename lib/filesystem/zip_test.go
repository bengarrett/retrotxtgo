package filesystem

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

func ExampleZip() {
	// Create a temporary directory
	tmpZip := tempFile("retrotxtgo_zip_directory_test")
	err := os.Mkdir(tmpZip, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpZip)

	// Create a temporary 1 byte file in the temporary directory
	tmpFile, err := SaveTemp(path.Join(tmpZip, "temp.zip"), []byte("x")...)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile)

	// Initialize the Zip archive file
	path := tempFile("exampleZip.zip")
	zip := Zip{
		Name:      path,
		Root:      tmpZip,
		Comment:   "",
		Overwrite: true,
		Quiet:     true,
	}

	// Create the Zip archive file
	if err := zip.Create(); err != nil {
		log.Fatal(err)
	}

	// Check the Zip archive exists
	s, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s, %d", s.Name(), s.Size())
	// Output: exampleZip.zip, 149
}

func ExampleUniqueName() {
	name := "retrotxtgo_uniquetest.txt"

	// Create a temporary 1 byte file in the temporary directory
	tmpFile, err := SaveTemp(tempFile(name), []byte("x")...)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile)

	// Use UniqueName to find a new unique filename
	// so not to conflict with the previously saved file
	u, err := UniqueName(tmpFile)
	if err != nil {
		log.Fatal(err)
	}

	// In Linux the new name will be retrotxtgo_uniquetest_1.txt
	// In Windows the name be retrotxtgo_uniquetest (1).txt
	newName := filepath.Base(u)

	// As the new unique names vary based on the host operating system
	// Compare the name lengths to confirm the creation of a new filename
	unique := bool(len(newName) > len(name))
	fmt.Print(unique)
	// Output: true
}
