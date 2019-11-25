//Package filesystem to handle the opening and reading of text files
package filesystem

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
)

//Read a text file and display its content
func Read(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	stat, _ := os.Stat(path)
	fmt.Printf("opening: %s\twhich is %v bytes\n", stat.Name(), stat.Size())
	data := make([]byte, stat.Size())
	count, err := file.Read(data)
	file.Close()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("read %d bytes: %q\n", count, data[:count])
	return fmt.Sprintf("%s", data[:count])
}

//TailBytes reads the name file from the offset position relative to the end of the file.
func TailBytes(name string, offset int64) ([]byte, error) {
	file, err := os.Open(name)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	// file.Seek(whence)
	// 0 means relative to the origin of the file
	// 1 means relative to the current offset
	// 2 means relative to the end
	_, err = file.Seek(offset, 2) // todo: have offset deal with runes not bytes
	if err != nil {
		return nil, err
	}

	var size int64 = int64(math.Abs(float64(offset)))
	stat, _ := os.Stat(name)
	if stat.Size() < size {
		return nil, fmt.Errorf("offset: value is %v too large for a %v byte file", offset, stat.Size())
	}

	buffer := make([]byte, size)
	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

//ReadAllBytes reads the named file and returns its content as a byte array.
func ReadAllBytes(name string) ([]byte, error) {
	file, err := os.Open(name)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	buffer, err := ioutil.ReadAll(file)
	return buffer, nil
}
