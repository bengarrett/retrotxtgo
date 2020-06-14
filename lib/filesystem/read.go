package filesystem

import (
	"bufio"
	"os"
)

// ReadBytes reads a named file location or a named temporary file and returns its byte content.
func ReadBytes(name string) (data []byte, err error) {
	var path = tempFile(name)
	file, err := os.OpenFile(path, os.O_RDONLY, permf)
	if err != nil {
		return data, err
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
