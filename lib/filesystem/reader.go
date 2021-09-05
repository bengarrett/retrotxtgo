package filesystem

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// LB is the text line break control represented as 2 runes.
type LB [2]rune

const (
	ansiEscape string = "\x1B\x5b" // esc[
	// Linefeed is a Linux/macOS line break.
	Linefeed rune = 10
	// CarriageReturn is a partial line break for Windows/DOS.
	CarriageReturn rune = 13
	// NewLine EBCDIC control.
	NewLine rune = 21
	// NextLine EBCDIC control in UTF-8 documents.
	NextLine rune = 133
)

// LF linefeed.
func LF() LB {
	return LB{Linefeed}
}

// CR carriage return.
func CR() LB {
	return LB{CarriageReturn}
}

// CRLF carriage return + linefeed.
func CRLF() LB {
	return LB{CarriageReturn, Linefeed}
}

// LFCR linefeed + carriage return.
func LFCR() LB {
	return LB{Linefeed, CarriageReturn}
}

// NL new line.
func NL() LB {
	return LB{NewLine}
}

// NEL next line.
func NEL() LB {
	return LB{NextLine}
}

var ErrLB = errors.New("linebreak runes cannot be empty")

// Columns counts the number of characters used per line in the reader interface.
func Columns(r io.Reader, lb LB) (int, error) {
	if reflect.DeepEqual(lb, LB{}) {
		return 0, ErrLB
	}
	var lineBreak = []byte{byte(lb[0]), byte(lb[1])}
	if lb[1] == 0 {
		lineBreak = []byte{byte(lb[0])}
	}
	buf, width := make([]byte, bufio.MaxScanTokenSize), 0
	for {
		size, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return -1, fmt.Errorf("columns could not read buffer: %w", err)
		}
		pos := 0
		for {
			if size == pos {
				break
			}
			i := bytes.Index(buf[pos:], lineBreak)
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
	lineBreak := []byte(ansiEscape)
	buf, count := make([]byte, bufio.MaxScanTokenSize), 0
	for {
		size, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, fmt.Errorf("controls could not read buffer: %w", err)
		}
		pos := 0
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
func Lines(r io.Reader, lb LB) (int, error) {
	lineBreak := []byte{byte(lb[0]), byte(lb[1])}
	if lb[1] == 0 {
		lineBreak = []byte{byte(lb[0])}
	}
	buf, count := make([]byte, bufio.MaxScanTokenSize), 0
	for {
		size, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, fmt.Errorf("lines could not read buffer: %w", err)
		}
		pos := 0
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
					return 1, nil // no line breaks = 1 line
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

// LineBreaks will try to guess the line break representation as a 2 byte value.
// A guess of Unix will return [10, 0], Windows [13, 10], otherwise a [0, 0] value is returned.
func LineBreaks(utf bool, runes ...rune) LB {
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
		case Linefeed:
			c[0].count = lfCnt(c[0].count, i, l, runes...)
		case CarriageReturn:
			if i < l && runes[i+1] == Linefeed {
				c[2].count++ // crlf
				continue
			}
			if i != 0 && runes[i-1] == Linefeed {
				// lfcr (already counted)
				continue
			}
			// carriage return on modern terminals will overwrite the existing line of text
			// todo: add flag or change behavor to replace CR (\r) with NL (\n)
			c[1].count++
		case NewLine, NextLine:
			if utf && r == NextLine {
				c[4].count++ // NL
				continue
			}
			if r == NewLine {
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

func abbr(utf bool, s string) LB {
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
	return LB{}
}

// LineBreak humanizes the value of LineBreaks().
func LineBreak(r LB, extraInfo bool) string {
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
	count := 0
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
	if err := scanner.Err(); err != nil {
		return -1, fmt.Errorf("words could not scan reader: %w", err)
	}
	return count, nil
}

// WordsEBCDIC counts the number of spaced words in the EBCDIC encoded reader interface.
func WordsEBCDIC(r io.Reader) (int, error) {
	// for the purposes of counting words, any EBCDIC codepage is fine
	c := transform.NewReader(r, charmap.CodePage037.NewDecoder())
	return Words(c)
}

// IsWord scans the content of a word for characters that are not digits,
// letters or punctuation and if discovered returns false.
// If a space or line break is encountered the scan will end.
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
