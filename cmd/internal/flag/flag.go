// Package flag provides the command flags handlers.
package flag

import (
	"errors"
	"fmt"
	"os"
	"slices"
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

var (
	ErrInput = errors.New("empty default encoding")
	ErrNames = errors.New("ignoring [filenames]")
)

// Args initializes the command arguments and flags.
func Args(cmd *cobra.Command, args ...string) (
	[]string, *convert.Convert, sample.Flags, error,
) {
	reset := []string{}
	converter := convert.Convert{}
	converter.Args = convert.Flag{
		Controls:  View().Controls,
		SwapChars: View().Swap,
		MaxWidth:  View().Width,
	}
	converter.Args = setFlags(cmd, converter.Args)
	pipeOW, err := fsys.IsPipe()
	if err != nil {
		logs.Fatal(err)
	}
	if pipeOW {
		converter.Input.Encoding = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
		args = reset
	} else {
		err := Help(cmd, args...)
		if err != nil {
			logs.Fatal(err)
		}
	}
	count := len(args)
	if pipeOW && count > 0 {
		err := fmt.Errorf("%w;%w for piped text", err, ErrNames)
		fmt.Fprintln(os.Stderr, logs.Sprint(err))
		args = reset
	}
	if count == 0 {
		args = reset
	}
	inputsOW, err := InputOverride(cmd, "")
	if err != nil {
		return nil, nil, inputsOW, err
	}
	// despite this default cp437 encoding being set,
	// the results of the --input flag will be used
	// when provided in the inputsOW struct.
	if converter.Input.Encoding == nil {
		converter.Input.Encoding = Default()
	}
	return args, &converter, inputsOW, nil
}

// setFlags applies the flag arguments to a convert flag struct.
func setFlags(cmd *cobra.Command, flag convert.Flag) convert.Flag {
	const (
		controls  = "controls"
		swapChars = "swap-chars"
		width     = "width"
	)
	if c := cmd.Flags().Lookup(controls); c != nil && c.Changed {
		const sep, minChrs = ",", 2
		val := c.Value.String()
		if len(val) > minChrs {
			val = val[1 : len(val)-1]
		}
		ctrls := strings.Split(val, sep)
		flag.Controls = ctrls
	}
	if s := cmd.Flags().Lookup(swapChars); s != nil && s.Changed {
		const sep, minChrs = ",", 2
		val := s.Value.String()
		if len(val) > minChrs {
			val = val[1 : len(val)-1]
		}
		swaps := strings.Split(val, sep)
		flag.SwapChars = swaps
	}
	if w := cmd.Flags().Lookup(width); w != nil && w.Changed {
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
	pipeOW, err := fsys.IsPipe()
	if err != nil {
		logs.Fatal(err)
	}
	if pipeOW {
		return unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
	}
	return charmap.CodePage437
}

// InputOverride applies the "input" encoding flag and the "original" bool flag.
func InputOverride(cmd *cobra.Command, fallback string) (sample.Flags, error) {
	none := sample.Flags{}
	if cmd == nil {
		return none, nil
	}
	const input, original = "input", "original"
	// handle encode flag or apply the default
	cp := cmd.Flags().Lookup(input)
	in, err := parseInput(cp.Changed, cp.Value.String(), fallback)
	// parse("input")
	if err != nil {
		if errors.Is(err, ErrInput) {
			return none, nil
		}
		return none, err
	}
	// handle the hidden original bool flag
	og := false
	l := cmd.Flags().Lookup(original)
	if l != nil && l.Changed {
		og = true
	}
	return sample.Flags{Input: in, Original: og}, nil
}

func parseInput(changed bool, value, fallback string) (encoding.Encoding, error) {
	name := fallback
	if changed {
		name = value
	}
	// 11-Aug-25, this was previously bugged with an OR statement,
	// it must always be an AND condition otherwise the logic will break.
	if name == "" && fallback == "" {
		return nil, ErrInput
	}
	return convert.Encoder(name)
}

// EndOfFile reports whether end-of-file control flag was requested.
func EndOfFile(flags convert.Flag) bool {
	return slices.Contains(flags.Controls, "eof")
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
