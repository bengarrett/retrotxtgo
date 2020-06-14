package samples

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

const (
	cp437hex = "\xCD\xB9\xB2\xCC\xCD" // `═╣░╠═`
	utf      = "═╣ ░ ╠═"
	// Newlines sample using operating system defaults
	Newlines = "a\nb\nc...\n"
	// Symbols for Unicode Wingdings
	Symbols = `[☠|☮|♺]`
	// Tabs and Unicode glyphs
	Tabs = "☠\tSkull and crossbones\n\n☮\tPeace symbol\n\n♺\tRecycling"
	// Escapes and control codes.
	Escapes = "bell:\a,back:\b,tab:\t,form:\f,vertical:\v,quote:\""
	// Digits in various formats
	Digits = "\xb0\260\u0170\U00000170"

	base = "rt_sample-"

	be   = unicode.BigEndian
	le   = unicode.LittleEndian
	bom  = unicode.UseBOM
	_bom = unicode.IgnoreBOM

	lf = "\x0a"
	cr = "\x0d"

	permf os.FileMode = 0644
)

// helpful references: How to convert from an encoding to UTF-8 in Go?
// https://stackoverflow.com/questions/32518432/how-to-convert-from-an-encoding-to-utf-8-in-go

// CP437Decode decodes IBM Code Page 437 encoded text.
func CP437Decode(s string) (result []byte, err error) {
	return CPDecode(s, *charmap.CodePage437)
}

// CP437Encode encodes text into IBM Code Page 437.
func CP437Encode(s string) (result []byte, err error) {
	return CPEncode(s, *charmap.CodePage437)
}

// CPDecode decodes simple character encoding text.
func CPDecode(s string, cp charmap.Charmap) (result []byte, err error) {
	decoder := cp.NewDecoder()
	reader := transform.NewReader(strings.NewReader(s), decoder)
	return ioutil.ReadAll(reader)
}

// CPEncode encodes text into a simple character encoding.
func CPEncode(s string, cp charmap.Charmap) (result []byte, err error) {
	if !utf8.Valid([]byte(s)) {
		return result, errors.New("cpencode: string is not encoded as utf8")
	}
	encoder := cp.NewEncoder()
	reader := transform.NewReader(strings.NewReader(s), encoder)
	return ioutil.ReadAll(reader)
}

// HexDecode ...
func HexDecode(hexadecimal string) (result []byte, err error) {
	src := []byte(hexadecimal)
	println(len(src))
	result = make([]byte, hex.DecodedLen(len(src)))
	_, err = hex.Decode(result, src)
	return result, err
}

// HexEncode ..
func HexEncode(text string) (result []byte) {
	src := []byte(text)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}

// ReadBytes reads a named file location or a named temporary file and returns its byte content.
func ReadBytes(name string) (data []byte, err error) {
	var path = tempFile(name)
	file, err := os.OpenFile(path, os.O_RDONLY, permf)
	if err != nil {
		return data, err
	}
	// bufio is the most performant
	scanner := bufio.NewScanner(file)
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

// ReadText reads a named file location or a named temporary file and returns its content.
func ReadText(name string) (text string, err error) {
	return ReadLine(name, "")
}

// nl returns a platform's newline character.
func nl(platform string) string {
	switch platform {
	case "dos", "windows":
		return cr + lf
	case "c64", "darwin", "mac":
		return cr
	case "amiga", "linux", "unix":
		return lf
	default: // use operating system default
		return "\n"
	}
}

// Save bytes to a named file location or a named temporary file.
func Save(b []byte, name string) (path string, err error) {
	path = tempFile(name)
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
	return filepath.Abs(file.Name())
}

func asciiFile() (path string) {
	var abs, err = filepath.Abs("ascii-logos.txt")
	if err != nil {
		log.Fatal(err)
	}
	return abs
}

func ansiFile() (path string) {
	var abs, err = filepath.Abs("ZII-RTXT.ans")
	if err != nil {
		log.Fatal(err)
	}
	return abs
}

func tempFile(name string) (path string) {
	path = name
	if filepath.Base(name) == name {
		path = filepath.Join(os.TempDir(), name)
	}
	return path
}

// these are reminders of how to implement the encoders

// Base64Decode decodes a base64 string.
func Base64Decode(s string) (result []byte, err error) {
	return base64.StdEncoding.DecodeString(s)
}

// Base64Encode encodes a string to base64.
func Base64Encode(s string) (result string) {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// UTF8 transforms string to UTF8 encoding.
func UTF8(s string) (result []byte, length int, err error) {
	return transform.Bytes(unicode.UTF8.NewEncoder(), []byte(s))
}

// UTF16BE transforms string to UTF16 big-endian encoding.
func UTF16BE(s string) (result []byte, length int, err error) {
	return transform.Bytes(unicode.UTF16(be, bom).NewEncoder(), []byte(s))
}

// UTF16LE transforms string to UTF16 little-endian encoding.
func UTF16LE(s string) (result []byte, length int, err error) {
	return transform.Bytes(unicode.UTF16(le, bom).NewEncoder(), []byte(s))
}
