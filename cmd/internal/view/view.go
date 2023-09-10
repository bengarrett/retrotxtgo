// Package view provides the view command run function.
package view

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/pkg/convert"
	"github.com/bengarrett/retrotxtgo/pkg/term"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding"
)

var ErrConv = errors.New("convert cannot be nil")

// Run parses the arguments supplied with the view command.
func Run(cmd *cobra.Command, args ...string) (*bytes.Buffer, error) {
	args, conv, samp, err := flag.Args(cmd, args...)
	if err != nil {
		return nil, err
	}
	w := new(bytes.Buffer)
	for i, arg := range args {
		if i == 0 && arg == "" {
			return w, nil
		}
		if i > 0 && i < len(arg) {
			const halfPage = 40
			fmt.Fprintln(w, term.HRLen(halfPage))
		}
		b, err := flag.ReadArgument(arg, conv, samp)
		if err != nil {
			return nil, err
		}
		r, err := Transform(samp.Input, samp.Output, conv, b...)
		if err != nil {
			return nil, err
		}
		fmt.Fprint(w, string(r))
	}
	return w, nil
}

// Transform bytes into Unicode runes.
// The optional in encoding argument is the bytes original character encoding.
// The optional out encoding argument is the encoding to replicate.
// When no encoding arguments are provided, UTF-8 unicode encoding is used.
func Transform(
	in encoding.Encoding,
	out encoding.Encoding,
	conv *convert.Convert,
	b ...byte,
) ([]rune, error) {
	if conv == nil {
		return nil, ErrConv
	}
	var err error
	if b == nil {
		return nil, nil
	}
	// handle input source encoding
	if in != nil {
		conv.Input.Encoding = in
	}
	b2 := b
	// handle any output re-encoding BEFORE converting to Unicode
	if out != nil {
		b2, err = out.NewDecoder().Bytes(b)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(os.Stdout, "%s\n", b2)
	}
	// convert the bytes into runes
	if flag.EndOfFile(conv.Flags) {
		return conv.Text(b2...)
	}
	return conv.Dump(b2...)
}
