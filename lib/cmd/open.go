package cmd

import (
	"fmt"
	"unicode/utf8"

	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/sample"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding"
)

func initArgs(cmd *cobra.Command, args ...string) ([]string, *convert.Convert, sample.Flags, error) {
	if len(args) == 0 {
		args = []string{""}
	}
	conv := convert.Convert{}
	conv.Flags = convert.Flag{
		Controls:  viewFlag.controls,
		SwapChars: viewFlag.swap,
		MaxWidth:  viewFlag.width,
	}
	if c := cmd.Flags().Lookup("controls"); !c.Changed {
		conv.Flags.Controls = []string{eof, tab}
	}
	if s := cmd.Flags().Lookup("swap-chars"); !s.Changed {
		conv.Flags.SwapChars = []int{null, verticalBar}
	}

	// html is a global flag, create.Args
	// if !changed {
	// 	only apply the sample font when the --font-family flag is unused
	// 	html.FontFamily.Value = p.Font.String()
	// }

	if filesystem.IsPipe() {
		if e := cmd.Flags().Lookup("encode"); e.Changed {
			fmt.Println("--encode flag is ignored when piped text is in use")
		}
		if args[0] != "" {
			fmt.Println("[filenames] are ignored when piped text is in use")
			args = []string{""}
		}
	} else {
		printUsage(cmd, args...)
	}
	samp, err := initEncodings(cmd, "")
	if err != nil {
		return nil, nil, samp, err
	}
	return args, &conv, samp, nil
}

// initEncodings applies the --encode and --to encoding values to embed sample data.
func initEncodings(cmd *cobra.Command, deft string) (sample.Flags, error) {
	encode := func(flag string) (encoding.Encoding, error) {
		cp := cmd.Flags().Lookup(flag)
		name := deft
		if cp.Changed {
			name = cp.Value.String()
		} else if deft == "" {
			return nil, nil
		}
		return convert.Encoder(name)
	}
	var (
		frm encoding.Encoding
		to  encoding.Encoding
	)
	if cmd == nil {
		return sample.Flags{}, nil
	}
	frm, err := encode("encode")
	if err != nil {
		return sample.Flags{}, err
	}
	return sample.Flags{From: frm, To: to}, err
}

func openArg(arg string, f sample.Flags, c *convert.Convert) ([]byte, error) {
	var (
		b   []byte
		err error
	)
	// if no argument, then assume the source is piped via stdin
	if arg == "" {
		b, err = filesystem.ReadPipe()
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	// first attempt to see if arg is a embed sample file request
	b, err = openSample(arg, c, f)
	if err != nil {
		return nil, err
	} else if b != nil {
		return b, nil
	}
	// otherwise, the arg should be a filepath
	b, err = openFile(arg)
	if err != nil {
		return nil, err
	}
	// handle flags
	//
	// return
	return b, nil
}

func openFile(arg string) ([]byte, error) {
	b, err := filesystem.Read(arg)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func openSample(arg string, c *convert.Convert, f sample.Flags) ([]byte, error) {
	if ok := sample.Valid(arg); !ok {
		return nil, nil
	}
	p, err := f.Open(arg, c)
	if err != nil {
		return nil, err
	}

	// TODO: handle encoding or use p.Encoding as fallback
	// this was used in staticTextfile() but seems to break things too
	//return create.Normalize(p.Encoding, p.Runes...), nil
	return []byte(string(p.Runes)), nil
}

func openBytes(f sample.Flags, conv *convert.Convert, b ...byte) ([]rune, error) {
	// make sure the file source isn't already encoded as UTF-8
	if utf8.Valid(b) {
		return []rune(string(b)), nil
	}
	// handle input source encoding
	if f.From != nil {
		conv.Input.Encoding = f.From
	}
	// convert text
	var (
		r   []rune
		err error
	)
	if endOfFile(conv.Flags) {
		r, err = conv.Text(b...)
	} else {
		r, err = conv.Dump(b...)
	}
	if err != nil {
		return nil, err
	}
	// return UTF-8 runes
	if f.To == nil {
		fmt.Println("returning utf-8")
		return r, nil
	}
	// re-encode text
	fmt.Println("re-encode text")
	newer, err := f.To.NewEncoder().String(string(r))
	if err != nil {
		return []rune(newer), err
	}
	return nil, nil
}
