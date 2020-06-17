package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/internal/pack"
	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"

	"github.com/spf13/cobra"
)

type viewFlags struct {
	codepage string
	width    int // TODO: not implemented
}

var viewFlag = viewFlags{
	codepage: "CP437",
	width:    80,
}

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view [filenames]",
	Short: "Print a legacy text file to the standard output",
	Example: `  retrotxt view file.txt -c latin1
  retrotxt view file1.txt file2.txt --codepage="iso-8859-1"`,
	Run: func(cmd *cobra.Command, args []string) {
		checkUse(cmd, args)
		for i, arg := range args {
			if ok, err := viewPackage(arg); err != nil {
				logs.CheckCont("view.pack", err)
				continue
			} else if ok {
				continue
			}
			b, err := filesystem.Read(arg)
			if ok := logs.CheckCont("view.codepage", err); !ok {
				continue
			}
			r, err := convert.Text(viewFlag.codepage, &b)
			if ok := logs.CheckCont("view.convert.text", err); !ok {
				continue
			}
			fmt.Println(string(r))
			if i < len(args) {
				fmt.Print("\n")
			}
		}
	},
}

func init() {
	// view cmd
	rootCmd.AddCommand(viewCmd)
	viewCmd.Flags().StringVarP(&viewFlag.codepage, "codepage", "c", viewFlag.codepage,
		"legacy character encoding used by the text file")
	viewCmd.Flags().IntVarP(&viewFlag.width, "width", "w", viewFlag.width, "document column character width")
	viewCmd.Flags().SortFlags = false
}

type viewPack struct {
	convert string
	name    string
}

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

func viewPackage(name string) (ok bool, err error) {
	var s = strings.ToLower(name)
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
		if r, err = convert.Dump(viewFlag.codepage, &b); err != nil {
			return false, err
		}
	case "", "t":
		if r, err = convert.Text(viewFlag.codepage, &b); err != nil {
			return false, err
		}
	}
	fmt.Println(string(r))
	return true, nil
}
