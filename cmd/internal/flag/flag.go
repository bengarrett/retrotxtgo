// Package flag provides the command flags handlers.
package flag

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bengarrett/retrotxtgo/convert"
	"github.com/bengarrett/retrotxtgo/fsys"
	"github.com/bengarrett/retrotxtgo/logs"
	"github.com/bengarrett/retrotxtgo/sample"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

var ErrNames = errors.New("ignoring [filenames]")

// Args initializes the command arguments and flags.
func Args(cmd *cobra.Command, args ...string) (
	[]string, *convert.Convert, sample.Flags, error,
) {
	reset := []string{}
	conv := convert.Convert{}
	conv.Args = convert.Flag{
		Controls:  View().Controls,
		SwapChars: View().Swap,
		MaxWidth:  View().Width,
	}
	conv.Args = setFlags(cmd, conv.Args)
	ok, err := fsys.IsPipe()
	if err != nil {
		logs.Fatal(err)
	}
	if ok {
		conv.Input.Encoding = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
		args = reset
	}
	if !ok {
		if err := Help(cmd, args...); err != nil {
			logs.Fatal(err)
		}
	}
	l := len(args)
	if ok && l > 0 {
		err := fmt.Errorf("%w;%w for piped text", err, ErrNames)
		fmt.Fprintln(os.Stderr, logs.Sprint(err))
		args = reset
	}
	if l == 0 {
		args = reset
	}
	samp, err := InputOriginal(cmd, "")
	if err != nil {
		return nil, nil, samp, err
	}
	if conv.Input.Encoding == nil {
		conv.Input.Encoding = Default()
	}
	return args, &conv, samp, nil
}

// setFlags applies the flag arguments to a convert flag struct.
func setFlags(cmd *cobra.Command, flag convert.Flag) convert.Flag {
	if c := cmd.Flags().Lookup("controls"); c != nil && c.Changed {
		const sep, minChrs = ",", 2
		val := c.Value.String()
		if len(val) > minChrs {
			val = val[1 : len(val)-1]
		}
		ctrls := strings.Split(val, sep)
		flag.Controls = ctrls
	}
	if s := cmd.Flags().Lookup("swap-chars"); s != nil && s.Changed {
		const sep, minChrs = ",", 2
		val := s.Value.String()
		if len(val) > minChrs {
			val = val[1 : len(val)-1]
		}
		swaps := strings.Split(val, sep)
		flag.SwapChars = swaps
	}
	if w := cmd.Flags().Lookup("width"); w != nil && w.Changed {
		i, err := strconv.Atoi(w.Value.String())
		if err != nil {
			logs.Fatal(err)
		}
		flag.MaxWidth = i
	}
	return flag
}

// Default returns a default encoding when the "input" flag is unused.
// If the input is a pipe, then the default encoding is UTF-16.
// Otherwise, the default encoding is CodePage437.
func Default() encoding.Encoding {
	ok, err := fsys.IsPipe()
	if err != nil {
		logs.Fatal(err)
	}
	if ok {
		return unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
	}
	return charmap.CodePage437
}

// InputOriginal applies the "input" and the (hidden) "original" encoding flag values
// to the sample data.
func InputOriginal(cmd *cobra.Command, dfault string) (sample.Flags, error) {
	parse := func(name string) (encoding.Encoding, error) {
		cp := cmd.Flags().Lookup(name)
		lookup := dfault
		if cp != nil && cp.Changed {
			lookup = cp.Value.String()
		}
		if dfault == "" || lookup == "" {
			return nil, nil
		}
		return convert.Encoder(lookup)
	}
	if cmd == nil {
		return sample.Flags{}, nil
	}
	// handle encode flag or apply the default
	in, err := parse("input")
	if err != nil {
		return sample.Flags{}, err
	}
	// handle the hidden original flag
	og := false
	l := cmd.Flags().Lookup("original")
	if l != nil && l.Changed {
		og = true
	}
	return sample.Flags{Input: in, Original: og}, nil
}

// EndOfFile reports whether end-of-file control flag was requested.
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
func OpenSample(name string, c *convert.Convert, f sample.Flags) ([]byte, error) {
	if ok := sample.Valid(name); !ok {
		return nil, nil
	}
	// return the sample with the original encoding
	if f.Original {
		p, err := sample.Open(name)
		if err != nil {
			return nil, err
		}
		return p, nil
	}
	// return the sample with utf-8 encoding
	r, err := f.Open(c, name)
	if err != nil {
		return nil, err
	}
	return []byte(string(r)), nil
}

// ReadArgument returns the content of argument supplied filepath, embed sample file or piped data.
func ReadArgument(arg string, c *convert.Convert, f sample.Flags) ([]byte, error) {
	// attempt to see if arg is a embed sample file request
	b, err := OpenSample(arg, c, f)
	if err != nil {
		return nil, err
	} else if b != nil {
		return b, nil
	}
	// the arg should be a filepath
	b, err = fsys.Read(arg)
	if err != nil {
		return nil, err
	}
	return b, nil
}
