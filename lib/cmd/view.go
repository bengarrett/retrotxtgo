package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/lib/transform"
	"github.com/bengarrett/retrotxtgo/samples"
	"golang.org/x/text/encoding/ianaindex"

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
	Short: "Print a legacy text file to the standard output",
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

var viewTablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "display tables showing known codepages and characters",
	Run: func(cmd *cobra.Command, args []string) {
		for _, e := range transform.Encodings() {
			name, err := ianaindex.MIME.Name(e)
			if err != nil {
				logs.Log(err)
			} else {
				table, err := transform.Table(name)
				logs.ChkErr(logs.Err{Issue: "table", Arg: name, Msg: err})
				fmt.Println(table.String())
			}
		}
	},
}

func init() {
	// view cmd
	rootCmd.AddCommand(viewCmd)
	viewCmd.Flags().StringVarP(&viewArgs.name, "name", "n", "",
		str.Required("text file to display")+"\n")
	viewCmd.Flags().StringVarP(&viewArgs.cp, "codepage", "c", "cp437",
		"legacy character encoding used by the text file")
	viewCmd.Flags().StringVarP(&viewArgs.format, "format", "f", "color",
		str.Options("output format", viewArgs.formats, true))
	viewCmd.Flags().IntVarP(&viewArgs.width, "width", "w", 80, "document column character width")
	err := viewCmd.MarkFlagFilename("name")
	logs.Check("view.filename", err)
	err = viewCmd.MarkFlagRequired("name")
	logs.Check("view.required", err)
	viewCmd.Flags().SortFlags = false
	// codepages cmd
	viewCmd.AddCommand(viewCodePagesCmd)
	// table cmd
	viewCmd.AddCommand(viewTableCmd)
	viewTableCmd.Flags().StringVarP(&viewArgs.cp, "codepage", "c", "cp437",
		"legacy character encoding table to display")
	// tables cmd
	viewCmd.AddCommand(viewTablesCmd)
}
