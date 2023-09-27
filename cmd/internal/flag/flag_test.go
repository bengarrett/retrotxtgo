package flag_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/convert"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding/charmap"
)

func ExampleEndOfFile() {
	var f convert.Flag
	f.Controls = []string{"eof"}
	fmt.Fprint(os.Stdout, flag.EndOfFile(f))
	// Output: true
}

func TestDefault(t *testing.T) {
	t.Parallel()
	e := flag.Default()
	assert.Equal(t, charmap.CodePage437, e)
}

func TestInputOriginal(t *testing.T) {
	t.Parallel()
	g, err := flag.InputOriginal(nil, "")
	assert.Empty(t, g)
	assert.Nil(t, err)
	g, err = flag.InputOriginal(nil, "CP437")
	assert.Empty(t, g)
	assert.Nil(t, err)
}
