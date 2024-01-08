package dump

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/fsys"
	"github.com/bengarrett/retrotxtgo/term"
	"github.com/spf13/cobra"
)

var ErrPipeRead = errors.New("could not read text stream from piped stdin (standard input)")

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
		return Pipe(w)
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
			const page = 76
			term.HR(w, page)
		}
		b, err := flag.ReadArgument(arg, c, samp)
		if err != nil {
			return err
		}
		fmt.Fprint(w, hex.Dump(b))
	}
	return nil
}

// Pipe parses a standard input (stdin) stream of data.
func Pipe(w io.Writer) error {
	if w == nil {
		w = io.Discard
	}
	b, err := fsys.ReadPipe()
	if err != nil {
		return fmt.Errorf("%w, %w", ErrPipeRead, err)
	}
	fmt.Fprint(w, hex.Dump(b))
	return nil
}
