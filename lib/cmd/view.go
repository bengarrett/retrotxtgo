package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/lib/transform"
	"github.com/bengarrett/retrotxtgo/samples"

	"github.com/bengarrett/retrotxtgo/lib/filesystem"

	"github.com/spf13/cobra"
)

type viewFlags struct {
	cp      string
	name    string
	formats []string
	format  string
	width   int
}

var viewArgs = viewFlags{
	cp:      "cp437",
	formats: []string{"color", "text"},
}

/*
TODO:
- reverse scan of file looking for EOF, SAUCE00 & COMNTT
- scan for unique color codes like 24-bit colors
- newline scanner to determine the maxWidth
*/

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err error
			t   transform.Set
		)
		switch viewArgs.name {
		case "ansi":
			t.B, err = samples.Base64Decode(samples.LogoANSI)
			logs.ChkErr(logs.Err{Issue: "logoansi is invalid", Arg: htmlArgs.Src, Msg: err})
		case "ascii":
			t.B, err = samples.Base64Decode(samples.LogoASCII)
			logs.ChkErr(logs.Err{Issue: "logoascii is invalid", Arg: htmlArgs.Src, Msg: err})
		case "":
			viewArgs.name = "textfiles/cp-437-all-characters.txt"
			fallthrough
		default:
			t.B, err = filesystem.Read(viewArgs.name)
			logs.Check("codepage", err)
		}
		_, err = t.Transform(viewArgs.cp)
		logs.Check("codepage", err)
		t.Newline = true
		t.Swap()
		fmt.Println(string(t.B))
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
		table, err := transform.Table(viewArgs.cp)
		logs.ChkErr(logs.Err{Issue: "table", Arg: viewArgs.cp, Msg: err})
		fmt.Println(table.String())
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
	viewCmd.Flags().StringVarP(&viewArgs.name, "name", "n", "",
		str.Required("text file to display")+"\n")
	viewCmd.Flags().StringVarP(&viewArgs.cp, "codepage", "c", "cp437", "legacy character encoding used by the text file")
	viewCmd.Flags().StringVarP(&viewArgs.format, "format", "f", "color",
		str.Options("output format", viewArgs.formats, true))
	viewCmd.Flags().IntVarP(&viewArgs.width, "width", "w", 80, "document column character width")
	// override ascii 0-F + 1-F || Control characters || IBM, ASCII, IBM+
	// example flag showing CP437 table
	_ = viewCmd.MarkFlagFilename("name")
	_ = viewCmd.MarkFlagRequired("name")
	viewCmd.Flags().SortFlags = false
	viewCmd.AddCommand(viewCodePagesCmd)
	viewCmd.AddCommand(viewTableCmd)
	viewTableCmd.Flags().StringVarP(&viewArgs.cp, "codepage", "c", "cp437", "legacy character encoding table to display")
	_ = viewTableCmd.MarkFlagRequired("name")
}
