// Package cmd handles the terminal interface, user flags and arguments.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/pkg/example"
	"github.com/bengarrett/retrotxtgo/cmd/pkg/version"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/pkg/logs"
	"github.com/spf13/cobra"
)

var (
	ErrHide  = errors.New("could not hide the flag")
	ErrUsage = errors.New("command usage could not display")
)

// Cmd represents the base command when called without any subcommands.
//
//nolint:gochecknoglobals
var Cmd = &cobra.Command{
	Use:   meta.Bin,
	Short: fmt.Sprintf("Use %s to print text, BBS and ANSI files", meta.Name),
	Long: fmt.Sprintf(`%s takes legacy encoded text, BBS, and ANSI files
	and print them to a modern UTF-8 terminal.`, meta.Name),
	Example: fmt.Sprint(example.Cmd),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Do nothing other than print the help.
		// This func must remain otherwise root command flags are ignored by Cobra.
		return flag.Help(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// TODO: MAKE EXECUTE return an error.
func Execute() {
	Cmd.CompletionOptions.DisableDefaultCmd = true
	Cmd.SilenceErrors = true // set to false to debug errors
	Cmd.Version = meta.Print()
	Cmd.SetVersionTemplate(version.Template())
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

func Init() {
	Cmd = Tester(Cmd)

	Cmd.AddGroup(&cobra.Group{ID: "listCmds", Title: "Codepages:"})
	Cmd.AddGroup(&cobra.Group{ID: "fileCmds", Title: "Files:"})
	Cmd.AddGroup(&cobra.Group{ID: "exaCmds", Title: "Testers:"})

	// create a version flag that only works on root.
	Cmd.LocalNonPersistentFlags().BoolP("version", "v", false, "")
	// hide the cobra introduced help command.
	// https://github.com/spf13/cobra/issues/587#issuecomment-810159087
	Cmd.SetHelpCommand(&cobra.Command{Hidden: true})
}

// Tester creates and hides a custom tester flag.
// It is its own function so it can also be applied to unit tests as well as init.
func Tester(c *cobra.Command) *cobra.Command {
	c.PersistentFlags().BoolVar(&flag.Command.Tester, "tester", false,
		"optional in-memory, tester config file")
	if err := c.PersistentFlags().MarkHidden("tester"); err != nil {
		logs.FatalMark("tester", ErrHide, err)
	}
	return c
}

//nolint:gochecknoinits
func init() {
	Init()
}
