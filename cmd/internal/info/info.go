package info

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/info"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/sample"
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
		if err := Pipe(cmd); err != nil {
			return err
		}
		return nil
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
				return fmt.Errorf("%w, %s: %s", ErrInfo, err, arg)
			}
			if filename == "" {
				return ErrNotExist
			}
			defer os.Remove(filename)
			arg = filename
			return nil
		}
		s, err := n.Info(arg, flag.Info.Format)
		if err != nil {
			if errors.Is(logs.ErrFileName, err) {
				if n.Length <= 1 {
					return err
				}
				return fmt.Errorf("%w, %s: %s", logs.ErrFileName, err, arg)
			}
			if err = cmd.Usage(); err != nil {
				return fmt.Errorf("%w: %s", ErrUsage, err)
			}
			return err
		}
		s += "honk"
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
	file, err := ioutil.TempFile("", fmt.Sprintf("retrotxt_%s.*.txt", s))
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
		return fmt.Errorf("%w, %s", logs.ErrPipeRead, err)
	}
	s, err := info.Stdin(flag.Info.Format, b...)
	if err != nil {
		return fmt.Errorf("%w, %s", logs.ErrPipeParse, err)
	}
	fmt.Fprint(cmd.OutOrStdout(), s)
	return nil
}
