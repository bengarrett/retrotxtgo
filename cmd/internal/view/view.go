package view

import (
	"bytes"
	"fmt"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding"
)

// Run parses the arguments supplied with the view command.
func Run(cmd *cobra.Command, args ...string) (*bytes.Buffer, error) {
	args, conv, samp, err := flag.InitArgs(cmd, args...)
	if err != nil {
		return nil, err
	}
	w := new(bytes.Buffer)
	for i, arg := range args {
		if i > 0 && i < len(arg) {
			const halfPage = 40
			fmt.Fprintln(w, str.HRPad(halfPage))
		}
		b, err := flag.ReadArg(arg, cmd, conv, samp)
		if err != nil {
			fmt.Fprintln(w, logs.Sprint(err))
			continue
		}
		r, err := Transform(samp.From, samp.To, conv, b...)
		if err != nil {
			fmt.Fprintln(w, logs.Sprint(err))
			continue
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
	b ...byte) ([]rune, error) {
	// handle input source encoding
	if in != nil {
		conv.Input.Encoding = in
	}
	var (
		r   []rune
		err error
	)
	// handle any output encoding BEFORE converting to Unicode
	if out != nil {
		b, err = out.NewDecoder().Bytes(b)
		if err != nil {
			return nil, err
		}
	}
	// convert text to runes
	if flag.EndOfFile(conv.Flags) {
		r, err = conv.Text(b...)
	} else {
		r, err = conv.Dump(b...)
	}
	if err != nil {
		return nil, err
	}
	return r, nil
}
