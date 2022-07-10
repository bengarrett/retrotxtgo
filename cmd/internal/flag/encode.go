package flag

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// Encode handles the --encode flag.
func Encode(p *string, cc *cobra.Command) {
	cc.Flags().StringVarP(p, "encode", "e", "",
		fmt.Sprintf("character encoding used by the filename(s) (default \"CP437\")\n%s\n%s%s\n",
			color.Info.Sprint("this flag has no effect for Unicode and EBCDIC samples"),
			"see the list of encode values ",
			str.Example(meta.Bin+" list codepages")))
}
