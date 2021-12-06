package usage

import (
	"os"

	"github.com/spf13/cobra"
)

// printUsage will print the help and exit when no arguments are supplied.
func Print(cmd *cobra.Command, args ...string) error {
	if len(args) == 0 {
		if err := cmd.Help(); err != nil {
			return err
		}
		os.Exit(0)
	}
	return nil
}
