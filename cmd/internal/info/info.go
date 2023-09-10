// Package info provides the info command run function.
package info

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"github.com/bengarrett/retrotxtgo/pkg/info"
	"github.com/bengarrett/retrotxtgo/pkg/sample"
	"github.com/bengarrett/retrotxtgo/static"
	"github.com/spf13/cobra"
)

var (
	ErrNotExist  = errors.New("no such file or directory")
	ErrNotSamp   = errors.New("no such example or sample file")
	ErrInfo      = errors.New("could not any obtain information")
	ErrUsage     = errors.New("command usage could not display")
	ErrPipeRead  = errors.New("could not read text stream from piped stdin (standard input)")
	ErrPipeParse = errors.New("could not parse the text stream from piped stdin (standard input)")
	ErrTmpClose  = errors.New("could not close the temporary file")
	ErrTmpOpen   = errors.New("could not open the temporary file")
	ErrTmpSave   = errors.New("could not save to the temporary file")
)

// Run parses the arguments supplied with the info command.
func Run(cmd *cobra.Command, args []string) error {
	// piped input from other programs and then exit
	ok, err := fsys.IsPipe()
	if err != nil {
		return err
	}
	if ok {
		return Pipe(cmd)
	}
	if err := flag.Help(cmd, args...); err != nil {
		return err
	}
	var n info.Names
	n.Length = len(args)
	for i, arg := range args {
		n.Index = i + 1
		_, err := os.Stat(arg)
		if os.IsNotExist(err) {
			// embed sample filename
			filename, err := Sample(arg)
			if errors.Is(err, ErrNotSamp) {
				return fmt.Errorf("%w, %w: %s", ErrInfo, err, arg)
			}
			if filename == "" {
				return ErrNotExist
			}
			defer os.Remove(filename)
			arg = filename
		}
		s, err := n.Info(arg, flag.Info.Format)
		if err != nil {
			if err := cmd.Usage(); err != nil {
				return fmt.Errorf("%w: %w", ErrUsage, err)
			}
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), s)
	}
	return nil
}

// Sample extracts and saves an embed sample file then returns its location.
func Sample(name string) (string, error) {
	s := strings.ToLower(name)
	samp, exist := sample.Map()[s]
	if !exist {
		return "", ErrNotSamp
	}
	b, err := static.File.ReadFile(samp.Name)
	if err != nil {
		return "", fmt.Errorf(" sample file %q: %w", samp.Name, err)
	}
	file, err := os.CreateTemp("", fmt.Sprintf("retrotxt_%s.*.txt", s))
	if err != nil {
		return "", fmt.Errorf(" sample file %q: %w", samp.Name, ErrTmpOpen)
	}
	if _, err = file.Write(b); err != nil {
		return "", fmt.Errorf(" sample file %q: %w", samp.Name, ErrTmpSave)
	}
	if err := file.Close(); err != nil {
		return "", fmt.Errorf(" sample file %q: %w", samp.Name, ErrTmpClose)
	}
	return file.Name(), nil
}

// Pipe parses a standard input (stdin) stream of data.
func Pipe(cmd *cobra.Command) error {
	b, err := fsys.ReadPipe()
	if err != nil {
		return fmt.Errorf("%w, %w", ErrPipeRead, err)
	}
	s, err := info.Stdin(flag.Info.Format, b...)
	if err != nil {
		return fmt.Errorf("%w, %w", ErrPipeParse, err)
	}
	fmt.Fprint(cmd.OutOrStdout(), s)
	return nil
}
