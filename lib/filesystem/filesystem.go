//Package filesystem to handle the opening and reading of text files
package filesystem

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// Read opens and returns the content of the named file.
func Read(name string) (data []byte, err error) {
	// check name is file not anything else
	return ReadAllBytes(name)
}

// ReadAllBytes reads the named file and returns the content as a byte array.
func ReadAllBytes(name string) (data []byte, err error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	// bufio is the most performant
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Bytes()...)
	}
	if err = scanner.Err(); err != nil {
		return data, err
	}
	if err = file.Close(); err != nil {
		return data, err
	}
	err = scanner.Err()
	return data, err
}

// ReadChunk reads a section of the named file and returns the result.
// TODO: finish
func ReadChunk(name string, size int) (data []byte, err error) {
	file, err := os.Open(name)
	if err != nil {
		return data, err
	}
	defer file.Close()
	read := bufio.NewReaderSize(file, size)
	buf := make([]byte, size)
	_, err = read.Read(buf)
	if err != nil {
		return data, err
	}
	//seek logic goes here
	return buf, err
}

// ReadTail reads the named file from the offset position relative to the end of the file.
func ReadTail(name string, offset int64) (data []byte, err error) {
	file, err := os.Open(name)
	if err != nil {
		return data, err
	}
	defer file.Close()
	var size int64 = int64(math.Abs(float64(offset)))
	if stat, err := os.Stat(name); err == nil && stat.Size() < size {
		return data, fmt.Errorf("offset: value is %v too large for a %v byte file", offset, stat.Size())
	}
	// file.Seek(whence)
	// 0 means relative to the origin of the file
	// 1 means relative to the current offset
	// 2 means relative to the end
	if _, err = file.Seek(offset, 2); err != nil {
		// todo: have offset deal with runes not bytes
		return data, err
	}
	data = make([]byte, size)
	if _, err = file.Read(data); err != nil {
		return data, err
	}
	return data, err
}
