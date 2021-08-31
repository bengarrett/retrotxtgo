// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/cobra"
)

type viewFlags struct {
	controls []string
	encode   string
	swap     []int
	to       string
	width    int
}

var viewFlag = viewFlags{
	controls: []string{eof, tab},
	encode:   "CP437",
	swap:     []int{null, verticalBar},
	to:       "",
	width:    0,
}

var viewExample = fmt.Sprintf("  %s\n%s\n%s",
	fmt.Sprintf("%s view file.txt -e latin1", meta.Bin),
	fmt.Sprintf("%s view file1.txt file2.txt --encode=\"iso-8859-1\"", meta.Bin),
	fmt.Sprintf("cat file.txt | %s view", meta.Bin))

// viewCmd represents the view command.
var viewCmd = &cobra.Command{
	Use:     fmt.Sprintf("view %s", filenames),
	Aliases: []string{"v"},
	Short:   "Print a text file to the terminal using standard output",
	Long:    "Print a text file to the terminal using standard output.",
	Example: exampleCmd(viewExample),
	Run: func(cmd *cobra.Command, args []string) {
		viewParseArgs(cmd, args...)
	},
}

// viewParseArgs parses the arguments supplied with the view command.
func viewParseArgs(cmd *cobra.Command, args ...string) {
	args, conv, samp, err := initArgs(cmd, args...)
	if err != nil {
		logs.Fatal(err)
	}
	for i, arg := range args {
		if i > 0 && i < len(arg) {
			const halfPage = 40 // output must be reset
			fmt.Println(str.HRPad(halfPage))
		}
		fmt.Printf("%d. temp string: %s\n\n", i, arg)
		b, err := openArg(arg, samp, conv)
		if err != nil {
			fmt.Println(logs.Printf(err))
			continue
		}
		r, err := openBytes(samp, conv, b...)
		if err != nil {
			fmt.Println(logs.Printf(err))
			continue
		}
		fmt.Println("arg end:\n", string(r))
	}
}

func init() {
	rootCmd.AddCommand(viewCmd)
	flagEncode(&viewFlag.encode, viewCmd)
	flagControls(&viewFlag.controls, viewCmd)
	flagRunes(&viewFlag.swap, viewCmd)
	flagTo(&viewFlag.to, viewCmd)
	flagWidth(&viewFlag.width, viewCmd)
	viewCmd.Flags().SortFlags = false
}
