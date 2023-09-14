package xud_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/xud"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding/charmap"
)

func ExampleEncoding() {
	fmt.Println(xud.XUserDefined1963)
	fmt.Println(xud.XUserDefined1965)
	fmt.Println(xud.XUserDefined1967)
	// Output: ASA X3.4 1963
	// ASA X3.4 1965
	// ANSI X3.4 1967/77/86
}

func ExampleName() {
	fmt.Println(xud.Name(xud.XUserDefined1963))
	fmt.Println(xud.Name(xud.XUserDefined1965))
	fmt.Println(xud.Name(xud.XUserDefined1967))
	// Output: ascii-63
	// ascii-65
	// ascii-67
}

func ExampleNumeric() {
	fmt.Println(xud.Numeric(xud.XUserDefined1963))
	fmt.Println(xud.Numeric(xud.XUserDefined1965))
	fmt.Println(xud.Numeric(xud.XUserDefined1967))
	// Output: 1963
	// 1965
	// 1967
}

func ExampleAlias() {
	fmt.Println(xud.Alias(xud.XUserDefined1963))
	fmt.Println(xud.Alias(xud.XUserDefined1965))
	fmt.Println(xud.Alias(xud.XUserDefined1967))
	// Output:
	//
	//
	// ansi
}

func ExampleChar() {
	const code = 94
	r := xud.Char(xud.XUserDefined1963, code)
	fmt.Printf("code %d is the rune %s (%v)\n", code, string(r), r)
	r = xud.Char(xud.XUserDefined1965, code)
	fmt.Printf("code %d is skipped (%v)\n", code, r)
	// Output: code 94 is the rune ↑ (8593)
	// code 94 is skipped (-1)
}

func TestCode7bit(t *testing.T) {
	b := xud.Code7bit(nil)
	assert.False(t, b)
	b = xud.Code7bit(charmap.CodePage037)
	assert.False(t, b)
	b = xud.Code7bit(xud.XUserDefined1963)
	assert.True(t, b)
}

func TestFootnote(t *testing.T) {
	s := &strings.Builder{}
	xud.Footnote(s, xud.XUserDefined1963)
	assert.Contains(t, s.String(), "* ASA X3.4 1963")
}

func TestCharX3465(t *testing.T) {
	const skip = int32(-1)
	r := xud.Char(xud.XUserDefined1965, 94)
	assert.Equal(t, skip, r)
	r = xud.Char(xud.XUserDefined1965, 92)
	assert.Equal(t, int32(126), r)
}

func TestCharX3467(t *testing.T) {
	const skip = int32(-1)
	r := xud.Char(xud.XUserDefined1967, 94)
	assert.Equal(t, skip, r)
	r = xud.Char(xud.XUserDefined1967, 130)
	assert.Equal(t, int32(32), r)
}
