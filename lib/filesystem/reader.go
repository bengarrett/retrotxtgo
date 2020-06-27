package filesystem

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
	"unicode/utf8"
)

// Columns counts the number of characters used per line in the reader interface.
func Columns(r io.Reader) (width int, err error) {
	const lineBreak = '\n'
	buf := make([]byte, bufio.MaxScanTokenSize)
	for {
		size, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return -1, err
		}
		var pos int
		for {
			i := bytes.IndexByte(buf[pos:], lineBreak)
			if i == -1 {
				// when no linebreaks are used
				return -1, err
			}
			if size == pos {
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

// Controls counts the number of ANSI escape controls in the reader interface.
func Controls(r io.Reader) (count int, err error) {
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

// Lines counts the number of lines in the interface.
func Lines(r io.Reader) (count int, err error) {
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
			if i == -1 {
				// when no linebreaks are used
				// return 0 lines for an empty buffer or 1 line of text
				if size == 0 {
					return 0, nil
				}
				return 1, nil
			}
			if size == pos {
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

// Runes returns the number of runes in the reader interface.
func Runes(r io.Reader) (count int, err error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		count++
	}
	return count, nil
}

// Words counts the number of spaced words in the reader interface.
func Words(r io.Reader) (count int, err error) {
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
		if word(t) {
			count++
		}
	}
	if err = scanner.Err(); err != nil {
		return -1, err
	}
	return count, err
}

// isWord scans the content of a word for characters that are not digits,
// letters or punctuation and if discovered returns false.
// If a space or newline is encountered the scan will end.
func word(s string) bool {
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
