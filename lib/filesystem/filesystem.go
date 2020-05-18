//Package filesystem to handle the opening and reading of text files
package filesystem

import (
	"bufio"
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
	defer file.Close()
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

// ReadChunk reads and returns the start of the named file.
func ReadChunk(name string, chars int) (data []byte, err error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// bufio is the most performant
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	count := 0
	for scanner.Scan() {
		count++
		if count > chars {
			break
		}
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

// ReadTail reads the named file from the offset position relative to the end of the file.
func ReadTail(name string, offset int) (data []byte, err error) {
	file, err := os.Open(name)
	if err != nil {
		return data, err
	}
	defer file.Close()
	count, total := 0, 0
	total, err = Runes(name)
	if err != nil {
		return data, err
	}
	// bufio is the most performant
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		count++
		if count <= (total - offset) {
			continue
		}
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

// Runes returns the number of runes in the named file.
func Runes(name string) (count int, err error) {
	file, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		count++
	}
	return count, err
}
