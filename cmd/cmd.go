// Package cmd handles the terminal interface, user flags and arguments.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/cmd/example"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/cmd/version"
	"github.com/bengarrett/retrotxtgo/logs"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/cobra"
)

var (
	ErrHide  = errors.New("could not hide the flag")
	ErrUsage = errors.New("command usage could not display")
)

// ID are the cobra command group IDs.
const (
	IDcodepage = "idcp"     // codepage group
	IDfile     = "idfile"   // file group
	IDsample   = "idsample" // sample group
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
		RunE: func(cmd *cobra.Command, _ []string) error {
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
	// disable the default "completion" command.
	Cmd.CompletionOptions.DisableDefaultCmd = true
	// hide the cobra introduced "help" command.
	// https://github.com/spf13/cobra/issues/587#issuecomment-810159087
	Cmd.SetHelpCommand(&cobra.Command{Hidden: true})
	// the help "command" is hidden, so it needs to be assigned to a
	// group otherwise it will display an empty "Additional Commands:".
	Cmd.SetHelpCommandGroupID(IDcodepage)
	// hide the cobra errors.
	Cmd.SilenceErrors = true // set to false to debug errors
	// build the version flag template.
	Cmd.Version = meta.String()
	s := strings.Builder{}
	if err := version.Template(&s); err != nil {
		return err
	}
	Cmd.SetVersionTemplate(s.String())
	if errE := Cmd.Execute(); errE != nil {
		const minArgs = 2
		if len(os.Args) < minArgs {
			if errU := Cmd.Usage(); errU != nil {
				logs.FatalS(ErrUsage, errU, "rootCmd")
			}
		}
		args := strings.Join(os.Args[1:], " ")
		return fmt.Errorf("%w: %s %s", errE, meta.Bin, strings.TrimSpace(args))
	}
	return nil
}

// Tester creates and hides a custom tester flag.
// It is its own function so it can also be applied to unit tests as well as init.
func Tester(c *cobra.Command) *cobra.Command {
	c.PersistentFlags().BoolVar(&flag.Cmd.Tester, "tester", false,
		"optional in-memory, tester config file")
	if err := c.PersistentFlags().MarkHidden("tester"); err != nil {
		logs.FatalS(ErrHide, err, "tester")
	}
	return c
}

func init() {
	Cmd = Tester(Cmd)
	c := &cobra.Group{ID: IDcodepage, Title: "Codepage:"}
	f := &cobra.Group{ID: IDfile, Title: "File:"}
	s := &cobra.Group{ID: IDsample, Title: "Sample:"}
	Cmd.AddGroup(c, f, s)
	// create a version flag that only works on root.
	Cmd.LocalNonPersistentFlags().BoolP("version", "v", false, "")
}
