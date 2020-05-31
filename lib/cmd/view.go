package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/lib/transform"

	"github.com/bengarrett/retrotxtgo/lib/filesystem"

	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/ianaindex"
)

var (
	viewCodePage string = "ibm437"
	viewFilename string
	viewFormats  = []string{"color", "text"}
	viewFormat   string
	viewWidth    int
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// todo reverse scan of file looking for SAUCE00 and COMNTT
		// todo scan for unique color codes like 24-bit color
		// todo scan for new lines or character counts and hard-code the width

		if viewFilename == "" {
			// todo: remove from final
			viewFilename = "textfiles/cp-437-all-characters.txt"
		}
		// todo handle unchanged viewCodePage, where UTF8 encoding will be checked otherwise use
		encoding, err := ianaindex.IANA.Encoding(viewCodePage)
		// TODO: check errors
		//CheckCodePage(ErrorFmt{"", viewCodePage, err})
		var d transform.Set

		data, err := filesystem.Read(viewFilename)
		//logs.ChkErr(logs.Err{"file open", viewFilename, err})

		err = d.Transform(data, encoding)
		logs.Check("codepage", err) // TODO: replace
		//logs.ChkErr(logs.Err{"Transform", "encoding", err})
		d.SwapAll(true)
		fmt.Printf("\n%s\n", d.Data)
		// todo: make an --example that auto generates
		// a table bytes 0 - 255 | lf every 16 characters
	},
}

var viewCodePagesCmd = &cobra.Command{
	Use:   "codepages",
	Short: "list available legacy codepages that RetroTxt can convert into UTF-8",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(transform.List())
	},
}

var viewTableCmd = &cobra.Command{
	Use:   "table",
	Short: "display a table showing the codepage and all its characters",
	Run: func(cmd *cobra.Command, args []string) {
		encoding, err := ianaindex.IANA.Encoding(viewCodePage)
		//CheckCodePage(ErrorFmt{"", viewCodePage, err})
		cp, err := ianaindex.IANA.Name(encoding)
		//CheckCodePage(ErrorFmt{"", viewCodePage, err})
		table, err := transform.Table(cp)
		logs.ChkErr(logs.Err{Issue: "table", Arg: cp, Msg: err})
		println(table)
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
	viewCmd.Flags().StringVarP(&viewFilename, "name", "n", "",
		str.Required("text file to display")+"\n")
	viewCmd.Flags().StringVarP(&viewCodePage, "codepage", "c", "cp437", "legacy character encoding used by the text file")
	viewCmd.Flags().StringVarP(&viewFormat, "format", "f", "color",
		str.Options("output format", viewFormats, true))
	viewCmd.Flags().IntVarP(&viewWidth, "width", "w", 80, "document column character width")
	// override ascii 0-F + 1-F || Control characters || IBM, ASCII, IBM+
	// example flag showing CP437 table
	_ = viewCmd.MarkFlagFilename("name")
	_ = viewCmd.MarkFlagRequired("name")
	viewCmd.Flags().SortFlags = false
	viewCmd.AddCommand(viewCodePagesCmd)
	viewCmd.AddCommand(viewTableCmd)
	viewTableCmd.Flags().StringVarP(&viewCodePage, "codepage", "c", "cp437", "legacy character encoding table to display")
	_ = viewTableCmd.MarkFlagRequired("name")
}
