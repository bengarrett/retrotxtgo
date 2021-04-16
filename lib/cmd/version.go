// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"github.com/spf13/cobra"

	"retrotxt.com/retrotxt/lib/config"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
	"retrotxt.com/retrotxt/lib/version"
)

type versionFlags struct {
	format string
}

var versionFlag = versionFlags{
	format: "color",
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"ver"},
	Example: "  retrotxt version --format=text\n  retrotxt ver -f t",
	Short:   "Version information for RetroTxt",
	Long: `Version information for Retrotxt

The shown ` + str.Cc("RetroTxt URL") + ` is the weblink to the application Github page.

` + str.Cc("Version") + ` number reflects ` + str.Ci("[major].[minor].[patch]") + `.
* Major is a generational iteration that may break backwards compatibility.
* Minor changes are increased whenever new features are added.
* Patch reflect hot fixes or bug corrections.

` + str.Cc("Go version") + ` reports the edition of Go used to build this application.
` + str.Cc("OS/Arch") + ` reports both the operating system and CPU architecture.

` + str.Cc("Binary") + ` should return the path name of this program. It maybe inaccurate
if it is launched through an operating system symlink.`,
	Run: func(cmd *cobra.Command, args []string) {
		if ok := version.Print(versionFlag.format); !ok {
			v := config.Format().Version
			logs.FlagFatal("format", versionFlag.format, v[:]...)
		}
	},
}

func init() {
	// cmds and flags
	rootCmd.AddCommand(versionCmd)
	v := config.Format().Version
	versionCmd.Flags().StringVarP(&versionFlag.format, "format", "f", versionFlag.format,
		str.Options("output format", true, v[:]...))
}
