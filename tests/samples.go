package tests

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var perm os.FileMode = 0644

const (
	utf8nl = `☠	Skull and crossbones

☮	Peace symbol

♺	Recycling`
	utf8      = `[☠|☮|♺]`
	toCP437   = "═╣ ░ ╠═"
	fromCP437 = "\xCD\xB9\xB2\xCC\xCD" // `═╣░╠═`

	be   = unicode.BigEndian
	le   = unicode.LittleEndian
	bom  = unicode.UseBOM
	_bom = unicode.IgnoreBOM
)

// helpful references: How to convert from an encoding to UTF-8 in Go?
// https://stackoverflow.com/questions/32518432/how-to-convert-from-an-encoding-to-utf-8-in-go

// BinaryIn ...
func BinaryIn(binary []byte) (text string) {
	return base64.StdEncoding.EncodeToString(binary)
}

// BinaryOut ...
func BinaryOut(encoded string) (result []byte, err error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// CP437Out ..
func CP437Out(text string) (result []byte, err error) {
	encoder := charmap.CodePage437.NewEncoder()
	nr := transform.NewReader(strings.NewReader(text), encoder)
	return ioutil.ReadAll(nr)
}

// CP437In ..
func CP437In(cp437 string) (result []byte, err error) {
	decoder := charmap.CodePage437.NewDecoder()
	nr := transform.NewReader(strings.NewReader(cp437), decoder)
	return ioutil.ReadAll(nr)
}

// HexOut ..
func HexOut(text string) (result []byte) {
	src := []byte(text)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}

// EmbedText ...
func EmbedText() (result string) {
	abs, err := filepath.Abs("../textfiles/ZII-RTXT.ans")
	if err != nil {
		log.Fatal(err)
	}
	print("file:", abs)
	d, err := ReadBinary(abs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("length: %dB\n", len(d))
	return BinaryIn(d)
}

// ReadBinary ..
func ReadBinary(name string) (data []byte, err error) {
	var path = name
	if filepath.Base(name) == name {
		path = filepath.Join(os.TempDir(), name)
	}
	file, err := os.OpenFile(path, os.O_RDONLY, perm)
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

// ReadFile ..
func ReadFile(name string) (text string, err error) {
	var path = name
	if filepath.Base(name) == name {
		path = filepath.Join(os.TempDir(), name)
	}
	file, err := os.OpenFile(path, os.O_RDONLY, perm)
	if err != nil {
		return text, err
	}
	// bufio is the most performant
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text += scanner.Text()
		//fmt.Printf("\n%v\n", text)
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

// SaveFile ..
func SaveFile(b []byte, name string) (path string, len int, err error) {
	if filepath.Base(name) == name {
		path = filepath.Join(os.TempDir(), name)
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return path, len, err
	}
	// bufio is the most performant
	writer := bufio.NewWriter(file)
	for len, c := range b {
		if err = writer.WriteByte(c); err != nil {
			return path, len, err
		}
	}
	if err = writer.Flush(); err != nil {
		return path, len, err
	}
	if err = file.Close(); err != nil {
		return path, len, err
	}
	path, err = filepath.Abs(file.Name())
	return path, len, err
}

// UTF8 ...
func UTF8(text string) (result []byte, len int, err error) {
	return transform.Bytes(unicode.UTF8.NewEncoder(), []byte(text))
}

// UTF16LE ..
func UTF16LE(text string) (result []byte, len int, err error) {
	return transform.Bytes(unicode.UTF16(le, bom).NewEncoder(), []byte(text))
}

// UTF16BE ...
func UTF16BE(text string) (result []byte, len int, err error) {
	return transform.Bytes(unicode.UTF16(be, bom).NewEncoder(), []byte(text))
}
