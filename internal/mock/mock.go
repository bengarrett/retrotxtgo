package mock

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	mr "math/rand"
	"os"
	"path/filepath"

	"github.com/bengarrett/retrotxtgo/internal/save"
	"github.com/bengarrett/retrotxtgo/internal/tmp"
)

var (
	ErrMax = errors.New("maxpow value cannot be less than or equal to zero")
	ErrZB  = errors.New("zero bytes written")
)

const (
	xPow = 1000
	yPow = 2
)

// Input returns a file pointer to a temporary file containing the input string.
func Input(input string) (*os.File, error) {
	s := []byte(input)
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	defer w.Close()
	if _, err = w.Write(s); err != nil {
		return nil, err
	}
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

// maxPow calculates and returns the maximum value of yPow raised to the power of xPow.
// It always confirms that the value is safe to use with rand.Int,
// that panics if the y value is less than or equal to zero.
func maxPow() *big.Int {
	y := big.NewInt(int64(math.Pow(xPow, yPow)))
	x := big.NewInt(1)
	if x.Cmp(y) == 1 {
		log.Fatalf("%s: %s", ErrMax, y)
	}
	return y
}

// FileExample saves the string to a randomized, threadsafe filename.
// The path to the file is returned.
func FileExample(s string) string {
	v, err := rand.Int(rand.Reader, maxPow())
	if err != nil {
		log.Fatal(err)
	}
	name := fmt.Sprintf("rt_fs_save%s.txt", v)
	path, err := SaveTemp(name, []byte(s)...)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

// LargeExample generates and saves a 800k file of filler us-ascii text
// to a randomized, threadsafe filename. The path to the file is returned.
func LargeExample() string {
	v, err := rand.Int(rand.Reader, maxPow())
	if err != nil {
		log.Fatal(err)
	}
	const sizeMB = 0.8
	name := fmt.Sprintf("rs_mega_example_save%s.txt", v)
	s := Filler(sizeMB)
	path, err := SaveTemp(name, []byte(s)...)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

// MegaExample generates and saves a 1.5MB file of filler us-ascii text
// to a randomized, threadsafe filename. The path to the file is returned.
func MegaExample() string {
	v, err := rand.Int(rand.Reader, maxPow())
	if err != nil {
		log.Fatal(err)
	}
	const sizeMB = 1.5
	name := fmt.Sprintf("rs_giga_mega_save%s.txt", v)
	s := Filler(sizeMB)
	path, err := SaveTemp(name, []byte(s)...)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

// ByteExample saves a string of multibyte Unicode characters, tabs and
// newlines to a randomized, threadsafe filename. The path to the file is
// returned.
func ByteExample() string {
	v, err := rand.Int(rand.Reader, maxPow())
	if err != nil {
		log.Fatal(err)
	}
	name := fmt.Sprintf("rs_byte_chars_save%s.txt", v)
	b := []byte(T()["Tabs"]) // Tabs and Unicode glyphs
	path, err := SaveTemp(name, b...)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

// Filler generates random us-ascii text.
func Filler(sizeMB float64) string {
	if sizeMB <= 0 {
		return ""
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
		s[i] = chars[mr.Intn(charsLen)]
	}
	return string(s)
}

type DirTests []struct {
	Name    string
	WantDir string
}

func WindowsTests(h, hp, s, w, wp string) DirTests {
	return DirTests{
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

func NixTests(h, hp, s, w, wp string) DirTests {
	return DirTests{
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

// SaveTemp saves bytes to a named temporary file.
func SaveTemp(name string, b ...byte) (string, error) {
	i, path, err := save.Save(tmp.File(name), b...)
	if err != nil {
		return path, fmt.Errorf("could not save the temporary file: %w", err)
	}
	if i == 0 && len(b) > 0 {
		return path, fmt.Errorf("%w: %s", ErrZB, path)
	}
	return path, nil
}
