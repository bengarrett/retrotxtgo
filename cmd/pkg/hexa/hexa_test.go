package hexa_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/pkg/hexa"
	"github.com/stretchr/testify/assert"
)

func ExampleTrimIdent() {
	// unicode example
	s := hexa.TrimIdent("U+00A9")
	fmt.Println(s)

	// retro hex example
	s = hexa.TrimIdent("$0F")
	fmt.Println(s)
	// Output: 00A9
	// 0F
}

func ExampleTrimNCR() {
	s := hexa.TrimNCR("&#169;")
	fmt.Println(s)
	// Output: 169
}

func ExampleTrimIndents() {
	s := hexa.TrimIndents("U+00A9", "$0F")
	fmt.Println(s)
	// Output: [00A9 0F]
}

func ExampleParse() {
	n := hexa.Parse(16, "0F", "01", "FF", "Z")
	fmt.Println(n)

	n = hexa.Parse(10, "15", "1", "255", "abc")
	fmt.Println(n)
	// Output: [15 1 255 -1]
	// [15 1 255 -1]
}

func ExampleParser() {
	_ = hexa.Parser(os.Stdout, 16, "0F", "01", "FF", "Z")
	// Output: 15 1 255 NaN
}

func ExampleWriter() {
	// hexadecimal example
	_ = hexa.Writer(os.Stdout, 16, "0F", "01", "FF", "Z")

	// decimal example
	_ = hexa.Writer(os.Stdout, 10, "15", "1", "255", "abc")
	// Output: 0F = 15  01 = 1  FF = 255  Z = invalid
	// 15 = F  1 = 1  255 = FF  ABC = invalid
}

func TestParse(t *testing.T) {
	t.Parallel()
	x := hexa.Parse(0, nil...)
	assert.Empty(t, x)
	x = hexa.Parse(0, "0F", "01", "FF", "Z")
	assert.Equal(t, []int64{-1, 1, -1, -1}, x)
	x = hexa.Parse(10, "0F", "01", "FF", "Z")
	assert.Equal(t, []int64{-1, 1, -1, -1}, x)
	x = hexa.Parse(16, "0F", "01", "FF", "Z")
	assert.Equal(t, []int64{15, 1, 255, -1}, x)
}

func TestParser(t *testing.T) {
	t.Parallel()
	err := hexa.Parser(nil, 0, nil...)
	assert.Nil(t, err)
	err = hexa.Parser(nil, 2, nil...)
	assert.Nil(t, err)
	err = hexa.Parser(nil, 10, nil...)
	assert.Nil(t, err)
	err = hexa.Parser(nil, 10, "1", "2", "3")
	assert.Nil(t, err)
}

func TestWriter(t *testing.T) {
	t.Parallel()
	err := hexa.Writer(nil, 0, nil...)
	assert.Nil(t, err)
	err = hexa.Writer(nil, 2, nil...)
	assert.Nil(t, err)
	err = hexa.Writer(nil, 10, nil...)
	assert.Nil(t, err)
	err = hexa.Writer(nil, 10, "1", "2", "3")
	assert.Nil(t, err)
}
