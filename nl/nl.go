// Package nl provides line break characters for multiple system and microcomputer platforms.
package nl

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var ErrReader = errors.New("the r reader cannot be nil")

// Common ASCII and EBCDIC control codes for new lines.
const (
	LF  rune = 10  // LF is the control code for a line feed.
	CR  rune = 13  // CR is the control code for a carriage return.
	NL  rune = 21  // NL is the control code for a new line, an EBCDIC control.
	NEL rune = 133 // NEL is the control code for a next line.
)

// LineBreak contains details on the line break sequence used to create a new line in a text file.
type LineBreak struct {
	Abbr    string  `json:"string"  xml:"string,attr"` // Abbr is the abbreviation for the line break.
	Escape  string  `json:"escape"  xml:"-"`           // Escape is the escape sequence for the line break.
	Decimal [2]rune `json:"decimal" xml:"decimal"`     // Decimal is the numeric character code for the line break.
}

// Find determines the new lines characters found in the rune pair.
func (lb *LineBreak) Find(r [2]rune) {
	a, e := "", ""
	switch r {
	case [2]rune{LF}:
		a = "lf"
		e = "\n"
	case [2]rune{CR}:
		a = "cr"
		e = "\r"
	case [2]rune{CR, LF}:
		a = "crlf"
		e = "\r\n"
	case [2]rune{LF, CR}:
		a = "lfcr"
		e = "\n\r"
	case [2]rune{NL}, [2]rune{NEL}:
		a = "nl"
		e = "\025"
	}
	lb.Decimal = r
	lb.Abbr = strings.ToUpper(a)
	lb.Escape = e
}

// Total counts the number of lines in the named file
// based on the provided line break sequence.
func (lb *LineBreak) Total(name string) (int, error) {
	f, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	l, err := Lines(f, lb.Decimal)
	if err != nil {
		return 0, err
	}
	return l, nil
}

// Lines counts the number of lines in the interface.
// The lb rune pair is the line break sequence.
// If the line break only has one rune, then the second rune should be 0.
func Lines(r io.Reader, lb [2]rune) (int, error) {
	if r == nil {
		return 0, ErrReader
	}
	lineBreak := []byte{byte(lb[0]), byte(lb[1])}
	if lb[1] == 0 {
		lineBreak = []byte{byte(lb[0])}
	}
	p, count := make([]byte, bufio.MaxScanTokenSize), 0
	for {
		size, err := r.Read(p)
		if err != nil && err != io.EOF {
			return 0, fmt.Errorf("lines could not read buffer: %w", err)
		}
		pos := 0
		for {
			i := bytes.Index(p[pos:], lineBreak)
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

// System is a micro or computer platform.
type System int

const (
	Host      System = iota // Host is the host operating system default.
	Acorn                   // Acorn is the Acorn Archimedes.
	Amiga                   // Amiga is the Commodore Amiga.
	Commodore               // Commodore is the Commodore 64 and compatibles.
	Darwin                  // Darwin is Apple macOS.
	Linux                   // Linux is the Linux operating system.
	Macintosh               // Macintosh is the Apple Macintosh.
	PCDos                   // PCDos is an IBM PC and Microsoft DOS and compatibles.
	Unix                    // Unix are both Unix and BSD.
	Windows                 // Windows is Microsoft Windows.
)

// NewLine returns a new line or line break characters for the system platform.
func NewLine(s System) string {
	switch s {
	case PCDos, Windows:
		return string(CR) + string(LF)
	case Commodore, Darwin, Macintosh:
		return string(CR)
	case Amiga, Linux, Unix:
		return string(LF)
	case Acorn:
		return string(LF) + string(CR)
	case Host:
		return "\n"
	}
	return ""
}
