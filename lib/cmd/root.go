// Package cmd handles the terminal interface, user flags and arguments.
// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/long"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/rootcmd"
	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/ver"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/cobra"
)

// const (
// 	eof         = "eof"
// 	tab         = "tab"
// 	null        = "null" // 0
// 	verticalBar = "bar"  // 124
// )

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:     meta.Bin,
	Short:   fmt.Sprintf("%s is the tool that turns ANSI, ASCII, NFO text into browser ready HTML", meta.Name),
	Long:    long.Root,
	Example: example.Print(example.Root),
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing other than print the help.
		// This func must remain otherwise root command flags are ignored by Cobra.
		if err := flag.PrintUsage(cmd); err != nil {
			logs.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SilenceErrors = true // set to false to debug errors
	rootCmd.Version = meta.Print()
	rootCmd.SetVersionTemplate(ver.Print())
	if err := rootCmd.Execute(); err != nil {
		const minArgs = 2
		if len(os.Args) < minArgs {
			if err1 := rootCmd.Usage(); err1 != nil {
				logs.FatalMark("rootCmd", ErrUsage, err1)
			}
		}
		logs.FatalExecute(err, os.Args[1:]...)
	}
}

func init() {
	cobra.OnInitialize(rootcmd.Init)
	// create and hide custom configuration file location flag.
	rootCmd.PersistentFlags().StringVar(&flag.RootFlag.Config, "config", "",
		"optional config file location")
	if err := rootCmd.PersistentFlags().MarkHidden("config"); err != nil {
		logs.FatalMark("config", ErrHide, err)
	}
	// create a version flag that only works on root.
	rootCmd.LocalNonPersistentFlags().BoolP("version", "v", false, "")
	// hide the cobra introduced help command.
	// https://github.com/spf13/cobra/issues/587#issuecomment-810159087
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}
