package filesystem

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func ExampleZip() {
	tmpZip := tempFile("retrotxtgo_zip_directory_test")
	err := os.Mkdir(tmpZip, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpZip)
	tmpFile, err := SaveTemp(path.Join(tmpZip, "temp.zip"), []byte("x")...)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile)
	path := tempFile("exampleZip.zip")
	zip := Zip{
		path, tmpZip, "", true, true,
	}
	if err := zip.Create(); err != nil {
		log.Fatal(err)
	}
	s, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s, %d", s.Name(), s.Size())
	// Output: exampleZip.zip, 149
}

func ExampleUniqueName() {
	name := "retrotxtgo_uniquetest.txt"
	tmpFile, err := SaveTemp(tempFile(name), []byte("x")...)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile)
	u, err := UniqueName(tmpFile)
	if err != nil {
		log.Fatal(err)
	}
	newName := bool(len(filepath.Base(u)) > len(name))
	fmt.Print(newName)
	// Output: true
}

func TestUniqueName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UniqueName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("UniqueName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UniqueName() = %v, want %v", got, tt.want)
			}
		})
	}
}
