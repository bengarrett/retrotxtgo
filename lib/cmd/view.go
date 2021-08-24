// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/sample"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
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
	controls: []string{eof, tab},
	encode:   "CP437",
	swap:     []int{null, verticalBar},
	to:       "",
	width:    0,
}

var viewExample = fmt.Sprintf("  %s\n%s\n%s",
	fmt.Sprintf("%s view file.txt -e latin1", meta.Bin),
	fmt.Sprintf("%s view file1.txt file2.txt --encode=\"iso-8859-1\"", meta.Bin),
	fmt.Sprintf("cat file.txt | %s view", meta.Bin))

// viewCmd represents the view command.
var viewCmd = &cobra.Command{
	Use:     "view [filenames]",
	Aliases: []string{"v"},
	Short:   "Print a text file to the terminal using standard output",
	Long:    "Print a text file to the terminal using standard output.",
	Example: exampleCmd(viewExample),
	Run: func(cmd *cobra.Command, args []string) {
		viewParsePipe(cmd)
		viewParseArgs(cmd, args...)
	},
}

// viewParseArgs parses the arguments supplied with the view command.
func viewParseArgs(cmd *cobra.Command, args ...string) {
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

// viewParseArg parses an argument supplied with the view command.
func viewParseArg(cmd *cobra.Command, conv *convert.Convert, i int, arg string) (skip bool, r []rune) {
	const halfPage = 40
	conv.Output = convert.Output{} // output must be reset
	f := sample.Flags{}
	// internal example file
	var err error
	if ok := sample.Valid(arg); ok {
		var p sample.File
		if p, err = f.Open(arg, conv); err != nil {
			logs.FatalMark(arg, logs.ErrSampView, err)
		}
		// --to flag is currently ignored
		if to := cmd.Flags().Lookup("to"); to.Changed {
			if viewToFlag(p.Runes...) {
				return true, nil
			}
		}
		if i > 0 {
			fmt.Println(str.HRPadded(halfPage))
		}
		fmt.Println(string(p.Runes))
		return true, nil
	}
	// read file
	b, err := filesystem.Read(arg)
	if err != nil {
		logs.FatalMark(arg, logs.ErrFileOpen, err)
	}
	if i > 0 {
		fmt.Println(str.HRPadded(halfPage))
	}
	return viewParseBytes(cmd, conv, arg, b...)
}

// viewParsePipe parses piped the standard input (stdin) for the view command.
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
			logs.FatalMark("pipe", logs.ErrEncode, err)
		}
		conv.Source.E = f.From
	}
	b, err := filesystem.ReadPipe()
	if err != nil {
		logs.FatalMark("view", logs.ErrPipe, err)
	}
	_, r := viewParseBytes(cmd, &conv, "piped", b...)
	fmt.Print(string(r))
	os.Exit(0)
}

// viewParseBytes parses byte data.
func viewParseBytes(cmd *cobra.Command, conv *convert.Convert, arg string, b ...byte) (skip bool, r []rune) {
	// configure source encoding
	name := "CP437"
	cp := cmd.Flags().Lookup("encode")
	if cp.Changed && cp.Value.String() != "" {
		name = cp.Value.String()
	}
	var f = sample.Flags{}
	var err error
	if f.From, err = convert.Encoding(name); err != nil {
		logs.FatalMark(arg, logs.ErrEncode, err)
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
		logs.FatalMark(arg, ErrUTF8, err)
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
	rootCmd.AddCommand(viewCmd)
	flagEncode(&viewFlag.encode, viewCmd)
	flagControls(&viewFlag.controls, viewCmd)
	flagRunes(&viewFlag.swap, viewCmd)
	flagTo(&viewFlag.to, viewCmd)
	flagWidth(&viewFlag.width, viewCmd)
	viewCmd.Flags().SortFlags = false
}

// viewToFlag prints the output of viewEncode.
func viewToFlag(r ...rune) (success bool) {
	newer, err := viewEncode(viewFlag.to, r...)
	if err != nil {
		logs.FatalMark(viewFlag.to, ErrEncode, err)
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
