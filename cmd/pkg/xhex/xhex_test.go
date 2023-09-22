package xhex_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/pkg/xhex"
	"github.com/stretchr/testify/assert"
)

func ExampleTrimIdent() {
	// unicode example
	s := xhex.TrimIdent("U+00A9")
	fmt.Println(s)

	// retro hex example
	s = xhex.TrimIdent("$0F")
	fmt.Println(s)
	// Output: 00A9
	// 0F
}

func ExampleTrimNCR() {
	s := xhex.TrimNCR("&#169;")
	fmt.Println(s)
	// Output: 169
}

func ExampleTrimIndents() {
	s := xhex.TrimIndents("U+00A9", "$0F")
	fmt.Println(s)
	// Output: [00A9 0F]
}

func ExampleParse() {
	n := xhex.Parse(16, "0F", "01", "FF", "Z")
	fmt.Println(n)

	n = xhex.Parse(10, "15", "1", "255", "abc")
	fmt.Println(n)
	// Output: [15 1 255 -1]
	// [15 1 255 -1]
}

func ExampleRaw() {
	_ = xhex.Raw(os.Stdout, 16, "0F", "01", "FF", "Z")
	// Output: 15 1 255 NaN
}

func ExampleWrite() {
	// hexadecimal example
	_ = xhex.Write(os.Stdout, 16, "0F", "01", "FF", "Z")

	// decimal example
	_ = xhex.Write(os.Stdout, 10, "15", "1", "255", "abc")
	// Output: 0F = 15  01 = 1  FF = 255  Z = invalid
	// 15 = F  1 = 1  255 = FF  ABC = invalid
}

func TestParse(t *testing.T) {
	t.Parallel()
	x := xhex.Parse(0, nil...)
	assert.Empty(t, x)
	x = xhex.Parse(0, "0F", "01", "FF", "Z")
	assert.Equal(t, []int64{-1, 1, -1, -1}, x)
	x = xhex.Parse(10, "0F", "01", "FF", "Z")
	assert.Equal(t, []int64{-1, 1, -1, -1}, x)
	x = xhex.Parse(16, "0F", "01", "FF", "Z")
	assert.Equal(t, []int64{15, 1, 255, -1}, x)
}

func TestRaw(t *testing.T) {
	t.Parallel()
	err := xhex.Raw(nil, 0, nil...)
	assert.Nil(t, err)
	err = xhex.Raw(nil, 2, nil...)
	assert.Nil(t, err)
	err = xhex.Raw(nil, 10, nil...)
	assert.Nil(t, err)
	err = xhex.Raw(nil, 10, "1", "2", "3")
	assert.Nil(t, err)
}

func TestWrite(t *testing.T) {
	t.Parallel()
	err := xhex.Write(nil, 0, nil...)
	assert.Nil(t, err)
	err = xhex.Write(nil, 2, nil...)
	assert.Nil(t, err)
	err = xhex.Write(nil, 10, nil...)
	assert.Nil(t, err)
	err = xhex.Write(nil, 10, "1", "2", "3")
	assert.Nil(t, err)
}
