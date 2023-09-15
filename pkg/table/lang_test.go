package table_test

import (
	"fmt"
	"strings"

	"github.com/bengarrett/retrotxtgo/pkg/table"
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

func ExampleListLanguage() {
	b := &strings.Builder{}
	_ = table.ListLanguage(b)
	fmt.Print(b)
	// Output:
}
