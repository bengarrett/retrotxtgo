package table_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/table"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding/charmap"
)

func ExampleLanguages() {
	l := table.Languages()
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
	assert.Nil(t, err)
	assert.Contains(t, b.String(), "ANSI X3.4 1967/77/86")
}
