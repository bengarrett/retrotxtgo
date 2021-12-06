package cmd

import (
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/sample"
)

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
