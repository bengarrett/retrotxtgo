package usage

import (
	"os"

	"github.com/spf13/cobra"
)

// Print the usage help and exit; but only when no arguments are given.
func Print(cmd *cobra.Command, args ...string) error {
	if len(args) == 0 {
		if err := cmd.Help(); err != nil {
			return err
		}
		os.Exit(0)
	}
	return nil
}
