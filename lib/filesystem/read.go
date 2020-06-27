package filesystem

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/logs"
)

// IsPipe determines if Stdin (standard input) is piped from another command.
func IsPipe() bool {
	// source: https://dev.to/napicella/linux-pipes-in-golang-2e8j
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		logs.LogCont(err)
	}
	return fileInfo.Mode()&os.ModeCharDevice == 0
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

// ReadColumns counts the number of characters used per line in the named file.
func ReadColumns(name string) (count int, err error) {
	file, err := os.Open(name)
	if err != nil {
		return -1, err
	}
	defer file.Close()
	nl, err := ReadNewlines(name)
	if err != nil {
		return -1, err
	}
	count, err = Columns(file, nl)
	return count, err
}

// ReadControls counts the number of ANSI escape sequences in the named file.
func ReadControls(name string) (count int, err error) {
	file, err := os.Open(name)
	if err != nil {
		return -1, err
	}
	defer file.Close()
	count, err = Controls(file)
	return count, err
}

// ReadLine reads a named file location or a named temporary file and returns its content.
func ReadLine(name, newline string) (text string, err error) {
	var path, n = tempFile(name), nl(newline)
	file, err := os.OpenFile(path, os.O_RDONLY, permf)
	if err != nil {
		return text, err
	}
	// bufio is the most performant
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text += fmt.Sprintf("%s%s", scanner.Text(), n)
	}
	if err = scanner.Err(); err != nil {
		return text, err
	}
	if err = file.Close(); err != nil {
		return text, err
	}
	err = scanner.Err()
	return text, err
}

// ReadLines counts the number of lines in the named file.
func ReadLines(name string) (count int, err error) {
	file, err := os.Open(name)
	if err != nil {
		return -1, err
	}
	defer file.Close()
	count, err = Lines(file)
	return count, err
}

// ReadNewlines scans the named file for the most commonly used newline method.
func ReadNewlines(name string) ([2]rune, error) {
	z := [2]rune{0, 0}
	file, err := os.Open(name)
	if err != nil {
		return z, err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return z, err
	}
	return Newlines(bytes.Runes(b)), nil
}

// ReadPipe reads data piped by the operating system's STDIN.
// If no data is detected the program will exit.
func ReadPipe() (b []byte, err error) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		b = append(b, scanner.Bytes()...)
	}
	if err = scanner.Err(); err != nil {
		return b, err
	}
	if len(b) == 0 {
		os.Exit(0)
	}
	return b, nil
}

// ReadRunes returns the number of runes in the named file.
func ReadRunes(name string) (count int, err error) {
	file, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	return Runes(file)
}

// ReadTail reads the named file from the offset position relative to the end of the file.
func ReadTail(name string, offset int) (data []byte, err error) {
	file, err := os.Open(name)
	if err != nil {
		return data, err
	}
	defer file.Close()
	count, total := 0, 0
	total, err = ReadRunes(name)
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

// ReadText reads a named file location or a named temporary file and returns its content.
func ReadText(name string) (text string, err error) {
	return ReadLine(name, "")
}

// ReadWords counts the number of spaced words in the named file.
func ReadWords(name string) (count int, err error) {
	file, err := os.Open(name)
	if err != nil {
		return -1, err
	}
	defer file.Close()
	return Words(file)
}
