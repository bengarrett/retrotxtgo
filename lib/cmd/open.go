package cmd

import (
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/sample"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding"
)

// initArgs initializes the command arguments and flags.
func initArgs(cmd *cobra.Command, args ...string) ([]string, *convert.Convert, sample.Flags, error) {
	conv := convert.Convert{}
	conv.Flags = convert.Flag{
		Controls:  viewFlag.controls,
		SwapChars: viewFlag.swap,
		MaxWidth:  viewFlag.width,
	}
	l := len(args)

	if c := cmd.Flags().Lookup("controls"); !c.Changed {
		conv.Flags.Controls = []string{eof, tab}
	}
	if s := cmd.Flags().Lookup("swap-chars"); !s.Changed {
		conv.Flags.SwapChars = []string{null, verticalBar}
	}
	if filesystem.IsPipe() {
		if e := cmd.Flags().Lookup("encode"); e.Changed {
			fmt.Println("--encode flag is ignored when piped text is in use")
		}
		if l == 0 {
			fmt.Println("[filenames] are ignored when piped text is in use")
		}
	} else if err := printUsage(cmd, args...); err != nil {
		logs.Fatal(err)
	}
	if l == 0 {
		args = []string{""}
	}
	samp, err := initEncodings(cmd, "")
	if err != nil {
		return nil, nil, samp, err
	}
	return args, &conv, samp, nil
}

// initEncodings applies the --encode and --to encoding values to embed sample data.
func initEncodings(cmd *cobra.Command, def string) (sample.Flags, error) {
	encode := func(flag string) (encoding.Encoding, error) {
		cp := cmd.Flags().Lookup(flag)
		name := def
		if cp.Changed {
			name = cp.Value.String()
		} else if def == "" {
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

// readArg returns the content of argument supplied filepath, embed sample file or piped data.
func readArg(arg string, cmd *cobra.Command, c *convert.Convert, f sample.Flags) ([]byte, error) {
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
	// attempt to see if arg is a embed sample file request
	if b, err = openSample(arg, cmd, c, f); err != nil {
		return nil, err
	} else if b != nil {
		return b, nil
	}
	// the arg should be a filepath
	b, err = openFile(arg)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func openFile(arg string) ([]byte, error) {
	b, err := filesystem.Read(arg)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func openSample(arg string, cmd *cobra.Command, c *convert.Convert, f sample.Flags) ([]byte, error) {
	if ok := sample.Valid(arg); !ok {
		return nil, nil
	}
	p, err := f.Open(arg, c)
	if err != nil {
		return nil, err
	}
	// handle flags
	if ff := cmd.Flags().Lookup("font-family"); ff != nil && !ff.Changed {
		// only apply the sample font when the --font-family flag is unused
		// html is a global flag, create.Args
		html.FontFamily.Value = p.Font.String()
	}
	// TODO: handle encoding or use p.Encoding as fallback
	//return create.Normalize(p.Encoding, p.Runes...), nil

	return []byte(string(p.Runes)), nil
}

func transform(conv *convert.Convert, f sample.Flags, b ...byte) ([]rune, error) {
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
		return r, nil
	}
	// re-encode text
	newer, err := f.To.NewEncoder().String(string(r))
	if err != nil {
		return []rune(newer), err
	}
	return nil, nil
}
