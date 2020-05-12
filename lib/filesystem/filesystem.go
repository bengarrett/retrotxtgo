//Package filesystem to handle the opening and reading of text files
package filesystem

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
)

// Read opens and returns the content of the name file.
func Read(name string) ([]byte, error) {
	// check name is file not anything else
	data, err := ReadAllBytes(name)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//TailBytes reads the name file from the offset position relative to the end of the file.
func TailBytes(name string, offset int64) ([]byte, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var size int64 = int64(math.Abs(float64(offset)))
	if stat, err := os.Stat(name); err == nil && stat.Size() < size {
		return nil, fmt.Errorf("offset: value is %v too large for a %v byte file", offset, stat.Size())
	}
	// file.Seek(whence)
	// 0 means relative to the origin of the file
	// 1 means relative to the current offset
	// 2 means relative to the end
	if _, err = file.Seek(offset, 2); err != nil {
		// todo: have offset deal with runes not bytes
		return nil, err
	}
	buffer := make([]byte, size)
	if _, err = file.Read(buffer); err != nil {
		return nil, err
	}
	return buffer, nil
}

//ReadAllBytes reads the named file and returns its content as a byte array.
func ReadAllBytes(name string) ([]byte, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

// ReadChunk is unused
func ReadChunk(name string, size int) ([]byte, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	read := bufio.NewReaderSize(file, size)
	buf := make([]byte, size)

	_, err = read.Read(buf)
	if err != nil {
		return nil, err
	}
	// do things
	return buf, nil
}
