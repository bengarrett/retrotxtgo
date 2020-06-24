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
	Use:     "view [filenames]",
	Aliases: []string{"v"},
	Short:   "Print a legacy text file to the standard output",
	Example: `  retrotxt view file.txt -c latin1
  retrotxt view file1.txt file2.txt --codepage="iso-8859-1"`,
	Run: func(cmd *cobra.Command, args []string) {
		checkUse(cmd, args)
		for i, arg := range args {
			if ok, err := viewPackage(cmd, arg); err != nil {
				logs.CheckCont("pack", err)
				continue
			} else if ok {
				continue
			}
			b, err := filesystem.Read(arg)
			if ok := logs.CheckCont("read file", err); !ok {
				continue
			}
			r, err := convert.Text(viewFlag.codepage, &b)
			if ok := logs.CheckCont("convert", err); !ok {
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

func viewPackage(cmd *cobra.Command, name string) (ok bool, err error) {
	var s = strings.ToLower(name)
	if _, err := os.Stat(s); !os.IsNotExist(err) {
		return false, nil
	}
	pkg, exist := internalPacks[s]
	if !exist {
		return false, nil
	}
	b := pack.Get(pkg.name)
	if b == nil {
		return false, errors.New("pkg.name is unknown: " + pkg.name)
	}
	// codepage defaults
	encoding := viewFlag.codepage
	if cp := cmd.Flags().Lookup("codepage"); !cp.Changed {
		encoding = pkg.encoding
	}
	// convert and print
	var r []rune
	switch pkg.convert {
	case "d":
		if r, err = convert.Dump(encoding, &b); err != nil {
			return false, err
		}
	case "", "t":
		if r, err = convert.Text(encoding, &b); err != nil {
			return false, err
		}
	}
	fmt.Println(string(r))
	return true, nil
}
