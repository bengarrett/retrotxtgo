package dump

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/bengarrett/retrotxtgo/fsys"
	"github.com/bengarrett/retrotxtgo/sample"
	"github.com/bengarrett/retrotxtgo/term"
	"github.com/spf13/cobra"
)

var ErrPipeRead = errors.New("could not read text stream from piped stdin (standard input)")

// Run parses the arguments supplied with the dump command.
func Run(w io.Writer, _ *cobra.Command, args ...string) error {
	if w == nil {
		w = io.Discard
	}
	// piped input from other programs and then exit
	ok, err := fsys.IsPipe()
	if err != nil {
		return err
	}
	if ok {
		return Pipe(w)
	}
	// read from files or samples
	for i, arg := range args {
		if i == 0 && arg == "" {
			return nil
		}
		if i > 0 && i < len(arg) {
			const page = 76
			term.HR(w, page)
		}
		// Try to read as sample first, then as regular file
		b, err := tryReadSample(arg)
		if err == nil && b != nil {
			fmt.Fprint(w, hex.Dump(b))
			continue
		}
		// Read as regular file
		b, err = fsys.Read(arg)
		if err != nil {
			return err
		}
		fmt.Fprint(w, hex.Dump(b))
	}
	return nil
}

// tryReadSample attempts to read a sample file if it exists.
func tryReadSample(name string) ([]byte, error) {
	if ok := sample.Valid(name); !ok {
		return nil, nil
	}
	p, err := sample.Open(name)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// Pipe parses a standard input (stdin) stream of data.
func Pipe(w io.Writer) error {
	if w == nil {
		w = io.Discard
	}
	data, err := fsys.ReadPipe()
	if err != nil {
		return fmt.Errorf("%w, %w", ErrPipeRead, err)
	}
	fmt.Fprint(w, hex.Dump(data))
	return nil
}
