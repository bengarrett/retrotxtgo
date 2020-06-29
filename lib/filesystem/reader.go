package filesystem

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	"unicode"
	"unicode/utf8"
)

// Columns counts the number of characters used per line in the reader interface.
func Columns(r io.Reader, nl [2]rune) (width int, err error) {
	var lineBreak = []byte{byte(nl[0]), byte(nl[1])}
	if nl[1] == 0 {
		lineBreak = []byte{byte(nl[0])}
	}
	buf := make([]byte, bufio.MaxScanTokenSize)
	for {
		size, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return -1, err
		}
		var pos int
		for {
			if size == pos {
				break
			}
			i := bytes.Index(buf[pos:], lineBreak)
			if i == -1 {
				break
			}
			pos += i + len(lineBreak)
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
			if size == pos {
				break
			}
			if i == -1 {
				return count, nil
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
			if size == pos {
				break
			}
			if i == -1 {
				// when no linebreaks are used
				// return 0 lines for an empty buffer or 1 line of text
				if size == 0 {
					return 0, nil
				}
				return 1, nil
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

// Newlines will try to guess the newline representation as a 2 byte value.
// A guess of Unix will return [10, 0], Windows [13, 10], otherwise a [0, 0] value is returned.
func Newlines(runes []rune) [2]rune {
	// scan data for possible newlines
	c := []struct {
		abbr  string
		count int
	}{
		{"lf", 0},   // linux, unix, amiga...
		{"cr", 0},   // 8-bit micros & legacy mac
		{"crlf", 0}, // windows, dos, cp/m...
		{"lfcr", 0}, // acorn bbc micro
		{"nl", 0},   // ibm ebcdic encodings
	}
	l := len(runes) - 1 // range limit
	for i, r := range runes {
		switch r {
		case 10:
			if i < l && runes[i+1] == 13 {
				c[3].count++ // lfcr
				continue
			}
			if i != 0 && runes[i-1] == 13 {
				// crlf (already counted)
				continue
			}
			c[0].count++
		case 13:
			if i < l && runes[i+1] == 10 {
				c[2].count++ // crlf
				continue
			}
			if i != 0 && runes[i-1] == 10 {
				// lfcr (already counted)
				continue
			}
			// carriage return on modern terminals will overwrite the existing line of text
			// todo: add flag or change behaviour to replace CR (\r) with NL (\n)
			c[1].count++
		case 21:
			c[4].count++
		case 155:
			// atascii (not currently used)
		}
	}
	// sort results
	sort.SliceStable(c, func(i, j int) bool {
		return c[i].count > c[j].count
	})
	fmt.Println("runes:", len(runes), c, c[0].abbr)
	switch c[0].abbr {
	case "lf":
		return [2]rune{10}
	case "cr":
		return [2]rune{13}
	case "crlf":
		return [2]rune{13, 10}
	case "lfcr":
		return [2]rune{10, 13}
	case "nl":
		return [2]rune{21}
	}
	return [2]rune{}
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
