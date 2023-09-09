package info

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/pkg/filesystem"
	"github.com/bengarrett/retrotxtgo/pkg/info"
	"github.com/bengarrett/retrotxtgo/pkg/logs"
	"github.com/bengarrett/retrotxtgo/pkg/sample"
	"github.com/bengarrett/retrotxtgo/static"
	"github.com/spf13/cobra"
)

var (
	ErrNotExist = errors.New("no such file or directory")
	ErrNotSamp  = errors.New("no such example or sample file")
	ErrInfo     = errors.New("could not any obtain information")
	ErrUsage    = errors.New("command usage could not display")
)

func Run(cmd *cobra.Command, args []string) error {
	// piped input from other programs and then exit
	if filesystem.IsPipe() {
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
			if errors.Is(logs.ErrFileName, err) {
				if n.Length <= 1 {
					return err
				}
				return fmt.Errorf("%w, %w: %s", logs.ErrFileName, err, arg)
			}
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
		return "", fmt.Errorf(" sample file %q: %w", samp.Name, logs.ErrTmpOpen)
	}
	if _, err = file.Write(b); err != nil {
		return "", fmt.Errorf(" sample file %q: %w", samp.Name, logs.ErrTmpSave)
	}
	if err := file.Close(); err != nil {
		return "", fmt.Errorf(" sample file %q: %w", samp.Name, logs.ErrTmpClose)
	}
	return file.Name(), nil
}

// Pipe parses a standard input (stdin) stream of data.
func Pipe(cmd *cobra.Command) error {
	b, err := filesystem.ReadPipe()
	if err != nil {
		return fmt.Errorf("%w, %w", logs.ErrPipeRead, err)
	}
	s, err := info.Stdin(flag.Info.Format, b...)
	if err != nil {
		return fmt.Errorf("%w, %w", logs.ErrPipeParse, err)
	}
	fmt.Fprint(cmd.OutOrStdout(), s)
	return nil
}
