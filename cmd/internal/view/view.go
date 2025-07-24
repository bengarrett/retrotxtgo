// Package view provides the view command run function.
package view

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/convert"
	"github.com/bengarrett/retrotxtgo/fsys"
	"github.com/bengarrett/retrotxtgo/term"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding"
)

var (
	ErrConv     = errors.New("convert cannot be nil")
	ErrPipeRead = errors.New("could not read text stream from piped stdin (standard input)")
)

// Run parses the arguments supplied with the view command.
func Run(w io.Writer, cmd *cobra.Command, args ...string) error {
	if w == nil {
		w = io.Discard
	}
	// piped input from other programs and then exit
	ok, err := fsys.IsPipe()
	if err != nil {
		return err
	}
	if ok {
		return Pipe(w, cmd, args...)
	}
	// read from files or samples
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
		// write out the sample with its original encoding
		// this could display poorly in a terminal
		if samp.Original {
			fmt.Fprint(w, string(b))
			continue
		}
		// write out the sample with the utf-8 encoding
		r, err := Transform(c, samp.Input, nil, b...)
		if err != nil {
			return err
		}
		fmt.Fprint(w, string(r))
	}
	fmt.Fprintln(w)
	return nil
}

// Pipe parses a standard input (stdin) stream of data.
func Pipe(w io.Writer, cmd *cobra.Command, args ...string) error {
	if w == nil {
		w = io.Discard
	}
	_, c, samp, err := flag.Args(cmd, args...)
	if err != nil {
		return err
	}
	b, err := fsys.ReadPipe()
	if err != nil {
		return fmt.Errorf("%w, %w", ErrPipeRead, err)
	}
	// write out the sample with the utf-8 encoding
	r, err := Transform(c, samp.Input, nil, b...)
	if err != nil {
		return err
	}
	fmt.Fprint(w, string(r))
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
	// handle any encoding BEFORE outputting to Unicode
	// we also make sure the bytes are not valid UTF-8
	// otherwise the bytes will become corrupted
	if out != nil && !utf8.Valid(b) {
		p, err = out.NewDecoder().Bytes(b)
		if err != nil {
			return nil, err
		}
	}
	// convert the bytes into runes
	if flag.EndOfFile(c.Args) {
		return c.Text(p...)
	}
	return c.Dump(p...)
}
