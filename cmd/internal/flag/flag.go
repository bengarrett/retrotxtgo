package flag

import (
	"errors"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/pkg/convert"
	"github.com/bengarrett/retrotxtgo/pkg/create"
	"github.com/bengarrett/retrotxtgo/pkg/filesystem"
	"github.com/bengarrett/retrotxtgo/pkg/logs"
	"github.com/bengarrett/retrotxtgo/pkg/sample"
	"github.com/bengarrett/sauce"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

const EncodingDefault = "CP437"

var ErrFilenames = errors.New("ignoring [filenames]")

// Args initializes the command arguments and flags.
func Args(cmd *cobra.Command, args ...string) ([]string, *convert.Convert, sample.Flags, error) {
	conv := convert.Convert{}
	conv.Flags = convert.Flag{
		Controls:  View().Controls,
		SwapChars: View().Swap,
		MaxWidth:  View().Width,
	}
	l := len(args)

	if c := cmd.Flags().Lookup("controls"); c != nil && !c.Changed {
		conv.Flags.Controls = []string{"eof", "tab"}
	}
	if s := cmd.Flags().Lookup("swap-chars"); s != nil && !s.Changed {
		conv.Flags.SwapChars = []string{"null", "bar"}
	}
	if filesystem.IsPipe() {
		var err error
		if l > 0 {
			err = fmt.Errorf("%w;%w for piped text", err, ErrFilenames)
			args = []string{""}
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, logs.Sprint(err))
		}
	} else if err := Help(cmd, args...); err != nil {
		logs.Fatal(err)
	}
	if l == 0 {
		args = []string{""}
	}
	samp, err := EncodeAndHide(cmd, "")
	if err != nil {
		return nil, nil, samp, err
	}
	if conv.Input.Encoding == nil {
		conv.Input.Encoding = Default()
	}
	return args, &conv, samp, nil
}

// Default returns the default encoding when the --encoding flag is unused.
func Default() encoding.Encoding { //nolint:ireturn
	if filesystem.IsPipe() {
		return unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
	}
	return charmap.CodePage437
}

// EncodeAndHide applies the public --encode and the hidden --to encoding values to embed sample data.
func EncodeAndHide(cmd *cobra.Command, dfault string) (sample.Flags, error) {
	parse := func(name string) (encoding.Encoding, error) {
		cp := cmd.Flags().Lookup(name)
		lookup := dfault
		if cp != nil && cp.Changed {
			lookup = cp.Value.String()
		} else if dfault == "" {
			return nil, nil
		}
		if lookup == "" {
			return nil, nil
		}
		return convert.Encoder(lookup)
	}
	var (
		frm encoding.Encoding
		to  encoding.Encoding
	)
	if cmd == nil {
		return sample.Flags{}, nil
	}
	// handle encode flag or apply the default
	frm, err := parse("encode")
	if err != nil {
		return sample.Flags{}, err
	}
	// handle the hidden reencode (--to) flag
	to, err = parse("to")
	if err != nil {
		return sample.Flags{}, err
	}
	return sample.Flags{From: frm, To: to}, err
}

// EndOfFile returns true if the end-of-file control flag was requested.
func EndOfFile(flags convert.Flag) bool {
	for _, c := range flags.Controls {
		if c == "eof" {
			return true
		}
	}
	return false
}

// Help will print the help and exit when no arguments are supplied.
func Help(cmd *cobra.Command, args ...string) error {
	if len(args) != 0 {
		return nil
	}
	return cmd.Help()
}

// OpenSample returns the content of the named embed sample file given via an argument.
func OpenSample(cmd *cobra.Command, arg string, c *convert.Convert, f sample.Flags) ([]byte, error) {
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
		Build.FontFamily.Value = p.Font.String()
	}
	return []byte(string(p.Runes)), nil
}

// ReadArgument returns the content of argument supplied filepath, embed sample file or piped data.
func ReadArgument(arg string, cmd *cobra.Command, c *convert.Convert, f sample.Flags) ([]byte, error) {
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
	if b, err = OpenSample(cmd, arg, c, f); err != nil {
		return nil, err
	} else if b != nil {
		return b, nil
	}
	// the arg should be a filepath
	b, err = filesystem.Read(arg)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// SAUCE returns any SAUCE metadata that is attached to src.
func SAUCE(src *[]byte) create.SAUCE {
	sr := sauce.Decode(*src)
	if !sr.Valid() {
		return create.SAUCE{}
	}
	return create.SAUCE{
		Use:         true,
		Title:       sr.Title,
		Author:      sr.Author,
		Group:       sr.Group,
		Description: sr.Desc,
		Width:       uint(sr.Info.Info1.Value),
		Lines:       uint(sr.Info.Info2.Value),
	}
}
