// Package cmd handles the terminal interface, user flags and arguments.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

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

const (
	IDcodepage = "idcp"
	IDfile     = "idfile"
	IDsample   = "idsample"
)

// Cmd represents the base command when called without any subcommands.
var Cmd = base()

func base() *cobra.Command {
	s := "Use " + meta.Name + " to print legacy text on modern terminals."
	l := "Text files and art created without Unicode often fail to display on modern systems. " +
		"\n" + s
	expl := strings.Builder{}
	example.Cmd.String(&expl)
	return &cobra.Command{
		Use:     meta.Bin,
		Short:   s,
		Long:    l,
		Example: expl.String(),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Do nothing other than print the help.
			// This func must remain otherwise root
			// command flags are ignored by Cobra.
			return flag.Help(cmd)
		},
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	Cmd.CompletionOptions.DisableDefaultCmd = true
	Cmd.SilenceErrors = true // set to false to debug errors
	Cmd.Version = meta.Print()
	s := strings.Builder{}
	if err := version.Template(&s); err != nil {
		return err
	}
	Cmd.SetVersionTemplate(s.String())
	if err := Cmd.Execute(); err != nil {
		const minArgs = 2
		if len(os.Args) < minArgs {
			if err1 := Cmd.Usage(); err1 != nil {
				logs.FatalS(ErrUsage, err1, "rootCmd")
			}
		}
		return fmt.Errorf("%w: %s", err, os.Args[1:])
	}
	return nil
}

func Init() {
	Cmd = Tester(Cmd)
	Cmd.AddGroup(&cobra.Group{ID: IDcodepage, Title: "Codepage:"})
	Cmd.AddGroup(&cobra.Group{ID: IDfile, Title: "File:"})
	Cmd.AddGroup(&cobra.Group{ID: IDsample, Title: "Sample:"})
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
		logs.FatalS(ErrHide, err, "tester")
	}
	return c
}

//nolint:gochecknoinits
func init() {
	Init()
}
