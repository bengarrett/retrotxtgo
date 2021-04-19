// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"retrotxt.com/retrotxt/lib/convert"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/sample"
	"retrotxt.com/retrotxt/lib/str"
)

type viewFlags struct {
	controls []string
	encode   string
	swap     []int
	to       string
	width    int
}

var viewFlag = viewFlags{
	controls: []string{eof, tab},
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
		var err error
		conv := convert.Convert{}
		conv.Flags = convert.Flags{
			Controls:  viewFlag.controls,
			SwapChars: viewFlag.swap,
			Width:     viewFlag.width,
		}
		f := sample.Flags{}
		// handle defaults that are left empty for usage formatting
		if c := cmd.Flags().Lookup("controls"); !c.Changed {
			conv.Flags.Controls = []string{eof, tab}
		}
		if s := cmd.Flags().Lookup("swap-chars"); !s.Changed {
			conv.Flags.SwapChars = []int{null, verticalBar}
		}
		// piped input from other programs
		if filesystem.IsPipe() {
			if cp := cmd.Flags().Lookup("encode"); cp.Changed {
				if f.From, err = convert.Encoding(cp.Value.String()); err != nil {
					logs.Fatal("encoding not known or supported", "pipe", err)
				}
				conv.Source.E = f.From
			}
			viewPipe(cmd, &conv)
		}
		// user arguments
		checkUse(cmd, args...)
		for i, arg := range args {
			conv.Output = convert.Output{} // output must be reset
			if i > 0 {
				fmt.Printf(" \n%s\n\n", str.Cb(strings.Repeat("\u2500", 40))) // horizontal bar
			}
			if cp := cmd.Flags().Lookup("encode"); cp.Changed {
				if f.From, err = convert.Encoding(cp.Value.String()); err != nil {
					logs.Fatal("encoding not known or supported", arg, err)
				}
				conv.Source.E = f.From
			}
			// internal example file
			if ok := sample.Valid(arg); ok {
				var p sample.File
				if p, err = f.Open(&conv, arg); err != nil {
					logs.Println("sample", arg, err)
					continue
				}
				// --to flag is currently ignored
				if to := cmd.Flags().Lookup("to"); to.Changed {
					if viewToFlag(p.Runes...) {
						continue
					}
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
			// convert text
			var r []rune
			if endOfFile(conv.Flags) {
				r, err = conv.Text(&b)
			} else {
				r, err = conv.Dump(&b)
			}
			if err != nil {
				logs.Println("convert text", arg, err)
				continue
			}
			// to flag
			if to := cmd.Flags().Lookup("to"); to.Changed {
				if viewToFlag(r...) {
					continue
				}
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

func viewPipe(cmd *cobra.Command, conv *convert.Convert) {
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
		if viewToFlag(r...) {
			os.Exit(0)
		}
	}
	fmt.Println(string(r))
	os.Exit(0)
}

func viewToFlag(r ...rune) (success bool) {
	newer, err := viewEncode(viewFlag.to, r...)
	if err != nil {
		logs.Println("using the original encoding and not", viewFlag.to, err)
		return false
	}
	fmt.Println(string(newer))
	return true
}

func viewEncode(name string, r ...rune) (b []byte, err error) {
	encode, err := convert.Encoding(name)
	if err != nil {
		return b, fmt.Errorf("encoding not known or supported %s: %w", encode, err)
	}
	b = []byte(string(r))
	newer, err := encode.NewEncoder().Bytes(b)
	if err != nil {
		if len(newer) == 0 {
			return b, fmt.Errorf("encoder could not convert bytes to %s: %w", encode, err)
		}
		return newer, nil
	}
	return newer, nil
}
