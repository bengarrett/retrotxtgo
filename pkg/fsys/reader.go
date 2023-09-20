package fsys

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/bengarrett/retrotxtgo/pkg/nl"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const ansiEscape string = "\x1B\x5b" // equals runes 27 and 91 or "ESC["

// LF linefeed.
func LF() [2]rune {
	return [2]rune{nl.LF}
}

// CR carriage return.
func CR() [2]rune {
	return [2]rune{nl.CR}
}

// CRLF carriage return + linefeed.
func CRLF() [2]rune {
	return [2]rune{nl.CR, nl.LF}
}

// LFCR linefeed + carriage return.
func LFCR() [2]rune {
	return [2]rune{nl.LF, nl.CR}
}

// NL new line.
func NL() [2]rune {
	return [2]rune{nl.NL}
}

// NEL next line.
func NEL() [2]rune {
	return [2]rune{nl.NEL}
}

// Columns counts the number of characters used per line in the reader interface.
func Columns(r io.Reader, lb [2]rune) (int, error) {
	if r == nil {
		return 0, ErrReader
	}
	if reflect.DeepEqual(lb, [2]rune{}) {
		return 0, ErrLB
	}
	lineBreak := []byte{byte(lb[0]), byte(lb[1])}
	if lb[1] == 0 {
		lineBreak = []byte{byte(lb[0])}
	}
	p, width := make([]byte, bufio.MaxScanTokenSize), 0
	for {
		size, err := r.Read(p)
		if err != nil && err != io.EOF {
			return -1, fmt.Errorf("columns could not read buffer: %w", err)
		}
		pos := 0
		for {
			if size == pos {
				break
			}
			i := bytes.Index(p[pos:], lineBreak)
			if i == -1 {
				width = size
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
func Controls(r io.Reader) (int, error) {
	if r == nil {
		return 0, ErrReader
	}
	lineBreak := []byte(ansiEscape)
	p, count := make([]byte, bufio.MaxScanTokenSize), 0
	for {
		size, err := r.Read(p)
		if err != nil && err != io.EOF {
			return 0, fmt.Errorf("controls could not read buffer: %w", err)
		}
		pos := 0
		for {
			i := bytes.Index(p[pos:], lineBreak)
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

// LineBreaks will try to guess the line break representation as a 2 byte value.
// A guess of Unix will return [10, 0], Windows [13, 10], otherwise a [0, 0] value is returned.
func LineBreaks(utf bool, runes ...rune) [2]rune {
	// scan data for possible line breaks
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
		case nl.LF:
			c[0].count = lfCnt(c[0].count, i, l, runes...)
		case nl.CR:
			if i < l && runes[i+1] == nl.LF {
				c[2].count++ // crlf
				continue
			}
			if i != 0 && runes[i-1] == nl.LF {
				// lfcr (already counted)
				continue
			}
			// carriage return on modern terminals will overwrite the existing line of text
			// todo: add flag or change behavor to replace CR (\r) with NL (\n)
			c[1].count++
		case nl.NL, nl.NEL:
			if utf && r == nl.NEL {
				c[4].count++ // NL
				continue
			}
			if r == nl.NL {
				c[4].count++ // NEL
				continue
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
	if i < l && runes[i+1] == nl.CR {
		c++ // lfcr
		return c
	}
	if i != 0 && runes[i-1] == nl.CR {
		// crlf (already counted)
		return c
	}
	c++
	return c
}

func abbr(utf bool, s string) [2]rune {
	switch s {
	case "lf":
		return LF()
	case "cr":
		return CR()
	case "crlf":
		return CRLF()
	case "lfcr":
		return LFCR()
	case "nl":
		if utf {
			return NEL()
		}
		return NL()
	}
	return [2]rune{}
}

// LineBreak humanizes the value of LineBreaks().
func LineBreak(r [2]rune, extraInfo bool) string {
	if !extraInfo {
		switch r {
		case LF():
			return "LF"
		case CR():
			return "CR"
		case CRLF():
			return "CRLF"
		case LFCR():
			return "LFCR"
		case NL():
			return "NL"
		case NEL():
			return "NEL"
		}
	}
	switch r {
	case LF():
		return fmt.Sprintf("LF (%s)", lfNix())
	case CR():
		return "CR (8-bit microcomputers)"
	case CRLF():
		return "CRLF (Windows, DOS)"
	case LFCR():
		return "LFCR (Acorn BBC)"
	case NL():
		return "NL (IBM EBCDIC)"
	case NEL():
		return "NEL (EBCDIC to Unicode)"
	}
	return "??"
}

func lfNix() string {
	const mac, linux, unix = "macOS", "Linux", "Unix"
	s := strings.Join([]string{mac, linux, unix}, ", ")
	switch runtime.GOOS {
	case "linux":
		s = linux
	case "darwin":
		s = mac
	case "dragonfly", "illumos", "solaris":
		s = unix
	default:
		if strings.HasSuffix(runtime.GOOS, "bsd") {
			s = unix
		}
	}
	return s
}

// Runes returns the number of runes in the reader interface.
func Runes(r io.Reader) (int, error) {
	if r == nil {
		return 0, ErrReader
	}
	count := 0
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
func Words(r io.Reader) (int, error) {
	if r == nil {
		return 0, ErrReader
	}
	count := 0
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		t := scanner.Text()
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
		if Word(t) {
			count++
		}
	}
	if err := scanner.Err(); err != nil {
		return -1, fmt.Errorf("words could not scan reader: %w", err)
	}
	return count, nil
}

// WordsEBCDIC counts the number of spaced words in the EBCDIC encoded reader interface.
func WordsEBCDIC(r io.Reader) (int, error) {
	if r == nil {
		return 0, ErrReader
	}
	// for the purposes of counting words, any EBCDIC code page is fine
	c := transform.NewReader(r, charmap.CodePage037.NewDecoder())
	return Words(c)
}

// Word reports whether content of s contains only characters
// that are comprised of digits, letters and punctuation.
// If a space or line break is encountered the scan ends.
func Word(s string) bool {
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
