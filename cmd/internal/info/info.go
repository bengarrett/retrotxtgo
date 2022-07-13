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
	ErrInfo     = errors.New("could not any obtain information")
	ErrUsage    = errors.New("command usage could not display")
)

func Run(cmd *cobra.Command, args []string) error {
	// piped input from other programs
	if filesystem.IsPipe() {
		Pipe(cmd)
	}
	if err := flag.PrintUsage(cmd, args...); err != nil {
		return err
	}
	var n info.Names
	n.Length = len(args)
	for i, arg := range args {
		if _, err := os.Stat(arg); os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrNotExist, arg)
		} else if err != nil {
			return err
		}
		n.Index = i + 1
		// embed sample filename
		filename, err := Sample(arg)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStdout(), logs.SprintMark(arg, ErrInfo, err))
			continue
		}
		if filename != "" {
			defer os.Remove(filename)
			arg = filename
		}
		s, err := n.Info(arg, flag.InfoFlag.Format)
		if err != nil {
			if errors.Is(logs.ErrFileName, err) {
				if n.Length <= 1 {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), logs.SprintMark(arg, logs.ErrFileName, err))
				continue
			}
			if err = cmd.Usage(); err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), logs.SprintWrap(ErrUsage, err))
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
	if _, err := os.Stat(s); !os.IsNotExist(err) {
		return "", nil
	}
	samp, exist := sample.Map()[s]
	if !exist {
		return "", nil
	}
	b, err := static.File.ReadFile(samp.Name)
	if err != nil {
		return "", fmt.Errorf("view package %q: %w", samp.Name, err)
	}
	file, err := ioutil.TempFile("", fmt.Sprintf("retrotxt_%s.*.txt", s))
	if err != nil {
		return "", fmt.Errorf("view package %q: %w", samp.Name, logs.ErrTmpOpen)
	}
	if _, err = file.Write(b); err != nil {
		return "", fmt.Errorf("view package %q: %w", samp.Name, logs.ErrTmpSave)
	}
	if err := file.Close(); err != nil {
		return "", fmt.Errorf("view package %q: %w", samp.Name, logs.ErrTmpClose)
	}
	return file.Name(), nil
}

// Pipe parses a standard input (stdin) stream of data.
func Pipe(cmd *cobra.Command) {
	b, err := filesystem.ReadPipe()
	if err != nil {
		logs.FatalMark("info", logs.ErrPipeRead, err)
	}
	s, err := info.Stdin(flag.InfoFlag.Format, b...)
	if err != nil {
		logs.FatalMark("info", logs.ErrPipeParse, err)
	}
	fmt.Fprint(cmd.OutOrStdout(), s)
	os.Exit(0)
}
