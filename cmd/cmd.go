// Package cmd handles the terminal interface, user flags and arguments.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/cmd/internal/example"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/internal/version"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/cobra"
)

var (
	ErrHide  = errors.New("could not hide the flag")
	ErrUsage = errors.New("command usage could not display")
)

var cmdShort = fmt.Sprintf("%s is the tool that turns ANSI, ASCII, NFO text into browser ready HTML",
	meta.Name)

var cmdLong = fmt.Sprintf(`Turn many pieces of ANSI art, ASCII and NFO texts into HTML5 using %s.
It is the platform agnostic tool that takes nostalgic text files and stylises
them into a more modern, useful format to view or copy in a web browser.`, meta.Name)

// Cmd represents the base command when called without any subcommands.
//nolint:gochecknoglobals
var Cmd = &cobra.Command{
	Use:     meta.Bin,
	Short:   cmdShort,
	Long:    cmdLong,
	Example: fmt.Sprint(example.Cmd),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Do nothing other than print the help.
		// This func must remain otherwise root command flags are ignored by Cobra.
		if err := flag.PrintUsage(cmd); err != nil {
			return err
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// TODO: MAKE EXECUTE return an error.
func Execute() {
	Cmd.CompletionOptions.DisableDefaultCmd = true
	Cmd.SilenceErrors = true // set to false to debug errors
	Cmd.Version = meta.Print()
	Cmd.SetVersionTemplate(version.Print())
	if err := Cmd.Execute(); err != nil {
		const minArgs = 2
		if len(os.Args) < minArgs {
			if err1 := Cmd.Usage(); err1 != nil {
				logs.FatalMark("rootCmd", ErrUsage, err1)
			}
		}
		logs.FatalExecute(err, os.Args[1:]...)
	}
}

func CmdInit() {
	cobra.OnInitialize(Load)
	// create and hide custom configuration file location flag.
	Cmd.PersistentFlags().StringVar(&flag.RootFlag.Config, "config", "",
		"optional config file location")
	if err := Cmd.PersistentFlags().MarkHidden("config"); err != nil {
		logs.FatalMark("config", ErrHide, err)
	}
	// create a version flag that only works on root.
	Cmd.LocalNonPersistentFlags().BoolP("version", "v", false, "")
	// hide the cobra introduced help command.
	// https://github.com/spf13/cobra/issues/587#issuecomment-810159087
	Cmd.SetHelpCommand(&cobra.Command{Hidden: true})
}

//nolint:gochecknoinits
func init() {
	CmdInit()
}
