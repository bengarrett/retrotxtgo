// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"
	"os"
	"unicode/utf8"

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

const viewExample = `  retrotxt view file.txt -e latin1
  retrotxt view file1.txt file2.txt --encode="iso-8859-1"
  cat file.txt | retrotxt view`

// viewCmd represents the view command.
var viewCmd = &cobra.Command{
	Use:     "view [filenames]",
	Aliases: []string{"v"},
	Short:   "Print a legacy text file to the standard output",
	Example: exampleCmd(viewExample),
	Run: func(cmd *cobra.Command, args []string) {
		viewParsePipe(cmd)
		viewParseArgs(cmd, args)
	},
}

func viewParseArgs(cmd *cobra.Command, args []string) {
	conv := convert.Convert{}
	conv.Flags = convert.Flags{
		Controls:  viewFlag.controls,
		SwapChars: viewFlag.swap,
		Width:     viewFlag.width,
	}
	if c := cmd.Flags().Lookup("controls"); !c.Changed {
		conv.Flags.Controls = []string{eof, tab}
	}
	if s := cmd.Flags().Lookup("swap-chars"); !s.Changed {
		conv.Flags.SwapChars = []int{null, verticalBar}
	}
	printUsage(cmd, args...)
	for i, arg := range args {
		skip, r := viewParseArg(cmd, &conv, i, arg)
		if skip {
			continue
		}
		fmt.Println(string(r))
		if i < len(args) {
			fmt.Println("")
		}
	}
}

func viewParseArg(cmd *cobra.Command, conv *convert.Convert, i int, arg string) (skip bool, r []rune) {
	var err error
	conv.Output = convert.Output{} // output must be reset
	f := sample.Flags{}
	// internal example file
	if ok := sample.Valid(arg); ok {
		var p sample.File
		if p, err = f.Open(arg, conv); err != nil {
			logs.MarkProblem(arg, ErrSampleView, err)
			return true, nil
		}
		// --to flag is currently ignored
		if to := cmd.Flags().Lookup("to"); to.Changed {
			if viewToFlag(p.Runes...) {
				return true, nil
			}
		}
		if i > 0 {
			str.HR(40)
		}
		fmt.Println(string(p.Runes))
		return true, nil
	}
	// read file
	b, err := filesystem.Read(arg)
	if err != nil {
		logs.MarkProblem(arg, logs.ErrOpenFile, err)
		return true, nil
	}
	if i > 0 {
		str.HR(40)
	}
	return viewParseBytes(cmd, conv, arg, b)
}

func viewParsePipe(cmd *cobra.Command) {
	if !filesystem.IsPipe() {
		return
	}
	conv := convert.Convert{}
	conv.Flags = convert.Flags{
		Controls:  viewFlag.controls,
		SwapChars: viewFlag.swap,
		Width:     viewFlag.width,
	}
	if c := cmd.Flags().Lookup("controls"); !c.Changed {
		conv.Flags.Controls = []string{eof, tab}
	}
	if s := cmd.Flags().Lookup("swap-chars"); !s.Changed {
		conv.Flags.SwapChars = []int{null, verticalBar}
	}
	// piped input from other programs
	var err error
	if cp := cmd.Flags().Lookup("encode"); cp.Changed {
		f := sample.Flags{}
		if f.From, err = convert.Encoding(cp.Value.String()); err != nil {
			logs.MarkProblemFatal("pipe", logs.ErrEncode, err)
		}
		conv.Source.E = f.From
	}
	b, err := filesystem.ReadPipe()
	if err != nil {
		logs.MarkProblemFatal("view", logs.ErrPipe, err)
	}
	_, r := viewParseBytes(cmd, &conv, "piped", b)
	fmt.Print(string(r))
	os.Exit(0)
}

func viewParseBytes(cmd *cobra.Command, conv *convert.Convert, arg string, b []byte) (skip bool, r []rune) {
	// configure source encoding
	name := "CP437"
	cp := cmd.Flags().Lookup("encode")
	if cp.Changed && cp.Value.String() != "" {
		name = cp.Value.String()
	}
	var f = sample.Flags{}
	var err error
	if f.From, err = convert.Encoding(name); err != nil {
		logs.MarkProblemFatal(arg, logs.ErrEncode, err)
	}
	conv.Source.E = f.From

	// make sure the file source isn't already encoded as UTF-8
	if utf8.Valid(b) {
		fmt.Println(string(b))
		return true, nil
	}
	// convert text
	if endOfFile(conv.Flags) {
		r, err = conv.Text(&b)
	} else {
		r, err = conv.Dump(&b)
	}
	if err != nil {
		logs.MarkProblem(arg, ErrViewUTF8, err)
		return true, nil
	}
	// to flag
	if to := cmd.Flags().Lookup("to"); to.Changed {
		if viewToFlag(r...) {
			return true, nil
		}
	}
	return false, r
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

// viewToFlag prints the results of viewEncode().
func viewToFlag(r ...rune) (success bool) {
	newer, err := viewEncode(viewFlag.to, r...)
	if err != nil {
		logs.MarkProblem(viewFlag.to, ErrEncode, err)
		return false
	}
	fmt.Println(string(newer))
	return true
}

// viewEncode encodes runes into the named encoding.
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
