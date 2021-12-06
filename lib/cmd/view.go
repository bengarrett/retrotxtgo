// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"bytes"
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/cobra"
)

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
		b, err := viewCmdRun(cmd, args...)
		if err != nil {
			logs.Fatal(err)
		}
		fmt.Print(b)
	},
}

// viewCmdRun parses the arguments supplied with the view command.
func viewCmdRun(cmd *cobra.Command, args ...string) (*bytes.Buffer, error) {
	args, conv, samp, err := flag.InitArgs(cmd, args...)
	if err != nil {
		return nil, err
	}
	w := new(bytes.Buffer)
	for i, arg := range args {
		if i > 0 && i < len(arg) {
			const halfPage = 40
			fmt.Fprintln(w, str.HRPad(halfPage))
		}
		b, err := flag.ReadArg(arg, cmd, conv, samp)
		if err != nil {
			fmt.Fprintln(w, logs.Sprint(err))
			continue
		}
		r, err := transform(conv, samp, b...)
		if err != nil {
			fmt.Fprintln(w, logs.Sprint(err))
			continue
		}
		fmt.Fprint(w, string(r))
	}
	return w, nil
}

func init() {
	rootCmd.AddCommand(viewCmd)
	flagEncode(&flag.ViewFlag.Encode, viewCmd)
	flagControls(&flag.ViewFlag.Controls, viewCmd)
	flagRunes(&flag.ViewFlag.Swap, viewCmd)
	flagTo(&flag.ViewFlag.To, viewCmd)
	flagWidth(&flag.ViewFlag.Width, viewCmd)
	viewCmd.Flags().SortFlags = false
}
