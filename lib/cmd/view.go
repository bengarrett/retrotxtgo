package cmd

import (
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
	swap     []int
	to       string
	width    int
}

var viewFlag = viewFlags{
	controls: []string{"tab"},
	encode:   "CP437",
	swap:     []int{0, 124},
	to:       "",
	width:    0,
}

// viewCmd represents the view command.
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
			Swap:     viewFlag.swap,
			Width:    viewFlag.width,
		}
		// handle defaults that are left empty for usage formatting
		if c := cmd.Flags().Lookup("controls"); !c.Changed {
			conv.Controls = []string{"tab"}
		}
		if s := cmd.Flags().Lookup("swap-chars"); !s.Changed {
			conv.Swap = []int{0, 124}
		}
		// piped input from other programs
		if filesystem.IsPipe() {
			viewPipe(cmd, conv)
		}
		// user arguments
		checkUse(cmd, args...)
		for i, arg := range args {
			// internal, packed example file
			if ok, err := viewPackage(cmd, conv, arg); err != nil {
				logs.Println("pack", arg, err)
				continue
			} else if ok {
				continue
			}
			// read file
			b, err := filesystem.Read(arg)
			if err != nil {
				logs.Println("read file", arg, err)
				continue
			}
			// convert text
			r, err := conv.Text(&b)
			if err != nil {
				logs.Println("convert text", arg, err)
				continue
			}
			// to flag
			if to := cmd.Flags().Lookup("to"); to.Changed {
				r, err = toDecode(viewFlag.to, r...)
				if err != nil {
					logs.Println("using utf8 encoding and not", viewFlag.to, err)
				}
				fmt.Println(string(r))
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
	flagRunes(&viewFlag.swap, viewCmd)
	flagTo(&viewFlag.to, viewCmd)
	flagWidth(&viewFlag.width, viewCmd)
	viewCmd.Flags().SortFlags = false
}

func viewPackage(cmd *cobra.Command, conv convert.Args, name string) (ok bool, err error) {
	var s = strings.ToLower(name)
	if _, err = os.Stat(s); !os.IsNotExist(err) {
		return false, nil
	}
	pkg, exist := internalPacks[s]
	if !exist {
		return false, nil
	}
	b := pack.Get(pkg.name)
	if b == nil {
		return false, fmt.Errorf("view package %q: %w", pkg.name, ErrPackGet)
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
	// to flag
	if to := cmd.Flags().Lookup("to"); to.Changed {
		r, err = toDecode(viewFlag.to, r...)
		if err != nil {
			logs.Println("using utf8 encoding and not", viewFlag.to, err)
		}
		fmt.Println(string(r))
		return true, nil
	}
	fmt.Println(string(r))
	return true, nil
}

func viewPipe(cmd *cobra.Command, conv convert.Args) {
	b, err := filesystem.ReadPipe()
	if err != nil {
		logs.Fatal("view", "stdin read", err)
	}
	r, err := conv.Text(&b)
	if err != nil {
		logs.Fatal("view", "stdin convert", err)
	}
	// to flag
	if to := cmd.Flags().Lookup("to"); to.Changed {
		r, err = toDecode(viewFlag.to, r...)
		if err != nil {
			logs.Println("using utf8 encoding and not", viewFlag.to, err)
		}
		fmt.Println(string(r))
		os.Exit(0)
	}
	fmt.Println(string(r))
	os.Exit(0)
}

func toDecode(name string, r ...rune) ([]rune, error) {
	encode, err := convert.Encoding(name)
	if err != nil {
		return r, fmt.Errorf("encoding not known or supported %s: %w", encode, err)
	}
	cp, err := encode.NewEncoder().String(string(r))
	if err != nil {
		if cp == "" {
			return r, fmt.Errorf("encoder could not convert runes to %s: %w", encode, err)
		}
		return []rune(cp), nil
	}
	return []rune(cp), nil
}
