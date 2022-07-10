package flag

import "github.com/spf13/cobra"

// Controls handles the --controls flag.
func Controls(p *[]string, cc *cobra.Command) {
	cc.Flags().StringSliceVarP(p, "controls", "c", []string{},
		`implement these control codes (default "eof,tab")
separate multiple controls with commas
  eof    end of file mark
  tab    horizontal tab
  bell   bell or terminal alert
  cr     carriage return
  lf     line feed
  bs backspace, del delete character, esc escape character
  ff formfeed, vt vertical tab
`)
}
