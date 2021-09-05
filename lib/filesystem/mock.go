package filesystem

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
)

const windows = "windows"

// MockInput uses the os pipe to mock the user input.
// os.Pipe() https://stackoverflow.com/questions/46365221/fill-os-stdin-for-function-that-reads-from-it
/*
	Usage:
	r, err := filesystem.MockInput(tt.args.input)
	if err != nil {
		t.Error(err)
	}
	stdin := os.Stdin
	defer func() {
		os.Stdin = stdin
	}()
	os.Stdin = r
*/
func MockInput(input string) (*os.File, error) {
	s := []byte(input)
	r, w, err := os.Pipe()
	if err != nil {
		return r, err
	}
	_, err = w.Write(s)
	if err != nil {
		return r, err
	}
	w.Close()
	return r, nil
}

// T are short ASCII and Unicode strings used in various unit tests.
func T() map[string]string {
	return map[string]string{
		// Newline sample using yjr operating system defaults
		"Newline": "a\nb\nc...\n",
		// Symbols for Unicode Wingdings
		"Symbols": `[☠|☮|♺]`,
		// Tabs and Unicode glyphs
		"Tabs": "☠\tSkull and crossbones\n\n☮\tPeace symbol\n\n♺\tRecycling",
		// Escapes and control codes.
		"Escapes": "bell:\a,back:\b,tab:\t,form:\f,vertical:\v,quote:\"",
		// Digits in various formats
		"Digits": "\xb0\260\u0170\U00000170",
	}
}

// fileExample saves the string to a numbered text file.
func fileExample(s string, i int) string {
	name := fmt.Sprintf("rt_fs_save%d.txt", i)
	path, err := SaveTemp(name, []byte(s)...)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

// largeExample generates and saves a 800k file of random us-ascii text.
func largeExample() string {
	const name, sizeMB = "rs_mega_example_save.txt", 0.8
	_, s := filler(sizeMB)
	path, err := SaveTemp(name, []byte(s)...)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

// megaExample generates and saves a 1.5MB file of random us-ascii text.
func megaExample() string {
	const name, sizeMB = "rs_giga_mega_save.txt", 1.5
	_, s := filler(sizeMB)
	path, err := SaveTemp(name, []byte(s)...)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

// filler generates random us-ascii text.
func filler(sizeMB float64) (length int, random string) {
	if sizeMB <= 0 {
		return length, random
	}
	// make characters to randomize
	const (
		// ascii code points (rune codes)
		start    = 33  // "!"
		end      = 122 // "z"
		charsLen = end - start + 1
	)
	chars := make([]rune, charsLen)
	for c, i := 0, start; i <= end; i++ {
		chars[c] = rune(i)
		c++
	}
	// initialize rune slice
	const base, exp = 1000, 2
	f := (math.Pow(base, exp) * sizeMB)
	s := make([]rune, int(f))
	// generate random string
	for i := range s {
		s[i] = chars[rand.Intn(charsLen)] // nolint:gosec
	}
	return len(s), string(s)
}

type dirTests []struct {
	name    string
	wantDir string
}

// nolint:dupl
func windowsTests(h, hp, s, w, wp string) dirTests {
	return dirTests{
		{fmt.Sprintf("C:%shome%suser", s, s), fmt.Sprintf("C:%shome%suser", s, s)},
		{"~", h},
		{filepath.Join("~", "foo"), filepath.Join(h, "foo")},
		{".", w},
		{fmt.Sprintf(".%sfoo", s), filepath.Join(w, "foo")},
		{fmt.Sprintf("..%sfoo", s), filepath.Join(wp, "foo")},
		{fmt.Sprintf("~%s..%sfoo", s, s), filepath.Join(hp, "foo")},
		{fmt.Sprintf("d:%sroot%sfoo%s..%sblah", s, s, s, s), fmt.Sprintf("D:%sroot%sblah", s, s)},
		{fmt.Sprintf("z:%sroot%sfoo%s.%sblah", s, s, s, s), fmt.Sprintf("Z:%sroot%sfoo%sblah", s, s, s)},
	}
}

// nolint:dupl
func nixTests(h, hp, s, w, wp string) dirTests {
	return dirTests{
		{fmt.Sprintf("%shome%suser", s, s), fmt.Sprintf("%shome%suser", s, s)},
		{"~", h},
		{filepath.Join("~", "foo"), filepath.Join(h, "foo")},
		{".", w},
		{fmt.Sprintf(".%sfoo", s), filepath.Join(w, "foo")},
		{fmt.Sprintf("..%sfoo", s), filepath.Join(wp, "foo")},
		{fmt.Sprintf("~%s..%sfoo", s, s), filepath.Join(hp, "foo")},
		{fmt.Sprintf("%sroot%sfoo%s..%sblah", s, s, s, s), fmt.Sprintf("%sroot%sblah", s, s)},
		{fmt.Sprintf("%sroot%sfoo%s.%sblah", s, s, s, s), fmt.Sprintf("%sroot%sfoo%sblah", s, s, s)},
	}
}
