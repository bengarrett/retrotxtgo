package flag

import (
	"errors"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/cobra"
)

var ErrHide = errors.New("could not hide the flag")

// To handles the hidden --to flag.
func To(p *string, cc *cobra.Command) {
	const name = "to"
	cc.Flags().StringVar(p, name, "",
		"alternative character encoding to print to stdout\nthis flag is unreliable and not recommended")
	if err := cc.Flags().MarkHidden(name); err != nil {
		logs.FatalMark(name, ErrHide, err)
	}
}
