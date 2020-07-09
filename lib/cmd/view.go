package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"retrotxt.com/retrotxt/internal/pack"
	"retrotxt.com/retrotxt/lib/convert"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/logs"

	"github.com/spf13/cobra"
)

type viewFlags struct {
	controls []string
	encode   string
	width    int // TODO: not implemented
}

var viewFlag = viewFlags{
	controls: nil,
	encode:   "CP437",
	width:    80,
}

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:     "view [filenames]",
	Aliases: []string{"v"},
	Short:   "Print a legacy text file to the standard output",
	Example: `  retrotxt view file.txt -e latin1
  retrotxt view file1.txt file2.txt --encode="iso-8859-1"
  cat file.txt | retrotxt view`,
	Run: func(cmd *cobra.Command, args []string) {
		var conv = convert.Args{
			Controls: viewFlag.controls,
			Encoding: viewFlag.encode,
			Width:    viewFlag.width,
		}
		// piped input from other programs
		if filesystem.IsPipe() {
			b, err := filesystem.ReadPipe()
			if err != nil {
				logs.Fatal("view", "stdin read", err)
			}
			r, err := conv.Text(&b)
			if err != nil {
				logs.Fatal("view", "stdin convert", err)
			}
			fmt.Println(string(r))
			os.Exit(0)
		}
		// user arguments
		checkUse(cmd, args)
		for i, arg := range args {
			if ok, err := viewPackage(cmd, arg); err != nil {
				logs.Println("pack", arg, err)
				continue
			} else if ok {
				continue
			}
			b, err := filesystem.Read(arg)
			if err != nil {
				logs.Println("read file", arg, err)
				continue
			}
			r, err := conv.Text(&b)
			if err != nil {
				logs.Println("convert text", arg, err)
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
	flagEncode(&viewFlag.encode, viewCmd)
	flagControls(&viewFlag.controls, viewCmd)
	flagWidth(&viewFlag.width, viewCmd)
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
	var conv = convert.Args{
		Controls: viewFlag.controls,
		Encoding: viewFlag.encode,
		Width:    viewFlag.width,
	}
	// encode defaults
	if cp := cmd.Flags().Lookup("encode"); !cp.Changed {
		conv.Encoding = pkg.encoding
	}
	// convert and print
	var r []rune
	switch pkg.convert {
	case "d":
		if r, err = conv.Dump(&b); err != nil {
			return false, err
		}
	case "", "t":
		if r, err = conv.Text(&b); err != nil {
			return false, err
		}
	}
	// TODO: optional char conversion
	// ab, err := charmap.CodePage437.NewDecoder().Bytes(b)
	// if err != nil {
	// 	logs.Println("conversion from", "CodePage437", err)
	// }
	// nr, err := charmap.CodePage437.NewEncoder().Bytes(ab)
	// if err != nil {
	// 	logs.Println("conversion to", "CodePage437", err)
	// }
	// fmt.Println(string(nr))
	// os.Exit(0)
	fmt.Println(string(r))
	return true, nil
}
