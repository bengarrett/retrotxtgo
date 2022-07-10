package viewcmd

import (
	"bytes"
	"fmt"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/sample"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/cobra"
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
		r, err := transform(conv, samp, b...)
		if err != nil {
			fmt.Fprintln(w, logs.Sprint(err))
			continue
		}
		fmt.Fprint(w, string(r))
	}
	return w, nil
}

// transform the bytes into Unicode runes.
func transform(conv *convert.Convert, f sample.Flags, b ...byte) ([]rune, error) {
	// handle input source encoding
	if f.From != nil {
		conv.Input.Encoding = f.From
	}
	var (
		r   []rune
		err error
	)
	// handle any output encoding BEFORE converting to Unicode
	if f.To != nil {
		b, err = f.To.NewDecoder().Bytes(b)
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
