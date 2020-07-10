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
	to       string
	width    int
}

var viewFlag = viewFlags{
	controls: nil,
	encode:   "CP437",
	to:       "",
	width:    0,
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
			if to := cmd.Flags().Lookup("to"); to.Changed {
				b, err = toDecode(viewFlag.to, &b)
				if err != nil {
					logs.Println("using UTF8 encoding as text could not convert to", viewFlag.to, err)
				}
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
			// read file
			b, err := filesystem.Read(arg)
			if err != nil {
				logs.Println("read file", arg, err)
				continue
			}
			// to flag
			if to := cmd.Flags().Lookup("to"); to.Changed {
				b, err = toDecode(viewFlag.to, &b)
				if err != nil {
					logs.Println("using UTF8 encoding as text could not convert to", viewFlag.to, err)
				}
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
	flagTo(&viewFlag.to, viewCmd)
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
	// to flag
	if to := cmd.Flags().Lookup("to"); to.Changed {
		b, err = toDecode(viewFlag.to, &b)
		if err != nil {
			logs.Println("using UTF8 encoding as text could not convert to", viewFlag.to, err)
		}
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
	fmt.Println(string(r))
	return true, nil
}

func toDecode(name string, b *[]byte) ([]byte, error) {
	encode, err := convert.Encoding(name)
	if err != nil {
		return *b, err
	}
	cp, err := encode.NewDecoder().Bytes(*b)
	if err != nil {
		return *b, err
	}
	return cp, nil
}
