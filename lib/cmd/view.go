package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/internal/pack"
	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
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

var viewTypes = []string{"chars", "dump", "text"}

type viewPack struct {
	convert string
	name    string
}

// TODO: replace --name with args.
// have a fileexist check for rare possible conflicting names.
// retrotxt view filename.txt
// retrotxt view rt.ansi <- ...

var viewPacks = map[string]viewPack{
	"437.cr":        {"d", "text/cp437-cr.txt"},
	"437.crlf":      {"d", "text/cp437-crlf.txt"},
	"437.lf":        {"d", "text/cp437-lf.txt"},
	"865":           {"", "text/cp865.txt"},
	"1252":          {"", "text/cp1252.txt"},
	"ascii":         {"", "text/retrotxt.asc"},
	"ansi":          {"", "text/retrotxt.ans"},
	"ansi.aix":      {"", "text/ansi-aixterm.ans"},
	"ansi.blank":    {"", "text/ansi-blank"},
	"ansi.cp":       {"", "text/ansi-cp.ans"},
	"ansi.cpf":      {"", "text/ansi-cpf.ans"},
	"ansi.hvp":      {"", "text/ansi-hvp.ans"},
	"ansi.proof":    {"", "text/ansi-proof.ans"},
	"ansi.rgb":      {"", "text/ansi-rgb.ans"},
	"ansi.setmodes": {"", "text/ansi-setmodes.ans"},
	"iso-1":         {"", "text/iso-8859-1.txt"},
	"iso-15":        {"", "text/iso-8859-15.txt"},
	"sauce":         {"", "text/sauce.txt"},
	"shiftjis":      {"", "text/shiftjis.txt"},
	"us-ascii":      {"", "text/us-ascii.txt"},
	"utf8":          {"", "text/utf-8.txt"},
	"utf8.bom":      {"", "text/utf-8-bom.txt"},
	"utf16.be":      {"", "text/utf-16-be.txt"},
	"utf16.le":      {"", "text/utf-16-le.txt"},
}

func viewPackage() (ok bool, err error) {
	var s = strings.ToLower(viewArgs.name)
	if _, err := os.Stat(s); !os.IsNotExist(err) {
		return false, nil
	}
	pkg, exist := viewPacks[s]
	println(fmt.Sprintf("%+v", pkg), exist)
	if !exist {
		return false, nil
	}
	b := pack.Get(pkg.name)
	if b == nil {
		return false, errors.New("pkg.name is unknown: " + pkg.name)
	}
	var r []rune
	switch pkg.convert {
	case "d":
		if r, err = convert.Dump(viewArgs.cp, &b); err != nil {
			return false, err
		}
	case "", "t":
		if r, err = convert.Text(viewArgs.cp, &b); err != nil {
			return false, err
		}
	}
	fmt.Println(string(r))
	return true, nil
}

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Print a legacy text file to the standard output",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("args: %+v\n", args)
		ok, err := viewPackage()
		logs.Check("view.pack", err)
		if ok {
			os.Exit(0)
		}
		b, err := filesystem.Read(viewArgs.name)
		logs.Check("view.codepage", err)
		r, err := convert.Text(viewArgs.cp, &b)
		logs.Check("view.convert.text", err)
		fmt.Println(string(r))
	},
}

var viewCodePagesCmd = &cobra.Command{
	Use:   "codepages",
	Short: "list available legacy codepages that RetroTxt can convert into UTF-8",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(convert.List())
	},
}

var viewTableCmd = &cobra.Command{
	Use:   "table",
	Short: "display a table showing the codepage and all its characters",
	Run: func(cmd *cobra.Command, args []string) {
		table, err := convert.Table(viewArgs.cp)
		logs.ChkErr(logs.Err{Issue: "table", Arg: viewArgs.cp, Msg: err})
		fmt.Println(table.String())
	},
}

var viewTablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "display tables showing known codepages and characters",
	Run: func(cmd *cobra.Command, args []string) {
		for _, e := range convert.Encodings() {
			name, err := ianaindex.MIME.Name(e)
			if err != nil {
				logs.Log(err)
			} else {
				// keep 0F,1F controls. blank other ?
				// tables -> Macintosh to list alt. names Mac OS Roman
				// Windows 874 is not showing different chars from ISO-11
				// https://en.wikipedia.org/wiki/ISO/IEC_8859-11#Vendor_extensions
				// japanese needs fixing
				table, err := convert.Table(name)
				logs.ChkErr(logs.Err{Issue: "tables", Arg: name, Msg: err})
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
