package infocmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/info"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/sample"
	"github.com/bengarrett/retrotxtgo/static"
	"github.com/spf13/cobra"
)

var (
	ErrInfo  = errors.New("could not any obtain information")
	ErrUsage = errors.New("command usage could not display")
)

func Run(cmd *cobra.Command, args []string) {
	// piped input from other programs
	if filesystem.IsPipe() {
		Pipe()
	}
	if err := flag.PrintUsage(cmd, args...); err != nil {
		logs.Fatal(err)
	}
	var n info.Names
	n.Length = len(args)
	for i, arg := range args {
		n.Index = i + 1
		// embed sample filename
		filename, err := Sample(arg)
		if err != nil {
			fmt.Println(logs.SprintMark(arg, ErrInfo, err))
			continue
		}
		if filename != "" {
			defer os.Remove(filename)
			arg = filename
		}
		if err := n.Info(arg, flag.InfoFlag.Format); err != nil {
			if errors.As(logs.ErrFileName, &err) {
				if n.Length <= 1 {
					logs.Fatal(err)
				}
				fmt.Println(logs.SprintMark(arg, logs.ErrFileName, err))
				continue
			}
			if err = cmd.Usage(); err != nil {
				fmt.Println(logs.SprintWrap(ErrUsage, err))
			}
			logs.Fatal(err)
		}
	}
}

// infoSample extracts and saves an embed sample file then returns its location.
func Sample(name string) (filename string, err error) {
	s := strings.ToLower(name)
	if _, err = os.Stat(s); !os.IsNotExist(err) {
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

// infoPipe parses a standard input (stdin) stream of data.
func Pipe() {
	b, err := filesystem.ReadPipe()
	if err != nil {
		logs.FatalMark("info", logs.ErrPipeRead, err)
	}
	if err = info.Stdin(flag.InfoFlag.Format, b...); err != nil {
		logs.FatalMark("info", logs.ErrPipeParse, err)
	}
	os.Exit(0)
}
