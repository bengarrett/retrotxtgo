package cmd

import (
	"fmt"
	"os"

	"retrotxt.com/retrotxt/lib/convert"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/pack"

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
	controls: []string{tab},
	encode:   "CP437",
	swap:     []int{null, verticalBar},
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
		conv := convert.Args{
			Controls: viewFlag.controls,
			Encoding: viewFlag.encode,
			Swap:     viewFlag.swap,
			Width:    viewFlag.width,
		}
		f := pack.Flags{}
		// handle defaults that are left empty for usage formatting
		if c := cmd.Flags().Lookup("controls"); !c.Changed {
			conv.Controls = []string{tab}
		}
		if s := cmd.Flags().Lookup("swap-chars"); !s.Changed {
			conv.Swap = []int{null, verticalBar}
		}
		// piped input from other programs
		if filesystem.IsPipe() {
			viewPipe(cmd, conv)
		}
		// user arguments
		checkUse(cmd, args...)
		var err error
		for i, arg := range args {
			// internal, packed example file
			if ok := pack.Valid(arg); ok {
				if cp := cmd.Flags().Lookup("encode"); cp.Changed {
					fmt.Println("cp changed", cp.Value.String())
					if f.Encode, err = convert.Encoding(cp.Value.String()); err != nil {
						logs.Fatal("encoding not known or supported", arg, err)
					}
				}
				if to := cmd.Flags().Lookup("to"); to.Changed {
					if f.To, err = convert.Encoding(to.Value.String()); err != nil {
						logs.Fatal("to not known or supported", arg, err)
					}
				}
				var p pack.Pack
				if p, err = f.Open(conv, arg); err != nil {
					logs.Println("pack", arg, err)
					continue
				}
				fmt.Println(string(p.Runes))
				continue
			}
			// read file
			b, err := filesystem.Read(arg)
			if err != nil {
				logs.Println("read file", arg, err)
				continue
			}
			// to flag
			if to := cmd.Flags().Lookup("to"); to.Changed {
				viewToFlag(b...) // todo: move to root.go and return a value.
				continue
			}
			// convert text
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
	flagRunes(&viewFlag.swap, viewCmd)
	flagTo(&viewFlag.to, viewCmd)
	flagWidth(&viewFlag.width, viewCmd)
	viewCmd.Flags().SortFlags = false
}

func viewPipe(cmd *cobra.Command, conv convert.Args) {
	b, err := filesystem.ReadPipe()
	if err != nil {
		logs.Fatal("view", "stdin read", err)
	}
	// to flag
	if to := cmd.Flags().Lookup("to"); to.Changed {
		viewToFlag(b...)
		os.Exit(0)
	}
	r, err := conv.Text(&b)
	if err != nil {
		logs.Fatal("view", "stdin convert", err)
	}
	fmt.Println(string(r))
	os.Exit(0)
}

func viewToFlag(b ...byte) {
	b, err := viewEncode(viewFlag.to, b...)
	if err != nil {
		logs.Println("using the original encoding and not", viewFlag.to, err)
	}
	fmt.Println(string(b))
}

func viewEncode(name string, b ...byte) ([]byte, error) {
	encode, err := convert.Encoding(name)
	if err != nil {
		return b, fmt.Errorf("encoding not known or supported %s: %w", encode, err)
	}
	newer, err := encode.NewEncoder().Bytes(b)
	if err != nil {
		if len(newer) == 0 {
			return b, fmt.Errorf("encoder could not convert bytes to %s: %w", encode, err)
		}
		return newer, nil
	}
	return newer, nil
}
