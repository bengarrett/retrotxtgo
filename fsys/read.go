package fsys

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bengarrett/retrotxtgo/internal/save"
	"github.com/bengarrett/retrotxtgo/nl"
)

// IsPipe reports whether Stdin (standard input) is piped from another command.
func IsPipe() (bool, error) {
	// source: https://dev.to/napicella/linux-pipes-in-golang-2e8j
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false, fmt.Errorf("could not stat stdin: %w", err)
	}
	return fi.Mode()&os.ModeCharDevice == 0, nil
}

// Read opens and returns the content of the named file.
func Read(name string) ([]byte, error) {
	return ReadAllBytes(name)
}

// ReadAllBytes reads the named file and returns the content as a byte array.
// Create a word and random character generator to make files larger than 64k.
func ReadAllBytes(name string) ([]byte, error) {
	file, err := os.Open(name)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// bufio is the most performant way to scan streamed data
	scanner := bufio.NewScanner(file)
	// optional adjustment to the token size
	// Go by default will scan 64 * 1024 bytes (64KB) per iteration
	const size = 64 * 1024
	buf := []byte{}
	scanner.Buffer(buf, size)
	// required, split scan into Buffer(data, x) sized byte chuncks
	// otherwise scanner will panic on files larger than 64 * 1024 bytes
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		buf = append(buf, scanner.Bytes()...)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner %q: %w", name, err)
	}
	return buf, nil
}

// ReadChunk reads and returns the start of the named file.
func ReadChunk(name string, chars int) ([]byte, error) {
	file, err := os.Open(name)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf := []byte{}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	count := 0
	for scanner.Scan() {
		count++
		if count > chars {
			break
		}
		buf = append(buf, scanner.Bytes()...)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read chunk could not scan file: %q: %w", name, err)
	}
	return buf, nil
}

// ReadColumns counts the number of characters used per line in the named file.
func ReadColumns(name string) (int, error) {
	return readLineBreaks(name, true)
}

func readLineBreaks(name string, cols bool) (int, error) {
	file, err := os.Open(name)
	if errors.Is(err, os.ErrNotExist) {
		return -1, fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	if err != nil {
		return -1, err
	}
	defer file.Close()
	lb, err := ReadLineBreaks(name)
	if err != nil {
		return -1, fmt.Errorf("could not find the line break method: %w", err)
	}
	if !cols {
		cnt, err := nl.Lines(file, lb)
		if err != nil {
			return -1, fmt.Errorf("read lines count the file: %q: %w", name, err)
		}
		return cnt, nil
	}
	cnt, err := Columns(file, lb)
	if err != nil {
		return -1, fmt.Errorf("read lines count the file: %q: %w", name, err)
	}
	return cnt, nil
}

// ReadControls counts the number of ANSI escape sequences in the named file.
func ReadControls(name string) (int, error) {
	file, err := os.Open(name)
	if errors.Is(err, os.ErrNotExist) {
		return -1, fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	if err != nil {
		return -1, err
	}
	defer file.Close()
	cnt, err := Controls(file)
	if err != nil {
		return -1, fmt.Errorf("read controls could not parse the file: %q: %w", name, err)
	}
	return cnt, nil
}

// ReadLine reads a named file location or a named temporary file and returns its content.
func ReadLine(name string, sys nl.System) (string, error) {
	path, n := temp(name), nl.NewLine(sys)
	file, err := os.OpenFile(path, os.O_RDONLY, save.LogFileMode)
	if errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	if err != nil {
		return "", err
	}
	defer file.Close()
	// bufio is the most performant
	s := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s += scanner.Text() + n
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("read line could not scan file: %w", err)
	}
	return s, nil
}

// ReadLines counts the number of lines in the named file.
func ReadLines(name string) (int, error) {
	return readLineBreaks(name, false)
}

// ReadLineBreaks scans the named file for the most commonly used line break method.
func ReadLineBreaks(name string) ([2]rune, error) {
	z := [2]rune{0, 0}
	file, err := os.Open(name)
	if errors.Is(err, os.ErrNotExist) {
		return z, fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	if err != nil {
		return z, err
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		return z, fmt.Errorf("read line breaks could not read the file: %q: %w", name, err)
	}
	return LineBreaks(true, bytes.Runes(b)...), nil
}

// ReadPipe reads data piped by the operating system's STDIN.
// If no data is detected the program will exit.
func ReadPipe() ([]byte, error) {
	b := []byte{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		b = append(b, scanner.Bytes()...)
		b = append(b, []byte("\n")...)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read pipe could not scan stdin: %w", err)
	}
	if len(b) == 0 {
		return nil, ErrPipeEmpty
	}
	return b, nil
}

// ReadRunes returns the number of runes in the named file.
func ReadRunes(name string) (int, error) {
	file, err := os.Open(name)
	if errors.Is(err, os.ErrNotExist) {
		return -1, fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	if err != nil {
		return -1, err
	}
	defer file.Close()
	cnt, err := Runes(file)
	if err != nil {
		return 0, fmt.Errorf("read runes could not calculate this file: %q: %w", name, err)
	}
	return cnt, nil
}

// ReadTail reads the named file from the offset position relative to the end of the file.
func ReadTail(name string, offset int) ([]byte, error) {
	file, err := os.Open(name)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	total, err := ReadRunes(name)
	if err != nil {
		return nil, fmt.Errorf("read tail could not read runes: %q: %w", name, err)
	}
	// bufio is the most performant
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	buf, cnt := []byte{}, 0
	for scanner.Scan() {
		cnt++
		if cnt <= (total - offset) {
			continue
		}
		buf = append(buf, scanner.Bytes()...)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read tail could scan file bytes: %q: %w", name, err)
	}
	return buf, nil
}

// ReadText reads a named file location or a named temporary file and returns its content.
func ReadText(name string) (string, error) {
	return ReadLine(name, nl.Host)
}

// ReadWords counts the number of spaced words in the named file.
func ReadWords(name string) (int, error) {
	file, err := os.Open(name)
	if errors.Is(err, os.ErrNotExist) {
		return -1, fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	if err != nil {
		return -1, err
	}
	defer file.Close()
	cnt, err := Words(file)
	if err != nil {
		return -1, fmt.Errorf("read words failed to count words: %q: %w", name, err)
	}
	return cnt, nil
}
