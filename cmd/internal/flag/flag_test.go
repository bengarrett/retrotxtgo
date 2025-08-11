package flag_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/convert"
	"github.com/bengarrett/retrotxtgo/sample"
	"github.com/nalgeon/be"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

func ExampleEndOfFile() {
	var f convert.Flag
	f.Controls = []string{"eof"}
	fmt.Fprint(os.Stdout, flag.EndOfFile(f))
	// Output: true
}

var cp437 encoding.Encoding = charmap.CodePage437

func TestDefault(t *testing.T) {
	t.Parallel()
	e := flag.Default()
	be.True(t, cp437 == e)
}

func TestInputOriginal(t *testing.T) {
	t.Parallel()
	g, err := flag.InputOverride(nil, "")
	be.Equal(t, g, sample.Flags{})
	be.Err(t, err, nil)
	g, err = flag.InputOverride(nil, "CP437")
	be.Equal(t, g, sample.Flags{})
	be.Err(t, err, nil)
}
