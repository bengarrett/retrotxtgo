//Package filesystem to handle the opening and reading of text files
package filesystem

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"unicode"
	"unicode/utf8"
)

const (
	// PermF is posix permission bits for files
	permf os.FileMode = 0660
	// PermD is posix permission bits for directories
	permd os.FileMode = 0700
)

// Clean removes the named file or directory.
func Clean(name string) {
	if err := os.RemoveAll(name); err != nil {
		fmt.Fprintln(os.Stderr, "removing path:", err)
	}
}

// Columns counts the number of characters used per line in the named file.
func Columns(name string) (count int, err error) {
	file, err := os.Open(name)
	if err != nil {
		return -1, err
	}
	defer file.Close()
	count, err = columnsCounter(file)
	return count, err
}

func columnsCounter(r io.Reader) (width int, err error) {
	const lineBreak = '\n'
	buf := make([]byte, bufio.MaxScanTokenSize)
	for {
		size, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}
		var pos int
		for {
			i := bytes.IndexByte(buf[pos:], lineBreak)
			if i == -1 || size == pos {
				break
			}
			pos += i + 1
			if i > width {
				width = i
			}
		}
		if err == io.EOF {
			break
		}
	}
	return width, nil
}

// Controls counts the number of lines in the named file.
func Controls(name string) (count int, err error) {
	file, err := os.Open(name)
	if err != nil {
		return -1, err
	}
	defer file.Close()
	count, err = ctrlCounter(file)
	return count, err
}

func ctrlCounter(r io.Reader) (count int, err error) {
	lineBreak := []byte("\x1B\x5b") // esc[
	buf := make([]byte, bufio.MaxScanTokenSize)
	for {
		size, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}
		var pos int
		for {
			i := bytes.Index(buf[pos:], lineBreak)
			if i == -1 || size == pos {
				break
			}
			pos += i + 1
			if i > count {
				count = i
			}
		}
		if err == io.EOF {
			break
		}
	}
	return count, nil
}

// Lines counts the number of lines in the named file.
func Lines(name string) (count int, err error) {
	file, err := os.Open(name)
	if err != nil {
		return -1, err
	}
	defer file.Close()
	count, err = lineCounter(file)
	return count, err
}

func lineCounter(r io.Reader) (count int, err error) {
	const lineBreak = '\n'
	buf := make([]byte, bufio.MaxScanTokenSize)
	for {
		size, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}
		var pos int
		for {
			i := bytes.IndexByte(buf[pos:], lineBreak)
			if i == -1 || size == pos {
				break
			}
			pos += i + 1
			count++
		}
		if err == io.EOF {
			break
		}
	}
	return count, nil
}

// Words counts the number of spaced words in the named file.
func Words(name string) (count int, err error) {
	file, err := os.Open(name)
	if err != nil {
		return -1, err
	}
	defer file.Close()
	count, err = wordCounter(file)
	return count, err
}

func wordCounter(r io.Reader) (count int, err error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		var t = scanner.Text()
		r, _ := utf8.DecodeRuneInString(t)
		if r >= 65533 {
			continue
		}
		// scan single chars
		if len(t) == 1 {
			if unicode.IsDigit(r) || unicode.IsLetter(r) {
				count++
			}
			continue
		}
		// scan chars within each word
		if IsWord(t) {
			count++
		}
	}
	if err = scanner.Err(); err != nil {
		return -1, err
	}
	return count, err
}

// IsWord scans the content of a word for characters that are not digits,
// letters or punctuation and if discovered returns false.
// If a space or newline is encountered the scan will end.
func IsWord(s string) bool {
	if len(s) == 0 {
		return false
	}
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		if unicode.IsSpace(r) || s == "\n" {
			break
		}
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) && !unicode.IsPunct(r) {
			return false
		}
		s = s[size:]
	}
	return true
}

// Read opens and returns the content of the named file.
func Read(name string) (data []byte, err error) {
	return ReadAllBytes(name)
}

// ReadAllBytes reads the named file and returns the content as a byte array.
// Create a word and random character generator to make files larger than 64k.
func ReadAllBytes(name string) (data []byte, err error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// bufio is the most performant way to scan streamed data
	scanner := bufio.NewScanner(file)
	// optional adjustment to the token size
	// Go by default will scan 64 * 1024 bytes (64KB) per iteration
	scanner.Buffer(data, 64*1024)
	// required, split scan into Buffer(data, x) sized byte chuncks
	// otherwise scanner will panic on files larger than 64 * 1024 bytes
	scanner.Split(bufio.ScanBytes)
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

func dir(name string) (path string, err error) {
	path = filepath.Dir(name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, permd)
	}
	return path, err
}

// Save bytes to a named file location.
func Save(filename string, b []byte) (path string, err error) {
	path, err = dir(filename)
	if err != nil {
		return path, err
	}
	path = filename
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, permf)
	if err != nil {
		return path, err
	}
	// bufio is the most performant
	writer := bufio.NewWriter(file)
	for _, c := range b {
		if err = writer.WriteByte(c); err != nil {
			return path, err
		}
	}
	if err = writer.Flush(); err != nil {
		return path, err
	}
	if err = file.Close(); err != nil {
		return path, err
	}
	//ioutil.WriteFile(filename,data,perm)
	return filepath.Abs(file.Name())
}

// SaveTemp saves bytes to a named temporary file.
func SaveTemp(filename string, b []byte) (path string, err error) {
	return Save(tempFile(filename), b)
}

func tempFile(name string) (path string) {
	path = name
	if filepath.Base(name) == name {
		path = filepath.Join(os.TempDir(), name)
	}
	return path
}

// Touch creates an empty file at the named location.
func Touch(name string) (path string, err error) {
	return Save(name, nil)
}
