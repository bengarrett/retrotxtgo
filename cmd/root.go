// Package cmd handles the terminal interface, user flags and arguments.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/internal/long"
	"github.com/bengarrett/retrotxtgo/cmd/internal/root"
	"github.com/bengarrett/retrotxtgo/cmd/internal/ver"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/cobra"
)

var (
	ErrHide  = errors.New("could not hide the flag")
	ErrUsage = errors.New("command usage could not display")
)

// rootCmd represents the base command when called without any subcommands.
//nolint:gochecknoglobals
var rootCmd = &cobra.Command{
	Use:     meta.Bin,
	Short:   fmt.Sprintf("%s is the tool that turns ANSI, ASCII, NFO text into browser ready HTML", meta.Name),
	Long:    long.Root.String(),
	Example: example.Root.Print(),
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

//nolint:gochecknoinits
func init() {
	cobra.OnInitialize(root.Init)
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
