package filesystem

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/logs"
)

type lineBreaks int

const (
	nl lineBreaks = iota
	dos
	win
	c64
	darwin
	mac
	amiga
	linux
	unix
)

// IsPipe determines if Stdin (standard input) is piped from another command.
func IsPipe() bool {
	// source: https://dev.to/napicella/linux-pipes-in-golang-2e8j
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		logs.Save(err)
	}
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

// Read opens and returns the content of the named file.
func Read(name string) ([]byte, error) {
	return ReadAllBytes(name)
}

// ReadAllBytes reads the named file and returns the content as a byte array.
// Create a word and random character generator to make files larger than 64k.
func ReadAllBytes(name string) (data []byte, err error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("read all bytes could not open file: %q: %w", name, err)
	}
	defer file.Close()
	// bufio is the most performant way to scan streamed data
	scanner := bufio.NewScanner(file)
	// optional adjustment to the token size
	// Go by default will scan 64 * 1024 bytes (64KB) per iteration
	const KB = 1024
	const max = 64 * KB
	scanner.Buffer(data, max)
	// required, split scan into Buffer(data, x) sized byte chuncks
	// otherwise scanner will panic on files larger than 64 * 1024 bytes
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		data = append(data, scanner.Bytes()...)
	}
	if err = scanner.Err(); err != nil {
		return data, fmt.Errorf("read all bytes could not scan file: %q: %w", name, err)
	}
	return data, file.Close()
}

// ReadChunk reads and returns the start of the named file.
func ReadChunk(name string, chars int) (data []byte, err error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("read chunk could not open file: %q: %w", name, err)
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
		return data, fmt.Errorf("read chunk could not scan file: %q: %w", name, err)
	}
	return data, file.Close()
}

// ReadColumns counts the number of characters used per line in the named file.
func ReadColumns(name string) (count int, err error) {
	file, err := os.Open(name)
	if err != nil {
		return -1, fmt.Errorf("read columns could not open file: %q: %w", name, err)
	}
	defer file.Close()
	nl, err := ReadLineBreaks(name)
	if err != nil {
		return -1, fmt.Errorf("read columns could not find the line break method: %w", err)
	}
	count, err = Columns(file, nl)
	if err != nil {
		return -1, fmt.Errorf("read columns count the file: %q: %w", name, err)
	}
	return count, file.Close()
}

// ReadControls counts the number of ANSI escape sequences in the named file.
func ReadControls(name string) (count int, err error) {
	file, err := os.Open(name)
	if err != nil {
		return -1, fmt.Errorf("read countrols could not open file: %q: %w", name, err)
	}
	defer file.Close()
	count, err = Controls(file)
	if err != nil {
		return -1, fmt.Errorf("read countrols could not parse the file: %q: %w", name, err)
	}
	return count, file.Close()
}

// ReadLine reads a named file location or a named temporary file and returns its content.
func ReadLine(name string, lb lineBreaks) (text string, err error) {
	var path, n = tempFile(name), lineBreak(lb)
	file, err := os.OpenFile(path, os.O_RDONLY, filemode)
	if err != nil {
		return text, fmt.Errorf("read line could not open file: %q: %w", name, err)
	}
	defer file.Close()
	// bufio is the most performant
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text += fmt.Sprintf("%s%s", scanner.Text(), n)
	}
	if err = scanner.Err(); err != nil {
		return text, fmt.Errorf("read line could not scan file: %w", err)
	}
	return text, file.Close()
}

// ReadLines counts the number of lines in the named file.
func ReadLines(name string) (count int, err error) {
	file, err := os.Open(name)
	if err != nil {
		return -1, fmt.Errorf("read lines could not open file: %q: %w", name, err)
	}
	defer file.Close()
	nl, err := ReadLineBreaks(name)
	if err != nil {
		return -1, fmt.Errorf("read lines could not scan the file: %w", err)
	}
	count, err = Lines(file, nl)
	if err != nil {
		return -1, fmt.Errorf("read tail could not open file: %q: %w", name, err)
	}
	return count, file.Close()
}

// ReadLineBreaks scans the named file for the most commonly used line break method.
func ReadLineBreaks(name string) ([2]rune, error) {
	z := [2]rune{0, 0}
	file, err := os.Open(name)
	if err != nil {
		return z, fmt.Errorf("read line breaks could not open file: %q: %w", name, err)
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return z, fmt.Errorf("read line breaks could not read the file: %q: %w", name, err)
	}
	return LineBreaks(true, bytes.Runes(b)...), file.Close()
}

// ReadPipe reads data piped by the operating system's STDIN.
// If no data is detected the program will exit.
func ReadPipe() (b []byte, err error) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		b = append(b, scanner.Bytes()...)
		b = append(b, []byte("\n")...)
	}
	if err = scanner.Err(); err != nil {
		return b, fmt.Errorf("read pipe could not scan stdin: %w", err)
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
		return 0, fmt.Errorf("read runes could not open file: %q: %w", name, err)
	}
	defer file.Close()
	count, err = Runes(file)
	if err != nil {
		return 0, fmt.Errorf("read runes could not calculate this file: %q: %w", name, err)
	}
	return count, file.Close()
}

// ReadTail reads the named file from the offset position relative to the end of the file.
func ReadTail(name string, offset int) (data []byte, err error) {
	file, err := os.Open(name)
	if err != nil {
		return data, fmt.Errorf("read tail could not open file: %q: %w", name, err)
	}
	defer file.Close()
	count, total := 0, 0
	total, err = ReadRunes(name)
	if err != nil {
		return data, fmt.Errorf("read tail could not read runes: %q: %w", name, err)
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
		return data, fmt.Errorf("read tail could scan file bytes: %q: %w", name, err)
	}
	return data, file.Close()
}

// ReadText reads a named file location or a named temporary file and returns its content.
func ReadText(name string) (text string, err error) {
	text, err = ReadLine(name, nl)
	if err != nil {
		return text, fmt.Errorf("read text: %q: %w", name, err)
	}
	return text, nil
}

// ReadWords counts the number of spaced words in the named file.
func ReadWords(name string) (count int, err error) {
	file, err := os.Open(name)
	if err != nil {
		return -1, fmt.Errorf("read words could not open: %q: %w", name, err)
	}
	defer file.Close()
	count, err = Words(file)
	if err != nil {
		return count, fmt.Errorf("read words failed to count words: %q: %w", name, err)
	}
	return count, file.Close()
}
