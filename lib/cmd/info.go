// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/info"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/sample"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/static"
	"github.com/spf13/cobra"
)

var infoFlag struct {
	format string
}

var infoExample = fmt.Sprintf("  %s %s\n%s %s",
	meta.Bin, "info text.asc logo.jpg # print the information of multiple files",
	meta.Bin, "info file.txt --format=json # print the information using a structured syntax")

// infoCmd represents the info command.
var infoCmd = &cobra.Command{
	Use:     fmt.Sprintf("info %s", filenames),
	Aliases: []string{"i"},
	Short:   "Information on a text file",
	Long:    "Discover details and information about any text or text art file.",
	Example: exampleCmd(infoExample),
	Run: func(cmd *cobra.Command, args []string) {
		// piped input from other programs
		if filesystem.IsPipe() {
			infoPipe()
		}
		printUsage(cmd, args...)
		var n info.Names
		n.Length = len(args)
		for i, arg := range args {
			n.Index = i + 1
			// embed sample filename
			filename, err := infoSample(arg)
			if err != nil {
				fmt.Println(logs.PrintfMark(arg, ErrInfo, err))
				continue
			}
			if filename != "" {
				defer os.Remove(filename)
				arg = filename
			}
			if err := n.Info(arg, infoFlag.format); err != nil {
				if errors.As(logs.ErrFileName, &err) {
					if n.Length <= 1 {
						logs.Fatal(err)
					}
					fmt.Println(logs.PrintfMark(arg, logs.ErrFileName, err))
					continue
				}
				if err = cmd.Usage(); err != nil {
					fmt.Println(logs.PrintfWrap(ErrUsage, err))
				}
				logs.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	i := config.Format().Info
	infoCmd.Flags().StringVarP(&infoFlag.format, "format", "f", "color",
		str.Options("print format or syntax", true, true, i[:]...))
}

// infoSample extracts and saves an embed sample file then returns its location.
func infoSample(name string) (filename string, err error) {
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
func infoPipe() {
	b, err := filesystem.ReadPipe()
	if err != nil {
		logs.FatalMark("info", logs.ErrPipeRead, err)
	}
	if err = info.Stdin(infoFlag.format, b...); err != nil {
		logs.FatalMark("info", logs.ErrPipeParse, err)
	}
	os.Exit(0)
}
