package flag

import "github.com/spf13/cobra"

// Width handles the --width flag.
func Width(p *int, cc *cobra.Command) {
	cc.Flags().IntVarP(p, "width", "w", ViewFlag.Width,
		"maximum document character/column width")
}
