package dump

import (
	"encoding/hex"
	"fmt"
	"io"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/pkg/term"
	"github.com/spf13/cobra"
)

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
