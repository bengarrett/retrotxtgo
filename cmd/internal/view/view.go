// Package view provides the view command run function.
package view

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/pkg/convert"
	"github.com/bengarrett/retrotxtgo/pkg/term"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding"
)

var ErrConv = errors.New("convert cannot be nil")

// Run parses the arguments supplied with the view command.
func Run(w io.Writer, cmd *cobra.Command, args ...string) error {
	if w == nil {
		w = io.Discard
	}
	args, c, samp, err := flag.Args(cmd, args...)
	if err != nil {
		return err
	}
	for i, arg := range args {
		if i == 0 && arg == "" {
			return nil
		}
		if i > 0 && i < len(arg) {
			const halfPage = 40
			fmt.Fprintln(w)
			term.HR(w, halfPage)
		}
		b, err := flag.ReadArgument(arg, c, samp)
		if err != nil {
			return err
		}
		r, err := Transform(c, samp.Input, samp.Output, b...)
		if err != nil {
			return err
		}
		fmt.Fprint(w, string(r))
	}
	fmt.Fprintln(w)
	return nil
}

// Transform bytes into Unicode runes.
// The optional in encoding argument is the bytes original character encoding.
// The optional out encoding argument is the encoding to replicate.
// When no encoding arguments are provided, UTF-8 unicode encoding is used.
func Transform(c *convert.Convert, in, out encoding.Encoding, b ...byte,
) ([]rune, error) {
	if c == nil {
		return nil, ErrConv
	}
	var err error
	if b == nil {
		return nil, nil
	}
	// handle input source encoding
	if in != nil {
		c.Input.Encoding = in
	}
	p := b
	// handle any output re-encoding BEFORE converting to Unicode
	if out != nil {
		p, err = out.NewDecoder().Bytes(b)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(os.Stdout, "%s\n", p)
	}
	// convert the bytes into runes
	if flag.EndOfFile(c.Args) {
		return c.Text(p...)
	}
	return c.Dump(p...)
}
