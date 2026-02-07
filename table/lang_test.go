package table_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/table"
	"github.com/nalgeon/be"
	"golang.org/x/text/encoding/charmap"
)

func ExampleLanguages() {
	l := *table.Languages()
	fmt.Print(l[charmap.CodePage437])
	// Output: US English
}

func ExampleLanguage() {
	l := table.Language(charmap.CodePage437)
	fmt.Print(l)
	// Output: US English
}

func TestListLanguage(t *testing.T) {
	t.Parallel()
	b := &strings.Builder{}
	err := table.ListLanguage(b)
	be.Err(t, err, nil)
	s := strings.Contains(b.String(), "ANSI X3.4 1967")
	be.True(t, s)
}
