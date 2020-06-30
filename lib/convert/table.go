package convert

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"text/tabwriter"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/gookit/color"
)

// Table prints out all the characters in the named 8-bit character set.
func Table(name string) (*bytes.Buffer, error) {
	cp, err := Encoding(name)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	w := new(tabwriter.Writer).Init(&buf, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, " "+color.OpFuzzy.Sprint(strings.Repeat("\u2015", 67)))
	fmt.Fprintln(w, color.Primary.Sprint(str.Center(fmt.Sprintf("%s", cp), 67)))
	for i := 0; i < 16; i++ {
		switch {
		case i == 0:
			fmt.Fprintf(w, "%s", color.OpFuzzy.Sprintf("     %X  ", i))
		case i == 15:
			fmt.Fprintf(w, "%s", color.OpFuzzy.Sprintf(" %X  \n", i))
		default:
			fmt.Fprintf(w, "%s", color.OpFuzzy.Sprintf(" %X  ", i))
		}
	}
	var b, row = MakeBytes(), 0
	runes, err := Chars(name, &b)
	logs.Check("convert.table.chars", err)
	for i, r := range runes {
		char := string(r)
		switch {
		case i == 0:
			fmt.Fprintf(w, " %s %s %s %s",
				color.OpFuzzy.Sprint("0"),
				color.OpFuzzy.Sprint("|"),
				char, color.OpFuzzy.Sprint("|"))
		case i == 255:
			fmt.Fprintf(w, " %s %s\n", char,
				color.OpFuzzy.Sprint("|"))
		case math.Mod(float64(i+1), 16) == 0:
			// every 16th loop
			row++
			fmt.Fprintf(w, " %s %s\n %s %s", char,
				color.OpFuzzy.Sprint("|"),
				color.OpFuzzy.Sprintf("%X", row),
				color.OpFuzzy.Sprint("|"))
		default:
			fmt.Fprintf(w, " %s %s", char,
				color.OpFuzzy.Sprint("|"))
		}
	}
	fmt.Fprint(w, "\n")
	w.Flush()
	return &buf, nil
}
