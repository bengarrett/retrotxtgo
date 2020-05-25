package cmd

import (
	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/lib/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var versionFmt string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"ver", "v"},
	Example: "  retrotxt version --format=text",
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
		if ok := version.Print(viper.GetString("format.version")); !ok {
			logs.CheckFlag("format", versionFmt, config.Format.Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringVarP(&versionFmt, "format", "f",
		viper.GetString("format.version"),
		str.Options("output format", config.Format.Version, true))
	err := viper.BindPFlag("format.version", versionCmd.Flags().Lookup("format"))
	logs.Check("format.version", err)
}
