package filesystem

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
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
			return -1, fmt.Errorf("columns could not read buffer: %w", err)
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
			return 0, fmt.Errorf("controls could not read buffer: %w", err)
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
func Lines(r io.Reader, nl [2]rune) (count int, err error) {
	var lineBreak = []byte{byte(nl[0]), byte(nl[1])}
	if nl[1] == 0 {
		lineBreak = []byte{byte(nl[0])}
	}
	buf := make([]byte, bufio.MaxScanTokenSize)
	for {
		size, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, fmt.Errorf("lines could not read buffer: %w", err)
		}
		var pos int
		for {
			i := bytes.Index(buf[pos:], lineBreak)
			if size == pos {
				break
			}
			if i == -1 {
				if size == 0 {
					return 0, nil // empty file
				}
				if count == 0 {
					return 1, nil // no newlines = 1 line
				}
				count++
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

// Newlines will try to guess the newline representation as a 2 byte value.
// A guess of Unix will return [10, 0], Windows [13, 10], otherwise a [0, 0] value is returned.
func Newlines(utf bool, runes ...rune) [2]rune {
	const (
		lf     = 10
		cr     = 13
		nl     = 21
		nlutf8 = 133
	)
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
		case lf:
			c[0].count = lfCnt(c[0].count, i, l, runes...)
		case cr:
			if i < l && runes[i+1] == lf {
				c[2].count++ // crlf
				continue
			}
			if i != 0 && runes[i-1] == lf {
				// lfcr (already counted)
				continue
			}
			// carriage return on modern terminals will overwrite the existing line of text
			// todo: add flag or change behaviour to replace CR (\r) with NL (\n)
			c[1].count++
		case nl, nlutf8:
			if utf && r == nlutf8 {
				c[4].count++ // NL as utf8
			} else if r == nl {
				c[4].count++ // NL as ebcdic
			}
		}
	}
	// sort results
	sort.SliceStable(c, func(i, j int) bool {
		return c[i].count > c[j].count
	})
	return abbr(utf, c[0].abbr)
}

func lfCnt(c, i, l int, runes ...rune) int {
	const cr = 13
	if i < l && runes[i+1] == cr {
		c++ // lfcr
		return c
	}
	if i != 0 && runes[i-1] == cr {
		// crlf (already counted)
		return c
	}
	c++
	return c
}

func abbr(utf bool, s string) [2]rune {
	switch s {
	case "lf":
		return [2]rune{10}
	case "cr":
		return [2]rune{13}
	case "crlf":
		return [2]rune{13, 10}
	case "lfcr":
		return [2]rune{10, 13}
	case "nl":
		if utf {
			return [2]rune{133}
		}
		return [2]rune{21}
	}
	return [2]rune{}
}

// Newline humanizes the value of Newlines().
func Newline(r [2]rune, extraInfo bool) string {
	if !extraInfo {
		switch r {
		case [2]rune{10}:
			return "LF"
		case [2]rune{13}:
			return "CR"
		case [2]rune{13, 10}:
			return "CRLF"
		case [2]rune{10, 13}:
			return "LFCR"
		case [2]rune{21}, [2]rune{133}:
			return "NL"
		}
	}
	switch r {
	case [2]rune{10}:
		return "LF (Linux, macOS, Unix)"
	case [2]rune{13}:
		return "CR (8-bit microcomputers)"
	case [2]rune{13, 10}:
		return "CRLF (Windows, DOS)"
	case [2]rune{10, 13}:
		return "LFCR (Acorn BBC)"
	case [2]rune{21}, [2]rune{133}:
		return "NL (IBM EBCDIC)"
	}
	return "??"
}

// Runes returns the number of runes in the reader interface.
func Runes(r io.Reader) (count int, err error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return -1, fmt.Errorf("runes could not scan reader: %w", err)
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
		const max = 65533
		if r >= max {
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
		return -1, fmt.Errorf("words could not scan reader: %w", err)
	}
	return count, nil
}

// WordsEBCDIC counts the number of spaced words in the EBCDIC encoded reader interface.
func WordsEBCDIC(r io.Reader) (count int, err error) {
	// for the purposes of counting words, any EBCDIC codepage is fine
	c := transform.NewReader(r, charmap.CodePage037.NewDecoder())
	return Words(c)
}

// isWord scans the content of a word for characters that are not digits,
// letters or punctuation and if discovered returns false.
// If a space or newline is encountered the scan will end.
func word(s string) bool {
	if s == "" {
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
